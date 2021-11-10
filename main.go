package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

const ESTIMATED_YEARS = 70
const PROGRESS_FILL_COUNT = 50
const MAX_PRINT_WIDTH = 90

type Habit struct {
	Action    string    `json:"action"`
	CreatedAt time.Time `json:"createAt"`
}

type Tracker struct {
	Start  time.Time `json:"start"`
	Habits []Habit   `json:"habits"`
	End    time.Time `json:"end"`
}

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
	fmt.Print("usage: days <track|since|reset|list|life <start|end [-v]> > [habit]")
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
		file.Write([]byte("{\"start\":null, \"habits\":[], \"end\": null}"))
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

func findHabitPosition(action string, habits []Habit) (int, error) {
	for i, habit := range habits {
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
	newHabit := &Habit{
		habit,
		time.Now(),
	}

	return Tracker{
		Start:  tracker.Start,
		Habits: append(tracker.Habits, *newHabit),
	}
}

func list(tracker Tracker) {
	for i, habit := range tracker.Habits {
		fmt.Printf("%d) %s", i+1, getHabitDescription(habit))
	}
}

func since(action string, tracker Tracker) {
	p, err := findHabitPosition(action, tracker.Habits)
	check(err)
	fmt.Println(getHabitDescription(tracker.Habits[p]))
}

func reset(action string, tracker Tracker) Tracker {
	p, err := findHabitPosition(action, tracker.Habits)
	check(err)
	tracker.Habits[p].CreatedAt = time.Now()
	return tracker
}

func lifeStart(startDate string, tracker Tracker) Tracker {
	t, err := time.Parse("2006-01-02", startDate)
	estimatedEnd := t.AddDate(ESTIMATED_YEARS, 0, 0)
	check(err)
	return Tracker{
		Start:  t,
		Habits: tracker.Habits,
		End:    estimatedEnd,
	}
}

func printSummary(daysExpected int, daysSinceStart int, daysToEnd int) {
	fmt.Println()
	fmt.Println("Summary:")
	fmt.Printf(`days expected: %d
days exhausted: %d
days till the END: %d
`, daysExpected, daysSinceStart, daysToEnd)
}

func printProgressBar(daysExpected int, daysSinceStart int) {
	fractionCompleted := float32(daysSinceStart) / float32(daysExpected)
	percentageComplete := int(fractionCompleted * 100)

	fillEnd := int(fractionCompleted * PROGRESS_FILL_COUNT)

	fmt.Println()
	fmt.Println("Progress:")
	for i := 0; i < PROGRESS_FILL_COUNT; i++ {
		if i > fillEnd {
			fmt.Print("=")
			continue
		}
		fmt.Print("#")
	}

	fmt.Printf(" %d%%\n", percentageComplete)
}

func printDetailGrid(daysExpected int, daysSinceStart int) {
	fmt.Println()
	fmt.Println("Details:")

	for i := 0; i < MAX_PRINT_WIDTH+7; i++ {
		fmt.Print("=")
	}

	fmt.Println()
	for i := 0; i < daysExpected; i++ {
		if i%MAX_PRINT_WIDTH == 0 {
			fmt.Println()
			fmt.Printf("%05d| ", i)
		}
		if i > daysSinceStart {
			fmt.Print("o")
			continue
		}
		fmt.Print("*")
	}
}

func lifeEnd(verbose bool, tracker Tracker) {
	daysSinceStart := getDaysSince(tracker.Start)
	daysExpected := ESTIMATED_YEARS * 365
	daysToEnd := daysExpected - daysSinceStart

	printSummary(daysExpected, daysSinceStart, daysToEnd)

	printProgressBar(daysExpected, daysSinceStart)

	if verbose {
		printDetailGrid(daysExpected, daysSinceStart)
	}
}

func init() {
	log.SetFlags(0) // remove timestamp
}

func main() {
	if len(os.Args) < 2 {
		usage("")
	}

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
	case "life":
		switch os.Args[2] {
		case "start":
			log.Printf("setting life start date...")
			updatedTracker := lifeStart(os.Args[3], tracker)
			writeTrackerTofile(updatedTracker, file, encoder)
		case "end":
			log.Printf("calculating approximately how long you may have till the end...")
			verbose := false
			if len(os.Args) == 4 && os.Args[3] == "-v" {
				verbose = true
			}
			lifeEnd(verbose, tracker)
		default:
			usage("")
		}
	default:
		usage("")
	}
}
