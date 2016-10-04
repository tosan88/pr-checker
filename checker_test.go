package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseLinkHeader(t *testing.T) {
	assert := assert.New(t)
	link := `<https://api.github.com/organizations/3502508/repos?page=3>; rel="next",
 <https://api.github.com/organizations/3502508/repos?page=22>; rel="last",
 <https://api.github.com/organizations/3502508/repos?page=1>; rel="first",
 <https://api.github.com/organizations/3502508/repos?page=1>; rel="prev"`
	nextPageLink := parseLinkHeader(link)
	fmt.Printf("%v\n", nextPageLink)
	assert.Equal("https://api.github.com/organizations/3502508/repos?page=3", nextPageLink)
}
