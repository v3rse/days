package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

type Habit struct {
	Action    string    `json:"action"`
	CreatedAt time.Time `json:"createAt"`
}

type Tracker []Habit

func NewTracker(decoder json.Decoder) (Tracker, error) {
	var tracker Tracker
	err := decoder.Decode(&tracker)
	if err != nil {
		err = fmt.Errorf("parsing error: %v", err)
	}
	return tracker, err
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func usage(message string) {
	fmt.Print("usage: days <track|since|reset|list> [habit]")
	if message != "" {
		fmt.Print(": ", message)
	}
	fmt.Println("")
	os.Exit(1)
}

func initialize() *os.File {
	file, err := os.OpenFile("days.json", os.O_RDWR|os.O_CREATE, 0666)
	check(err)

	info, err := file.Stat()
	check(err)

	if info.Size() == 0 {
		file.Write([]byte("[]"))
		file.Seek(0, 0)
	}

	return file
}

func getDaysSince(date time.Time) int {
	return int(time.Since(date).Hours() / 24)
}

func getHabitDescription(habit Habit) string {
	return fmt.Sprintf("%s (%d days since)\n", habit.Action, getDaysSince(habit.CreatedAt))
}

func findHabitPosition(action string, tracker Tracker) (int, error) {
	for i, habit := range tracker {
		if habit.Action == action {
			return i, nil
		}
	}
	return -1, fmt.Errorf("habit not found")
}

func writeTrackerTofile(tracker Tracker, file *os.File, encoder *json.Encoder) {
	file.Seek(0, 0)
	err := encoder.Encode(tracker)
	check(err)
}

func track(habit string, tracker Tracker) Tracker {
	newHabit := Habit{
		habit,
		time.Now(),
	}
	newTracker := append(tracker, newHabit)

	return newTracker
}

func list(tracker Tracker) {
	for i, habit := range tracker {
		fmt.Printf("%d) %s", i+1, getHabitDescription(habit))
	}
}

func since(action string, tracker Tracker) {
	p, err := findHabitPosition(action, tracker)
	check(err)
	fmt.Println(getHabitDescription(tracker[p]))
}

func reset(action string, tracker Tracker) Tracker {
	p, err := findHabitPosition(action, tracker)
	check(err)
	tracker[p].CreatedAt = time.Now()
	return tracker
}

func init() {
	log.SetFlags(0) // remove timestamp
}

func main() {
	command := os.Args[1]
	if command != "list" && len(os.Args) <= 2 {
		usage("expected habit name")
	}

	file := initialize()
	defer file.Close()

	decoder := json.NewDecoder(file)
	encoder := json.NewEncoder(file)

	tracker, err := NewTracker(*decoder)
	check(err)

	switch command {
	case "track":
		log.Printf("tracking '%s'...", os.Args[2])
		updatedTracker := track(os.Args[2], tracker)
		writeTrackerTofile(updatedTracker, file, encoder)
	case "since":
		log.Printf("reading days since %s...", os.Args[2])
		since(os.Args[2], tracker)
	case "list":
		log.Println("listing days since all habits...")
		list(tracker)
	case "reset":
		log.Printf("reseting count for %s...", os.Args[2])
		updatedTracker := reset(os.Args[2], tracker)
		writeTrackerTofile(updatedTracker, file, encoder)
	default:
		usage("")
	}
}
