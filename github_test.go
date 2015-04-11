package github

import (
	"testing"
	"time"
)

func TestPullCommitsForDay(t *testing.T) {
	// cs4d := PullCommitsForDay("jaffee", "gogurt", 27, 3, 2015)
	// fmt.Println(cs4d)
	// for i := range cs4d {
	// 	fmt.Println(string(cs4d[i]))
	// }
}

func TestOneDayAhead(t *testing.T) {
	startdate = time.Date(2015, time.Month(4), 3, 0, 0, 0, 0, time.Local)
	startminus1date = time.Date(2015, time.Month(4), 2, 0, 0, 0, 0, time.Local)
	boo := oneDayAhead(startminus1date, startdate)
	if !boo {
		t.Fail()
	}
}

func TestOneDayAhead2(t *testing.T) {
	d1 := time.Date(2015, time.Month(1), 1, 0, 0, 0, 0, time.Local)
	d2 := time.Date(2014, time.Month(12), 31, 7, 0, 0, 0, time.Local)
	boo := oneDayAhead(d2, d1)
	if !boo {
		t.Fail()
	}

}
