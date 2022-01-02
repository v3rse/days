package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"time"
  "path/filepath"

	"github.com/v3rse/days/store"
	"github.com/v3rse/days/tracker"
	"github.com/v3rse/days/utils"
)


func usage(message string) {
	fmt.Print("usage: days <track|since|reset|list|life <start|end [-v]> >|journal <write|read> [habit]")
	if message != "" {
		fmt.Print(": ", message)
	}
	fmt.Println("")
	os.Exit(1)
}

func init() {
	log.SetFlags(0) // remove timestamp
}

type Entry struct {
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"createdAt"`
}

type Journal struct {
	Entries   []Entry   `json:"entries"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`
}

func (j *Journal) write(text string) {
	entry := Entry{
		text,
		time.Now(),
	}
	j.Entries = append(j.Entries, entry)
	j.UpdatedAt = time.Now()
}

func (j *Journal) list(start string, end string) []Entry {
	//assumption: no modification is made to entry order in the data file
	// if start not set use today
	startDate := time.Now().Truncate(24 * time.Hour)

	if start != "" {
		var err error
		startDate, err = time.Parse("2006-01-02", start)
		utils.Check(err)
	}

	endDate := startDate
	if end != "" {
		var err error
		endDate, err = time.Parse("2006-01-02", end)
		utils.Check(err)
	}

	// find index of entry with earliest entry with created date
  // matching or after start date
	startIndex := -1
	for i, entry := range j.Entries {
		if startDate.Equal(entry.CreatedAt.Truncate(24 * time.Hour)) ||
    startDate.Before(entry.CreatedAt.Truncate(24 * time.Hour)) {
			startIndex = i
			break
		} 
	}

	if startIndex == -1 {
		return []Entry{}
	}

  // find index of entry from the bottom with created date
  // matching or before end date 
	endIndex := len(j.Entries)
	for i := len(j.Entries) - 1; i >= 0; i-- {
		entry := j.Entries[i]
		if endDate.Equal(entry.CreatedAt.Truncate(24 * time.Hour)) ||
    endDate.After(entry.CreatedAt.Truncate(24 * time.Hour)) {
			endIndex = i + 1
			break
    }
	}

	return j.Entries[startIndex:endIndex]
}

func NewJournal(store store.FileStore) Journal {
	var journal Journal

	store.Load(&journal)

	if (journal.CreatedAt == time.Time{}) {
		journal.CreatedAt = time.Now()
	}

	return journal
}

func promptTextfield(message string) string {
	log.Printf(message)
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	text = strings.Replace(text, "\n", "", -1)

	utils.Check(err)

	return text
}

func militaryTime(t time.Time) string {
	return fmt.Sprintf("%02d%02d", t.Hour(), t.Minute())
}

func printEntries(entries []Entry) {
	var currentPrintDate time.Time

	for _, entry := range entries {
		if reflect.ValueOf(currentPrintDate).IsZero() {
			fmt.Println("")
			currentPrintDate = entry.CreatedAt.Truncate(24 * time.Hour)
			fmt.Println(currentPrintDate.Format("2006-01-02"))
			fmt.Println("----------")
		}

		currentDate := entry.CreatedAt.Truncate(24 * time.Hour)

		if !currentPrintDate.Equal(currentDate) {
			fmt.Println("")
			currentPrintDate = currentDate
			fmt.Println(currentPrintDate.Format("2006-01-02"))
			fmt.Println("----------")
		}

		fmt.Println("")
		fmt.Println(militaryTime(entry.CreatedAt))
		fmt.Println("----")
		fmt.Println(entry.Text)
	}
}

func loadData() (tracker.Tracker, Journal, store.FileStore, store.FileStore) {
 homeDirPath, _ := os.UserHomeDir()
 appDir, _ := filepath.Abs(homeDirPath + "/.days")

  if _, err := os.Stat(appDir); os.IsNotExist(err) {
    err := os.Mkdir(appDir, os.ModePerm)
    utils.Check(err)
  }

	trackInit := []byte("{\"start\":null, \"habits\":[], \"end\": null}")
	trackerStore := store.NewFileStore(appDir+"/track.json", trackInit)
	trk := tracker.NewTracker(trackerStore)

	journalInit := []byte("{\"entries\":[], \"createdAt\":null, \"updatedAt\": null}")
	journalStore := store.NewFileStore(appDir+"/journal.json", journalInit)
	jnl := NewJournal(journalStore)

  return trk, jnl, trackerStore, journalStore
}

func main() {
	if len(os.Args) < 2 {
		usage("")
	}

	command := os.Args[1]
	if command != "list" && len(os.Args) <= 2 {
		usage("expected habit name")
	}

	if command == "life" && os.Args[2] == "start" && len(os.Args) <= 3 {
		usage("expected life start date in the form YYYY-MM-DD")
	}

  trk, jnl, trackerStore, journalStore := loadData()

	defer trackerStore.Close()
	defer journalStore.Close()

	switch command {
	case "track":
		log.Printf("tracking '%s'...", os.Args[2])
		trk.Track(os.Args[2])
		trackerStore.Save(trk)
	case "since":
		log.Printf("reading days since %s...", os.Args[2])
		trk.Since(os.Args[2])
	case "list":
		log.Println("listing days since all habits...")
		trk.List()
	case "reset":
		log.Printf("reseting count for %s...", os.Args[2])
		trk.Reset(os.Args[2])
		trackerStore.Save(trk)
	case "life":
		switch os.Args[2] {
		case "start":
			log.Printf("setting life start date...")
			trk.LifeStart(os.Args[3])
			trackerStore.Save(trk)
		case "end":
			log.Printf("calculating approximately how long you may have till the end...")
			verbose := false
			if len(os.Args) == 4 && os.Args[3] == "-v" {
				verbose = true
			}
			trk.LifeEnd(verbose)
		default:
			usage("")
		}
	case "journal":
		switch os.Args[2] {
		case "write":
			entry := promptTextfield("write your entry below. Hit [Enter] when done:")
			jnl.write(entry)
			journalStore.Save(jnl)
		case "read":
			var entries []Entry
			if len(os.Args) == 5 {
				entries = jnl.list(os.Args[3], os.Args[4])
			} else if len(os.Args) == 4 {
				entries = jnl.list(os.Args[3], "")
			} else {
				entries = jnl.list("", "")
			}
			printEntries(entries)
		default:
			usage("")
		}
	default:
		usage("")
	}
}
