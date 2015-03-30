package main

import (
	"fmt"
	"testing"
)

func TestPullCommitsForDay(t *testing.T) {
	cs4d := PullCommitsForDay("jaffee", "gogurt", 27, 3, 2015)
	fmt.Println(cs4d)
	// for i := range cs4d {
	// 	fmt.Println(string(cs4d[i]))
	// }
}
