package terminal

import (
    "fmt"
   _ "strings"
)

type TerminalBox struct {
    buffer []string
    x      int
    y      int
}

func NewTerminalBox(size string) *TerminalBox {
    return &TerminalBox{
        buffer: []string{},
        x:      0,
        y:      0,
    }
}

func (tb *TerminalBox) SetPosition(x, y int) {
    tb.x = x
    tb.y = y
}

func (tb *TerminalBox) AddText(text string) {
    tb.buffer = append(tb.buffer, text)
    tb.Print()
}

func (tb *TerminalBox) RemoveText() {
    if len(tb.buffer) > 0 {
        tb.buffer = tb.buffer[:len(tb.buffer)-1]
        tb.Print()
    }
}

func (tb *TerminalBox) Print() {
    maxLen := 0
    for _, line := range tb.buffer {
        if len(line) > maxLen {
            maxLen = len(line)
        }
    }

    moveCursor(tb.x, tb.y)

    fmt.Print("╔")
    for i := 0; i < maxLen+2; i++ {
        fmt.Print("═")
    }
    fmt.Println("╗")

    for _, line := range tb.buffer {
        fmt.Printf("║ %-*s ║\n", maxLen, line)
    }

    fmt.Print("╚")
    for i := 0; i < maxLen+2; i++ {
        fmt.Print("═")
    }
    fmt.Println("╝")
}

func moveCursor(x, y int) {
    fmt.Printf("\033[%d;%dH", y, x)
}

func clearScreen() {
    fmt.Print("\033[2J\033[H")
}

func ClearScreenFull() {
    clearScreen()
}