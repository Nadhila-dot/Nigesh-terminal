package main

import (
	"fmt"
	"nadhi/do-t/server"
	"nadhi/do-t/terminal"
	"os"
	"os/signal"
	"time"
)

func main() {
	//ui := NewConversationUI()
	//ui.Start()
	//Dynamic_demo()
	//Panel_test()
	text := terminal.CreateWrappedText("Your Message")
	text.Set(`This is a sample text that will be wrapped according to the terminal width. 
	It should properly handle line breaks and ensure that the text fits within the specified width.`)
	
	
	text3 := terminal.CreateWrappedText("Notes")
	text3.Set(`Hello there!`)
	loader := terminal.NewLoader()

	// loaders
	loader.Set("Loading using Model Gemini 2.5-flash")
	
	
	time.Sleep(2 * time.Second)
	loader.Set("Finished Generation")
	time.Sleep(1 * time.Second)
	loader.Stop()
	text2 := terminal.CreateWrappedText("Generated code")
	text2.Set(`for i in range(5):
    pint(f"This is iteration number: {i}")

fruits = ["apple", "banana", "cherry"]
print("\nIterating through a list:")
for fruit in fruits:
    print(f"I like {fruit}")`)
	
	time.Sleep(1 * time.Second)
	fmt.Println("\n Uh oh.. Let me try something new.\n")
	loader2 := terminal.NewLoader()
	facts := []string{
		"Fun Fact: Honey never spoils.",
		"Fun Fact: Bananas are berries, but strawberries aren't.",
		"Fun Fact: Octopuses have three hearts.",
		"Fun Fact: A group of flamingos is called a 'flamboyance'.",
	}
	for _, fact := range facts {
		loader2.Set(fact)
		time.Sleep(1 * time.Second)
	}
	loader2.Stop()
	text4 := terminal.CreateWrappedText("Generated code")
	text4.Set(`numbers = [1, 2, 3, 4, 5]
for num in numbers:
	print(f"Number: {num}")`)

	// server tests

	if server.HasWorkspaceFiles() {
    files, _ := server.ListWorkspaceFiles()
	file_text := terminal.CreateWrappedText("Workspace Files")
	file_text.Set(fmt.Sprintf("\n%s", files))
}
if server.HasPublicFiles() {
    files, _ := server.ListPublicFiles()
    file_text2 := terminal.CreateWrappedText("Public Files")
	file_text2.Set(fmt.Sprintf("\n%s", files))
}

 // start file server
go server.StartFileServer()
servertext := terminal.CreateWrappedText("File Server")
servertext.Set("File server running at http://localhost:8087 \n Serving files from ./workspace and ./public directories. \n\n Run EXPORT NIGESH_PORT=your_port to change the port. \n Press Ctrl+C to stop.")

c := make(chan struct{})
go func() {
	// Listen for interrupt signal (Ctrl+C)
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig
	//fmt.Println("\nShutting down file server...")
	close(c)
	os.Exit(0)
}()

<-c
}