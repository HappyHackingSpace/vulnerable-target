# VHS Demo Script for Vulnerable Target (vt)
# Generate with: vhs demo.tape
# Requires: https://github.com/charmbracelet/vhs

Output demo.gif

Set Shell "bash"
Set FontSize 16
Set Width 1200
Set Height 600
Set Theme "Dracula"
Set Padding 20

# Title
Type "# Vulnerable Target (vt) - Demo"
Enter
Sleep 500ms

# Clear and start demo
Type "clear"
Enter
Sleep 500ms

# Show available templates
Type "vt template --list"
Sleep 500ms
Enter
Sleep 3s

# Filter by tag
Type "vt template --list --filter dvwa"
Sleep 500ms
Enter
Sleep 2s

# Start a vulnerable environment
Type "vt start --id vt-dvwa"
Sleep 500ms
Enter
Sleep 4s

# Show running environments
Type "vt ps"
Sleep 500ms
Enter
Sleep 2s

# Stop the environment
Type "vt stop --id vt-dvwa"
Sleep 500ms
Enter
Sleep 3s

# Final message
Type "# Happy Hacking!"
Enter
Sleep 2s
