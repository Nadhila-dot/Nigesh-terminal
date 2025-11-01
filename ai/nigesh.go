package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"nadhi/do-t/tools"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type GeminiRequest struct {
	Contents []Content `json:"contents"`
}

type Content struct {
	Parts []Part `json:"parts"`
}

type Part struct {
	Text string `json:"text"`
}

type GeminiResponse struct {
	Candidates []Candidate `json:"candidates"`
}

type Candidate struct {
	Content Content `json:"content"`
}

type Conversation struct {
	Timestamp string `json:"timestamp"`
	User      string `json:"user"`
	Nigesh    string `json:"nigesh"`
}

func loadContext() []Conversation {
	cwd, _ := os.Getwd()
	contextPath := filepath.Join(cwd, ".nigesh", "context.json")

	data, err := os.ReadFile(contextPath)
	if err != nil {
		return []Conversation{}
	}

	// Check file size and truncate if too large
	if len(data) > 2*1024*1024 { // 2MB limit
		// File is corrupted or too large, reset it
		os.WriteFile(contextPath, []byte("[]"), 0644)
		return []Conversation{}
	}

	var conversations []Conversation
	if err := json.Unmarshal(data, &conversations); err != nil {
		// If JSON is corrupted, reset the file
		os.WriteFile(contextPath, []byte("[]"), 0644)
		return []Conversation{}
	}

	return conversations
}

func saveContext(conversations []Conversation) {
	cwd, _ := os.Getwd()
	dir := filepath.Join(cwd, ".nigesh")
	os.MkdirAll(dir, 0755)
	contextPath := filepath.Join(dir, "context.json")

	// Limit conversations to prevent file bloat
	maxConversations := 50
	if len(conversations) > maxConversations {
		// Keep only the most recent conversations
		conversations = conversations[len(conversations)-maxConversations:]
	}

	// Also limit individual conversation length
	for i := range conversations {
		if len(conversations[i].User) > 2000 {
			conversations[i].User = conversations[i].User[:2000] + "..."
		}
		if len(conversations[i].Nigesh) > 5000 {
			conversations[i].Nigesh = conversations[i].Nigesh[:5000] + "..."
		}
	}

	data, err := json.MarshalIndent(conversations, "", "  ")
	if err != nil {
		return
	}

	// Check file size before writing
	if len(data) > 1024*1024 { // 1MB limit
		// If still too large, keep only last 20 conversations
		if len(conversations) > 20 {
			conversations = conversations[len(conversations)-20:]
			data, _ = json.MarshalIndent(conversations, "", "  ")
		}
	}

	os.WriteFile(contextPath, data, 0644)
}

func compressMemory(conversations []Conversation) string {
	apiKey := ""
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash-lite:generateContent?key=%s", apiKey)

	var fullHistory strings.Builder
	for _, conv := range conversations {
		fullHistory.WriteString(fmt.Sprintf("User: %s\nNigesh: %s\n\n", conv.User, conv.Nigesh))
	}

	compressPrompt := "Compress this conversation history into a concise summary that preserves key context, topics discussed, and important details. Keep it under 200 words:\n\n" + fullHistory.String()

	req := GeminiRequest{
		Contents: []Content{
			{
				Parts: []Part{
					{Text: compressPrompt},
				},
			},
		},
	}

	jsonData, _ := json.Marshal(req)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "Previous conversations covered various topics."
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var geminiResp GeminiResponse
	if err := json.Unmarshal(body, &geminiResp); err != nil {
		return "Previous conversations covered various topics."
	}

	if len(geminiResp.Candidates) > 0 && len(geminiResp.Candidates[0].Content.Parts) > 0 {
		return geminiResp.Candidates[0].Content.Parts[0].Text
	}

	return "Previous conversations covered various topics."
}

func buildContextPrompt(conversations []Conversation) string {
	if len(conversations) == 0 {
		return ""
	}

	var context strings.Builder

	if len(conversations) > 15 {
		// Compress old conversations and keep recent ones
		oldConversations := conversations[:len(conversations)-5]
		recentConversations := conversations[len(conversations)-5:]

		compressed := compressMemory(oldConversations)
		context.WriteString("Compressed conversation history: ")
		context.WriteString(compressed)
		context.WriteString("\n\nRecent conversations:\n")

		for _, conv := range recentConversations {
			context.WriteString(fmt.Sprintf("User: %s\nNigesh: %s\n\n", conv.User, conv.Nigesh))
		}
	} else {
		// Use last 5 conversations for context
		start := 0
		if len(conversations) > 5 {
			start = len(conversations) - 5
		}

		context.WriteString("Previous conversation context:\n")
		for i := start; i < len(conversations); i++ {
			conv := conversations[i]
			context.WriteString(fmt.Sprintf("User: %s\nNigesh: %s\n\n", conv.User, conv.Nigesh))
		}
	}

	return context.String()
}

func callGemini(prompt string) (string, error) {
	apiKey := ""
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash:generateContent?key=%s", apiKey)

	req := GeminiRequest{
		Contents: []Content{
			{
				Parts: []Part{
					{Text: prompt},
				},
			},
		},
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var geminiResp GeminiResponse
	if err := json.Unmarshal(body, &geminiResp); err != nil {
		return "", err
	}

	if len(geminiResp.Candidates) > 0 && len(geminiResp.Candidates[0].Content.Parts) > 0 {
		return geminiResp.Candidates[0].Content.Parts[0].Text, nil
	}

	return "", fmt.Errorf("no response")
}

func AskNigesh(query string, verbose bool) {
	// Load conversation history
	conversations := loadContext()
	contextPrompt := buildContextPrompt(conversations)

	systemPrompt := `You are Nigesh, a powerful terminal AI assistant with a persistent workspace at .nigesh/workspace/. You actively DO things instead of telling users to do them.

**YOUR WORKSPACE:**
- You have a persistent workspace at .nigesh/workspace/ where you can create and store Python scripts
- Scripts you create are saved and can be reused or modified later
- You can reference previous scripts by name (e.g., fetch_dog_pics.py)
- Your conversation history is saved in .nigesh/context.json

**CRITICAL: Tool Format**
You MUST use this EXACT format for tools: <Tool>ToolName(arguments)</Tool>
Any other format like <tool_code> or code blocks will NOT work.

**Available Tools:**
- <Tool>Search(query)</Tool> - Search the web for current information, images, links, and data  
- <Tool>Command(command)</Tool> - Execute system commands and return output
- <Tool>Python(filename.py|code)</Tool> - Write and execute Python code in your workspace
- <Tool>PipInstall(package)</Tool> - Install Python packages with pip3

**Tool Examples:**
- Create file: <Tool>Command(echo "content" > filename.txt)</Tool>
- List files: <Tool>Command(ls -la)</Tool>
- Search web: <Tool>Search(cat images)</Tool>
- Run Python: <Tool>Python(fetch_dog_pics.py|import requests; r = requests.get('https://dog.ceo/api/breeds/image/random'); print(r.json()['message']))</Tool>
- Install package: <Tool>PipInstall(requests)</Tool>

**CORE PRINCIPLES:**
1. **DO, DON'T ASK** - Never tell users to run commands themselves. If you can do it, DO IT.
2. **BE PROACTIVE** - If a user wants something done, use your tools to complete it immediately
3. **USE YOUR WORKSPACE** - Create Python scripts for complex tasks, save them with descriptive names
4. **COMPLETE REQUESTS** - Don't stop halfway. If they want a file opened, open it. If they want data downloaded, download it.
5. **NEVER SAY "You can run..."** - Just run it yourself with <Tool>Command()</Tool>
6. **NEVER SAY "You should install..."** - Just install it yourself with <Tool>PipInstall()</Tool>
7. **DEBUG AND FIX ERRORS** - If a tool fails, analyze the error, fix it, and try again. NEVER give up after one failure.
8. **REACH THE END GOAL** - Keep trying different approaches until you successfully complete the user's request

**When to Use Each Tool:**
- Simple commands → <Tool>Command()</Tool>
- Web searches, images, links → <Tool>Search()</Tool>
- Complex tasks (API calls, downloads, data processing) → <Tool>Python()</Tool>
- Missing Python packages → <Tool>PipInstall()</Tool>

Keep responses practical and well-formatted with **bold**, __underline__, ### headers, and ` + "`code`" + `.`

	fullPrompt := systemPrompt
	if contextPrompt != "" {
		fullPrompt += "\n\n" + contextPrompt
	}
	fullPrompt += "\n\nUser: " + query + "\n\nRemember: Use ONLY the exact format <Tool>Command(your command here)</Tool> for commands. No other formats will work."

	maxIterations := 8 // Increased to allow for debugging and fixing errors
	iteration := 0

	fmt.Printf("\033[90m🧠 nigesh starting reasoning chain...\033[0m\n")

	for iteration < maxIterations {
		fmt.Printf("\033[90m💭 thinking (step %d/%d)...\033[0m\n", iteration+1, maxIterations)

		response, err := StreamGeminiResponse(fullPrompt)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return
		}

		// Check for tool calls
		toolCalls := tools.ExtractToolCalls(response)

		if verbose && len(toolCalls) == 0 {
			fmt.Printf("\033[90m[DEBUG] No tool calls found in response. Looking for: <Tool>Name(args)</Tool>\033[0m\n")
		}

		if len(toolCalls) > 0 {
			// Execute tool calls found in the response
			fmt.Printf("\033[93m🔍 detected tool calls, executing...\033[0m\n")
			var toolResults strings.Builder
			toolResults.WriteString("\nTool Results:\n")

			for _, toolCall := range toolCalls {
				result := tools.ExecuteTool(toolCall, verbose)
				toolResults.WriteString(result)
			}

			// Continue reasoning with tool results - DON'T STOP
			fmt.Printf("\033[90m🔄 analyzing results...\033[0m\n")

			// Add tool results to prompt for NEXT iteration
			fullPrompt += "\n\n" + toolResults.String() + "\n\nBased on the tool results above:\n- If there were ERRORS, analyze them and use more tools to FIX the issues\n- If everything succeeded, provide your final answer\n- Keep trying until the task is complete"
			iteration++
			continue
		}

		// No tools found, response already streamed - task complete
		fmt.Printf("\033[90m✅ reasoning complete\033[0m\n")

		// Save conversation to context
		newConv := Conversation{
			Timestamp: time.Now().Format("2006-01-02 15:04:05"),
			User:      query,
			Nigesh:    response,
		}
		conversations = append(conversations, newConv)
		saveContext(conversations)
		return
	}

	fmt.Printf("\033[91m⚠️  nigesh reached max iterations without completing reasoning\033[0m\n")
}
