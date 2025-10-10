package users

import (
	"bytes"
	"fmt"
	"os/exec"
	"os/user"
	"strconv"
	"strings"
	"time"
)

// UserManager interface para gerenciamento de usuários
type UserManager interface {
	List(options ListOptions) ([]*UserInfo, error)
	Info(username string) (*UserDetail, error)
	Add(options AddUserOptions) error
	Remove(username string, removeHome bool) error
	Modify(username string, options ModifyOptions) error
	ListGroups() ([]*GroupInfo, error)
	AddToGroup(username, group string) error
	RemoveFromGroup(username, group string) error
}

// ListOptions opções para listar usuários
type ListOptions struct {
	SystemUsers bool   // incluir usuários de sistema
	Filter      string // filtro por nome
	Group       string // filtro por grupo
}

// UserInfo informações básicas do usuário
type UserInfo struct {
	Username string
	UID      string
	GID      string
	HomeDir  string
	Shell    string
	FullName string
}

// UserDetail informações detalhadas do usuário
type UserDetail struct {
	*UserInfo
	Groups        []string
	LastLogin     time.Time
	PasswordSet   bool
	Locked        bool
	ExpiryDate    string
	DaysSincePass int
}

// GroupInfo informações de grupo
type GroupInfo struct {
	Name    string
	GID     string
	Members []string
}

// AddUserOptions opções para adicionar usuário
type AddUserOptions struct {
	Username   string
	FullName   string
	HomeDir    string
	Shell      string
	Groups     []string
	CreateHome bool
	System     bool
}

// ModifyOptions opções para modificar usuário
type ModifyOptions struct {
	FullName    string
	HomeDir     string
	Shell       string
	Lock        bool
	Unlock      bool
	ExpireDate  string
	AddGroups   []string
	RemoveGroup string
}

// SystemUserManager implementação padrão
type SystemUserManager struct{}

// NewUserManager cria um novo user manager
func NewUserManager() UserManager {
	return &SystemUserManager{}
}

// List lista usuários
func (m *SystemUserManager) List(options ListOptions) ([]*UserInfo, error) {
	// Lê /etc/passwd
	cmd := exec.Command("getent", "passwd")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %v", err)
	}

	var users []*UserInfo
	lines := strings.Split(string(output), "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}

		fields := strings.Split(line, ":")
		if len(fields) < 7 {
			continue
		}

		username := fields[0]
		uid := fields[2]
		gid := fields[3]
		fullName := fields[4]
		homeDir := fields[5]
		shell := fields[6]

		// Parse UID
		uidInt, err := strconv.Atoi(uid)
		if err != nil {
			continue
		}

		// Filtrar usuários de sistema se necessário
		if !options.SystemUsers && uidInt < 1000 {
			continue
		}

		// Aplicar filtros
		if options.Filter != "" {
			if !strings.Contains(strings.ToLower(username), strings.ToLower(options.Filter)) &&
				!strings.Contains(strings.ToLower(fullName), strings.ToLower(options.Filter)) {
				continue
			}
		}

		if options.Group != "" {
			// Verifica se usuário pertence ao grupo
			inGroup, err := m.isUserInGroup(username, options.Group)
			if err != nil || !inGroup {
				continue
			}
		}

		users = append(users, &UserInfo{
			Username: username,
			UID:      uid,
			GID:      gid,
			HomeDir:  homeDir,
			Shell:    shell,
			FullName: fullName,
		})
	}

	return users, nil
}

// Info obtém informações detalhadas de um usuário
func (m *SystemUserManager) Info(username string) (*UserDetail, error) {
	// Obtém info básica
	u, err := user.Lookup(username)
	if err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}

	detail := &UserDetail{
		UserInfo: &UserInfo{
			Username: u.Username,
			UID:      u.Uid,
			GID:      u.Gid,
			HomeDir:  u.HomeDir,
		},
	}

	// Obtém grupos
	groups, err := u.GroupIds()
	if err == nil {
		for _, gid := range groups {
			g, err := user.LookupGroupId(gid)
			if err == nil {
				detail.Groups = append(detail.Groups, g.Name)
			}
		}
	}

	// Obtém informações adicionais do sistema
	// Last login
	lastCmd := exec.Command("lastlog", "-u", username)
	if output, err := lastCmd.Output(); err == nil {
		lines := strings.Split(string(output), "\n")
		if len(lines) > 1 {
			// Parse lastlog output
			detail.LastLogin = time.Now() // Simplified
		}
	}

	// Password status
	passwdCmd := exec.Command("passwd", "-S", username)
	if output, err := passwdCmd.Output(); err == nil {
		statusLine := string(output)
		if strings.Contains(statusLine, "P") {
			detail.PasswordSet = true
		}
		if strings.Contains(statusLine, "L") {
			detail.Locked = true
		}
	}

	// Expiry information
	chageCmd := exec.Command("chage", "-l", username)
	if output, err := chageCmd.Output(); err == nil {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if strings.Contains(line, "Account expires") {
				parts := strings.Split(line, ":")
				if len(parts) > 1 {
					detail.ExpiryDate = strings.TrimSpace(parts[1])
				}
			}
		}
	}

	return detail, nil
}

// Add adiciona um novo usuário
func (m *SystemUserManager) Add(options AddUserOptions) error {
	args := []string{}

	if options.System {
		args = append(args, "--system")
	}

	if options.FullName != "" {
		args = append(args, "--gecos", options.FullName)
	}

	if options.HomeDir != "" {
		args = append(args, "--home", options.HomeDir)
	}

	if options.Shell != "" {
		args = append(args, "--shell", options.Shell)
	}

	if options.CreateHome {
		args = append(args, "--create-home")
	} else {
		args = append(args, "--no-create-home")
	}

	if len(options.Groups) > 0 {
		args = append(args, "--groups", strings.Join(options.Groups, ","))
	}

	args = append(args, options.Username)

	cmd := exec.Command("useradd", args...)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to add user: %v - %s", err, string(output))
	}

	return nil
}

// Remove remove um usuário
func (m *SystemUserManager) Remove(username string, removeHome bool) error {
	args := []string{}

	if removeHome {
		args = append(args, "--remove")
	}

	args = append(args, username)

	cmd := exec.Command("userdel", args...)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to remove user: %v - %s", err, string(output))
	}

	return nil
}

// Modify modifica um usuário existente
func (m *SystemUserManager) Modify(username string, options ModifyOptions) error {
	// Lock/Unlock
	if options.Lock {
		cmd := exec.Command("usermod", "--lock", username)
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to lock user: %v - %s", err, string(output))
		}
	}

	if options.Unlock {
		cmd := exec.Command("usermod", "--unlock", username)
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to unlock user: %v - %s", err, string(output))
		}
	}

	// Modify attributes
	args := []string{}

	if options.FullName != "" {
		args = append(args, "--comment", options.FullName)
	}

	if options.HomeDir != "" {
		args = append(args, "--home", options.HomeDir)
	}

	if options.Shell != "" {
		args = append(args, "--shell", options.Shell)
	}

	if options.ExpireDate != "" {
		args = append(args, "--expiredate", options.ExpireDate)
	}

	if len(options.AddGroups) > 0 {
		args = append(args, "--append", "--groups", strings.Join(options.AddGroups, ","))
	}

	if len(args) > 0 {
		args = append(args, username)
		cmd := exec.Command("usermod", args...)
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to modify user: %v - %s", err, string(output))
		}
	}

	return nil
}

// ListGroups lista grupos do sistema
func (m *SystemUserManager) ListGroups() ([]*GroupInfo, error) {
	cmd := exec.Command("getent", "group")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list groups: %v", err)
	}

	var groups []*GroupInfo
	lines := strings.Split(string(output), "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}

		fields := strings.Split(line, ":")
		if len(fields) < 4 {
			continue
		}

		name := fields[0]
		gid := fields[2]
		members := []string{}

		if fields[3] != "" {
			members = strings.Split(fields[3], ",")
		}

		groups = append(groups, &GroupInfo{
			Name:    name,
			GID:     gid,
			Members: members,
		})
	}

	return groups, nil
}

// AddToGroup adiciona usuário a um grupo
func (m *SystemUserManager) AddToGroup(username, group string) error {
	cmd := exec.Command("usermod", "--append", "--groups", group, username)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to add user to group: %v - %s", err, string(output))
	}

	return nil
}

// RemoveFromGroup remove usuário de um grupo
func (m *SystemUserManager) RemoveFromGroup(username, group string) error {
	cmd := exec.Command("gpasswd", "--delete", username, group)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to remove user from group: %v - %s", err, string(output))
	}

	return nil
}

// isUserInGroup verifica se usuário pertence a um grupo
func (m *SystemUserManager) isUserInGroup(username, group string) (bool, error) {
	cmd := exec.Command("id", "-Gn", username)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return false, err
	}

	groups := strings.Fields(stdout.String())
	for _, g := range groups {
		if g == group {
			return true, nil
		}
	}

	return false, nil
}
