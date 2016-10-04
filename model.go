package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

type User struct {
	User string `json:"login"`
}

type GHURL struct {
	Url string
}

type Repo struct {
	Name            string `json:"name"`
	ContributorsURL string `json:"contributors_url"`
	PullsURL        GHURL  `json:"pulls_url"`
}

type PullReq struct {
	HTMLURL   string `json:"html_url"`
	Title     string `json:"title"`
	User      string `json:"user.login"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (url *GHURL) UnmarshalJSON(data []byte) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	var rawURL string
	if err := dec.Decode(&rawURL); err != nil {
		return fmt.Errorf("Error Decoding GHURL: %v", err)
	}
	url.Url = strings.Split(rawURL, "{")[0]
	return nil
}
