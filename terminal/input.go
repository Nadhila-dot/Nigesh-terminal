package terminal

import (
    "fmt"
)

type InputBox struct {
    placeholder string
    x           int
    y           int
    width       int
}

func NewInputBox(placeholder string) *InputBox {
    return &InputBox{
        placeholder: placeholder,
        x:           0,
        y:           0,
        width:       50,
    }
}

func (ib *InputBox) SetPosition(x, y int) {
    ib.x = x
    ib.y = y
}

func (ib *InputBox) SetWidth(width int) {
    ib.width = width
}

func (ib *InputBox) Print(input string) {
    moveCursor(ib.x, ib.y)

    fmt.Print("╭─Input ")
    for i := 0; i < ib.width-7; i++ {
        fmt.Print("─")
    }
    fmt.Println("╮")

    displayText := input
    if input == "" {
        displayText = ib.placeholder
    }

    if len(displayText) > ib.width-4 {
        displayText = displayText[:ib.width-4]
    }

    fmt.Printf("│ %-*s │\n", ib.width-4, displayText)

    fmt.Print("╰")
    for i := 0; i < ib.width-2; i++ {
        fmt.Print("─")
    }
    fmt.Println("╯")
}