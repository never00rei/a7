package ui

import (
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m AppModel) contentWidth() int {
	frameWidth := m.frameWidth()
	if frameWidth > 4 {
		return frameWidth - 2
	}
	return 80
}

func (m AppModel) frameWidth() int {
	if m.width > 3 {
		return m.width - 2
	}
	return 80
}

func (m AppModel) frameHeight() int {
	if m.height > 3 {
		return m.height - 2
	}
	return 24
}

func (m AppModel) twoPane(left, right string) string {
	total := m.contentWidth()
	gapWidth := 2
	available := total - gapWidth - 4
	if available < 0 {
		available = 0
	}
	leftWidth := available / 2
	rightWidth := available - leftWidth

	leftStyle := lipgloss.NewStyle().
		Width(leftWidth).
		Padding(1, 2).
		Border(lipgloss.RoundedBorder())
	rightStyle := lipgloss.NewStyle().
		Width(rightWidth).
		Padding(1, 2).
		Border(lipgloss.RoundedBorder())
	gap := lipgloss.NewStyle().Width(gapWidth).Render("")

	return lipgloss.JoinHorizontal(lipgloss.Top, leftStyle.Render(left), gap, rightStyle.Render(right))
}

func (m AppModel) singlePane(content string) string {
	width := m.contentWidth() - 2
	if width < 0 {
		width = 0
	}
	style := lipgloss.NewStyle().Width(width).Padding(1, 2).Border(lipgloss.RoundedBorder())
	return style.Render(content)
}

func (m AppModel) footer(help string) string {
	if help == "" {
		help = "enter: continue  b: back  q: quit"
	}
	help = help + "  |  size: " + strconv.Itoa(m.width) + "x" + strconv.Itoa(m.height)
	footerStyle := lipgloss.NewStyle().Width(m.contentWidth()).Padding(0, 2)
	return "\n" + footerStyle.Render(strings.TrimSpace(help))
}

func (m AppModel) frame(content string) string {
	frameWidth := m.frameWidth()
	frameHeight := m.frameHeight()
	innerWidth := frameWidth - 2
	innerHeight := frameHeight - 2

	if innerWidth < 0 {
		innerWidth = 0
	}
	if innerHeight < 0 {
		innerHeight = 0
	}

	placed := lipgloss.Place(
		innerWidth,
		innerHeight,
		lipgloss.Left,
		lipgloss.Top,
		content,
	)

	frameStyle := lipgloss.NewStyle().
		Width(frameWidth).
		Height(frameHeight).
		Border(lipgloss.NormalBorder())

	rendered := frameStyle.Render(placed)
	return m.injectFrameTitle(rendered, "A7")
}

func (m AppModel) injectFrameTitle(frame, title string) string {
	lines := strings.Split(frame, "\n")
	if len(lines) == 0 {
		return frame
	}

	border := lipgloss.NormalBorder()
	targetWidth := lipgloss.Width(lines[0])
	if targetWidth < 2 {
		return frame
	}
	innerWidth := targetWidth - 2

	paddedTitle := " " + title + " "
	titleWidth := lipgloss.Width(paddedTitle)
	if titleWidth > innerWidth {
		paddedTitle = paddedTitle[:innerWidth]
		titleWidth = lipgloss.Width(paddedTitle)
	}

	fillWidth := innerWidth - titleWidth
	if fillWidth < 0 {
		fillWidth = 0
	}

	topLine := border.TopLeft + paddedTitle + strings.Repeat(border.Top, fillWidth) + border.TopRight
	lines[0] = topLine
	return strings.Join(lines, "\n")
}
