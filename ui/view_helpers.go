package ui

import (
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m AppModel) contentWidth() int {
	frameWidth := m.frameWidth()
	paddingX := m.contentPaddingX()
	if frameWidth > paddingX*2 {
		return frameWidth - paddingX*2
	}
	return 80
}

func (m AppModel) frameWidth() int {
	if m.width > 2 {
		return m.width - 2
	}
	return 80
}

func (m AppModel) frameHeight() int {
	if m.height > 2 {
		return m.height - 2
	}
	return 24
}

func (m AppModel) contentPaddingX() int {
	return 0
}

func (m AppModel) contentPaddingY() int {
	return 0
}

func (m AppModel) twoPane(left, right string) string {
	return m.twoPaneWithRatio(left, right, 0.5)
}

func (m AppModel) twoPaneWithRatio(left, right string, leftRatio float64) string {
	return m.twoPaneWithRatioAndTitles("", "", left, right, leftRatio)
}

func (m AppModel) twoPaneWithRatioAndTitles(leftTitle, rightTitle, left, right string, leftRatio float64) string {
	return m.twoPaneWithRatioAndTitlesAndWidth(leftTitle, rightTitle, left, right, leftRatio, m.contentWidth())
}

func (m AppModel) twoPaneWithRatioAndTitlesAndWidth(leftTitle, rightTitle, left, right string, leftRatio float64, totalWidth int) string {
	theme := currentTheme()
	total := totalWidth
	if total <= 0 {
		total = m.contentWidth()
	}
	gapWidth := 2
	paddingX := 2
	borderX := m.cardBorderX()
	extra := paddingX + borderX
	available := total - gapWidth - (extra * 2)
	if available < 0 {
		available = 0
	}
	availableHeight := m.bodyHeight()
	if leftRatio < 0 {
		leftRatio = 0
	}
	if leftRatio > 1 {
		leftRatio = 1
	}
	leftWidth := int(float64(available) * leftRatio)
	rightWidth := available - leftWidth

	leftStyle := lipgloss.NewStyle().
		Width(leftWidth).
		Padding(1, paddingX).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.PaneBorder).
		Foreground(theme.Text)
	rightStyle := lipgloss.NewStyle().
		Width(rightWidth).
		Padding(1, paddingX).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.PaneBorder).
		Foreground(theme.Text)
	if availableHeight > 0 {
		boxHeight := m.paneBoxHeight(availableHeight)
		leftStyle = leftStyle.Height(boxHeight)
		rightStyle = rightStyle.Height(boxHeight)
	}

	leftPane := leftStyle.Render(left)
	rightPane := rightStyle.Render(right)
	if leftTitle != "" {
		leftPane = m.injectBorderTitle(leftPane, leftTitle, lipgloss.RoundedBorder(), theme.PaneBorder)
	}
	if rightTitle != "" {
		rightPane = m.injectBorderTitle(rightPane, rightTitle, lipgloss.RoundedBorder(), theme.PaneBorder)
	}
	gap := lipgloss.NewStyle().Width(gapWidth).Render("")

	return lipgloss.JoinHorizontal(lipgloss.Top, leftPane, gap, rightPane)
}

func (m AppModel) singlePane(content string) string {
	return m.singlePaneWithWidth(content, m.contentWidth())
}

func (m AppModel) singlePaneWithWidth(content string, totalWidth int) string {
	theme := currentTheme()
	paddingX := 2
	if totalWidth <= 0 {
		totalWidth = m.contentWidth()
	}
	width := m.paneContentWidth(totalWidth)
	height := m.paneBoxHeight(m.bodyHeight())
	style := lipgloss.NewStyle().
		Width(width).
		Padding(1, paddingX).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.PaneBorder).
		Foreground(theme.Text)
	if height > 0 {
		style = style.Height(height)
	}
	return style.Render(content)
}

func (m AppModel) cardBorderX() int {
	return lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).GetHorizontalBorderSize()
}

func (m AppModel) cardBorderY() int {
	return lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).GetVerticalBorderSize()
}

func (m AppModel) centerContent(content string) string {
	return lipgloss.Place(m.contentWidth(), lipgloss.Height(content), lipgloss.Center, lipgloss.Top, content)
}

func (m AppModel) paneHeader(title string, paneWidth int) string {
	paddingX := 2
	headerWidth := paneWidth - paddingX*2
	if headerWidth < 0 {
		headerWidth = 0
	}
	return lipgloss.NewStyle().
		Bold(true).
		Padding(0, paddingX).
		Width(headerWidth).
		Align(lipgloss.Left).
		Render(title)
}

func (m AppModel) paneContentWidth(totalWidth int) int {
	paddingX := 2
	borderX := m.cardBorderX()
	width := totalWidth - (paddingX*2 + borderX)
	if width < 0 {
		return 0
	}
	return width
}

func (m AppModel) paneBoxHeight(totalHeight int) int {
	borderY := m.cardBorderY()
	height := totalHeight - borderY
	if height < 0 {
		return 0
	}
	return height
}

func (m AppModel) paneContentHeight(totalHeight int) int {
	paddingY := 1
	borderY := m.cardBorderY()
	height := totalHeight - (paddingY*2 + borderY)
	if height < 0 {
		return 0
	}
	return height
}

func (m AppModel) bodyHeight() int {
	height := m.frameHeight()
	if height <= 0 {
		return 0
	}
	if height-1 < 0 {
		return 0
	}
	return height - 1
}

func (m AppModel) splitPaneContentWidths(leftRatio float64) (int, int) {
	return m.splitPaneContentWidthsForTotal(m.contentWidth(), leftRatio)
}

func (m AppModel) splitPaneContentWidthsForTotal(totalWidth int, leftRatio float64) (int, int) {
	total := totalWidth
	if total <= 0 {
		total = m.contentWidth()
	}
	gapWidth := 2
	paddingX := 2
	borderX := m.cardBorderX()
	extra := (paddingX * 2) + borderX
	available := total - gapWidth - (extra * 2)
	if available < 0 {
		available = 0
	}
	if leftRatio < 0 {
		leftRatio = 0
	}
	if leftRatio > 1 {
		leftRatio = 1
	}
	leftWidth := int(float64(available) * leftRatio)
	rightWidth := available - leftWidth
	return leftWidth, rightWidth
}

func (m AppModel) helpLine(help string, width int) string {
	theme := currentTheme()
	if help != "" {
		help = help + "  |  size: " + strconv.Itoa(m.width) + "x" + strconv.Itoa(m.height)
	}
	helpStyle := lipgloss.NewStyle().Width(width).Padding(0, 2).Foreground(theme.Help)
	return helpStyle.Render(strings.TrimSpace(help))
}

func padToHeight(s string, height int) string {
	if height <= 0 {
		return ""
	}
	lines := strings.Split(s, "\n")
	if len(lines) > height {
		lines = lines[:height]
	}
	if len(lines) < height {
		lines = append(lines, make([]string, height-len(lines))...)
	}
	return strings.Join(lines, "\n")
}

func fitToHeight(s string, height int) string {
	if height <= 0 {
		return ""
	}
	return padToHeight(s, height)
}

func (m AppModel) frame(content, help string) string {
	theme := currentTheme()
	frameWidth := m.frameWidth()
	frameHeight := m.frameHeight()
	paddingX := m.contentPaddingX()
	paddingY := m.contentPaddingY()
	innerWidth := frameWidth - paddingX*2
	innerHeight := frameHeight - paddingY*2

	if innerWidth < 0 {
		innerWidth = 0
	}
	if innerHeight < 0 {
		innerHeight = 0
	}

	content = strings.TrimRight(content, "\n")
	bodyHeight := innerHeight
	if help != "" && innerHeight > 0 {
		bodyHeight = innerHeight - 1
	}

	body := lipgloss.Place(
		innerWidth,
		bodyHeight,
		lipgloss.Left,
		lipgloss.Top,
		content,
	)
	body = padToHeight(body, bodyHeight)

	placed := body
	if help != "" && innerHeight > 0 {
		helpLine := m.helpLine(help, innerWidth)
		if bodyHeight > 0 {
			placed = body + "\n" + helpLine
		} else {
			placed = helpLine
		}
	}

	if paddingX > 0 || paddingY > 0 {
		placed = lipgloss.NewStyle().Padding(paddingY, paddingX).Render(placed)
	}

	frameStyle := lipgloss.NewStyle().
		Width(frameWidth).
		Height(frameHeight).
		Border(lipgloss.NormalBorder()).
		BorderForeground(theme.FrameBorder)

	rendered := frameStyle.Render(placed)
	return m.injectFrameTitle(rendered, "A7", theme.FrameBorder)
}

func (m AppModel) injectFrameTitle(frame, title string, color lipgloss.Color) string {
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
	if color != "" {
		topLine = lipgloss.NewStyle().Foreground(color).Bold(true).Render(topLine)
	}
	lines[0] = topLine
	return strings.Join(lines, "\n")
}

func (m AppModel) injectBorderTitle(frame, title string, border lipgloss.Border, color lipgloss.Color) string {
	lines := strings.Split(frame, "\n")
	if len(lines) == 0 {
		return frame
	}

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
	if color != "" {
		topLine = lipgloss.NewStyle().Foreground(color).Bold(true).Render(topLine)
	}
	lines[0] = topLine
	return strings.Join(lines, "\n")
}

func (m AppModel) titledPane(title, content string) string {
	return m.titledPaneWithWidth(title, content, m.contentWidth())
}

func (m AppModel) titledPaneWithWidth(title, content string, totalWidth int) string {
	theme := currentTheme()
	pane := m.singlePaneWithWidth(content, totalWidth)
	return m.injectBorderTitle(pane, title, lipgloss.RoundedBorder(), theme.PaneBorder)
}

func (m AppModel) titledPaneWithWidthAndHeight(title, content string, totalWidth, totalHeight int) string {
	theme := currentTheme()
	paddingX := 2
	width := m.paneContentWidth(totalWidth)
	height := m.paneContentHeight(totalHeight)
	if width < 0 {
		width = 0
	}
	style := lipgloss.NewStyle().
		Width(width).
		Padding(1, paddingX).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.PaneBorder).
		Foreground(theme.Text)
	if height > 0 {
		style = style.Height(height)
	}
	pane := style.Render(content)
	return m.injectBorderTitle(pane, title, lipgloss.RoundedBorder(), theme.PaneBorder)
}

func (m AppModel) primaryPaneWidth() int {
	width := m.contentWidth() / 2
	if width <= 0 {
		return m.contentWidth()
	}
	return width
}

func (m AppModel) formWidth() int {
	return m.paneContentWidth(m.primaryPaneWidth())
}

func (m AppModel) editorPaneWidth() int {
	minWidth := m.primaryPaneWidth()
	target := m.contentWidth()
	if target < minWidth {
		return minWidth
	}
	return target
}
