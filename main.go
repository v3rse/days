package main

import (
	"fmt"
	"log"
	"os"

	"github.com/v3rse/days/store"
	"github.com/v3rse/days/tracker"
)

func usage(message string) {
	fmt.Print("usage: days <track|since|reset|list|life <start|end [-v]> > [habit]")
	if message != "" {
		fmt.Print(": ", message)
	}
	fmt.Println("")
	os.Exit(1)
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

	if command == "life" && os.Args[2] == "start" && len(os.Args) <= 3 {
		usage("expected life start date in the form YYYY-MM-DD")
	}

	trackerStore := store.NewFileStore("track.json")
	tracker := tracker.NewTracker(trackerStore)

	defer trackerStore.Close()

	switch command {
	case "track":
		log.Printf("tracking '%s'...", os.Args[2])
		tracker.Track(os.Args[2])
		trackerStore.Save(tracker)
	case "since":
		log.Printf("reading days since %s...", os.Args[2])
		tracker.Since(os.Args[2])
	case "list":
		log.Println("listing days since all habits...")
		tracker.List()
	case "reset":
		log.Printf("reseting count for %s...", os.Args[2])
		tracker.Reset(os.Args[2])
		trackerStore.Save(tracker)
	case "life":
		switch os.Args[2] {
		case "start":
			log.Printf("setting life start date...")
			tracker.LifeStart(os.Args[3])
			trackerStore.Save(tracker)
		case "end":
			log.Printf("calculating approximately how long you may have till the end...")
			verbose := false
			if len(os.Args) == 4 && os.Args[3] == "-v" {
				verbose = true
			}
			tracker.LifeEnd(verbose)
		default:
			usage("")
		}
	default:
		usage("")
	}
}
