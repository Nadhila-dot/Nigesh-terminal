package terminal

import (
	"fmt"
)

type Panel struct {
	title   string
	content []string
	x       int
	y       int
	width   int
	height  int
	style   PanelStyle
}

type PanelStyle struct {
	BorderColor string
	TitleColor  string
	TextColor   string
}

var (
	DefaultStyle = PanelStyle{
		BorderColor: "\033[36m",   // Cyan
		TitleColor:  "\033[1;37m", // Bold white
		TextColor:   "\033[0m",    // Reset
	}

	UserStyle = PanelStyle{
		BorderColor: "\033[34m", // Blue
		TitleColor:  "\033[1;34m",
		TextColor:   "\033[0m",
	}

	AssistantStyle = PanelStyle{
		BorderColor: "\033[35m", // Magenta
		TitleColor:  "\033[1;35m",
		TextColor:   "\033[0m",
	}
)

func NewPanel(title string, width, height int) *Panel {
	return &Panel{
		title:   title,
		content: []string{},
		x:       0,
		y:       0,
		width:   width,
		height:  height,
		style:   DefaultStyle,
	}
}

func (p *Panel) SetPosition(x, y int) {
	p.x = x
	p.y = y
}

func (p *Panel) SetStyle(style PanelStyle) {
	p.style = style
}

func (p *Panel) SetContent(lines []string) {
	p.content = lines
}

func (p *Panel) AppendContent(line string) {
	p.content = append(p.content, line)
}

func (p *Panel) Clear() {
	p.content = []string{}
}

func (p *Panel) Render() {
	moveCursor(p.x, p.y)

	// Top border with title
	fmt.Print(p.style.BorderColor + "╔")
	if p.title != "" {
		titleText := " " + p.title + " "
		fmt.Print(p.style.TitleColor + titleText + p.style.BorderColor)
		remaining := p.width - len(titleText) - 2
		for i := 0; i < remaining; i++ {
			fmt.Print("═")
		}
	} else {
		for i := 0; i < p.width-2; i++ {
			fmt.Print("═")
		}
	}
	fmt.Println("╗" + p.style.TextColor)

	// Content
	for i := 0; i < p.height-2; i++ {
		fmt.Print(p.style.BorderColor + "║ " + p.style.TextColor)
		if i < len(p.content) {
			line := p.content[i]
			if len(line) > p.width-4 {
				line = line[:p.width-4]
			}
			fmt.Printf("%-*s", p.width-4, line)
		} else {
			fmt.Printf("%-*s", p.width-4, "")
		}
		fmt.Println(p.style.BorderColor + " ║" + p.style.TextColor)
	}

	// Bottom border
	fmt.Print(p.style.BorderColor + "╚")
	for i := 0; i < p.width-2; i++ {
		fmt.Print("═")
	}
	fmt.Println("╝" + p.style.TextColor)
}

func (p *Panel) RenderAt(x, y int) {
	p.SetPosition(x, y)
	p.Render()
}
