package infra

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"regexp"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

// CronModule provides cron job management
type CronModule struct {
	L *lua.LState
}

// NewCronModule creates a new Cron module instance
func NewCronModule(L *lua.LState) *CronModule {
	return &CronModule{L: L}
}

// Register registers the Cron module with the Lua state
func (m *CronModule) Register(L *lua.LState) {
	cronTable := L.NewTable()

	L.SetField(cronTable, "add", L.NewFunction(m.add))
	L.SetField(cronTable, "remove", L.NewFunction(m.remove))
	L.SetField(cronTable, "list", L.NewFunction(m.list))
	L.SetField(cronTable, "exists", L.NewFunction(m.exists))
	L.SetField(cronTable, "enable", L.NewFunction(m.enable))
	L.SetField(cronTable, "disable", L.NewFunction(m.disable))
	L.SetField(cronTable, "validate_schedule", L.NewFunction(m.validateSchedule))

	L.SetGlobal("cron", cronTable)
}

func (m *CronModule) execCommand(args ...string) (string, error) {
	cmd := exec.Command(args[0], args[1:]...)
	output, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(output)), err
}

func (m *CronModule) needsSudo() bool {
	output, err := m.execCommand("id", "-u")
	if err != nil {
		return true
	}
	return output != "0"
}

func (m *CronModule) prependSudo(args []string) []string {
	if m.needsSudo() {
		return append([]string{"sudo"}, args...)
	}
	return args
}

// getCurrentUser returns the current user
func (m *CronModule) getCurrentUser() string {
	currentUser, err := user.Current()
	if err != nil {
		return "root"
	}
	return currentUser.Username
}

// validateCronSchedule validates a cron schedule expression
func (m *CronModule) validateCronSchedule(schedule string) error {
	// Cron format: minute hour day month weekday
	// Also support special formats: @hourly, @daily, @weekly, @monthly, @yearly, @reboot
	specialFormats := []string{"@hourly", "@daily", "@weekly", "@monthly", "@yearly", "@reboot", "@annually", "@midnight"}
	for _, special := range specialFormats {
		if schedule == special {
			return nil
		}
	}

	parts := strings.Fields(schedule)
	if len(parts) != 5 {
		return fmt.Errorf("invalid cron schedule format (expected 5 fields: minute hour day month weekday)")
	}

	// Validate each field with basic regex
	patterns := []string{
		`^(\*|[0-5]?[0-9]|(\*\/[0-9]+)|([0-5]?[0-9]-[0-5]?[0-9])|([0-5]?[0-9](,[0-5]?[0-9])*))$`, // minute
		`^(\*|[01]?[0-9]|2[0-3]|(\*\/[0-9]+)|([01]?[0-9]-2[0-3])|([01]?[0-9](,[01]?[0-9])*))$`,     // hour
		`^(\*|[01]?[0-9]|2[0-9]|3[01]|(\*\/[0-9]+)|([0-2]?[0-9]-3[01])|([0-2]?[0-9](,[0-2]?[0-9])*))$`, // day
		`^(\*|[0-9]|1[0-2]|(\*\/[0-9]+)|([0-9]-1[0-2])|([0-9](,[0-9])*))$`,                         // month
		`^(\*|[0-6]|(\*\/[0-9]+)|([0-6]-[0-6])|([0-6](,[0-6])*))$`,                                 // weekday
	}

	for i, part := range parts {
		matched, err := regexp.MatchString(patterns[i], part)
		if err != nil || !matched {
			return fmt.Errorf("invalid cron field %d: %s", i+1, part)
		}
	}

	return nil
}

// getCrontab reads crontab for a user
func (m *CronModule) getCrontab(username string) (string, error) {
	var cmd *exec.Cmd
	if username == "root" || username == "" {
		args := m.prependSudo([]string{"crontab", "-l"})
		cmd = exec.Command(args[0], args[1:]...)
	} else {
		args := m.prependSudo([]string{"crontab", "-u", username, "-l"})
		cmd = exec.Command(args[0], args[1:]...)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		// Empty crontab returns error, which is okay
		if strings.Contains(string(output), "no crontab") {
			return "", nil
		}
		return "", err
	}

	return string(output), nil
}

// setCrontab writes crontab for a user
func (m *CronModule) setCrontab(username, content string) error {
	// Write content to temp file
	tmpFile, err := os.CreateTemp("", "crontab-*.txt")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(content); err != nil {
		return fmt.Errorf("failed to write temp file: %w", err)
	}
	tmpFile.Close()

	// Install crontab from temp file
	var args []string
	if username == "root" || username == "" {
		args = m.prependSudo([]string{"crontab", tmpFile.Name()})
	} else {
		args = m.prependSudo([]string{"crontab", "-u", username, tmpFile.Name()})
	}

	cmd := exec.Command(args[0], args[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to install crontab: %s\n%s", err, string(output))
	}

	return nil
}

// parseCronEntry parses a cron entry line
func (m *CronModule) parseCronEntry(line string) (name, schedule, command string, disabled bool) {
	line = strings.TrimSpace(line)

	// Check if disabled (commented out)
	if strings.HasPrefix(line, "#") {
		disabled = true
		line = strings.TrimPrefix(line, "#")
		line = strings.TrimSpace(line)
	}

	// Check for name marker
	if strings.Contains(line, "# SLOTH_NAME:") {
		parts := strings.Split(line, "# SLOTH_NAME:")
		if len(parts) == 2 {
			return strings.TrimSpace(parts[1]), "", "", disabled
		}
	}

	// Parse schedule and command
	fields := strings.Fields(line)
	if len(fields) < 6 {
		return "", "", "", disabled
	}

	// Check for special schedule format
	if strings.HasPrefix(fields[0], "@") {
		schedule = fields[0]
		command = strings.Join(fields[1:], " ")
	} else {
		// Standard 5-field schedule
		schedule = strings.Join(fields[0:5], " ")
		command = strings.Join(fields[5:], " ")
	}

	return "", schedule, command, disabled
}

// findCronJob finds a cron job by name in crontab content
func (m *CronModule) findCronJob(content, name string) (found bool, entry string, lineNum int) {
	lines := strings.Split(content, "\n")

	for i, line := range lines {
		if strings.Contains(line, "# SLOTH_NAME:"+name) {
			// Found name marker, next line should be the job
			if i+1 < len(lines) {
				return true, lines[i+1], i + 1
			}
		}
	}

	return false, "", -1
}

// add adds a cron job
func (m *CronModule) add(L *lua.LState) int {
	options := L.CheckTable(1)

	name := options.RawGetString("name").String()
	schedule := options.RawGetString("schedule").String()
	command := options.RawGetString("command").String()
	username := options.RawGetString("user").String()

	if name == "" {
		L.Push(lua.LFalse)
		L.Push(lua.LString("Job name is required"))
		return 2
	}

	if schedule == "" {
		L.Push(lua.LFalse)
		L.Push(lua.LString("Schedule is required"))
		return 2
	}

	if command == "" {
		L.Push(lua.LFalse)
		L.Push(lua.LString("Command is required"))
		return 2
	}

	// Default to current user
	if username == "" {
		username = m.getCurrentUser()
	}

	// Validate schedule
	if err := m.validateCronSchedule(schedule); err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Invalid schedule: %s", err)))
		return 2
	}

	// Get current crontab
	currentCrontab, err := m.getCrontab(username)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to read crontab: %s", err)))
		return 2
	}

	// Check if job already exists (idempotency)
	found, existingEntry, lineNum := m.findCronJob(currentCrontab, name)
	if found {
		// Parse existing entry to compare
		_, existingSchedule, existingCommand, existingDisabled := m.parseCronEntry(existingEntry)

		if existingSchedule == schedule && existingCommand == command && !existingDisabled {
			L.Push(lua.LTrue)
			L.Push(lua.LString(fmt.Sprintf("Cron job '%s' already exists with same schedule and command (idempotent)", name)))
			return 2
		}

		// Update existing job
		lines := strings.Split(currentCrontab, "\n")
		lines[lineNum] = fmt.Sprintf("%s %s", schedule, command)
		newCrontab := strings.Join(lines, "\n")

		if err := m.setCrontab(username, newCrontab); err != nil {
			L.Push(lua.LFalse)
			L.Push(lua.LString(fmt.Sprintf("Failed to update cron job: %s", err)))
			return 2
		}

		L.Push(lua.LTrue)
		L.Push(lua.LString(fmt.Sprintf("Cron job '%s' updated for user %s", name, username)))
		return 2
	}

	// Add new job
	var newCrontab string
	if currentCrontab != "" {
		newCrontab = currentCrontab + "\n"
	}
	newCrontab += fmt.Sprintf("# SLOTH_NAME:%s\n%s %s\n", name, schedule, command)

	if err := m.setCrontab(username, newCrontab); err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to add cron job: %s", err)))
		return 2
	}

	L.Push(lua.LTrue)
	L.Push(lua.LString(fmt.Sprintf("Cron job '%s' added for user %s", name, username)))
	return 2
}

// remove removes a cron job
func (m *CronModule) remove(L *lua.LState) int {
	name := L.CheckString(1)
	username := ""
	if L.GetTop() >= 2 {
		username = L.CheckString(2)
	}

	if username == "" {
		username = m.getCurrentUser()
	}

	// Get current crontab
	currentCrontab, err := m.getCrontab(username)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to read crontab: %s", err)))
		return 2
	}

	// Find job
	found, _, lineNum := m.findCronJob(currentCrontab, name)
	if !found {
		L.Push(lua.LTrue)
		L.Push(lua.LString(fmt.Sprintf("Cron job '%s' does not exist (idempotent)", name)))
		return 2
	}

	// Remove job and its name marker
	lines := strings.Split(currentCrontab, "\n")
	newLines := []string{}

	for i, line := range lines {
		// Skip name marker and the job line
		if i == lineNum-1 || i == lineNum {
			continue
		}
		newLines = append(newLines, line)
	}

	newCrontab := strings.Join(newLines, "\n")
	newCrontab = strings.TrimSpace(newCrontab)

	if err := m.setCrontab(username, newCrontab); err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to remove cron job: %s", err)))
		return 2
	}

	L.Push(lua.LTrue)
	L.Push(lua.LString(fmt.Sprintf("Cron job '%s' removed from user %s", name, username)))
	return 2
}

// list lists cron jobs
func (m *CronModule) list(L *lua.LState) int {
	username := ""
	if L.GetTop() >= 1 {
		username = L.CheckString(1)
	}

	if username == "" {
		username = m.getCurrentUser()
	}

	// Get current crontab
	currentCrontab, err := m.getCrontab(username)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("Failed to read crontab: %s", err)))
		return 2
	}

	// Parse crontab entries
	jobsTable := L.NewTable()
	lines := strings.Split(currentCrontab, "\n")

	var currentName string
	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])

		if line == "" || strings.HasPrefix(line, "# ") && !strings.Contains(line, "SLOTH_NAME:") {
			continue
		}

		// Check for name marker
		if strings.Contains(line, "# SLOTH_NAME:") {
			currentName = strings.TrimSpace(strings.TrimPrefix(line, "# SLOTH_NAME:"))
			continue
		}

		// Parse job entry
		_, schedule, command, disabled := m.parseCronEntry(line)

		if schedule != "" && command != "" {
			jobInfo := L.NewTable()

			if currentName != "" {
				jobInfo.RawSetString("name", lua.LString(currentName))
				currentName = "" // Reset
			} else {
				jobInfo.RawSetString("name", lua.LString(fmt.Sprintf("unnamed_%d", i)))
			}

			jobInfo.RawSetString("schedule", lua.LString(schedule))
			jobInfo.RawSetString("command", lua.LString(command))
			jobInfo.RawSetString("user", lua.LString(username))
			jobInfo.RawSetString("enabled", lua.LBool(!disabled))

			jobsTable.Append(jobInfo)
		}
	}

	L.Push(jobsTable)
	L.Push(lua.LNil)
	return 2
}

// exists checks if a cron job exists
func (m *CronModule) exists(L *lua.LState) int {
	name := L.CheckString(1)
	username := ""
	if L.GetTop() >= 2 {
		username = L.CheckString(2)
	}

	if username == "" {
		username = m.getCurrentUser()
	}

	// Get current crontab
	currentCrontab, err := m.getCrontab(username)
	if err != nil {
		L.Push(lua.LFalse)
		return 1
	}

	// Check if job exists
	found, _, _ := m.findCronJob(currentCrontab, name)
	L.Push(lua.LBool(found))
	return 1
}

// enable enables a cron job (uncomments it)
func (m *CronModule) enable(L *lua.LState) int {
	name := L.CheckString(1)
	username := ""
	if L.GetTop() >= 2 {
		username = L.CheckString(2)
	}

	if username == "" {
		username = m.getCurrentUser()
	}

	// Get current crontab
	currentCrontab, err := m.getCrontab(username)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to read crontab: %s", err)))
		return 2
	}

	// Find job
	found, entry, lineNum := m.findCronJob(currentCrontab, name)
	if !found {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Cron job '%s' not found", name)))
		return 2
	}

	// Check if already enabled
	if !strings.HasPrefix(strings.TrimSpace(entry), "#") {
		L.Push(lua.LTrue)
		L.Push(lua.LString(fmt.Sprintf("Cron job '%s' is already enabled (idempotent)", name)))
		return 2
	}

	// Enable job (remove leading #)
	lines := strings.Split(currentCrontab, "\n")
	lines[lineNum] = strings.TrimPrefix(strings.TrimSpace(entry), "#")
	lines[lineNum] = strings.TrimSpace(lines[lineNum])
	newCrontab := strings.Join(lines, "\n")

	if err := m.setCrontab(username, newCrontab); err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to enable cron job: %s", err)))
		return 2
	}

	L.Push(lua.LTrue)
	L.Push(lua.LString(fmt.Sprintf("Cron job '%s' enabled", name)))
	return 2
}

// disable disables a cron job (comments it out)
func (m *CronModule) disable(L *lua.LState) int {
	name := L.CheckString(1)
	username := ""
	if L.GetTop() >= 2 {
		username = L.CheckString(2)
	}

	if username == "" {
		username = m.getCurrentUser()
	}

	// Get current crontab
	currentCrontab, err := m.getCrontab(username)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to read crontab: %s", err)))
		return 2
	}

	// Find job
	found, entry, lineNum := m.findCronJob(currentCrontab, name)
	if !found {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Cron job '%s' not found", name)))
		return 2
	}

	// Check if already disabled
	if strings.HasPrefix(strings.TrimSpace(entry), "#") {
		L.Push(lua.LTrue)
		L.Push(lua.LString(fmt.Sprintf("Cron job '%s' is already disabled (idempotent)", name)))
		return 2
	}

	// Disable job (add leading #)
	lines := strings.Split(currentCrontab, "\n")
	lines[lineNum] = "# " + strings.TrimSpace(entry)
	newCrontab := strings.Join(lines, "\n")

	if err := m.setCrontab(username, newCrontab); err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to disable cron job: %s", err)))
		return 2
	}

	L.Push(lua.LTrue)
	L.Push(lua.LString(fmt.Sprintf("Cron job '%s' disabled", name)))
	return 2
}

// validateSchedule validates a cron schedule expression
func (m *CronModule) validateSchedule(L *lua.LState) int {
	schedule := L.CheckString(1)

	err := m.validateCronSchedule(schedule)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LTrue)
	L.Push(lua.LString("Valid cron schedule"))
	return 2
}
