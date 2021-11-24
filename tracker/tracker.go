package tracker

import (
	"fmt"
	"strconv"
	"time"

	"github.com/v3rse/days/store"
	"github.com/v3rse/days/utils"
)

const ESTIMATED_YEARS = 70
const PROGRESS_FILL_COUNT = 50
const MAX_PRINT_WIDTH = 90
const MARGIN_OFFSET = 8

type Habit struct {
	Action    string    `json:"action"`
	CreatedAt time.Time `json:"createAt"`
}

type Tracker struct {
	Start  time.Time `json:"start"`
	Habits []Habit   `json:"habits"`
	End    time.Time `json:"end"`
}

func NewTracker(trackerStore store.FileStore) Tracker {
	var tracker Tracker

	trackerStore.Load(&tracker)

	return tracker
}

func (tracker *Tracker) Track(habit string) {
	newHabit := &Habit{
		habit,
		time.Now(),
	}

	tracker.Habits = append(tracker.Habits, *newHabit)
}

func (tracker *Tracker) List() {
	for i, habit := range tracker.Habits {
		fmt.Printf("%d) %s", i+1, getHabitDescription(habit))
	}
}

func (tracker *Tracker) Since(action string) {
	p := findHabitPosition(action, tracker.Habits)
	fmt.Println(getHabitDescription(tracker.Habits[p]))
}

func (tracker *Tracker) Reset(action string) {
	p := findHabitPosition(action, tracker.Habits)
	tracker.Habits[p].CreatedAt = time.Now()
}

func (tracker *Tracker) LifeStart(startDate string) {
	t, err := time.Parse("2006-01-02", startDate)
	estimatedEnd := t.AddDate(ESTIMATED_YEARS, 0, 0)
	utils.Check(err)
	tracker.Start = t
	tracker.End = estimatedEnd
}

func (tracker *Tracker) LifeEnd(verbose bool) {
	daysSinceStart := getDaysSince(tracker.Start)
	daysExpected := ESTIMATED_YEARS * 365
	daysToEnd := daysExpected - daysSinceStart

	printSummary(daysExpected, daysSinceStart, daysToEnd)

	printProgressBar(daysExpected, daysSinceStart)

	if verbose {
		printDetailGrid(daysExpected, daysSinceStart)
	}
}

func getDaysSince(date time.Time) int {
	return int(time.Since(date).Hours() / 24)
}

func getHabitDescription(habit Habit) string {
	return fmt.Sprintf("%s (%d days since)\n", habit.Action, getDaysSince(habit.CreatedAt))
}

func findHabitPositionByAction(action string, habits []Habit) (int, error) {
	for i, habit := range habits {
		if habit.Action == action {
			return i, nil
		}
	}
	return -1, fmt.Errorf("habit not found")
}

func findHabitPosition(action string, habits []Habit) int {
	var p int
	index, err := strconv.ParseInt(action, 10, 64)
	if err != nil {
		var err error
		p, err = findHabitPositionByAction(action, habits)
		utils.Check(err)
	} else {
		p = int(index) - 1
	}

	return p
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

	for i := 1; i < MAX_PRINT_WIDTH; i++ {
		fmt.Print("=")
	}

	fmt.Println()
	for i := 0; i < daysExpected; i++ {
		if i%(MAX_PRINT_WIDTH-MARGIN_OFFSET) == 0 {
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
