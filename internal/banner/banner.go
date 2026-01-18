// Package banner provides ASCII art banner and colorful text utilities for the vulnerable target application.
package banner

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"strings"
	"time"

	"golang.org/x/term"
)

const (
	// AppName is the application name.
	AppName = "vt"
	// AppVersion is the application version.
	AppVersion = "v0.0.1"
)

// ANSI color codes
const (
	yellow  = "\033[38;2;255;216;43m" // HHS Yellow #ffd82b
	cyan    = "\033[38;2;0;255;255m"  // Cyan for text
	magenta = "\033[38;2;255;0;255m"  // Magenta for accents
	gray    = "\033[90m"              // Gray for dim text
	reset   = "\033[0m"
)

// Quote represents a motivational quote with its author.
type Quote struct {
	Text   string
	Author string
}

var quotesList = []Quote{
	{Text: "Pirêze Hayat, Doxrî Yașanmaz.", Author: "Pișo Meheme"},
	{Text: "Talk is cheap. Show me the code.", Author: "Linus Torvalds"},
	{Text: "Given enough eyeballs, all bugs are shallow.", Author: "Eric S. Raymond"},
	{Text: "The quieter you become, the more you are able to hear.", Author: "Anonymous"},
	{Text: "Hack the planet!", Author: "Hackers (1995)"},
	{Text: "Code is poetry.", Author: "WP Community"},
	{Text: "Think like a hacker, act like an engineer.", Author: "Security Community"},
	{Text: "Open source is power.", Author: "Open Source Advocates"},
	{Text: "Information wants to be free.", Author: "Stewart Brand"},
}

// originalLogo is the full HHS ASCII art logo
var originalLogo = `                                                   ::::::::::::::
                             :::::::::::::::::::::::::::::::::::::::
                         ::::::::::::::::::::::::::::::-===-::::::::::
                      :::::::::::::::::::::::::::::+%@@@%##%@@%-::::::::
                    ::::::::-=*#@@@@@@@@@%#####%@@#*=:::::::::+%@+:::::::
                   :::::::*@@#-:::::::::::::::::::::::::::::::::-%@=::::::::
                  ::::::+@#-::::::::::::::::::::::::::::::::::::::+@%=::::::::::
                 :::::-%@=::::::::::::::::::::::::::::::::::::::::::+@%+:::::::::::
                ::::::*@-:::::::::::::::::::::::::::::::::::::::::::::-#@%+-::::::::::
                :::::-@*:::::::::::::::::::::::::::::::::::::::::::::::::-*@@#=:::::::::
                :::::=@=:::::::::::::::::::::::::::::::::::::::::::::::::::::+%@#=::::::::
               ::::::+@=:::::::::::::=*+:::::::::::::::::::+*=::::::::::::::::::+%@#-:::::::
            ::::::::=@#::::::::::::::@@@%-::::::::::::::::+@@@*::::::::::::::::::::#@#:::::::
          ::::::::=#@+:::::::::::::-#@@@@=:::::::::::::::-%@@@#:::::::::::::::::::::-%@-::::::
         :::::::+@@+::::::::::::::+@@@@@%-:::::::::::::-%@@@@@+:::::::::::::::::::::::#@-:::::
  ::::::::::::+@%-::::::::::::::::%@@@@@%-:::::::::::::=@@@@@@*:::::::::::::::::::::::-@#::::::::
 ::::::::::::+@+::::::::::::::::::=@@@@@@@+:::::::::::::*@@@@@@%-::::::::::::::::::::::#@-:::::::::
:::::::::::::--:::::::::::::::::::::--------::::::::::::::-------::::::::::::::::::::::--:::::::::::
:::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
::::::::::::  HAPPY HACKING  :::::::::::  ** SPACE **  ::::::::::::::::::::::::::::::::::::::::::::
:::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
:::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
::::::::::::::::@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@:::::::::::::::::
:::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
 ::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
  :::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
   :::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
     :::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
      ::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
         ::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
           :::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
              :::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
                  :::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
                       :::::::::::::::::::::::::::::::::::::::::::::::::::::::
                              :::::::::::::::::::::::::::::::::::::::::
                                       :::::::::::::::::::::::::                                `

// HHS mascot frames - derived from originalLogo
var mascotFrames = [][]string{
	// Frame 0 - eyes open (original)
	{
		"                                                   ::::::::::::::                                   ",
		"                             :::::::::::::::::::::::::::::::::::::::                                ",
		"                         ::::::::::::::::::::::::::::::-===-::::::::::                              ",
		"                      :::::::::::::::::::::::::::::+%@@@%##%@@%-::::::::                            ",
		"                    ::::::::-=*#@@@@@@@@@%#####%@@#*=:::::::::+%@+:::::::                           ",
		"                   :::::::*@@#-:::::::::::::::::::::::::::::::::-%@=::::::::                        ",
		"                  ::::::+@#-::::::::::::::::::::::::::::::::::::::+@%=::::::::::                    ",
		"                 :::::-%@=::::::::::::::::::::::::::::::::::::::::::+@%+:::::::::::                 ",
		"                ::::::*@-:::::::::::::::::::::::::::::::::::::::::::::-#@%+-::::::::::              ",
		"                :::::-@*:::::::::::::::::::::::::::::::::::::::::::::::::-*@@#=:::::::::            ",
		"                :::::=@=:::::::::::::::::::::::::::::::::::::::::::::::::::::+%@#=::::::::          ",
		"               ::::::+@=:::::::::::::=*+:::::::::::::::::::+*=::::::::::::::::::+%@#-:::::::        ",
		"            ::::::::=@#::::::::::::::@@@%-::::::::::::::::+@@@*::::::::::::::::::::#@#:::::::       ",
		"          ::::::::=#@+:::::::::::::-#@@@@=:::::::::::::::-%@@@#:::::::::::::::::::::-%@-::::::      ",
		"         :::::::+@@+::::::::::::::+@@@@@%-:::::::::::::-%@@@@@+:::::::::::::::::::::::#@-:::::      ",
		"  ::::::::::::+@%-::::::::::::::::%@@@@@%-:::::::::::::=@@@@@@*:::::::::::::::::::::::-@#::::::::   ",
		" ::::::::::::+@+::::::::::::::::::=@@@@@@@+:::::::::::::*@@@@@@%-::::::::::::::::::::::#@-::::::::: ",
		":::::::::::::--:::::::::::::::::::::--------::::::::::::::-------::::::::::::::::::::::--:::::::::::"},
	// Frame 1 - eyes half closed
	{
		"                                                   ::::::::::::::                                   ",
		"                             :::::::::::::::::::::::::::::::::::::::                                ",
		"                         ::::::::::::::::::::::::::::::-===-::::::::::                              ",
		"                      :::::::::::::::::::::::::::::+%@@@%##%@@%-::::::::                            ",
		"                    ::::::::-=*#@@@@@@@@@%#####%@@#*=:::::::::+%@+:::::::                           ",
		"                   :::::::*@@#-:::::::::::::::::::::::::::::::::-%@=::::::::                        ",
		"                  ::::::+@#-::::::::::::::::::::::::::::::::::::::+@%=::::::::::                    ",
		"                 :::::-%@=::::::::::::::::::::::::::::::::::::::::::+@%+:::::::::::                 ",
		"                ::::::*@-:::::::::::::::::::::::::::::::::::::::::::::-#@%+-::::::::::              ",
		"                :::::-@*:::::::::::::::::::::::::::::::::::::::::::::::::-*@@#=:::::::::            ",
		"                :::::=@=:::::::::::::::::::::::::::::::::::::::::::::::::::::+%@#=::::::::          ",
		"               ::::::+@=:::::::::::::----::::::::::::::::::----:::::::::::::::::+%@#-:::::::        ",
		"            ::::::::=@#::::::::::::::@@@%-::::::::::::::::+@@@*::::::::::::::::::::#@#:::::::       ",
		"          ::::::::=#@+:::::::::::::-#@@@@=:::::::::::::::-%@@@#:::::::::::::::::::::-%@-::::::      ",
		"         :::::::+@@+::::::::::::::+@@@@@%-:::::::::::::-%@@@@@+:::::::::::::::::::::::#@-:::::      ",
		"  ::::::::::::+@%-::::::::::::::::%@@@@@%-:::::::::::::=@@@@@@*:::::::::::::::::::::::-@#::::::::   ",
		" ::::::::::::+@+::::::::::::::::::=@@@@@@@+:::::::::::::*@@@@@@%-::::::::::::::::::::::#@-::::::::: ",
		":::::::::::::--:::::::::::::::::::::--------::::::::::::::-------::::::::::::::::::::::--:::::::::::"},
	// Frame 2 - eyes closed
	{
		"                                                   ::::::::::::::                                   ",
		"                             :::::::::::::::::::::::::::::::::::::::                                ",
		"                         ::::::::::::::::::::::::::::::-===-::::::::::                              ",
		"                      :::::::::::::::::::::::::::::+%@@@%##%@@%-::::::::                            ",
		"                    ::::::::-=*#@@@@@@@@@%#####%@@#*=:::::::::+%@+:::::::                           ",
		"                   :::::::*@@#-:::::::::::::::::::::::::::::::::-%@=::::::::                        ",
		"                  ::::::+@#-::::::::::::::::::::::::::::::::::::::+@%=::::::::::                    ",
		"                 :::::-%@=::::::::::::::::::::::::::::::::::::::::::+@%+:::::::::::                 ",
		"                ::::::*@-:::::::::::::::::::::::::::::::::::::::::::::-#@%+-::::::::::              ",
		"                :::::-@*:::::::::::::::::::::::::::::::::::::::::::::::::-*@@#=:::::::::            ",
		"                :::::=@=:::::::::::::::::::::::::::::::::::::::::::::::::::::+%@#=::::::::          ",
		"               ::::::+@=:::::::::::::----::::::::::::::::::----:::::::::::::::::+%@#-:::::::        ",
		"            ::::::::=@#::::::::::::::----::::::::::::::::::----::::::::::::::::::::#@#:::::::       ",
		"          ::::::::=#@+:::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::-%@-::::::     ",
		"         :::::::+@@+:::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::#@-:::::      ",
		"  ::::::::::::+@%-::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::-@#::::::::   ",
		" ::::::::::::+@+:::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::#@-::::::::: ",
		":::::::::::::--:::::::::::::::::::::--------::::::::::::::-------::::::::::::::::::::::--:::::::::::"},
}

func randomQuote() string {
	if len(quotesList) == 0 {
		return ""
	}
	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(quotesList))))
	if err != nil {
		return ""
	}
	q := quotesList[n.Int64()]
	return fmt.Sprintf("\033[3m%s\033[0m — %s", q.Text, q.Author)
}

// isTerminal checks if stdout is a terminal.
func isTerminal() bool {
	return term.IsTerminal(int(os.Stdout.Fd()))
}

// flush forces stdout to display buffered content immediately.
func flush() {
	if err := os.Stdout.Sync(); err != nil {
		return
	}
}

// moveCursor moves cursor to specific position
func moveCursor(row, col int) {
	fmt.Printf("\033[%d;%dH", row, col)
}

// hideCursor hides the terminal cursor
func hideCursor() {
	fmt.Print("\033[?25l")
}

// showCursor shows the terminal cursor
func showCursor() {
	fmt.Print("\033[?25h")
}

// printColored prints text with specified color
func printColored(text, color string) {
	fmt.Print(color + text + reset)
}

// Banner returns the ASCII art banner with application information (non-animated).
func Banner() string {
	var b strings.Builder

	// Color the original logo - yellow for @#%*+= characters
	for _, ch := range originalLogo {
		switch ch {
		case '@', '#', '%', '*', '+', '=', '-':
			b.WriteString(yellow + string(ch) + reset)
		case ':':
			b.WriteString(gray + string(ch) + reset)
		default:
			b.WriteString(string(ch))
		}
	}

	b.WriteString("\n")
	b.WriteString(cyan + "  Vulnerable Target" + reset)
	b.WriteString(gray + " " + AppVersion + " - Next-gen vuln testing platform" + reset + "\n\n")
	b.WriteString("  " + randomQuote() + "\n")
	b.WriteString(gray + strings.Repeat("─", 100) + reset + "\n")

	return b.String()
}

// PrintAnimated displays the banner with frame-by-frame animation like GitHub Copilot CLI.
func PrintAnimated() {
	if !isTerminal() {
		fmt.Print(Banner())
		return
	}

	hideCursor()
	defer showCursor()

	quote := randomQuote()

	startRow := 1
	startCol := 1

	// Phase 1: Draw border corners
	corners := []struct{ row, col int }{
		{startRow, startCol},
		{startRow, 105},
		{startRow + len(mascotFrames[0]) + 3, startCol},
		{startRow + len(mascotFrames[0]) + 3, 105},
	}
	cornerChars := []string{"┌", "┐", "└", "┘"}

	for i, c := range corners {
		moveCursor(c.row, c.col)
		printColored(cornerChars[i], magenta)
		flush()
		time.Sleep(30 * time.Millisecond)
	}

	// Phase 2: Sparkles around the corners
	sparkles := []struct{ row, col int }{
		{startRow + 1, 4},
		{startRow + 2, 102},
		{startRow + len(mascotFrames[0]), 5},
		{startRow + len(mascotFrames[0]) - 1, 101},
	}
	for _, s := range sparkles {
		moveCursor(s.row, s.col)
		printColored("✦", yellow)
		flush()
		time.Sleep(30 * time.Millisecond)
	}

	// Phase 3: Draw mascot line by line with colors
	for lineIdx, line := range mascotFrames[0] {
		moveCursor(startRow+1+lineIdx, startCol+2)
		for _, ch := range line {
			switch ch {
			case '@', '#', '%', '*', '+', '=', '-':
				printColored(string(ch), yellow)
			case ':':
				printColored(string(ch), gray)
			default:
				fmt.Print(string(ch))
			}
			flush()
			time.Sleep(300 * time.Microsecond)
		}
	}

	// Phase 4: Eye blink animation (blink twice)
	time.Sleep(200 * time.Millisecond)
	blinkSequence := []int{0, 1, 2, 1, 0, 0, 0, 1, 2, 1, 0}
	for _, frameIdx := range blinkSequence {
		for lineIdx, line := range mascotFrames[frameIdx] {
			moveCursor(startRow+1+lineIdx, startCol+2)
			for _, ch := range line {
				switch ch {
				case '@', '#', '%', '*', '+', '=', '-':
					printColored(string(ch), yellow)
				case ':':
					printColored(string(ch), gray)
				default:
					fmt.Print(string(ch))
				}
			}
		}
		flush()
		time.Sleep(80 * time.Millisecond)
	}

	// Phase 5: Info line below logo
	time.Sleep(100 * time.Millisecond)
	infoRow := startRow + len(mascotFrames[0]) + 1
	moveCursor(infoRow, 35)

	title := "Vulnerable Target"
	for _, ch := range title {
		printColored(string(ch), cyan)
		flush()
		time.Sleep(20 * time.Millisecond)
	}
	printColored(" "+AppVersion+" - Next-gen vuln testing platform", gray)
	flush()

	// Phase 6: Quote with typewriter
	time.Sleep(100 * time.Millisecond)
	moveCursor(infoRow+1, 30)
	for _, ch := range quote {
		fmt.Print(string(ch))
		flush()
		time.Sleep(6 * time.Millisecond)
	}

	// Phase 7: Bottom border line
	moveCursor(infoRow+2, startCol)
	for range 105 {
		printColored("─", gray)
		flush()
		time.Sleep(1 * time.Millisecond)
	}

	moveCursor(infoRow+4, 1)
	fmt.Println()
}

// Print displays the banner to stdout with animation if in a terminal.
func Print() {
	PrintAnimated()
}
