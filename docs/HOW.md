# This is the document on how this terminal AI must work. 

# Working Directory
To keep things clean, the terminal AI is given a directory / workspace to work in.
The directory is located at .nigesh/workspace/, this gives nigesh a full workspace to write code, debug, do actions and what not. 

# Using the working directory

Nigesh must use the working directory for basically, behind the sences work. Since the terminal ai is given a special directory, imagine the user requests to get pictures of dogs and download them to the current directory.

Nigesh must break it down, current directory is NOT .nigesh/workspace/, so it should run commands to mv from nigesh's working directory to the current directory the files, this is just an example. 

# Tool calling. 
Since terminal ai is a terminal ai we must give it the ability call tools and get data.

How a tool is called, using simple XML. Inside that is the tool type.

<Tool>Search()<Tool/>

The tools nigesh should have, 

The ability to search
The ability to run commands
The ability to get current system. 
The ability to get current directory.
The ability to get where nigesh's workspace is
The ability to save code and stuff to Nigesh's workspace

We prefer nigesh use the the ability to run commands without having TOO many speecial tool calls to use nigesh's workspace


# Reasoning chain

Nigesh is given an input like this. 

<env>
<current enviroment data here>
<env/>

<The memory compressed>

<System prompt>

<Goal / what the user said nigesh to do>

<additonal context on how to use tools n shit>

nigesh must then basically think, on what it should do and basically reprompt itself again and again until it has decided the End goal has been reached. 

if a -pro flag is passed, use gemini2.5-pro

Graphical looks in the terminal 

without
⚡ running: rm -f hello.txt config.json
┌─ Command Output ─────────────────────────────────────────┐
└──────────────────────────────────────────────────────────┘

do

Executing: rm -f hello.txt config.json
┌─ Command Output ─────────────────────────────────────────┐
| No output was given from this command                    |
└──────────────────────────────────────────────────────────┘

and also No more

┌─ Command Output ─────────────────────────────────────────┐
│ Replaced all occurrences of 'first' with 'initial' in 'hello.txt'.
└──────────────────────────────────────────────────────────┘

properly fix to

┌─ Command Output ───────────────────────────────────────────────────┐
│ Replaced all occurrences of 'first' with 'initial' in 'hello.txt'. |
└────────────────────────────────────────────────────────────────────┘

without

🔄 analyzing results...
💭 thinking (step 7/8)...
⚙️  nigesh is working...
🔍 detected tool calls, executing...

do

🔄 Analyzing..
┌─ Thinking (7/8).. ───────────────────────────────────────┐
| <Thought chain output here or something?>                |
└──────────────────────────────────────────────────────────┘
⚙️  Working...
┌─ Tool call detected ───────────────────────────────────────────────┐
│ Nigesh wanted Search with inputs Thomas+the+tank+engine            |
└────────────────────────────────────────────────────────────────────┘
# Change the tool call message based on tool called. 
show the working part when u parse a tool call in the streama nd show tool calle  detected when u fully parse.

When creating the prompt please tell nigesh to finish it's goal, make it rember to install dependencies make it know that it's context is saved in a json file in .nigesh/context.json

