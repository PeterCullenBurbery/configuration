package main

import (
	"fmt"
	"log"

	"github.com/PeterCullenBurbery/go_functions_002"
	"github.com/PeterCullenBurbery/go_functions_002/date_time_functions"
	"github.com/PeterCullenBurbery/go_functions_002/system_management_functions"
)

func main() {
	go_functions_002.SayHello("Peter")

	ts := date_time_functions.Format_now()
	fmt.Println("⏰ Now:", ts)

	// Install Notepad++.Notepad++ using Winget
	if err := system_management_functions.Winget_install("Notepad++", "Notepad++.Notepad++"); err != nil {
		log.Printf("⚠️ Winget install warning: %v", err)
	} else {
		log.Println("✅ Winget install completed for Notepad++.Notepad++.")
	}
}
