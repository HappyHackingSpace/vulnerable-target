package utils

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const (
	AppName    = "vt"
	AppVersion = "v1.0.0"
)

type Quote struct {
	Text   string
	Author string
}

var quotesList = []Quote{
	{Text: "\033[1;3mTalk is cheap. Show me the code.\033[0m", Author: "Linus Torvalds"},
	{Text: "\033[1;3mGiven enough eyeballs, all bugs are shallow.\033[0m", Author: "Eric S. Raymond"},
	{Text: "\033[1;3mThe quieter you become, the more you are able to hear.\033[0m", Author: "Anonymous"},
	{Text: "\033[1;3mHack the planet!\033[0m", Author: "Hackers (1995)"},
	{Text: "\033[1;3mCode is poetry.\033[0m", Author: "WP Community"},
	{Text: "\033[1;3mThink like a hacker, act like an engineer.\033[0m", Author: "Security Community"},
	{Text: "\033[1;3mOpen source is power.\033[0m", Author: "Open Source Advocates"},
	{Text: "\033[1;3mInformation wants to be free.\033[0m", Author: "Stewart Brand"},
	//{Text: "\033[1;3mPreze Hayat Dogru Yasanmaz.\033[0m", Author: "Piso Meheme"}, :>
}

var rainbowColors = []string{
	"\033[31m", // Red
	"\033[33m", // Yellow
	"\033[32m", // Green
	"\033[36m", // Cyan
	"\033[34m", // Blue
	"\033[35m", // Magenta
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RainbowText(text string) string {
	runes := []rune(text)
	var b strings.Builder
	for i, r := range runes {
		color := rainbowColors[i%len(rainbowColors)]
		fmt.Fprintf(&b, "%s%c", color, r)
	}
	b.WriteString("\033[0m") // reset
	return b.String()
}

func randomQuote() string {
	if len(quotesList) == 0 {
		return ""
	}
	q := quotesList[rand.Intn(len(quotesList))]
	return fmt.Sprintf("%s — %s", q.Text, q.Author)
}

func Banner() string {
	title := RainbowText("Next-generation vuln-focused testing platform")
	quote := randomQuote()
	return fmt.Sprintf(`
 HHS     HHS HHSHHSHHSHHS
 HHS     HHS     HHS      %-40s
 HHS     HHS     HHS     
 HHSx   xHHS     HHS     
  xHHS xHHS      HHS      %-40s
   HHSHHS        HHS    
    HHHH         HHS    
     HHS         HHS       %s
%s
`, title, quote, AppVersion, strings.Repeat("-", 62))
}

func PrintBanner(_ string) {
	fmt.Println(Banner())
}
