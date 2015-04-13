package main

import (
	"encoding/json"
	"fmt"
	"github.com/jaffee/github"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

var startdate = time.Date(2015, time.Month(3), 19, 0, 0, 0, 0, l)
var startminus1date = time.Date(2015, time.Month(3), 18, 0, 0, 0, 0, l)

const activityPath = "/Users/jaffee/go/src/github.com/jaffee/github/"

var repos = []string{"gogurt", "goplait", "robpike.io", "github"}
var l, _ = time.LoadLocation("America/Chicago")

const username = "jaffee"

func main() {
	args := os.Args[1:]
	if len(args) == 3 {
		date := argsToDate(args)
		activity := github.GetDailyActivity(username, repos, argsToDate(args))
		fname := activityPath + fmt.Sprintf("%04v%02v%02v.activity", date.Year(), int(date.Month()), date.Day())
		writeActivity(activity, fname)

	} else {
		for {
			fmt.Println("Start new figurin")
			ds := figureDays()
			fmt.Printf("DaysNeeded: %v\n", ds)
			for _, d := range ds {
				fmt.Printf("Now pulling %v\n", d)
				activity := github.GetDailyActivity(username, repos, d)
				fname := activityPath + fmt.Sprintf("%04v%02v%02v.activity", d.Year(), int(d.Month()), d.Day())
				writeActivity(activity, fname)
			}
			fmt.Println("Done with this round... sleeping\n")
			time.Sleep(12 * time.Minute)
		}
	}
}

func writeActivity(activity []github.RepoActivity, fname string) {
	activityBytes, err := json.Marshal(activity)
	check(err)
	err = ioutil.WriteFile(fname, activityBytes, 0644)
	check(err)

}

func argsToDate(args []string) time.Time {
	year, err := strconv.Atoi(args[0])
	month, err := strconv.Atoi(args[1])
	day, err := strconv.Atoi(args[2])
	check(err)

	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, l)
	return date
}

func figureDays() []time.Time {
	files, err := ioutil.ReadDir(activityPath)
	check(err)
	now := time.Now()
	// TODO make now only down to the day, not hours etc.
	prevDate := startdate.Add(time.Hour * -24)
	var daysNeeded []time.Time

	for _, f := range files {
		name := f.Name()
		if strings.HasSuffix(name, ".activity") {
			loc := strings.LastIndex(name, ".activity")
			datestr := name[:loc]
			if len(datestr) != 8 {
				fmt.Printf("You have a .activity file with a malformatted name: %v\n", name)
				continue
			}
			year, err := strconv.Atoi(datestr[:4])
			month, err := strconv.Atoi(datestr[4:6])
			day, err := strconv.Atoi(datestr[6:8])
			if err != nil {
				continue // we don't hafta deal with these crappy filenames
			}
			date := time.Date(year, time.Month(month), day,
				0, 0, 0, 0, l)
			fmt.Printf("Got the date %v\n", date)
			if date.After(startminus1date) && date.Before(now) {
				for !oneDayAhead(prevDate, date) {
					prevDate = prevDate.Add(time.Hour * 24)
					daysNeeded = append(daysNeeded, prevDate)
				}
				prevDate = date
			}
		}
	}
	for now.Sub(prevDate) >= time.Hour*24 {
		prevDate = prevDate.Add(time.Hour * 24)
		daysNeeded = append(daysNeeded, prevDate)
	}
	return daysNeeded

}

func oneDayAhead(prevDate time.Time, date time.Time) bool {
	nd := time.Date(prevDate.Year(), prevDate.Month(), prevDate.Day()+1, 0, 0, 0, 0, l)
	return nd.Year() == date.Year() && nd.Month() == date.Month() && nd.Day() == date.Day()
}
