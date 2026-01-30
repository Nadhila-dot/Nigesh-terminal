# Nigesh's workspace is a place to privately craft things and give users outputs. 

Think of it as a workbench where Nigesh crafts things and Nigesh can give outputs

There should be a tool called, 

MoveToPublic(<array of files in nigesh's private workspace>)

these will be moved to the place the command was executed into a folder called nigesh-public

When nigesh is run there will always be a file server running to serve files in nigesh-public and .nigesh/workspace


The file server will run on port (8087) or EXPORT PORT_NIGESH=int

/private/<file paths> for private 

/public/<file path> for public

# Nigesh having access to files.

Best if you ask the user if they consent on letting nigesh search through there files.
Nigesh will then go through files using.

GetContents(<file name>)
, it will send with line numbers so nigesh can change line numbers

GetLineContents(<file name>, <line number>)

WriteFile(<file name>,<file contents>) // change the whole file with new contents

ChangeLine(<file name>, <line number>, <new line content>) // Change a specific line




