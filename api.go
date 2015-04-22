package github

import (
	"encoding/json"
	"fmt"
)

const githubApiUrl = "https://api.github.com/"

type Api struct {
	Username string
}

func (a *Api) User() (*User, error) {
	bytes, err := httpGetBody(githubApiUrl + "users/" + a.Username)
	switch err := err.(type) {
	case RateLimitError:
		// maybe wait and try again
		return &User{}, err
	}
	if err != nil {
		return &User{}, err
	}

	user := &User{}
	fmt.Printf("%v\n", bytes)
	err = json.Unmarshal(bytes, user)
	if err != nil {
		return &User{}, err
	}
	return user, nil
}

// TODO functions on Api for accessing github stuff
