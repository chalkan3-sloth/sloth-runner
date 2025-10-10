package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// ConfigValidator interface para validação de configuração
type ConfigValidator interface {
	ValidateFile(path string) (*ValidationResult, error)
	ValidateYAML(data []byte) (*ValidationResult, error)
}

// ValidationResult resultado da validação
type ValidationResult struct {
	Valid    bool
	Errors   []ValidationError
	Warnings []ValidationWarning
	FilePath string
}

// ValidationError erro de validação
type ValidationError struct {
	Field   string
	Message string
}

// ValidationWarning aviso de validação
type ValidationWarning struct {
	Field   string
	Message string
}

// SystemValidator implementação padrão
type SystemValidator struct{}

// NewValidator cria um novo validador
func NewValidator() ConfigValidator {
	return &SystemValidator{}
}

// ValidateFile valida um arquivo de configuração
func (v *SystemValidator) ValidateFile(path string) (*ValidationResult, error) {
	result := &ValidationResult{
		Valid:    true,
		Errors:   []ValidationError{},
		Warnings: []ValidationWarning{},
		FilePath: path,
	}

	// Verifica se arquivo existe
	if _, err := os.Stat(path); os.IsNotExist(err) {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:   "file",
			Message: fmt.Sprintf("Configuration file not found: %s", path),
		})
		return result, nil
	}

	// Lê arquivo
	data, err := os.ReadFile(path)
	if err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:   "file",
			Message: fmt.Sprintf("Failed to read file: %v", err),
		})
		return result, nil
	}

	// Valida conteúdo baseado na extensão
	ext := filepath.Ext(path)
	switch ext {
	case ".yaml", ".yml":
		return v.ValidateYAML(data)
	case ".json":
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:   "format",
			Message: "JSON validation not yet implemented",
		})
	default:
		result.Warnings = append(result.Warnings, ValidationWarning{
			Field:   "format",
			Message: fmt.Sprintf("Unknown file format: %s. Treating as YAML.", ext),
		})
		return v.ValidateYAML(data)
	}

	return result, nil
}

// ValidateYAML valida YAML
func (v *SystemValidator) ValidateYAML(data []byte) (*ValidationResult, error) {
	result := &ValidationResult{
		Valid:    true,
		Errors:   []ValidationError{},
		Warnings: []ValidationWarning{},
	}

	// Parse YAML
	var config map[string]interface{}
	if err := yaml.Unmarshal(data, &config); err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:   "syntax",
			Message: fmt.Sprintf("YAML syntax error: %v", err),
		})
		return result, nil
	}

	// Validações básicas
	v.validateBasicStructure(config, result)

	return result, nil
}

// validateBasicStructure valida estrutura básica
func (v *SystemValidator) validateBasicStructure(config map[string]interface{}, result *ValidationResult) {
	// Verifica campos comuns
	commonFields := []string{"server", "database", "logging", "agents"}

	for _, field := range commonFields {
		if _, ok := config[field]; ok {
			// Campo existe, validações específicas podem ser adicionadas aqui
			continue
		}
	}

	// Se nenhum campo comum foi encontrado, adiciona aviso
	hasAnyField := false
	for _, field := range commonFields {
		if _, ok := config[field]; ok {
			hasAnyField = true
			break
		}
	}

	if !hasAnyField && len(config) > 0 {
		result.Warnings = append(result.Warnings, ValidationWarning{
			Field:   "structure",
			Message: "Configuration doesn't contain expected fields (server, database, logging, agents)",
		})
	}

	// Valida valores específicos se existirem
	if server, ok := config["server"].(map[string]interface{}); ok {
		v.validateServer(server, result)
	}

	if database, ok := config["database"].(map[string]interface{}); ok {
		v.validateDatabase(database, result)
	}

	if logging, ok := config["logging"].(map[string]interface{}); ok {
		v.validateLogging(logging, result)
	}
}

// validateServer valida configurações de servidor
func (v *SystemValidator) validateServer(server map[string]interface{}, result *ValidationResult) {
	// Valida porta
	if port, ok := server["port"].(int); ok {
		if port < 1 || port > 65535 {
			result.Errors = append(result.Errors, ValidationError{
				Field:   "server.port",
				Message: fmt.Sprintf("Invalid port number: %d (must be between 1-65535)", port),
			})
			result.Valid = false
		}
	}

	// Valida host
	if host, ok := server["host"].(string); ok {
		if host == "" {
			result.Warnings = append(result.Warnings, ValidationWarning{
				Field:   "server.host",
				Message: "Empty host field, will use default",
			})
		}
	}
}

// validateDatabase valida configurações de banco de dados
func (v *SystemValidator) validateDatabase(database map[string]interface{}, result *ValidationResult) {
	// Valida path
	if path, ok := database["path"].(string); ok {
		if path == "" {
			result.Errors = append(result.Errors, ValidationError{
				Field:   "database.path",
				Message: "Database path cannot be empty",
			})
			result.Valid = false
		}
	}
}

// validateLogging valida configurações de logging
func (v *SystemValidator) validateLogging(logging map[string]interface{}, result *ValidationResult) {
	// Valida level
	if level, ok := logging["level"].(string); ok {
		validLevels := []string{"debug", "info", "warn", "error", "fatal"}
		valid := false
		for _, vl := range validLevels {
			if level == vl {
				valid = true
				break
			}
		}
		if !valid {
			result.Errors = append(result.Errors, ValidationError{
				Field:   "logging.level",
				Message: fmt.Sprintf("Invalid log level '%s', must be one of: %v", level, validLevels),
			})
			result.Valid = false
		}
	}
}

// GetDefaultConfigPaths retorna paths padrão de configuração
func GetDefaultConfigPaths() []string {
	home := os.Getenv("HOME")
	return []string{
		filepath.Join(home, ".config", "sloth-runner", "config.yaml"),
		filepath.Join(home, ".config", "sloth-runner", "config.yml"),
		filepath.Join("/etc", "sloth-runner", "config.yaml"),
		"./config.yaml",
		"./config.yml",
	}
}

// FindConfigFile procura por arquivo de configuração nos paths padrão
func FindConfigFile() (string, error) {
	for _, path := range GetDefaultConfigPaths() {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("no configuration file found in default locations")
}
