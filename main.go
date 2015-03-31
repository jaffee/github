package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

const username = "jaffee"

var repos = []string{"gogurt", "goplait", "robpike.io"}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	args := os.Args[1:]
	if len(args) < 3 {
		fmt.Println("Usage: github YYYY MM DD")
		return
	}
	year, err := strconv.Atoi(args[0])
	month, err := strconv.Atoi(args[1])
	day, err := strconv.Atoi(args[2])

	activity := GetDailyActivity(username, repos, day, month, year)

	fname := fmt.Sprintf("%v%v%v.activity", year, month, day)

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
	Url string
}

type CommitDiff struct {
	Metadata Commit
	Diff     string
}

func GetDailyActivity(username string, repos []string, day, month, year int) []RepoActivity {
	repoActivities := make([]RepoActivity, len(repos))
	for i := range repos {
		repoActivities[i] = PullCommitsForDay(username, repos[i], day, month, year)
	}
	return repoActivities
}

func PullCommitsForDay(username string, repo string, day, month, year int) RepoActivity {
	loc, _ := time.LoadLocation("Local")
	begOfDay := time.Date(year, time.Month(month), day, 0, 0, 0, 0, loc)
	endOfDay := time.Date(year, time.Month(month), day, 23, 59, 59, 999999999, loc)
	commitsdiffs := GetCommits(username, repo, begOfDay, endOfDay)
	var ra RepoActivity
	ra.Name = repo
	ra.Commits = commitsdiffs
	return ra
}

func getBody(url string) []byte {
	resp, err := http.Get(url)
	if r := handle_url_err(url, err); r != nil {
		return r
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if r := handle_url_err(url, err); r != nil {
		return r
	}
	//	ppJson(body)
	// fmt.Printf("Got body for URL %v\n%v\n", url, string(body))
	return body
}

func getBodyDiff(url string) []byte {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	check(err)
	req.Header.Del("Accept")
	req.Header.Add("Accept", "application/vnd.github.diff")
	resp, err := client.Do(req)
	check(err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if r := handle_url_err(url, err); r != nil {
		return r
	}
	return body
}

func handle_url_err(url string, err error) []byte {
	if err != nil {
		fmt.Printf("Error with URL: %v\n", url)
		return make([]byte, 0)
	}
	return nil
}

func GetCommits(username string, repo string, since, until time.Time) []CommitDiff {
	// TODO get diffs instead of just the commit object, this will
	// involve sending a special header "Accept:
	// application/vnd.github.diff" which in order to do you have to
	// create a custom client and request object and set the headers
	// on it.

	base_url := "https://api.github.com/repos/" + username + "/" + repo + "/commits"
	time_layout := "2006-01-02T15:04:05Z"
	full_url := base_url + "?since=" + since.Format(time_layout) + "&until=" + until.Format(time_layout)
	fmt.Println(full_url)
	body := getBody(full_url)
	var commits []Commit
	json.Unmarshal(body, &commits)
	ret := make([]CommitDiff, len(commits))
	for i := range commits {
		body = getBodyDiff(commits[i].Url)
		ret[i].Metadata = commits[i]
		ret[i].Diff = string(body)
	}
	return ret //make([]byte, 1) //
}
