package scaffolding

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/pterm/pterm"
)

// WorkflowTemplate represents a workflow template
type WorkflowTemplate struct {
	Name        string
	Description string
	Category    string
	Complexity  string
	Author      string
	Version     string
	Template    string
}

// TemplateData contains data for template rendering
type TemplateData struct {
	WorkflowName string
	Description  string
	Author       string
	Version      string
	Category     string
	Complexity   string
	CreatedAt    string
	ProjectName  string
	Tasks        []TaskTemplate
}

// TaskTemplate represents a task template
type TaskTemplate struct {
	Name        string
	Description string
	Command     string
	DependsOn   []string
}

// WorkflowScaffolder handles workflow scaffolding
type WorkflowScaffolder struct {
	templates map[string]WorkflowTemplate
}

// NewWorkflowScaffolder creates a new workflow scaffolder
func NewWorkflowScaffolder() *WorkflowScaffolder {
	scaffolder := &WorkflowScaffolder{
		templates: make(map[string]WorkflowTemplate),
	}
	scaffolder.loadBuiltinTemplates()
	return scaffolder
}

// loadBuiltinTemplates loads the built-in workflow templates
func (ws *WorkflowScaffolder) loadBuiltinTemplates() {
	// Basic workflow template
	ws.templates["basic"] = WorkflowTemplate{
		Name:        "basic",
		Description: "Basic workflow with a single task",
		Category:    "general",
		Complexity:  "beginner",
		Author:      "Sloth Runner Team",
		Version:     "1.0.0",
		Template:    basicWorkflowTemplate,
	}

	// CI/CD Pipeline template
	ws.templates["cicd"] = WorkflowTemplate{
		Name:        "cicd",
		Description: "Complete CI/CD pipeline with build, test, and deploy",
		Category:    "devops",
		Complexity:  "intermediate",
		Author:      "Sloth Runner Team",
		Version:     "1.0.0",
		Template:    cicdWorkflowTemplate,
	}

	// Infrastructure template
	ws.templates["infrastructure"] = WorkflowTemplate{
		Name:        "infrastructure",
		Description: "Infrastructure deployment workflow",
		Category:    "iac",
		Complexity:  "advanced",
		Author:      "Sloth Runner Team",
		Version:     "1.0.0",
		Template:    infrastructureWorkflowTemplate,
	}

	// Microservices template
	ws.templates["microservices"] = WorkflowTemplate{
		Name:        "microservices",
		Description: "Microservices deployment workflow",
		Category:    "kubernetes",
		Complexity:  "advanced",
		Author:      "Sloth Runner Team",
		Version:     "1.0.0",
		Template:    microservicesWorkflowTemplate,
	}

	// Data Pipeline template
	ws.templates["data-pipeline"] = WorkflowTemplate{
		Name:        "data-pipeline",
		Description: "Data processing pipeline workflow",
		Category:    "data",
		Complexity:  "intermediate",
		Author:      "Sloth Runner Team",
		Version:     "1.0.0",
		Template:    dataPipelineWorkflowTemplate,
	}
}

// InitWorkflow initializes a new workflow from a template
func (ws *WorkflowScaffolder) InitWorkflow(workflowName string, templateName string, interactive bool) error {
	// Show banner
	pterm.DefaultCenter.Printf("%s\n", pterm.LightCyan("🦥 Sloth Runner Workflow Scaffolder"))
	pterm.Printf("\n")

	var template WorkflowTemplate
	var exists bool

	if templateName == "" || interactive {
		// Interactive template selection
		template, exists = ws.selectTemplateInteractively()
		if !exists {
			return fmt.Errorf("no template selected")
		}
	} else {
		template, exists = ws.templates[templateName]
		if !exists {
			return fmt.Errorf("template '%s' not found", templateName)
		}
	}

	// Get workflow name if not provided
	if workflowName == "" || interactive {
		workflowName = ws.getWorkflowNameInteractively()
	}

	// Create project directory
	projectDir := workflowName
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}

	// Gather template data
	templateData := ws.gatherTemplateData(workflowName, template, interactive)

	// Generate workflow file
	workflowFile := filepath.Join(projectDir, fmt.Sprintf("%s.lua", workflowName))
	if err := ws.generateWorkflowFile(workflowFile, template, templateData); err != nil {
		return fmt.Errorf("failed to generate workflow file: %w", err)
	}

	// Generate additional files
	if err := ws.generateAdditionalFiles(projectDir, templateData); err != nil {
		return fmt.Errorf("failed to generate additional files: %w", err)
	}

	// Show success message
	ws.showSuccessMessage(workflowName, projectDir, template)

	return nil
}

// selectTemplateInteractively allows user to select a template interactively
func (ws *WorkflowScaffolder) selectTemplateInteractively() (WorkflowTemplate, bool) {
	var templateNames []string
	var templateDescriptions []string

	for name, template := range ws.templates {
		templateNames = append(templateNames, name)
		desc := fmt.Sprintf("%s - %s (%s, %s)", 
			template.Description, 
			template.Category, 
			template.Complexity,
			template.Version)
		templateDescriptions = append(templateDescriptions, desc)
	}

	prompt := &survey.Select{
		Message: "Select a workflow template:",
		Options: templateDescriptions,
		Help:    "Choose the type of workflow you want to create",
	}

	var selectedIndex int
	if err := survey.AskOne(prompt, &selectedIndex); err != nil {
		return WorkflowTemplate{}, false
	}

	selectedName := templateNames[selectedIndex]
	return ws.templates[selectedName], true
}

// getWorkflowNameInteractively gets workflow name from user
func (ws *WorkflowScaffolder) getWorkflowNameInteractively() string {
	var workflowName string
	prompt := &survey.Input{
		Message: "Enter workflow name:",
		Default: "my-workflow",
		Help:    "The name of your workflow (will be used as directory and file name)",
	}
	survey.AskOne(prompt, &workflowName)
	return workflowName
}

// gatherTemplateData gathers data for template rendering
func (ws *WorkflowScaffolder) gatherTemplateData(workflowName string, template WorkflowTemplate, interactive bool) TemplateData {
	data := TemplateData{
		WorkflowName: workflowName,
		Description:  template.Description,
		Author:       template.Author,
		Version:      "1.0.0",
		Category:     template.Category,
		Complexity:   template.Complexity,
		CreatedAt:    time.Now().Format("2006-01-02 15:04:05"),
		ProjectName:  workflowName,
	}

	if interactive {
		// Ask for additional details
		var description string
		survey.AskOne(&survey.Input{
			Message: "Enter workflow description:",
			Default: data.Description,
		}, &description)
		data.Description = description

		var author string
		survey.AskOne(&survey.Input{
			Message: "Enter author name:",
			Default: "Developer",
		}, &author)
		data.Author = author
	}

	return data
}

// generateWorkflowFile generates the main workflow file
func (ws *WorkflowScaffolder) generateWorkflowFile(filename string, workflowTemplate WorkflowTemplate, data TemplateData) error {
	tmpl, err := template.New("workflow").Parse(workflowTemplate.Template)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}

// generateAdditionalFiles generates additional project files
func (ws *WorkflowScaffolder) generateAdditionalFiles(projectDir string, data TemplateData) error {
	// Generate README.md
	readmeFile := filepath.Join(projectDir, "README.md")
	if err := ws.generateReadme(readmeFile, data); err != nil {
		return err
	}

	// Generate .gitignore
	gitignoreFile := filepath.Join(projectDir, ".gitignore")
	if err := ws.generateGitignore(gitignoreFile); err != nil {
		return err
	}

	// Generate sloth-runner.yaml (config file)
	configFile := filepath.Join(projectDir, "sloth-runner.yaml")
	if err := ws.generateConfig(configFile, data); err != nil {
		return err
	}

	return nil
}

// generateReadme generates a README.md file
func (ws *WorkflowScaffolder) generateReadme(filename string, data TemplateData) error {
	content := fmt.Sprintf(readmeTemplate, 
		data.WorkflowName, 
		data.Description, 
		data.Author, 
		data.CreatedAt,
		data.WorkflowName,
		data.WorkflowName)

	return os.WriteFile(filename, []byte(content), 0644)
}

// generateGitignore generates a .gitignore file
func (ws *WorkflowScaffolder) generateGitignore(filename string) error {
	return os.WriteFile(filename, []byte(gitignoreTemplate), 0644)
}

// generateConfig generates a sloth-runner.yaml config file
func (ws *WorkflowScaffolder) generateConfig(filename string, data TemplateData) error {
	content := fmt.Sprintf(configTemplate, data.WorkflowName, data.WorkflowName, data.Description, data.Version)
	return os.WriteFile(filename, []byte(content), 0644)
}

// showSuccessMessage displays success information
func (ws *WorkflowScaffolder) showSuccessMessage(workflowName, projectDir string, template WorkflowTemplate) {
	pterm.Printf("\n")
	pterm.DefaultHeader.WithFullWidth(false).WithBackgroundStyle(pterm.NewStyle(pterm.BgGreen)).WithTextStyle(pterm.NewStyle(pterm.FgBlack)).Printf("Workflow Created Successfully")
	
	pterm.Printf("\n%s %s\n", pterm.Green("✓"), pterm.DefaultBasicText.WithStyle(pterm.NewStyle(pterm.Bold)).Sprint("Project created:"))
	pterm.Printf("  Name: %s\n", pterm.Cyan(workflowName))
	pterm.Printf("  Template: %s (%s)\n", pterm.Cyan(template.Name), template.Description)
	pterm.Printf("  Directory: %s\n", pterm.Cyan(projectDir))
	
	pterm.Printf("\n%s %s\n", pterm.Blue("📁"), pterm.DefaultBasicText.WithStyle(pterm.NewStyle(pterm.Bold)).Sprint("Generated files:"))
	pterm.Printf("  %s/%s.lua     - Main workflow file\n", projectDir, workflowName)
	pterm.Printf("  %s/README.md        - Project documentation\n", projectDir)
	pterm.Printf("  %s/.gitignore       - Git ignore rules\n", projectDir)
	pterm.Printf("  %s/sloth-runner.yaml - Configuration file\n", projectDir)
	
	pterm.Printf("\n%s %s\n", pterm.Yellow("🚀"), pterm.DefaultBasicText.WithStyle(pterm.NewStyle(pterm.Bold)).Sprint("Next steps:"))
	pterm.Printf("  1. cd %s\n", workflowName)
	pterm.Printf("  2. sloth-runner run -f %s.lua\n", workflowName)
	pterm.Printf("  3. Edit the workflow to suit your needs\n")
	
	pterm.Printf("\n%s %s\n", pterm.Magenta("📖"), pterm.DefaultBasicText.WithStyle(pterm.NewStyle(pterm.Bold)).Sprint("Documentation:"))
	pterm.Printf("  https://github.com/chalkan3-sloth/sloth-runner/docs\n")
}

// ListTemplates lists available templates
func (ws *WorkflowScaffolder) ListTemplates() {
	pterm.DefaultHeader.WithFullWidth(false).WithBackgroundStyle(pterm.NewStyle(pterm.BgBlue)).WithTextStyle(pterm.NewStyle(pterm.FgWhite)).Printf("Available Workflow Templates")
	
	pterm.Printf("\n")
	for templateName, template := range ws.templates {
		pterm.Printf("%s %s\n", pterm.Cyan("├─"), pterm.DefaultBasicText.WithStyle(pterm.NewStyle(pterm.Bold)).Sprint(templateName))
		pterm.Printf("  %s %s\n", pterm.Gray("│"), template.Description)
		pterm.Printf("  %s Category: %s | Complexity: %s | Version: %s\n", 
			pterm.Gray("│"), 
			pterm.Green(template.Category), 
			pterm.Yellow(template.Complexity), 
			pterm.Blue(template.Version))
		pterm.Printf("  %s\n", pterm.Gray("│"))
	}
	pterm.Printf("%s\n", pterm.Gray("└─"))
	
	pterm.Printf("\n%s %s\n", pterm.Blue("💡"), pterm.DefaultBasicText.WithStyle(pterm.NewStyle(pterm.Bold)).Sprint("Usage:"))
	pterm.Printf("  sloth-runner workflow init <name> --template <template>\n")
	pterm.Printf("  sloth-runner workflow init <name> --interactive\n")
}