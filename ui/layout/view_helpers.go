package layout

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/never00rei/a7/ui/theme"
)

type Layout struct {
	Width       int
	Height      int
	ContentPadX int
	ContentPadY int
	ActiveTheme theme.Theme
}

func New(width, height int) Layout {
	return Layout{
		Width:       width,
		Height:      height,
		ContentPadX: 0,
		ContentPadY: 0,
		ActiveTheme: theme.CurrentTheme(),
	}
}

func (l Layout) ContentWidth() int {
	frameWidth := l.FrameWidth()
	paddingX := l.ContentPaddingX()
	if frameWidth > paddingX*2 {
		return frameWidth - paddingX*2
	}
	return 80
}

func (l Layout) FrameWidth() int {
	if l.Width > 2 {
		return l.Width - 2
	}
	return 80
}

func (l Layout) FrameHeight() int {
	if l.Height > 2 {
		return l.Height - 2
	}
	return 24
}

func (l Layout) ContentPaddingX() int {
	return l.ContentPadX
}

func (l Layout) ContentPaddingY() int {
	return l.ContentPadY
}

func (l Layout) TwoPane(left, right string) string {
	return l.TwoPaneWithRatio(left, right, 0.5)
}

func (l Layout) TwoPaneWithRatio(left, right string, leftRatio float64) string {
	return l.TwoPaneWithRatioAndTitles("", "", left, right, leftRatio)
}

func (l Layout) TwoPaneWithRatioAndTitles(leftTitle, rightTitle, left, right string, leftRatio float64) string {
	return l.TwoPaneWithRatioAndTitlesAndWidth(leftTitle, rightTitle, left, right, leftRatio, l.ContentWidth())
}

func (l Layout) TwoPaneWithRatioAndTitlesAndWidth(leftTitle, rightTitle, left, right string, leftRatio float64, totalWidth int) string {
	theme := l.ActiveTheme
	total := totalWidth
	if total <= 0 {
		total = l.ContentWidth()
	}
	gapWidth := 2
	paddingX := 2
	borderX := l.CardBorderX()
	extra := paddingX + borderX
	available := total - gapWidth - (extra * 2)
	if available < 0 {
		available = 0
	}
	availableHeight := l.BodyHeight()
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
		boxHeight := l.PaneBoxHeight(availableHeight)
		leftStyle = leftStyle.Height(boxHeight)
		rightStyle = rightStyle.Height(boxHeight)
	}

	leftPane := leftStyle.Render(left)
	rightPane := rightStyle.Render(right)
	if leftTitle != "" {
		leftPane = l.injectBorderTitle(leftPane, leftTitle, lipgloss.RoundedBorder(), theme.PaneBorder)
	}
	if rightTitle != "" {
		rightPane = l.injectBorderTitle(rightPane, rightTitle, lipgloss.RoundedBorder(), theme.PaneBorder)
	}
	gap := lipgloss.NewStyle().Width(gapWidth).Render("")

	return lipgloss.JoinHorizontal(lipgloss.Top, leftPane, gap, rightPane)
}

func (l Layout) SinglePane(content string) string {
	return l.SinglePaneWithWidth(content, l.ContentWidth())
}

func (l Layout) SinglePaneWithWidth(content string, totalWidth int) string {
	theme := l.ActiveTheme
	paddingX := 2
	if totalWidth <= 0 {
		totalWidth = l.ContentWidth()
	}
	width := l.PaneContentWidth(totalWidth)
	height := l.PaneBoxHeight(l.BodyHeight())
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

func (l Layout) CardBorderX() int {
	return lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).GetHorizontalBorderSize()
}

func (l Layout) CardBorderY() int {
	return lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).GetVerticalBorderSize()
}

func (l Layout) CenterContent(content string) string {
	return lipgloss.Place(l.ContentWidth(), lipgloss.Height(content), lipgloss.Center, lipgloss.Top, content)
}

func (l Layout) PaneHeader(title string, paneWidth int) string {
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

func (l Layout) PaneContentWidth(totalWidth int) int {
	paddingX := 2
	borderX := l.CardBorderX()
	width := totalWidth - (paddingX*2 + borderX)
	if width < 0 {
		return 0
	}
	return width
}

func (l Layout) PaneBoxHeight(totalHeight int) int {
	borderY := l.CardBorderY()
	height := totalHeight - borderY
	if height < 0 {
		return 0
	}
	return height
}

func (l Layout) PaneContentHeight(totalHeight int) int {
	paddingY := 1
	borderY := l.CardBorderY()
	height := totalHeight - (paddingY*2 + borderY)
	if height < 0 {
		return 0
	}
	return height
}

func (l Layout) BodyHeight() int {
	height := l.FrameHeight()
	if height <= 0 {
		return 0
	}
	if height-1 < 0 {
		return 0
	}
	return height - 1
}

func (l Layout) SplitPaneContentWidths(leftRatio float64) (int, int) {
	return l.SplitPaneContentWidthsForTotal(l.ContentWidth(), leftRatio)
}

func (l Layout) SplitPaneContentWidthsForTotal(totalWidth int, leftRatio float64) (int, int) {
	total := totalWidth
	if total <= 0 {
		total = l.ContentWidth()
	}
	gapWidth := 2
	paddingX := 2
	borderX := l.CardBorderX()
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

func (l Layout) helpLine(help string, width int) string {
	theme := l.ActiveTheme
	contentWidth := width - 4
	if contentWidth < 0 {
		contentWidth = 0
	}
	help = ansi.Truncate(strings.TrimSpace(help), contentWidth, "")
	helpStyle := lipgloss.NewStyle().Width(width).Padding(0, 2).Foreground(theme.Help)
	return helpStyle.Render(help)
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

func (l Layout) Frame(content, help string) string {
	theme := l.ActiveTheme
	frameWidth := l.FrameWidth()
	frameHeight := l.FrameHeight()
	paddingX := l.ContentPaddingX()
	paddingY := l.ContentPaddingY()
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
		helpLine := l.helpLine(help, innerWidth)
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
	return l.injectFrameTitle(rendered, "A7", theme.FrameBorder)
}

func (l Layout) injectFrameTitle(frame, title string, color lipgloss.Color) string {
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

func (l Layout) injectBorderTitle(frame, title string, border lipgloss.Border, color lipgloss.Color) string {
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

func (l Layout) TitledPane(title, content string) string {
	return l.TitledPaneWithWidth(title, content, l.ContentWidth())
}

func (l Layout) TitledPaneWithWidth(title, content string, totalWidth int) string {
	theme := l.ActiveTheme
	pane := l.SinglePaneWithWidth(content, totalWidth)
	return l.injectBorderTitle(pane, title, lipgloss.RoundedBorder(), theme.PaneBorder)
}

func (l Layout) TitledPaneWithWidthAndHeight(title, content string, totalWidth, totalHeight int) string {
	theme := l.ActiveTheme
	paddingX := 2
	width := l.PaneContentWidth(totalWidth)
	height := l.PaneContentHeight(totalHeight)
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
	return l.injectBorderTitle(pane, title, lipgloss.RoundedBorder(), theme.PaneBorder)
}

func (l Layout) PrimaryPaneWidth() int {
	width := l.ContentWidth() / 2
	if width <= 0 {
		return l.ContentWidth()
	}
	return width
}

func (l Layout) FormWidth() int {
	return l.PaneContentWidth(l.PrimaryPaneWidth())
}

func (l Layout) EditorPaneWidth() int {
	minWidth := l.PrimaryPaneWidth()
	target := l.ContentWidth()
	if target < minWidth {
		return minWidth
	}
	return target
}
