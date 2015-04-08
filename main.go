package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const username = "jaffee"

var l, _ = time.LoadLocation("America/Chicago")
var repos = []string{"gogurt", "goplait", "robpike.io", "github"}

var startdate = time.Date(2015, time.Month(3), 19, 0, 0, 0, 0, l)
var startminus1date = time.Date(2015, time.Month(3), 18, 0, 0, 0, 0, l)

const activityPath = "/Users/jaffee/go/src/github.com/jaffee/github/"

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	args := os.Args[1:]
	if len(args) == 3 {
		PullDate(args)
	} else {
		for {
			fmt.Println("Start new figurin")
			ds := figureDays()
			fmt.Printf("DaysNeeded: %v\n", ds)
			argz := make([]string, 3)
			for _, d := range ds {
				argz[0] = strconv.FormatInt(int64(d.Year()), 10)
				argz[1] = strconv.FormatInt(int64(d.Month()), 10)
				argz[2] = strconv.FormatInt(int64(d.Day()), 10)
				fmt.Printf("Now pulling %v %v %v\n", argz[0], argz[1], argz[2])
				PullDate(argz)
			}
			fmt.Println("Done with this round... sleeping\n")
			time.Sleep(12 * time.Minute)
		}
	}
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

func PullDate(args []string) {
	year, err := strconv.Atoi(args[0])
	check(err)
	month, err := strconv.Atoi(args[1])
	check(err)
	day, err := strconv.Atoi(args[2])
	check(err)
	fmt.Printf("Pulling date for args %v\n", args)

	activity := GetDailyActivity(username, repos, day, month, year)

	fname := activityPath + fmt.Sprintf("%04v%02v%02v.activity", year, month, day)

	activityBytes, err := json.Marshal(activity)
	check(err)
	err = ioutil.WriteFile(fname, activityBytes, 0644)
	check(err)
}

type RepoActivity struct {
	Name    string
	Commits []CommitDiff
}

type Commit struct {
	Url     string
	Message string
}

type CommitDiff struct {
	Metadata Commit
	Diff     string
}

func GetDailyActivity(username string, repos []string, day, month, year int) []RepoActivity {
	fmt.Printf("GetDailyActivity uname:%v repos:%v day:%v month:%v year:%v\n", username, repos, day, month, year)
	repoActivities := make([]RepoActivity, len(repos))
	for i := range repos {
		c4d := PullCommitsForDay(username, repos[i], day, month, year)
		repoActivities[i] = c4d
	}
	return repoActivities
}

func PullCommitsForDay(username string, repo string, day, month, year int) RepoActivity {
	fmt.Printf("PullCommitsForDay uname:%v repo:%v\n", username, repo)
	loc, err := time.LoadLocation("Local")
	check(err)
	begOfDay := time.Date(year, time.Month(month), day, 0, 0, 0, 0, loc)
	endOfDay := time.Date(year, time.Month(month), day, 23, 59, 59, 999999999, loc)
	commitsdiffs := GetCommits(username, repo, begOfDay, endOfDay)

	var ra RepoActivity
	ra.Name = repo
	ra.Commits = commitsdiffs
	return ra
}

func GetCommits(username string, repo string, since, until time.Time) []CommitDiff {
	fmt.Printf("GetCommits uname:%v repo:%v since:%v, until:%v\n", username, repo, since, until)
	base_url := "https://api.github.com/repos/" + username + "/" + repo + "/commits"
	time_layout := "2006-01-02T15:04:05Z"
	full_url := base_url + "?since=" + since.Format(time_layout) + "&until=" + until.Format(time_layout)

	body := getBody(full_url)
	var commits []Commit
	json.Unmarshal(body, &commits)
	ret := make([]CommitDiff, len(commits))
	for i := range commits {
		body := getBodyDiff(commits[i].Url)

		ret[i].Metadata = commits[i]
		ret[i].Diff = string(body)
	}
	return ret
}

func getBody(url string) []byte {
	fmt.Printf("getBody %v\n", url)
	resp, err := http.Get(url)
	check(err)
	fmt.Printf("%v\n", resp)
	if resp.StatusCode == 403 { // Forbidden
		waitForRateLimitReset(resp)
		return getBody(url)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	check(err)

	fmt.Printf("getBody returning %v\n", string(body))
	return body
}

type Resp403 struct {
	Message           string
	Documentation_url string
}

func getBodyDiff(url string) []byte {
	fmt.Printf("getBodyDiff %v\n", url)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	check(err)
	req.Header.Del("Accept")
	req.Header.Add("Accept", "application/vnd.github.diff")
	resp, err := client.Do(req)
	check(err)
	if resp.StatusCode == 403 { // Forbidden
		waitForRateLimitReset(resp)
		return getBodyDiff(url)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	check(err)
	fmt.Printf("getBodyDiff returning %v\n", string(body))
	return body
}

func waitForRateLimitReset(resp *http.Response) {
	fmt.Printf("waitForRateLimitReset resp:%v\n", resp)
	remstr := resp.Header.Get("X-RateLimit-Remaining")
	resetstr := resp.Header.Get("X-RateLimit-Reset")
	rem, err := strconv.Atoi(remstr)
	check(err)
	if rem > 0 {
		panic("Thought we hit the rate limit, but we have " + remstr + " remaining")
	}
	reset, err := strconv.ParseInt(resetstr, 10, 0)
	check(err)
	// resetTime := time.Unix(reset, 0)
	diff := reset - time.Now().Unix()
	if diff > 0 {
		dur := time.Duration((diff + 3) * int64(time.Second))
		time.Sleep(dur)
	}
}
