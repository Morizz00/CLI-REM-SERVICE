package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/olebedev/when"
	"github.com/olebedev/when/rules/common"
	"github.com/olebedev/when/rules/en"
)

const (
	markName  = "Golang CLI Reminder"
	markValue = "1"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("Usage:%s<hh:mm> <text message\n>[priority=low|med|high]\n", os.Args[0])
		os.Exit(1)
	}
	now := time.Now()
	w := when.New(nil)
	w.Add(en.All...)
	w.Add(common.All...)
	t, err := w.Parse(os.Args[1], now)
	if err != nil || t == nil {
		fmt.Println("unable to parse time:", os.Args[1])
		os.Exit(2)
	}
	if now.After(t.Time) {
		fmt.Println("set a future time")
		os.Exit(3)
	}
	diff := t.Time.Sub(now)
	priority := "low"
	messagePri := []string{}
	arg := os.Args[2:]
	for i := 0; i < len(arg); i++ {
		if arg[i] == "-priority" && i+1 < len(arg) {
			priority = arg[i+1]
			i++
		} else {
			messagePri = append(messagePri, arg[i])
		}
	}
	message := strings.Join(messagePri, " ")
	icon := "assets/information.png"
	title := "Reminder"
	switch priority {
	case "medium":
		icon = "assets/warning.png"
		title = "Reminder"
	case "high":
		icon = "assets/warning.png"
		title = "Urgent Reminder"
	}
	if os.Getenv(markName) == markValue {
		time.Sleep(diff)
		err := beeep.Alert(title, message, icon)
		if err != nil {
			fmt.Println(err)
			fmt.Println("Failed to show notif:", err)
			os.Exit(4)
		}
	} else {
		cmd := exec.Command(os.Args[0], os.Args[1:]...)
		cmd.Env = append(os.Environ(), fmt.Sprintf("%s=%s", markName, markValue))
		if err := cmd.Start(); err != nil {
			fmt.Println("Couldn't set backrgoudn process:", err)
			os.Exit(5)
		}
		fmt.Printf("Reminder set for: %s\n", t.Time.Format("2006-01-02 15:04:05"))
		fmt.Println("Reminder will be displayed after", diff.Round(time.Second))
		os.Exit(0)
	}
}

