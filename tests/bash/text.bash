#!/bin/bash

# Get terminal width
cols=$(tput cols)
padding=4                 # spaces for borders
max_width=$((cols - padding))

# Read user input
read -p "You: " input

# Wrap text if too long
wrapped=$(echo "$input" | fold -s -w $((max_width - 4)))

# Compute longest line length (for box width)
longest=$(echo "$wrapped" | awk '{ if (length > max) max = length } END { print max }')

# Top border
echo "╭─ You $(printf '─%.0s' $(seq 1 $((longest - 2))))╮"

# Content lines
while IFS= read -r line; do
    printf "│ %-*s   │\n" "$longest" "$line" #Important these spaces
done <<< "$wrapped"

# Bottom border
echo "╰$(printf '─%.0s' $(seq 1 $((longest + 4))))╯"
