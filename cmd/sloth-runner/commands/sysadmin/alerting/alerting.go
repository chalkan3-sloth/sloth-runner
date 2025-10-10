package alerting

import (
	"fmt"
	"strings"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func NewAlertingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "alerting",
		Short:   "System alerting and monitoring",
		Long:    `Configure and manage system alerts based on thresholds and conditions.`,
		Aliases: []string{"alert", "alerts"},
	}

	// list command
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List alert rules",
		Run: func(cmd *cobra.Command, args []string) {
			if err := runList(); err != nil {
				pterm.Error.Printf("Failed to list rules: %v\n", err)
			}
		},
	}
	cmd.AddCommand(listCmd)

	// add command
	addCmd := &cobra.Command{
		Use:   "add",
		Short: "Add a new alert rule",
		Run: func(cmd *cobra.Command, args []string) {
			name, _ := cmd.Flags().GetString("name")
			ruleType, _ := cmd.Flags().GetString("type")
			threshold, _ := cmd.Flags().GetFloat64("threshold")
			severity, _ := cmd.Flags().GetString("severity")
			target, _ := cmd.Flags().GetString("target")
			description, _ := cmd.Flags().GetString("description")

			if name == "" || ruleType == "" {
				pterm.Error.Println("Name and type are required")
				return
			}

			if err := runAdd(name, ruleType, threshold, severity, target, description); err != nil {
				pterm.Error.Printf("Failed to add rule: %v\n", err)
			}
		},
	}
	addCmd.Flags().StringP("name", "n", "", "Rule name")
	addCmd.Flags().StringP("type", "t", "", "Alert type (cpu, memory, disk, service, process)")
	addCmd.Flags().Float64P("threshold", "T", 0, "Threshold value")
	addCmd.Flags().StringP("severity", "s", "warning", "Severity (info, warning, critical)")
	addCmd.Flags().StringP("target", "r", "", "Target (disk path, service name, process name)")
	addCmd.Flags().StringP("description", "d", "", "Rule description")
	cmd.AddCommand(addCmd)

	// remove command
	removeCmd := &cobra.Command{
		Use:   "remove",
		Short: "Remove an alert rule",
		Run: func(cmd *cobra.Command, args []string) {
			id, _ := cmd.Flags().GetString("id")

			if id == "" {
				pterm.Error.Println("Rule ID is required (use --id)")
				return
			}

			if err := runRemove(id); err != nil {
				pterm.Error.Printf("Failed to remove rule: %v\n", err)
			}
		},
	}
	removeCmd.Flags().StringP("id", "i", "", "Rule ID")
	cmd.AddCommand(removeCmd)

	// check command
	checkCmd := &cobra.Command{
		Use:   "check",
		Short: "Check all alert rules now",
		Run: func(cmd *cobra.Command, args []string) {
			if err := runCheck(); err != nil {
				pterm.Error.Printf("Failed to check rules: %v\n", err)
			}
		},
	}
	cmd.AddCommand(checkCmd)

	// history command
	historyCmd := &cobra.Command{
		Use:   "history",
		Short: "View alert history",
		Run: func(cmd *cobra.Command, args []string) {
			limit, _ := cmd.Flags().GetInt("limit")

			if err := runHistory(limit); err != nil {
				pterm.Error.Printf("Failed to get history: %v\n", err)
			}
		},
	}
	historyCmd.Flags().IntP("limit", "l", 50, "Limit number of alerts to show")
	cmd.AddCommand(historyCmd)

	// docs subcommand
	cmd.AddCommand(&cobra.Command{
		Use:   "docs",
		Short: "Show detailed documentation (man page style)",
		Long:  `Display comprehensive documentation for the alerting command with examples and usage patterns.`,
		Run: func(cmd *cobra.Command, args []string) {
			showAlertingDocs()
		},
	})

	return cmd
}

func showAlertingDocs() {
	title := "SLOTH-RUNNER SYSADMIN ALERTING(1)"
	description := "sloth-runner sysadmin alerting - System alerting and monitoring"
	synopsis := "sloth-runner sysadmin alerting [subcommand] [options]"

	options := [][]string{
		{"list", "List all configured alert rules."},
		{"add", "Add a new alert rule with specified parameters."},
		{"remove", "Remove an existing alert rule by ID."},
		{"check", "Manually check all enabled alert rules now."},
		{"history", "View history of triggered alerts."},
		{"docs", "Show this documentation page."},
	}

	examples := [][]string{
		{
			"Add CPU alert",
			"sloth-runner sysadmin alert add --name \"High CPU\" --type cpu --threshold 80 --severity warning",
			"Alert when CPU exceeds 80%",
		},
		{
			"Add memory alert",
			"sloth-runner sysadmin alerts add --name \"Low Memory\" --type memory --threshold 90 --severity critical",
			"Critical alert when memory exceeds 90%",
		},
		{
			"Add disk space alert",
			"sloth-runner sysadmin alerting add --name \"Disk Full\" --type disk --threshold 85 --target /data",
			"Alert when /data partition exceeds 85%",
		},
		{
			"Add process monitoring",
			"sloth-runner sysadmin alert add --name \"Nginx Down\" --type process --threshold 0 --target nginx",
			"Alert when nginx process is not running",
		},
		{
			"List all rules",
			"sloth-runner sysadmin alerting list",
			"Shows all configured alert rules",
		},
		{
			"Check rules manually",
			"sloth-runner sysadmin alert check",
			"Checks all rules and displays any triggered alerts",
		},
		{
			"View alert history",
			"sloth-runner sysadmin alerting history --limit 100",
			"Shows last 100 triggered alerts",
		},
		{
			"Remove a rule",
			"sloth-runner sysadmin alert remove --id rule-123456",
			"Removes the specified alert rule",
		},
	}

	seeAlso := []string{
		"sloth-runner sysadmin performance - System performance monitoring",
		"sloth-runner sysadmin process - Process management",
	}

	showDocs(title, description, synopsis, options, examples, seeAlso)
}

// showDocs displays formatted documentation similar to man pages
func showDocs(title, description, synopsis string, options [][]string, examples [][]string, seeAlso []string) {
	// Header
	pterm.DefaultHeader.WithFullWidth().Println(title)
	fmt.Println()

	// Name and Description
	pterm.DefaultSection.Println("NAME")
	fmt.Printf("    %s\n\n", description)

	// Synopsis
	if synopsis != "" {
		pterm.DefaultSection.Println("SYNOPSIS")
		fmt.Printf("    %s\n\n", synopsis)
	}

	// Options
	if len(options) > 0 {
		pterm.DefaultSection.Println("OPTIONS")
		for _, opt := range options {
			if len(opt) >= 2 {
				pterm.FgCyan.Printf("    %s\n", opt[0])
				fmt.Printf("        %s\n\n", opt[1])
			}
		}
	}

	// Examples
	if len(examples) > 0 {
		pterm.DefaultSection.Println("EXAMPLES")
		for i, ex := range examples {
			if len(ex) >= 2 {
				pterm.FgYellow.Printf("    Example %d: %s\n", i+1, ex[0])
				pterm.FgGreen.Printf("    $ %s\n", ex[1])
				if len(ex) >= 3 {
					fmt.Printf("        %s\n", ex[2])
				}
				fmt.Println()
			}
		}
	}

	// See Also
	if len(seeAlso) > 0 {
		pterm.DefaultSection.Println("SEE ALSO")
		for _, item := range seeAlso {
			fmt.Printf("    • %s\n", item)
		}
		fmt.Println()
	}

	// Footer
	pterm.FgGray.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	pterm.FgGray.Println("Documentation generated for sloth-runner sysadmin v2.0")
	pterm.FgGray.Println("For more information: sloth-runner sysadmin alerting --help")
}

// runList lista regras de alerta
func runList() error {
	manager := NewAlertManager()

	pterm.DefaultHeader.WithFullWidth().Println("Alert Rules")
	pterm.Println()

	spinner, _ := pterm.DefaultSpinner.Start("Loading alert rules...")

	rules, err := manager.ListRules()
	if err != nil {
		spinner.Fail("Failed to load rules")
		return err
	}

	spinner.Success(fmt.Sprintf("✅ Found %d alert rules", len(rules)))
	pterm.Println()

	if len(rules) == 0 {
		pterm.Info.Println("No alert rules configured")
		pterm.Info.Println("Use 'sloth-runner sysadmin alerting add' to create your first rule")
		return nil
	}

	// Rules table
	tableData := pterm.TableData{
		{"ID", "Name", "Type", "Threshold", "Severity", "Target", "Enabled"},
	}

	for _, rule := range rules {
		enabled := "Yes"
		if !rule.Enabled {
			enabled = pterm.FgRed.Sprint("No")
		} else {
			enabled = pterm.FgGreen.Sprint("Yes")
		}

		severity := string(rule.Severity)
		if rule.Severity == SeverityCritical {
			severity = pterm.FgRed.Sprint(severity)
		} else if rule.Severity == SeverityWarning {
			severity = pterm.FgYellow.Sprint(severity)
		}

		threshold := fmt.Sprintf("%.1f", rule.Threshold)
		if rule.Type == AlertTypeCPU || rule.Type == AlertTypeMemory || rule.Type == AlertTypeDisk {
			threshold += "%"
		}

		target := rule.Target
		if target == "" {
			target = "-"
		}

		tableData = append(tableData, []string{
			truncate(rule.ID, 20),
			truncate(rule.Name, 25),
			string(rule.Type),
			threshold,
			severity,
			truncate(target, 20),
			enabled,
		})
	}

	pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
	pterm.Println()

	// Statistics
	enabled := 0
	for _, rule := range rules {
		if rule.Enabled {
			enabled++
		}
	}

	pterm.Info.Printf("Total rules: %d | Enabled: %d | Disabled: %d\n", len(rules), enabled, len(rules)-enabled)

	return nil
}

// runAdd adiciona uma regra de alerta
func runAdd(name, ruleType string, threshold float64, severity, target, description string) error {
	manager := NewAlertManager()

	pterm.DefaultHeader.WithFullWidth().Println("Add Alert Rule")
	pterm.Println()

	// Valida tipo
	alertType := AlertType(strings.ToLower(ruleType))
	if alertType != AlertTypeCPU && alertType != AlertTypeMemory && alertType != AlertTypeDisk &&
		alertType != AlertTypeService && alertType != AlertTypeProcess {
		return fmt.Errorf("invalid alert type: %s (must be cpu, memory, disk, service, or process)", ruleType)
	}

	// Valida severidade
	sev := Severity(strings.ToLower(severity))
	if sev != SeverityInfo && sev != SeverityWarning && sev != SeverityCritical {
		return fmt.Errorf("invalid severity: %s (must be info, warning, or critical)", severity)
	}

	rule := &AlertRule{
		Name:        name,
		Type:        alertType,
		Enabled:     true,
		Threshold:   threshold,
		Severity:    sev,
		Target:      target,
		Description: description,
	}

	spinner, _ := pterm.DefaultSpinner.Start("Adding alert rule...")

	err := manager.AddRule(rule)
	if err != nil {
		spinner.Fail("Failed to add rule")
		return err
	}

	spinner.Success(fmt.Sprintf("✅ Alert rule '%s' added successfully", name))
	pterm.Println()

	// Show created rule details
	pterm.Info.Printf("Rule ID: %s\n", rule.ID)
	pterm.Info.Printf("Type: %s | Threshold: %.1f | Severity: %s\n", rule.Type, rule.Threshold, rule.Severity)
	if target != "" {
		pterm.Info.Printf("Target: %s\n", target)
	}

	return nil
}

// runRemove remove uma regra de alerta
func runRemove(id string) error {
	manager := NewAlertManager()

	pterm.DefaultHeader.WithFullWidth().Println("Remove Alert Rule")
	pterm.Println()

	spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Removing rule %s...", id))

	err := manager.RemoveRule(id)
	if err != nil {
		spinner.Fail("Failed to remove rule")
		return err
	}

	spinner.Success(fmt.Sprintf("✅ Alert rule %s removed successfully", id))
	pterm.Println()

	return nil
}

// runCheck verifica todas as regras
func runCheck() error {
	manager := NewAlertManager()

	pterm.DefaultHeader.WithFullWidth().Println("Check Alert Rules")
	pterm.Println()

	spinner, _ := pterm.DefaultSpinner.Start("Checking all alert rules...")

	alerts, err := manager.CheckRules()
	if err != nil {
		spinner.Fail("Failed to check rules")
		return err
	}

	if len(alerts) == 0 {
		spinner.Success("✅ All rules checked - No alerts triggered")
		pterm.Println()
		pterm.Success.Println("System is within all defined thresholds")
		return nil
	}

	spinner.Warning(fmt.Sprintf("⚠️  %d alert(s) triggered", len(alerts)))
	pterm.Println()

	// Alerts table
	tableData := pterm.TableData{
		{"Severity", "Rule", "Type", "Message", "Value", "Threshold"},
	}

	for _, alert := range alerts {
		severity := string(alert.Severity)
		if alert.Severity == SeverityCritical {
			severity = pterm.FgRed.Sprint("🔴 " + severity)
		} else if alert.Severity == SeverityWarning {
			severity = pterm.FgYellow.Sprint("🟡 " + severity)
		} else {
			severity = pterm.FgCyan.Sprint("🔵 " + severity)
		}

		value := fmt.Sprintf("%.1f", alert.Value)
		threshold := fmt.Sprintf("%.1f", alert.Threshold)

		tableData = append(tableData, []string{
			severity,
			truncate(alert.RuleName, 25),
			string(alert.Type),
			truncate(alert.Message, 40),
			value,
			threshold,
		})
	}

	pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
	pterm.Println()

	// Summary by severity
	critical, warning, info := 0, 0, 0
	for _, alert := range alerts {
		switch alert.Severity {
		case SeverityCritical:
			critical++
		case SeverityWarning:
			warning++
		case SeverityInfo:
			info++
		}
	}

	if critical > 0 {
		pterm.Error.Printf("Critical: %d\n", critical)
	}
	if warning > 0 {
		pterm.Warning.Printf("Warning: %d\n", warning)
	}
	if info > 0 {
		pterm.Info.Printf("Info: %d\n", info)
	}

	return nil
}

// runHistory exibe histórico de alertas
func runHistory(limit int) error {
	manager := NewAlertManager()

	pterm.DefaultHeader.WithFullWidth().Println("Alert History")
	pterm.Println()

	spinner, _ := pterm.DefaultSpinner.Start("Loading alert history...")

	history, err := manager.GetHistory(limit)
	if err != nil {
		spinner.Fail("Failed to load history")
		return err
	}

	spinner.Success(fmt.Sprintf("✅ Found %d alert(s) in history", len(history)))
	pterm.Println()

	if len(history) == 0 {
		pterm.Info.Println("No alerts in history")
		return nil
	}

	// History table
	tableData := pterm.TableData{
		{"Time", "Severity", "Rule", "Type", "Message"},
	}

	for _, alert := range history {
		severity := string(alert.Severity)
		if alert.Severity == SeverityCritical {
			severity = pterm.FgRed.Sprint(severity)
		} else if alert.Severity == SeverityWarning {
			severity = pterm.FgYellow.Sprint(severity)
		}

		tableData = append(tableData, []string{
			alert.TriggeredAt.Format("2006-01-02 15:04:05"),
			severity,
			truncate(alert.RuleName, 25),
			string(alert.Type),
			truncate(alert.Message, 50),
		})
	}

	pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
	pterm.Println()

	pterm.Info.Printf("Showing %d most recent alerts\n", len(history))

	return nil
}

// truncate trunca uma string
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
