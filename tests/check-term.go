package main

import (
    "bufio"
    "fmt"
    "nadhi/do-t/terminal"
    "os"
)

func Check_term() {
    terminal.ClearScreenFull()
    fmt.Println("Adding text to terminal box...\n")

    tb := terminal.NewTerminalBox("medium")
    tb.SetPosition(0, 3)

    inputBox := terminal.NewInputBox("Type something...")
    inputBox.SetPosition(0, 12)
    inputBox.SetWidth(50)

    reader := bufio.NewReader(os.Stdin)

    var input string

    for {
        inputBox.Print(input)

        char, _, err := reader.ReadRune()
        if err != nil {
            break
        }

        if char == '\n' {
            if input != "" {
                tb.AddText(input)
                input = ""
            }
        } else if char == 127 {
            if len(input) > 0 {
                input = input[:len(input)-1]
            }
        } else if char >= 32 {
            input += string(char)
        }
    }
}