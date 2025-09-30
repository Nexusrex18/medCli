package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Nexusrex18/medCli/internal/client"
	"github.com/Nexusrex18/medCli/internal/config"
	"github.com/Nexusrex18/medCli/internal/models"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	// "github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	logoColor     = lipgloss.Color("#37AC88")
	medStyle      = lipgloss.NewStyle().Foreground(logoColor).Bold(true)
	bridgeStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Bold(true)
	logoStyle     = lipgloss.NewStyle().Foreground(logoColor).Bold(true)
	block2Style   = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Bold(true)
	titleStyle    = lipgloss.NewStyle().Bold(true).Foreground(logoColor).Padding(1, 0)
	selectedStyle = lipgloss.NewStyle().Foreground(logoColor).Bold(true)

	// New UI Theme Colors
	primaryColor    = lipgloss.Color("#37AC88")
	secondaryColor  = lipgloss.Color("#2D9C78")
	accentColor     = lipgloss.Color("#45C9A3")
	backgroundColor = lipgloss.Color("#0A0F0D")
	surfaceColor    = lipgloss.Color("#151A17")
	borderColor     = lipgloss.Color("#2A3A34")
	textColor       = lipgloss.Color("#E8F3EF")
	mutedTextColor  = lipgloss.Color("#8A9E96")

	// Input Styles - Fixed BorderRadius issues
	inputStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(borderColor).
			Padding(0, 1).
			Background(surfaceColor).
			Foreground(textColor).
			Height(1).
			Width(60)

	inputFocusedStyle = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder()).
				BorderForeground(primaryColor).
				Padding(0, 1).
				Background(surfaceColor).
				Foreground(textColor).
				Height(1).
				Width(60)

	// Result Box Styles - Fixed BorderRadius issues
	resultBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(borderColor).
			Padding(1, 2).
			Background(surfaceColor).
			Foreground(textColor).
			Width(64)

	resultTitleStyle = lipgloss.NewStyle().
				Foreground(primaryColor).
				Bold(true)

	resultSubtitleStyle = lipgloss.NewStyle().
				Foreground(secondaryColor)

	resultTextStyle = lipgloss.NewStyle().
			Foreground(textColor)

	resultMutedStyle = lipgloss.NewStyle().
				Foreground(mutedTextColor)

	resultSeparator = lipgloss.NewStyle().
			Foreground(borderColor).
			SetString("┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈")

	// Status Bar Style
	statusStyle = lipgloss.NewStyle().
			Foreground(mutedTextColor).
			Italic(true)

	// Menu Styles - Fixed BorderRadius issues
	menuStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(borderColor).
			Padding(1).
			Background(surfaceColor)

	menuTitleStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			Padding(0, 1)

	menuItemStyle = lipgloss.NewStyle().
			Foreground(textColor).
			Padding(0, 1)

	menuSelectedStyle = lipgloss.NewStyle().
				Foreground(primaryColor).
				Background(lipgloss.Color("#1A2A24")).
				Bold(true).
				Padding(0, 1)

	popupStyle = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(primaryColor).
			Padding(1, 2).
			Background(surfaceColor).
			Foreground(textColor).
			Width(70).
			Height(20)

	popupTitleStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			Underline(true).
			Width(66).
			Align(lipgloss.Center)

	popupSectionStyle = lipgloss.NewStyle().
				Foreground(secondaryColor).
				Bold(true)

	popupTextStyle = lipgloss.NewStyle().
			Foreground(textColor)

	popupButtonStyle = lipgloss.NewStyle().
				Foreground(textColor).
				Background(primaryColor).
				Padding(0, 3).
				Bold(true)

	popupButtonInactiveStyle = lipgloss.NewStyle().
					Foreground(mutedTextColor).
					Border(lipgloss.NormalBorder()).
					BorderForeground(borderColor).
					Padding(0, 3)
)

var block1 = []string{
	" ██████   ██████ ██████████ ██████████  ",
	" ░░██████ ██████ ░░███░░░░░█░░███░░░░███ ",
	"  ░███░█████░███  ░███  █ ░  ░███   ░░███",
	"  ░███░░███ ░███  ░██████    ░███    ░███",
	"  ░███ ░░░  ░███  ░███░░█    ░███    ░███",
	"  ░███      ░███  ░███ ░   █ ░███    ███ ",
	"  █████     █████ ██████████ ██████████  ",
	" ░░░░░     ░░░░░ ░░░░░░░░░░ ░░░░░░░░░░   ",
	"                                          ",
	"                                          ",
	"                                          ",
}

var block2 = []string{
	" ███████████  ███████████   █████ ██████████     █████████  ██████████",
	" ░░███░░░░░███░░███░░░░░███ ░░███ ░░███░░░░███   ███░░░░░███░░███░░░░░█",
	"  ░███    ░███ ░███    ░███  ░███  ░███   ░░███ ███     ░░░  ░███  █ ░ ",
	"  ░██████████  ░██████████   ░███  ░███    ░███░███          ░██████   ",
	"  ░███░░░░░███ ░███░░░░░███  ░███  ░███    ░███░███    █████ ░███░░█   ",
	"  ░███    ░███ ░███    ░███  ░███  ░███    ███ ░░███  ░░███  ░███ ░   █",
	"  ███████████  █████   █████ █████ ██████████   ░░█████████  ██████████",
	" ░░░░░░░░░░░  ░░░░░   ░░░░░ ░░░░░ ░░░░░░░░░░     ░░░░░░░░░  ░░░░░░░░░░ ",
	"                                                                       ",
	"                                                                       ",
	"                                                                       ",
}

type model struct {
	showIntro      bool
	menu           list.Model
	config         *config.Config
	state          AppState
	client         *client.TM2Client
	err            error
	input          textinput.Model
	// viewport       viewport.Model
	results        string
	startTime      time.Time
	selectedRecord *models.MedicineRecord
	showPopup      bool
	// New fields for selection
	currentRecords []models.MedicineRecord
	selectedIndex  int
	lastSearchType string
	viewingResults bool
}

type AppState int

const (
	StateIntro AppState = iota
	StateMenu
	StateSearch
	StateSymptoms
	StateHealth
	StateSettings
	StateHelp
	StateError
	StatePopup
)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

func initialModel() model {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	items := []list.Item{
		item{title: "Search Traditional Medicine Codes", desc: "Search by TM2 or traditional codes"},
		item{title: "Search by Symptoms", desc: "Find diseases by symptoms"},
		item{title: "Health Status Dashboard", desc: "View system health"},
		item{title: "Configuration & Settings", desc: "Manage CLI settings"},
		item{title: "Help & Documentation", desc: "Access documentation"},
	}

	// Create a custom delegate for the menu
	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = menuSelectedStyle
	delegate.Styles.SelectedDesc = menuSelectedStyle.Copy().Foreground(secondaryColor)
	delegate.Styles.NormalTitle = menuItemStyle
	delegate.Styles.NormalDesc = menuItemStyle.Copy().Foreground(mutedTextColor)

	menu := list.New(items, delegate, 60, 14)
	menu.Title = "🌿 TM2 Traditional Medicine CLI"
	menu.SetShowStatusBar(false)
	menu.SetShowFilter(false)
	menu.SetShowHelp(false)
	menu.Styles.Title = menuTitleStyle
	menu.Styles.NoItems = menuItemStyle.Copy().Italic(true)

	tm2Client, err := client.NewTM2Client(cfg)
	if err != nil {
		return model{
			showIntro: false,
			menu:      menu,
			config:    cfg,
			state:     StateError,
			err:       err,
		}
	}

	// Initialize styled text input
	ti := textinput.New()
	ti.Placeholder = "Enter code or symptoms..."
	ti.PlaceholderStyle = lipgloss.NewStyle().Foreground(mutedTextColor)
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 58
	ti.Prompt = "❯ "
	ti.PromptStyle = lipgloss.NewStyle().Foreground(primaryColor)
	ti.TextStyle = lipgloss.NewStyle().Foreground(textColor)

	return model{
		showIntro:      true,
		menu:           menu,
		config:         cfg,
		state:          StateIntro,
		client:         tm2Client,
		err:            nil,
		input:          ti,
		results:        "", // Empty results initially
		startTime:      time.Now(),
		currentRecords: nil,
		selectedIndex:  0,
		lastSearchType: "",
		viewingResults: false,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch m.state {
	case StateError:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if msg.String() == "q" {
				return m, tea.Quit
			}
		}

	case StateIntro:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "q":
				return m, tea.Quit
			case "enter", " ":
				if m.err != nil {
					m.state = StateError
				} else {
					m.state = StateMenu
				}
				m.showIntro = false
			}
		}

	case StateMenu:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "q":
				return m, tea.Quit
			case "enter":
				selectedItem := m.menu.SelectedItem().(item)
				switch selectedItem.title {
				case "Search Traditional Medicine Codes":
					m.state = StateSearch
					m.input.Reset()
					m.input.Placeholder = "Enter TM2 code..."
					m.input.Focus()
				case "Search by Symptoms":
					m.state = StateSymptoms
					m.input.Reset()
					m.input.Placeholder = "Enter symptoms (comma-separated)..."
					m.input.Focus()
				case "Health Status Dashboard":
					m.state = StateHealth
					m.results = m.getHealthStatus()
					// m.viewport.SetContent(m.results)
				case "Configuration & Settings":
					m.state = StateSettings
				case "Help & Documentation":
					m.state = StateHelp
				}
			}
		}
		m.menu, cmd = m.menu.Update(msg)
		return m, cmd

	case StateSearch:
        switch msg := msg.(type) {
        case tea.KeyMsg:
            switch msg.String() {
            case "q":
                m.state = StateMenu
                m.currentRecords = nil
                m.selectedIndex = 0
                m.viewingResults = false
                m.results = ""
            case "up", "k":
                if m.viewingResults && len(m.currentRecords) > 0 {
                    m.selectedIndex--
                    if m.selectedIndex < 0 {
                        m.selectedIndex = len(m.currentRecords) - 1
                    }
                    m.results = formatSearchResults(m.currentRecords, m.selectedIndex)
                }
            case "down", "j":
                if m.viewingResults && len(m.currentRecords) > 0 {
                    m.selectedIndex++
                    if m.selectedIndex >= len(m.currentRecords) {
                        m.selectedIndex = 0
                    }
                    m.results = formatSearchResults(m.currentRecords, m.selectedIndex)
                }
            case "enter":
                if m.viewingResults && len(m.currentRecords) > 0 {
                    // Show popup for selected record
                    m.selectedRecord = &m.currentRecords[m.selectedIndex]
                    m.showPopup = true
                    m.state = StatePopup
                } else {
                    // Perform new search
                    query := m.input.Value()
                    ctx := context.Background()
                    result, err := m.client.SearchByCode(ctx, query, "both")
                    if err != nil {
                        m.results = fmt.Sprintf("Error: %v", err)
                        m.currentRecords = nil
                        m.viewingResults = false
                    } else {
                        m.currentRecords = result.Records
                        m.selectedIndex = 0
                        m.lastSearchType = "code"
                        m.viewingResults = true
                        m.results = formatSearchResults(m.currentRecords, m.selectedIndex)
                    }
                }
            }
        }
        m.input, cmd = m.input.Update(msg)
        return m, cmd

	case StateSymptoms:
        switch msg := msg.(type) {
        case tea.KeyMsg:
            switch msg.String() {
            case "q":
                m.state = StateMenu
                m.currentRecords = nil
                m.selectedIndex = 0
                m.viewingResults = false
                m.results = ""
            case "up", "k":
                if m.viewingResults && len(m.currentRecords) > 0 {
                    m.selectedIndex--
                    if m.selectedIndex < 0 {
                        m.selectedIndex = len(m.currentRecords) - 1
                    }
                    m.results = formatSymptomResults(m.currentRecords, m.selectedIndex)
                }
            case "down", "j":
                if m.viewingResults && len(m.currentRecords) > 0 {
                    m.selectedIndex++
                    if m.selectedIndex >= len(m.currentRecords) {
                        m.selectedIndex = 0
                    }
                    m.results = formatSymptomResults(m.currentRecords, m.selectedIndex)
                }
            case "enter":
                if m.viewingResults && len(m.currentRecords) > 0 {
                    // Show popup for selected record
                    m.selectedRecord = &m.currentRecords[m.selectedIndex]
                    m.showPopup = true
                    m.state = StatePopup
                } else {
                    // Perform new search
                    symptoms := strings.Split(m.input.Value(), ",")
                    for i := range symptoms {
                        symptoms[i] = strings.TrimSpace(symptoms[i])
                    }
                    ctx := context.Background()
                    result, err := m.client.SearchBySymptoms(ctx, symptoms)
                    if err != nil {
                        m.results = fmt.Sprintf("Error: %v", err)
                        m.currentRecords = nil
                        m.viewingResults = false
                    } else {
                        m.currentRecords = result.Records
                        m.selectedIndex = 0
                        m.lastSearchType = "symptoms"
                        m.viewingResults = true
                        m.results = formatSymptomResults(m.currentRecords, m.selectedIndex)
                    }
                }
            }
        }
        m.input, cmd = m.input.Update(msg)
        return m, cmd

	case StateHealth:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if msg.String() == "q" {
				m.state = StateMenu
			}
			if msg.String() == "r" { // Refresh health status
				m.results = m.getHealthStatus()
				// m.viewport.SetContent(m.results)
			}
		}
		// m.viewport, cmd = m.viewport.Update(msg)
		return m, cmd

	case StatePopup:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "q", "esc", "enter":
				m.showPopup = false
				m.selectedRecord = nil
				// Return to previous state based on last search type
				if m.lastSearchType == "code" {
					m.state = StateSearch
				} else {
					m.state = StateSymptoms
				}
				// Keep viewingResults as true since we're still viewing results
			}
		}
		return m, nil

	default:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if msg.String() == "q" {
				m.state = StateMenu
			}
		}
		m.input, cmd = m.input.Update(msg)
		// m.viewport, cmd = m.viewport.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m model) View() string {
	if m.state == StateError {
		errorBox := resultBoxStyle.Copy().
			BorderForeground(lipgloss.Color("#FF6B6B")).
			Width(50).
			Height(8).
			Align(lipgloss.Center).
			Render(
				lipgloss.JoinVertical(lipgloss.Center,
					"🚨 Error",
					"",
					lipgloss.NewStyle().Foreground(textColor).Render(m.err.Error()),
					"",
					statusStyle.Render("Press 'q' to quit"),
				),
			)
		return lipgloss.Place(80, 20, lipgloss.Center, lipgloss.Center, errorBox)
	}

	if m.showIntro {
		return renderIntro()
	}

	// Show popup on top of everything if active
	if m.showPopup && m.selectedRecord != nil {
		return m.renderPopup()
	}

	switch m.state {
	case StateMenu:
		menuContainer := lipgloss.NewStyle().
			Padding(1, 0).
			Render(
				lipgloss.JoinVertical(lipgloss.Center,
					titleStyle.Render("🌿 TM2 Traditional Medicine CLI"),
					"",
					menuStyle.Render(m.menu.View()),
					"",
					statusStyle.Render("[↑↓] Navigate • [Enter] Select • [q] Quit"),
				),
			)
		return lipgloss.Place(80, 24, lipgloss.Center, lipgloss.Center, menuContainer)

	case StateSearch:
        inputDisplay := m.input.View()
        if m.input.Focused() {
            inputDisplay = inputFocusedStyle.Render(inputDisplay)
        } else {
            inputDisplay = inputStyle.Render(inputDisplay)
        }

        statusMsg := "[Enter] Search • [q] Back to Menu"
        if m.viewingResults {
            statusMsg = "[↑↓] Navigate • [Enter] View Details • [q] Back to Menu"
        }

        // Create results area with fixed height
        resultsArea := resultBoxStyle.Copy().Height(14).Render(m.results)

        content := lipgloss.JoinVertical(lipgloss.Left,
            titleStyle.Render("🔍 Search Traditional Medicine Codes"),
            "",
            inputDisplay,
            "",
            resultsArea, // Use the results directly
            "",
            statusStyle.Render(statusMsg),
        )
        return lipgloss.Place(80, 24, lipgloss.Center, lipgloss.Center, content)

	case StateSymptoms:
        inputDisplay := m.input.View()
        if m.input.Focused() {
            inputDisplay = inputFocusedStyle.Render(inputDisplay)
        } else {
            inputDisplay = inputStyle.Render(inputDisplay)
        }

        statusMsg := "[Enter] Search • [q] Back to Menu"
        if m.viewingResults {
            statusMsg = "[↑↓] Navigate • [Enter] View Details • [q] Back to Menu"
        }

        // Create results area with fixed height
        resultsArea := resultBoxStyle.Copy().Height(14).Render(m.results)

        content := lipgloss.JoinVertical(lipgloss.Left,
            titleStyle.Render("🤒 Search by Symptoms"),
            "",
            inputDisplay,
            "",
            resultsArea, // Use the results directly
            "",
            statusStyle.Render(statusMsg),
        )
        return lipgloss.Place(80, 24, lipgloss.Center, lipgloss.Center, content)
        
    case StateHealth:
		content := lipgloss.JoinVertical(lipgloss.Left,
			titleStyle.Render("📊 Health Status Dashboard"),
			"",
			// m.viewport.View(),
			"",
			statusStyle.Render("[r] Refresh • [q] Back to Menu"),
		)
		return lipgloss.Place(80, 24, lipgloss.Center, lipgloss.Center, content)

	case StateSettings:
		settingsBox := resultBoxStyle.Copy().
			Height(10).
			Width(50).
			Align(lipgloss.Center).
			Render(
				lipgloss.JoinVertical(lipgloss.Center,
					"⚙️ Configuration & Settings",
					"",
					resultMutedStyle.Render("Settings panel coming soon..."),
					"",
					statusStyle.Render("[q] Back to Menu"),
				),
			)
		return lipgloss.Place(80, 24, lipgloss.Center, lipgloss.Center, settingsBox)

	case StateHelp:
		helpBox := resultBoxStyle.Copy().
			Height(10).
			Width(50).
			Align(lipgloss.Center).
			Render(
				lipgloss.JoinVertical(lipgloss.Center,
					"📖 Help & Documentation",
					"",
					resultMutedStyle.Render("Documentation coming soon..."),
					"",
					statusStyle.Render("[q] Back to Menu"),
				),
			)
		return lipgloss.Place(80, 24, lipgloss.Center, lipgloss.Center, helpBox)
	}
	return ""
}

func (s AppState) String() string {
	switch s {
	case StateSearch:
		return "Search Traditional Medicine Codes"
	case StateSymptoms:
		return "Search by Symptoms"
	case StateHealth:
		return "Health Status Dashboard"
	default:
		return "Unknown State"
	}
}

// Update health status to show CSV stats
func (m *model) getHealthStatus() string {
	stats := m.client.GetRepoStats()
	hits, misses, items := m.client.GetCacheStats()
	uptime := time.Since(m.startTime).Truncate(time.Second)

	healthBox := lipgloss.JoinVertical(lipgloss.Left,
		resultTitleStyle.Render("📊 System Health Status"),
		"",
		resultSubtitleStyle.Render("📈 Data Statistics:"),
		resultTextStyle.Render(fmt.Sprintf("   📚 Total Records: %d", stats["total_records"])),
		resultTextStyle.Render(fmt.Sprintf("   🏷️  Unique Codes: %d", stats["unique_codes"])),
		resultTextStyle.Render(fmt.Sprintf("   🔢 Unique TM2 Codes: %d", stats["unique_tm2_codes"])),
		"",
		resultSubtitleStyle.Render("⚡ Performance:"),
		resultTextStyle.Render(fmt.Sprintf("   💾 Cache Hits: %d", hits)),
		resultTextStyle.Render(fmt.Sprintf("   🔍 Cache Misses: %d", misses)),
		resultTextStyle.Render(fmt.Sprintf("   📦 Cache Items: %d", items)),
		resultTextStyle.Render(fmt.Sprintf("   ⏱️  Uptime: %s", uptime)),
		"",
		resultMutedStyle.Render("💡 Tip: Press 'r' to refresh stats"),
	)

	return resultBoxStyle.Copy().Height(16).Render(healthBox)
}

func formatSearchResults(records []models.MedicineRecord, selectedIndex int) string {
	if len(records) == 0 {
		noResults := resultBoxStyle.Copy().
			BorderForeground(mutedTextColor).
			Height(6).
			Align(lipgloss.Center).
			Render(
				lipgloss.JoinVertical(lipgloss.Center,
					"🔍 No results found",
					resultMutedStyle.Render("Try a different code or term"),
				),
			)
		return noResults
	}

	// Limit to first 10 results to prevent overload
	displayRecords := records
	if len(records) > 10 {
		displayRecords = records[:10]
	}

	var results []string
	results = append(results,
		resultTitleStyle.Render(fmt.Sprintf("📋 Found %d results (showing first 10)", len(records))),
		resultMutedStyle.Render("↑↓ to navigate • Enter to view details • q to back"),
		"",
	)

	for i, record := range displayRecords {
		// Highlight selected item
		titleStyle := resultTitleStyle
		subtitleStyle := resultSubtitleStyle
		textStyle := resultTextStyle
		
		if i == selectedIndex {
			titleStyle = titleStyle.Copy().Foreground(accentColor).Bold(true)
			subtitleStyle = subtitleStyle.Copy().Foreground(accentColor)
			textStyle = textStyle.Copy().Foreground(accentColor)
		}

		resultBox := lipgloss.NewStyle().
			Padding(0, 0).
			Render(
				lipgloss.JoinVertical(lipgloss.Left,
					titleStyle.Render(fmt.Sprintf("%d. %s", i+1, record.TM2Title)),
					subtitleStyle.Render(fmt.Sprintf("   🏷️  TM2 Code: %s • Traditional: %s", record.TM2Code, record.Code)),
					subtitleStyle.Render(fmt.Sprintf("   📁 Type: %s • Confidence: %.1f%%", record.Type, record.ConfidenceScore*100)),
					"",
					textStyle.Render("   📖 Definition:"),
					textStyle.Render("   "+wrapText(record.TM2Definition, 56, "   ")),
				),
			)

		results = append(results, resultBox)

		// Add separator between results (except for the last one)
		if i < len(displayRecords)-1 {
			results = append(results, resultSeparator.String())
		}
	}

	return lipgloss.JoinVertical(lipgloss.Left, results...)
}

func formatSymptomResults(records []models.MedicineRecord, selectedIndex int) string {
	if len(records) == 0 {
		noResults := resultBoxStyle.Copy().
			BorderForeground(mutedTextColor).
			Height(6).
			Align(lipgloss.Center).
			Render(
				lipgloss.JoinVertical(lipgloss.Center,
					"🤒 No matches found",
					resultMutedStyle.Render("Try different symptoms or terms"),
				),
			)
		return noResults
	}

	// Limit to first 10 results to prevent overload
	displayRecords := records
	if len(records) > 10 {
		displayRecords = records[:10]
	}

	var results []string
	results = append(results,
		resultTitleStyle.Render(fmt.Sprintf("🎯 Found %d matches (showing first 10)", len(records))),
		resultMutedStyle.Render("↑↓ to navigate • Enter to view details • q to back"),
		"",
	)

	for i, record := range displayRecords {
		// Highlight selected item
		titleStyle := resultTitleStyle
		subtitleStyle := resultSubtitleStyle
		
		if i == selectedIndex {
			titleStyle = titleStyle.Copy().Foreground(accentColor).Bold(true)
			subtitleStyle = subtitleStyle.Copy().Foreground(accentColor)
		}

		resultBox := lipgloss.NewStyle().
			Padding(0, 0).
			Render(
				lipgloss.JoinVertical(lipgloss.Left,
					titleStyle.Render(fmt.Sprintf("%d. %s", i+1, record.TM2Title)),
					subtitleStyle.Render(fmt.Sprintf("   🏷️  TM2 Code: %s • Traditional: %s", record.TM2Code, record.Code)),
					subtitleStyle.Render(fmt.Sprintf("   📁 Type: %s • Confidence: %.1f%%", record.Type, record.ConfidenceScore*100)),
				),
			)

		results = append(results, resultBox)

		// Add separator between results (except for the last one)
		if i < len(displayRecords)-1 {
			results = append(results, resultSeparator.String())
		}
	}

	return lipgloss.JoinVertical(lipgloss.Left, results...)
}

// Helper function to wrap text
func wrapText(text string, width int, prefix string) string {
	if text == "" {
		return prefix + "No description available"
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return ""
	}

	var lines []string
	currentLine := prefix

	for _, word := range words {
		if len(currentLine)+len(word)+1 > width {
			lines = append(lines, currentLine)
			currentLine = prefix + word
		} else {
			if currentLine != prefix {
				currentLine += " "
			}
			currentLine += word
		}
	}

	if currentLine != prefix {
		lines = append(lines, currentLine)
	}

	return strings.Join(lines, "\n")
}

func renderIntro() string {
	var b strings.Builder

	// Render blocks side by side with proper styling
	maxLines := max(len(block1), len(block2))
	for i := 0; i < maxLines; i++ {
		var line strings.Builder

		// Add block1 (green)
		if i < len(block1) {
			line.WriteString(logoStyle.Render(block1[i]))
		} else {
			line.WriteString(strings.Repeat(" ", len(block1[0])))
		}

		// Add spacing between blocks
		line.WriteString("  ")

		// Add block2 (white)
		if i < len(block2) {
			line.WriteString(block2Style.Render(block2[i]))
		}

		// Center the entire line
		fullLine := line.String()
		b.WriteString(lipgloss.PlaceHorizontal(120, lipgloss.Center, fullLine))
		b.WriteString("\n")
	}

	// Add "Medbridge" text - centered properly
	b.WriteString("\n")
	medbridgeText := medStyle.Render("Med") + bridgeStyle.Render("bridge")
	b.WriteString(lipgloss.PlaceHorizontal(120, lipgloss.Center, medbridgeText))
	b.WriteString("\n\n")

	// Styled press enter text
	pressEnterBox := resultBoxStyle.Copy().
		Width(40).
		Height(5).
		Align(lipgloss.Center).
		Render(
			lipgloss.JoinVertical(lipgloss.Center,
				"🚀 Welcome to MedBridge CLI",
				resultMutedStyle.Render("Traditional Medicine Terminology"),
				"",
				resultTextStyle.Render("Press Enter to continue..."),
			),
		)

	b.WriteString(lipgloss.PlaceHorizontal(120, lipgloss.Center, pressEnterBox))

	// Use a wider container to prevent cutoff
	return lipgloss.Place(120, 28, lipgloss.Center, lipgloss.Center, b.String())
}

func (m model) renderPopup() string {
	if m.selectedRecord == nil {
		return ""
	}

	record := m.selectedRecord

	content := lipgloss.JoinVertical(lipgloss.Left,
		popupTitleStyle.Render(record.TM2Title),
		resultMutedStyle.Render(fmt.Sprintf("Item %d of %d", m.selectedIndex+1, len(m.currentRecords))),
		"",
		popupSectionStyle.Render("Code Information:"),
		popupTextStyle.Render(fmt.Sprintf("   TM2 Code: %s", record.TM2Code)),
		popupTextStyle.Render(fmt.Sprintf("   Traditional Code: %s", record.Code)),
		popupTextStyle.Render(fmt.Sprintf("   Medicine Type: %s", record.Type)),
		popupTextStyle.Render(fmt.Sprintf("   Confidence Score: %.1f%%", record.ConfidenceScore*100)),
		"",
		popupSectionStyle.Render("TM2 Definition:"),
		popupTextStyle.Render(wrapText(record.TM2Definition, 66, "")),
		"",
		popupSectionStyle.Render("Traditional Description:"),
		popupTextStyle.Render(wrapText(record.Description, 66, "")),
		"",
		popupSectionStyle.Render("Code Title:"),
		popupTextStyle.Render(record.CodeTitle),
		"",
		lipgloss.PlaceHorizontal(66, lipgloss.Center,
			popupButtonStyle.Render("Press Enter/ESC to close")),
	)

	return lipgloss.Place(80, 24, lipgloss.Center, lipgloss.Center,
		popupStyle.Render(content))
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error starting TM2 CLI: %v\n", err)
		os.Exit(1)
	}
}