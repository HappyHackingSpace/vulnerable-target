package utils

import "fmt"

// Uygulama bilgileri
const (
	AppName    = "vt"
	AppVersion = "v1.0.0"
)

// Banner returns the application banner with colored text
func Banner() string {
	return `
 HHS     HHS HHSHHSHHSHHS
 HHS     HHS     HHS     
 HHS     HHS     HHS     
 HHSx   xHHS     HHS     ` + " " + RainbowText("- Create vulnerable environment") + `
  xHHS xHHS      HHS     ` + " " + "\033[1;3m- Happy Hacking!\033[0m" + `
   HHSHHS        HHS     
    HHHH         HHS    
     HHS         HHS     			   ` + fmt.Sprintf("\033[1m%s\033[0m", AppVersion) + `
`
}
