package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/beevik/ntp"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// Version information (set by GoReleaser)
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
	builtBy = "unknown"
)

// Styles
var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("86")).
			Bold(true).
			Padding(1, 2)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Italic(true)

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("10")).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("9")).
			Bold(true)

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("12"))

	labelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("8"))

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("86")).
			Padding(1, 2).
			Margin(1, 0)
)

type state int

const (
	configState state = iota
	testingState
	resultState
	versionState
)

type formData struct {
	socks5URL string
	ntpServer string
}

type model struct {
	state    state
	config   *config
	formData *formData
	form     *huh.Form
	spinner  spinner.Model
	result   *ntp.Response
	err      error
	quitting bool
}

type testCompleteMsg struct {
	result *ntp.Response
	err    error
}

func initialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("86"))

	cfg := &config{
		socks5Config: socks5Config{
			address:  "localhost:1080",
			username: "",
			password: "",
		},
		ntpAddress: "time.google.com:123",
	}

	// Create form data that will persist
	data := &formData{
		socks5URL: "socks5://localhost:1080",
		ntpServer: "time.google.com:123",
	}

	// Create the form and bind to the data fields
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("SOCKS5 Connection URL").
				Description("Enter SOCKS5 proxy URL (see supported formats below)").
				Value(&data.socks5URL).
				Placeholder("socks5://user:pass@host:port").
				Validate(func(s string) error {
					_, err := parseSocks5String(s)
					if err != nil {
						return fmt.Errorf("invalid format: %v", err)
					}
					return nil
				}),

			huh.NewNote().
				Title("Supported URL Formats:").
				Description(`• socks5://host:port (no authentication)
• socks5://username:password@host:port
• socks5://host:port:username:password
• socks5:host:port (no authentication)
• socks5:host:port:username:password`),

			huh.NewInput().
				Title("NTP Server").
				Description("Enter the NTP server to test UDP connectivity").
				Value(&data.ntpServer).
				Placeholder("time.google.com:123"),
		),
	)

	return model{
		state:    configState,
		config:   cfg,
		formData: data,
		form:     form,
		spinner:  s,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.form.Init(), m.spinner.Tick)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			if m.state == resultState {
				m.state = configState
				m.result = nil
				m.err = nil

				return m, m.form.Init()
			}
		case "esc":
			if m.state == versionState {
				m.state = configState
				return m, nil
			}
		case "v":
			if m.state == configState || m.state == resultState {
				// Show version info in TUI modal
				m.state = versionState
				return m, nil
			}
		}

	case testCompleteMsg:
		m.state = resultState
		m.result = msg.result
		m.err = msg.err
		return m, nil

	case spinner.TickMsg:
		if m.state == testingState {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
	}

	switch m.state {
	case configState:
		form, cmd := m.form.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.form = f
		}

		if m.form.State == huh.StateCompleted {
			// Update config from form values
			m.config.ntpAddress = m.formData.ntpServer

			// Parse the SOCKS5 URL and update the config
			if parsedConfig, err := parseSocks5String(m.formData.socks5URL); err == nil {
				m.config.socks5Config = parsedConfig
			} else {
				// If parsing fails, use defaults (this shouldn't happen due to validation)
				m.config.socks5Config = socks5Config{
					address:  "localhost:1080",
					username: "",
					password: "",
				}
			}
			m.state = testingState
			return m, tea.Batch(cmd, m.spinner.Tick, m.runTest())
		}

		return m, cmd

	case testingState:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m model) runTest() tea.Cmd {
	return func() tea.Msg {
		result, err := performNTPTest(m.config)
		return testCompleteMsg{result: result, err: err}
	}
}

func (m model) View() string {
	if m.quitting {
		return ""
	}

	var content strings.Builder

	// Header
	content.WriteString(titleStyle.Render(fmt.Sprintf("🌐 SOCKS5 UDP Checker v%s", version)))
	content.WriteString("\n")
	content.WriteString(subtitleStyle.Render("Test UDP connectivity through SOCKS5 proxy using NTP"))
	content.WriteString("\n\n")

	switch m.state {
	case configState:
		content.WriteString(m.form.View())
		content.WriteString("\n")
		content.WriteString(labelStyle.Render("Press Enter to start test • 'v' for version • Ctrl+C to quit"))

	case testingState:
		content.WriteString(fmt.Sprintf("%s Testing UDP connectivity through SOCKS5 proxy...\n\n", m.spinner.View()))

		// Show parsed SOCKS5 configuration
		if m.config.username != "" {
			content.WriteString(infoStyle.Render(fmt.Sprintf("• Connecting to SOCKS5 proxy: %s (authenticated)", m.config.address)))
		} else {
			content.WriteString(infoStyle.Render(fmt.Sprintf("• Connecting to SOCKS5 proxy: %s (no auth)", m.config.address)))
		}
		content.WriteString("\n")
		content.WriteString(infoStyle.Render(fmt.Sprintf("• Testing NTP server: %s", m.config.ntpAddress)))
		content.WriteString("\n\n")
		content.WriteString(labelStyle.Render("Please wait..."))

	case resultState:
		if m.err != nil {
			content.WriteString(errorStyle.Render("❌ Test Failed"))
			content.WriteString("\n\n")
			content.WriteString(boxStyle.Render(fmt.Sprintf("Error: %s", m.err.Error())))
		} else {
			content.WriteString(successStyle.Render("✅ Test Successful"))
			content.WriteString("\n\n")
			content.WriteString(m.formatNTPResponse(m.result))
		}
		content.WriteString("\n\n")
		content.WriteString(labelStyle.Render("Press Enter to run another test • Ctrl+C to quit"))

	case versionState:
		content.WriteString(infoStyle.Render("📋 Version Information"))
		content.WriteString("\n\n")

		versionInfo := fmt.Sprintf(`SOCKS5 UDP Checker
  Version:  %s
  Commit:   %s
  Built:    %s
  Built by: %s`,
			version, commit, date, builtBy)

		content.WriteString(boxStyle.Render(versionInfo))
		content.WriteString("\n\n")
		content.WriteString(labelStyle.Render("Press Esc to go back • Ctrl+C to quit"))
	}

	return content.String()
}

func (m model) formatNTPResponse(resp *ntp.Response) string {
	var content strings.Builder

	// Basic Information
	basicInfo := fmt.Sprintf(`NTP Response Summary:
  Server Time:     %v
  Round Trip Time: %v
  Clock Offset:    %v
  Stratum:         %d`,
		resp.Time.Format("2006-01-02 15:04:05.000 MST"),
		resp.RTT,
		resp.ClockOffset,
		resp.Stratum)

	content.WriteString(boxStyle.Render(basicInfo))

	// Kiss of Death warning (if applicable)
	if resp.IsKissOfDeath() {
		content.WriteString("\n")
		kissInfo := fmt.Sprintf(`⚠️  Kiss of Death Response:
  Kiss Code:       %s`, resp.KissCode)

		content.WriteString(boxStyle.Copy().
			BorderForeground(lipgloss.Color("9")).
			Render(kissInfo))
	}

	return content.String()
}

func main() {
	// Set up the program
	p := tea.NewProgram(
		initialModel(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	// Start the program
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
