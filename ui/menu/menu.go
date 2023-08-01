package menu

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zhengkyl/gol/game"
	"github.com/zhengkyl/gol/ui/common"
	"github.com/zhengkyl/gol/ui/keybinds"
)

const title = `
████████    ███████    ██
██    ██    ██ f ██    ██
████████    ███████    ██ ife
  ame ██               █████████
████████
`

type Model struct {
	playerId       int
	gm             *game.Manager
	common         common.Common
	lobbyInfos     []game.LobbyInfo
	options        []listItem
	activeIndex    int
	scrollIndex    int
	visibleOptions int
}

func New(common common.Common, gm *game.Manager, playerId int) *Model {
	options := make([]listItem, 0, 2)
	options = append(options,
		listItem{
			titleLeft:  "Play singleplayer game",
			titleRight: "",
			descLeft:   "Classic Conway's Game of Life",
			descRight:  "",
		},
		listItem{
			titleLeft:  "Create multiplayer lobby",
			titleRight: "",
			descLeft:   "Play with up to 10 players",
			descRight:  "",
		},
	)

	return &Model{common: common, gm: gm, options: options, playerId: playerId}
}

func (m *Model) SetSize(width, height int) {
	m.common.Width = width
	m.common.Height = height - 7

	m.visibleOptions = (m.common.Height - 1) / 4
}

func (m *Model) Init() tea.Cmd {
	return func() tea.Msg {
		return m.gm.LobbyInfos()
	}
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.SetSize(msg.Width, msg.Height)
	case []game.LobbyInfo:
		m.lobbyInfos = msg
		m.options = m.options[:2]
		for _, status := range msg {
			m.options = append(m.options, listItem{
				titleLeft:  fmt.Sprintf("Join lobby: %v", status.Name),
				titleRight: fmt.Sprintf("%v/%v players", status.PlayerCount, status.MaxPlayers),
				descLeft:   fmt.Sprintf("id: %v", status.Id),
			})
		}
		if m.activeIndex >= len(m.options) {
			m.activeIndex = len(m.options) - 1
			m.scrollIndex = m.activeIndex - m.visibleOptions + 1
			if m.scrollIndex < 0 {
				m.scrollIndex = 0
			}
		}
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keybinds.KeyBinds.Down):
			m.activeIndex = (m.activeIndex + 1 + len(m.options)) % len(m.options)
			if m.activeIndex < m.scrollIndex {
				m.scrollIndex = m.activeIndex
			}
			if m.activeIndex >= m.scrollIndex+m.visibleOptions {
				m.scrollIndex = m.activeIndex - m.visibleOptions + 1
			}
		case key.Matches(msg, keybinds.KeyBinds.Up):
			m.activeIndex = (m.activeIndex - 1 + len(m.options)) % len(m.options)
			if m.activeIndex < m.scrollIndex {
				m.scrollIndex = m.activeIndex
			}
			if m.activeIndex >= m.scrollIndex+m.visibleOptions {
				m.scrollIndex = m.activeIndex - m.visibleOptions + 1
			}
		case key.Matches(msg, keybinds.KeyBinds.Enter):
			switch m.activeIndex {
			case 0:
				return m, func() tea.Msg { return game.SoloGameMsg{} }
			case 1:
				lid := m.gm.CreateLobby()
				return m, func() tea.Msg { return m.gm.JoinLobby(lid, m.playerId) }
			default:
				activeId := m.lobbyInfos[m.activeIndex-2].Id
				return m, func() tea.Msg { return m.gm.JoinLobby(activeId, m.playerId) }
			}
		}
	}
	return m, nil
}

func alignLeftRight(left, right string, width int) string {
	leftW := lipgloss.Width(left)
	rightW := lipgloss.Width(right)

	spaces := width - (leftW + rightW)

	if spaces < 1 {
		if leftW > width {
			return left[:width-1] + "…"
		}
		return left
	} else {
		return left + strings.Repeat(" ", spaces) + right
	}
}

var (
	itemStyle        = lipgloss.NewStyle().Border(lipgloss.HiddenBorder(), true).Padding(0, 1)
	activeItemStyle  = lipgloss.NewStyle().Border(lipgloss.RoundedBorder(), true).Padding(0, 1)
	titleStyle       = lipgloss.NewStyle().Bold(true)
	activeTitleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("207"))
	descStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
)

func (m *Model) View() string {
	viewSb := strings.Builder{}
	itemSb := strings.Builder{}

	contentWidth := m.common.Width
	if contentWidth > 60 {
		contentWidth = 60
	}

	viewStyle := lipgloss.NewStyle().MarginLeft((m.common.Width - contentWidth) / 2)

	for i := m.scrollIndex; i < m.scrollIndex+m.visibleOptions && i < len(m.options); i++ {
		li := m.options[i]
		titleStyle := titleStyle
		itemStyle := itemStyle
		if i == m.activeIndex {
			titleStyle = activeTitleStyle
			itemStyle = activeItemStyle
		}
		// factor in border + margin
		itemSb.WriteString(titleStyle.Render(alignLeftRight(li.titleLeft, li.titleRight, contentWidth-4)))
		itemSb.WriteString("\n")
		itemSb.WriteString(descStyle.Render(alignLeftRight(li.descLeft, li.descRight, contentWidth-4)))

		viewSb.WriteString(itemStyle.Render(itemSb.String()))
		viewSb.WriteString("\n")

		itemSb.Reset()
	}

	titleStr := title
	titleLeftPad := (m.common.Width - 32) / 2
	if titleLeftPad > 0 {
		titleStr = lipgloss.NewStyle().MarginLeft(titleLeftPad).Render(titleStr)
	}

	return titleStr + "\n" + viewStyle.Render(viewSb.String())
}

type listItem struct {
	titleLeft  string
	titleRight string
	descLeft   string
	descRight  string
}
