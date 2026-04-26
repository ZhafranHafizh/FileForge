package interactive

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
)

type SummaryRow struct {
	Label string
	Value string
}

var (
	titleStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("45")).Bold(true)
	subtitleStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("250"))
	sectionStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("39")).Bold(true)
	helpStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	mutedStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	successStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Bold(true)
	errorStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("203")).Bold(true)
	warningStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("220")).Bold(true)
	summaryBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("39")).
			Padding(0, 1)
	headerBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("45")).
			Padding(0, 1)
)

func RenderHeader() string {
	body := lipgloss.JoinVertical(
		lipgloss.Left,
		titleStyle.Render("FileForge"),
		subtitleStyle.Render("Local-first file conversion toolkit"),
	)
	return headerBoxStyle.Render(body)
}

func RenderSection(title string) string {
	return sectionStyle.Render(title)
}

func RenderSummary(title string, rows []SummaryRow) string {
	lines := []string{sectionStyle.Render(title)}
	for _, row := range rows {
		lines = append(lines, formatSummaryRow(row))
	}
	return summaryBoxStyle.Render(strings.Join(lines, "\n"))
}

func RenderSuccess(outputPath string) string {
	body := "The process is complete. Please access the following path to view the results.\n" + outputPath
	return summaryBoxStyle.BorderForeground(lipgloss.Color("42")).Render(successStyle.Render("Success") + "\n" + body)
}

func RenderError(err error) string {
	if err == nil {
		return ""
	}
	return summaryBoxStyle.BorderForeground(lipgloss.Color("203")).Render(errorStyle.Render("Error") + "\n" + err.Error())
}

func RenderComingSoon(feature string) string {
	body := fmt.Sprintf("%s\n\nThis feature is not available yet.\nPlanned for a future release.", feature)
	return summaryBoxStyle.BorderForeground(lipgloss.Color("220")).Render(warningStyle.Render("Coming Soon") + "\n" + body)
}

func RenderHelp(text string) string {
	return helpStyle.Render(text)
}

func RenderMuted(text string) string {
	return mutedStyle.Render(text)
}

func formatSummaryRow(row SummaryRow) string {
	valueLines := strings.Split(row.Value, "\n")
	if len(valueLines) == 0 {
		return fmt.Sprintf("%-13s", row.Label)
	}

	lines := []string{fmt.Sprintf("%-13s %s", row.Label, valueLines[0])}
	for _, line := range valueLines[1:] {
		lines = append(lines, fmt.Sprintf("%-13s %s", "", line))
	}
	return strings.Join(lines, "\n")
}
