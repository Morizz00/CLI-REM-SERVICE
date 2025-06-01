package main

import (
	"container/heap"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/olebedev/when"
	"github.com/olebedev/when/rules/common"
	"github.com/olebedev/when/rules/en"
)

type Reminder struct {
	time     time.Time
	message  string
	priority string
	index    int
}

type PriorityQueue []*Reminder

func (pq PriorityQueue) Len() int { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].time.Before(pq[j].time)
}
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}
func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Reminder)
	item.index = n
	*pq = append(*pq, item)
}
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1
	*pq = old[0 : n-1]
	return item
}
func main() {
	if len(os.Args) < 3 {
		fmt.Printf("Usage:%s<hh:mm> <text message\n>[priority=low|med|high]\n", os.Args[0])
		os.Exit(1)
	}
	now := time.Now()
	w := when.New(nil)
	w.Add(en.All...)
	w.Add(common.All...)

	args := os.Args[1:]
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	i := 0
	for i < len(args) {
		tmStr := args[i]
		tmParsed, err := w.Parse(tmStr, now)
		if err != nil || tmParsed == nil {
			fmt.Printf("Unable to parse time: %s\n", tmStr)
			os.Exit(2)
		}
		if now.After(tmParsed.Time) {
			fmt.Printf("Set a future time, %s is in the past\n", tmStr)
			os.Exit(3)
		}
		i++
		msgParts := []string{}
		for i < len(args) {
			if strings.Contains(args[i], ":") { // next time detected
				break
			}
			if args[i] == "-priority" {
				break
			}
			msgParts = append(msgParts, args[i])
			i++
		}
		if len(msgParts) == 0 {
			fmt.Println("Missing message for reminder at", tmStr)
			os.Exit(1)
		}
		message := strings.Join(msgParts, " ")

		priority := "low"
		if i < len(args) && args[i] == "-priority" {
			i++
			if i < len(args) {
				priority = strings.ToLower(args[i])
				i++
			} else {
				fmt.Println("Priority flag provided but no priority value")
				os.Exit(1)
			}
		}
		heap.Push(&pq, &Reminder{
			time:     tmParsed.Time,
			message:  message,
			priority: priority,
		})

		fmt.Printf("Reminder scheduled for: %s [%s] %s\n", tmParsed.Time.Format("2006-01-02 15:04:05"), priority, message)
	}
	for pq.Len() > 0 {
		next := heap.Pop(&pq).(*Reminder)
		now = time.Now()
		if next.time.After(now) {
			sleepDuration := next.time.Sub(now)
			fmt.Printf("Waiting %s for next reminder...\n", sleepDuration.Round(time.Second))
			time.Sleep(sleepDuration)
		}
		icon := "assets/information.png"
		title := "Reminder"
		switch next.priority {
		case "medium":
			icon = "assets/warning.png"
			title = "Reminder"
		case "high":
			icon = "assets/warning.png"
			title = "Urgent Reminder"
		}
		fmt.Printf("Showing reminder: [%s] %s\n", next.priority, next.message)
		err := beeep.Alert(title, next.message, icon)
		if err != nil {
			fmt.Println("Failed to show notification:", err)
		}
	}
	fmt.Println("No more reminders. Exiting.")
}
