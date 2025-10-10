package security

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func NewSecurityCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "security",
		Short: "Security auditing and management",
		Long:  `Audit access logs, check permissions, scan for vulnerabilities, and manage security policies.`,
		Example: `  # Audit access logs
  sloth-runner sysadmin security audit --agent do-sloth-runner-01

  # Check file permissions
  sloth-runner sysadmin security permissions --path /var/lib/sloth-runner

  # Scan for vulnerabilities
  sloth-runner sysadmin security scan --agent do-sloth-runner-01`,
	}

	// audit command
	auditCmd := &cobra.Command{
		Use:   "audit",
		Short: "Audit access logs and security events",
		Run: func(cmd *cobra.Command, args []string) {
			since, _ := cmd.Flags().GetDuration("since")
			showFailedAuth, _ := cmd.Flags().GetBool("show-failed-auth")
			detectAnomalies, _ := cmd.Flags().GetBool("detect-anomalies")
			outputFormat, _ := cmd.Flags().GetString("output")

			if err := runAudit(since, showFailedAuth, detectAnomalies, outputFormat); err != nil {
				pterm.Error.Printf("Audit failed: %v\n", err)
			}
		},
	}
	auditCmd.Flags().Duration("since", 24*time.Hour, "Audit logs since duration ago")
	auditCmd.Flags().Bool("show-failed-auth", true, "Show failed authentication attempts")
	auditCmd.Flags().Bool("detect-anomalies", true, "Detect suspicious patterns")
	auditCmd.Flags().StringP("output", "o", "table", "Output format (table, json)")
	cmd.AddCommand(auditCmd)

	// scan command
	scanCmd := &cobra.Command{
		Use:   "scan",
		Short: "Scan for security vulnerabilities",
		Run: func(cmd *cobra.Command, args []string) {
			agent, _ := cmd.Flags().GetString("agent")
			full, _ := cmd.Flags().GetBool("full")
			cveOnly, _ := cmd.Flags().GetBool("cve-only")
			dependencyAudit, _ := cmd.Flags().GetBool("dependency-audit")

			if agent == "" {
				pterm.Error.Println("Agent is required (use --agent)")
				return
			}

			if err := runScan(agent, full, cveOnly, dependencyAudit); err != nil {
				pterm.Error.Printf("Scan failed: %v\n", err)
			}
		},
	}
	scanCmd.Flags().StringP("agent", "a", "", "Agent to scan")
	scanCmd.Flags().Bool("full", false, "Full security scan (includes all checks)")
	scanCmd.Flags().Bool("cve-only", false, "Scan only for CVE vulnerabilities")
	scanCmd.Flags().Bool("dependency-audit", false, "Audit package dependencies")
	cmd.AddCommand(scanCmd)

	// docs subcommand
	cmd.AddCommand(&cobra.Command{
		Use:   "docs",
		Short: "Show detailed documentation (man page style)",
		Long:  `Display comprehensive documentation for the security command with examples and usage patterns.`,
		Run: func(cmd *cobra.Command, args []string) {
			showSecurityDocs()
		},
	})

	return cmd
}

func showSecurityDocs() {
	title := "SLOTH-RUNNER SYSADMIN SECURITY(1)"
	description := "sloth-runner sysadmin security - Security auditing and management"
	synopsis := "sloth-runner sysadmin security [subcommand] [options]"

	options := [][]string{
		{"audit", "Audit access logs and security events including failed authentication attempts and suspicious activity."},
		{"scan", "Scan for security vulnerabilities including CVE scanning, dependency audits, and configuration validation."},
		{"docs", "Show this documentation page."},
	}

	examples := [][]string{
		{
			"Security audit",
			"sloth-runner sysadmin security audit --since 24h --show-failed-auth",
			"Audits security logs for the last 24 hours, highlighting failed logins",
		},
		{
			"Full vulnerability scan",
			"sloth-runner sysadmin security scan --agent web-01 --full",
			"Performs comprehensive vulnerability scan on web-01",
		},
		{
			"Audit all agents",
			"sloth-runner sysadmin security audit --all-agents --output report.json",
			"Audits all agents and exports report as JSON",
		},
		{
			"CVE scanning",
			"sloth-runner sysadmin security scan --agent db-01 --cve-only",
			"Scans only for known CVE vulnerabilities",
		},
		{
			"Dependency audit",
			"sloth-runner sysadmin security scan --dependency-audit --all-agents",
			"Audits all package dependencies for security issues",
		},
		{
			"Suspicious activity detection",
			"sloth-runner sysadmin security audit --detect-anomalies --since 7d",
			"Analyzes logs for suspicious patterns over the last week",
		},
	}

	seeAlso := []string{
		"sloth-runner sysadmin logs - Log management",
		"sloth-runner sysadmin packages - Package management",
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
	pterm.FgGray.Println("For more information: sloth-runner sysadmin security --help")
}

// runAudit executa auditoria de segurança
func runAudit(since time.Duration, showFailedAuth bool, detectAnomalies bool, outputFormat string) error {
	scanner := NewScanner()

	pterm.DefaultHeader.WithFullWidth().Println("Security Audit")
	pterm.Println()

	spinner, _ := pterm.DefaultSpinner.Start("Analyzing security logs...")

	options := AuditOptions{
		Since:           since,
		ShowFailedAuth:  showFailedAuth,
		DetectAnomalies: detectAnomalies,
		OutputFormat:    outputFormat,
	}

	report, err := scanner.Audit(options)
	if err != nil {
		spinner.Fail("Audit failed")
		return err
	}

	spinner.Success("✅ Security audit completed")
	pterm.Println()

	// Output format
	if outputFormat == "json" {
		data, _ := json.MarshalIndent(report, "", "  ")
		fmt.Println(string(data))
		return nil
	}

	// Summary box
	pterm.DefaultBox.WithTitle("Audit Summary").WithTitleTopCenter().Println(
		fmt.Sprintf("Period: Last %v\nTotal Events: %d\nFailed Auth Attempts: %d\nSuspicious Events: %d",
			since, report.TotalEvents, report.FailedAuthAttempts, len(report.SuspiciousEvents)),
	)
	pterm.Println()

	// Suspicious events
	if len(report.SuspiciousEvents) > 0 {
		pterm.DefaultSection.Println("Suspicious Events")

		tableData := pterm.TableData{
			{"Timestamp", "Type", "Source", "Severity", "Description"},
		}

		for _, event := range report.SuspiciousEvents {
			tableData = append(tableData, []string{
				event.Timestamp.Format("2006-01-02 15:04:05"),
				event.Type,
				event.Source,
				getSeverityColor(event.Severity).Sprint(event.Severity),
				event.Description,
			})
		}

		pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
		pterm.Println()
	} else {
		pterm.Success.Println("No suspicious events detected ✓")
		pterm.Println()
	}

	// Recommendations
	if len(report.Recommendations) > 0 {
		pterm.DefaultSection.Println("Recommendations")
		for _, rec := range report.Recommendations {
			if rec == "No security issues detected - system appears secure" {
				pterm.Success.Println("  • " + rec)
			} else {
				pterm.Info.Println("  • " + rec)
			}
		}
		pterm.Println()
	}

	return nil
}

// runScan executa scan de vulnerabilidades
func runScan(agent string, full bool, cveOnly bool, dependencyAudit bool) error {
	scanner := NewScanner()

	pterm.DefaultHeader.WithFullWidth().Println("Security Scan")
	pterm.Println()

	scanType := "Quick scan"
	if full {
		scanType = "Full security scan"
	} else if cveOnly {
		scanType = "CVE vulnerability scan"
	} else if dependencyAudit {
		scanType = "Dependency audit"
	}

	pterm.Info.Printf("Running %s on agent '%s'...\n", scanType, agent)
	pterm.Println()

	spinner, _ := pterm.DefaultSpinner.Start("Scanning for vulnerabilities...")

	options := ScanOptions{
		Agent:           agent,
		Full:            full,
		CVEOnly:         cveOnly,
		DependencyAudit: dependencyAudit,
	}

	report, err := scanner.Scan(options)
	if err != nil {
		spinner.Fail("Scan failed")
		return err
	}

	spinner.Success("✅ Security scan completed")
	pterm.Println()

	// Score box
	scoreColor := getScoreColor(report.Score)
	pterm.DefaultBox.WithTitle("Security Score").WithTitleTopCenter().Println(
		scoreColor.Sprintf("%d/100 - %s", report.Score, report.Severity),
	)
	pterm.Println()

	// Vulnerabilities
	if len(report.Vulnerabilities) > 0 {
		pterm.DefaultSection.Println("Vulnerabilities Found")

		tableData := pterm.TableData{
			{"CVE", "Package", "Version", "Severity", "Fix Version"},
		}

		for _, vuln := range report.Vulnerabilities {
			tableData = append(tableData, []string{
				vuln.CVE,
				vuln.Package,
				vuln.Version,
				getSeverityColor(vuln.Severity).Sprint(vuln.Severity),
				vuln.FixVersion,
			})
		}

		pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
		pterm.Println()
	}

	// Configuration issues
	if len(report.ConfigIssues) > 0 {
		pterm.DefaultSection.Println("Configuration Issues")

		for _, issue := range report.ConfigIssues {
			pterm.Warning.Printf("  • %s: %s (%s)\n", issue.File, issue.Issue, issue.Severity)
			pterm.Println("    " + issue.Description)
		}
		pterm.Println()
	}

	// Permission issues
	if len(report.PermissionIssues) > 0 {
		pterm.DefaultSection.Println("Permission Issues")

		tableData := pterm.TableData{
			{"Path", "Current", "Expected", "Severity", "Description"},
		}

		for _, issue := range report.PermissionIssues {
			tableData = append(tableData, []string{
				issue.Path,
				issue.CurrentPerm,
				issue.ExpectedPerm,
				getSeverityColor(issue.Severity).Sprint(issue.Severity),
				issue.Description,
			})
		}

		pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
		pterm.Println()
	}

	// Summary
	if len(report.Vulnerabilities) == 0 && len(report.ConfigIssues) == 0 && len(report.PermissionIssues) == 0 {
		pterm.Success.Println("No security issues found - system is secure ✓")
	} else {
		totalIssues := len(report.Vulnerabilities) + len(report.ConfigIssues) + len(report.PermissionIssues)
		pterm.Warning.Printf("Found %d security issue(s) - review and fix as needed\n", totalIssues)
	}

	return nil
}

// getSeverityColor retorna cor para severity
func getSeverityColor(severity SeverityLevel) pterm.Color {
	switch severity {
	case SeverityCritical:
		return pterm.FgRed
	case SeverityHigh:
		return pterm.FgLightRed
	case SeverityMedium:
		return pterm.FgYellow
	case SeverityLow:
		return pterm.FgCyan
	case SeverityInfo:
		return pterm.FgGreen
	default:
		return pterm.FgWhite
	}
}

// getScoreColor retorna cor para score
func getScoreColor(score int) pterm.Color {
	if score >= 80 {
		return pterm.FgGreen
	} else if score >= 60 {
		return pterm.FgCyan
	} else if score >= 40 {
		return pterm.FgYellow
	}
	return pterm.FgRed
}
