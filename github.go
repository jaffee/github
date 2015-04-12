package github

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

var l, _ = time.LoadLocation("America/Chicago")

func check(e error) {
	if e != nil {
		panic(e)
	}
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

func GetDailyActivity(username string, repos []string, date time.Time) []RepoActivity {
	fmt.Printf("GetDailyActivity uname:%v repos:%v date:%v\n", username, repos, date)
	repoActivities := make([]RepoActivity, len(repos))
	for i := range repos {
		c4d := PullCommitsForDay(username, repos[i], date)
		repoActivities[i] = c4d
	}
	return repoActivities
}

func PullCommitsForDay(username string, repo string, date time.Time) RepoActivity {
	fmt.Printf("PullCommitsForDay uname:%v repo:%v\n", username, repo)

	begOfDay := date
	endOfDay := time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 999999999, l)
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
