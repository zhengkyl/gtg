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

type model struct {
	gm           *game.Manager
	common       common.Common
	options      []listItem
	activeOption int
}

func New(common common.Common, gm *game.Manager) *model {
	return &model{common: common, gm: gm}
}

func (m *model) SetSize(width, height int) {
	m.common.Width = width
	m.common.Height = height
}

func (m *model) Init() tea.Cmd {
	return func() tea.Msg {
		return m.gm.LobbyStatuses()
	}
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.common.Width = msg.Width
		m.common.Height = msg.Height
	case []game.LobbyStatus:
		m.options = make([]listItem, len(msg))
		m.options = append(m.options, listItem{
			titleLeft:  "Play singleplayer game",
			titleRight: "",
			descLeft:   "Conway's game of life",
			descRight:  "",
		},
			listItem{
				titleLeft:  "Create multiplayer lobby",
				titleRight: "",
				descLeft:   "Play with up to 10 other players",
				descRight:  "",
			})
		for _, status := range msg {
			m.options = append(m.options, listItem{
				titleLeft:  status.Name,
				titleRight: fmt.Sprintf("%v/%v", status.PlayerCount, status.MaxPlayers),
				descLeft:   fmt.Sprint(status.Id),
			})
		}
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keybinds.KeyBinds.Down):
			m.activeOption = (m.activeOption + 1 + len(m.options)) % len(m.options)
		case key.Matches(msg, keybinds.KeyBinds.Up):
			m.activeOption = (m.activeOption - 1 + len(m.options)) % len(m.options)
		case key.Matches(msg, keybinds.KeyBinds.Enter):
			// m.gm.JoinLobby() // lobbyId + p *tea.Program
		}
	}
	return m, nil
}

func combine(left, right string, width int) string {
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

var itemStyle = lipgloss.NewStyle().Border(lipgloss.HiddenBorder(), true).Padding(0, 1)
var activeItemStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder(), true).Padding(0, 1)
var titleStyle = lipgloss.NewStyle().Bold(true)
var activeTitleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("207"))
var descStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("254"))

func (m *model) View() string {
	viewSb := strings.Builder{}
	itemSb := strings.Builder{}

	for i, li := range m.options {
		title := titleStyle
		item := itemStyle
		if i == m.activeOption {
			title = activeTitleStyle
			item = activeItemStyle
		}
		// factor in border + margin
		itemSb.WriteString(title.Render(combine(li.titleLeft, li.titleRight, m.common.Width-4)))
		itemSb.WriteString("\n")
		itemSb.WriteString(descStyle.Render(combine(li.descLeft, li.descRight, m.common.Width-4)))

		viewSb.WriteString(item.Render(itemSb.String()))
		viewSb.WriteString("\n")
		itemSb.Reset()
	}

	return viewSb.String()
}

type listItem struct {
	titleLeft  string
	titleRight string
	descLeft   string
	descRight  string
}
