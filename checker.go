package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
)

var linkRex = regexp.MustCompile(`^.*<([^>]+)>; rel="next".*`)

func (chkr *checker) collectPullRequests(firstPageLink string) (prs []PullReq) {
	prsPart, nextPageLink := chkr.getPullRequests(firstPageLink)
	prs = append(prs, prsPart...)
	for nextPageLink != "" {
		prsPart, nextPageLink = chkr.getPullRequests(nextPageLink)
		prs = append(prs, prsPart...)
	}
	log.Printf("DEBUG - # of PRs: %v\n", len(prs))
	return
}

func (chkr *checker) collectContributors(firstPageLink string) (contributors []User) {
	contributorsPart, nextPageLink := chkr.getContributorsOfRepo(firstPageLink)
	contributors = append(contributors, contributorsPart...)
	for nextPageLink != "" {
		contributorsPart, nextPageLink = chkr.getContributorsOfRepo(nextPageLink)
		contributors = append(contributors, contributorsPart...)
	}
	//log.Printf("DEBUG - # of contributors: %v\n", len(contributors))
	return
}

func (chkr *checker) collectRepos() (repos []Repo) {
	reposPart, nextPageLink := chkr.getRepos("https://api.github.com/orgs/Financial-Times/repos")
	repos = append(repos, reposPart...)
	for nextPageLink != "" {
		reposPart, nextPageLink = chkr.getRepos(nextPageLink)
		repos = append(repos, reposPart...)
	}
	log.Printf("DEBUG - # of repos: %v\n", len(repos))
	return
}

func (chkr *checker) decideCoreContributors() {
	if chkr.conf.contributors != "" {
		chkr.coreContributors = chkr.getPassedContributors()
		return
	}

	users, nextPageLink := chkr.getContributorsOfRepo("https://api.github.com/repos/Financial-Times/up-service-files/contributors")
	for _, user := range users {
		chkr.coreContributors[user.User] = user
	}
	for nextPageLink != "" {
		users, nextPageLink = chkr.getContributorsOfRepo(nextPageLink)
		for _, user := range users {
			chkr.coreContributors[user.User] = user
		}
	}
	chkr.filter()
}

func (chkr *checker) isAnyCoreContributor(users []User) bool {
	for _, contr := range users {
		if user, found := chkr.coreContributors[contr.User]; found {
			log.Printf("Found core contributor: %v\n", user.User)
			return true
		}
	}
	return false
}

func (chkr *checker) getPassedContributors() map[string]User {
	contributors := strings.Split(chkr.conf.contributors, ",")
	var users map[string]User
	for _, contributor := range contributors {
		users[contributor] = User{contributor}
	}

	return users
}

func (chkr *checker) getContributorsOfRepo(url string) ([]User, string) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("ERROR http.NewRequest - %v\n", err)
	}

	if chkr.conf.token != "" {
		req.Header.Add("Authorization", fmt.Sprintf("token %v", chkr.conf.token))
	}

	resp, err := chkr.client.Do(req)
	if err != nil {
		log.Fatalf("ERROR client.Do - %v\n", err)
	}
	defer resp.Body.Close()

	var users []User
	err = json.NewDecoder(resp.Body).Decode(&users)
	if err != nil {
		log.Fatalf("ERROR Decode - %v\n", err)
	}

	return users, parseLinkHeader(resp.Header.Get("Link"))
}

// Example:
// Link: <https://api.github.com/organizations/3502508/repos?page=3>; rel="next",
// <https://api.github.com/organizations/3502508/repos?page=22>; rel="last",
// <https://api.github.com/organizations/3502508/repos?page=1>; rel="first",
// <https://api.github.com/organizations/3502508/repos?page=1>; rel="prev"
func parseLinkHeader(link string) string {
	submatch := linkRex.FindStringSubmatch(link)
	if len(submatch) != 2 {
		return ""
	}

	return submatch[1]
}

func (chkr *checker) getRepos(url string) ([]Repo, string) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("ERROR http.NewRequest - %v\n", err)
	}

	if chkr.conf.token != "" {
		req.Header.Add("Authorization", fmt.Sprintf("token %v", chkr.conf.token))
	}

	resp, err := chkr.client.Do(req)
	if err != nil {
		log.Fatalf("ERROR client.Do - %v\n", err)
	}
	defer resp.Body.Close()

	var repos []Repo
	err = json.NewDecoder(resp.Body).Decode(&repos)
	if err != nil {
		log.Fatalf("ERROR Decode - %v\n", err)
	}

	return repos, parseLinkHeader(resp.Header.Get("Link"))
}

func (chkr *checker) getPullRequests(url string) ([]PullReq, string) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("ERROR http.NewRequest - %v\n", err)
	}

	if chkr.conf.token != "" {
		req.Header.Add("Authorization", fmt.Sprintf("token %v", chkr.conf.token))
	}

	resp, err := chkr.client.Do(req)
	if err != nil {
		log.Fatalf("ERROR client.Do - %v\n", err)
	}
	defer resp.Body.Close()

	var pullReqs []PullReq
	err = json.NewDecoder(resp.Body).Decode(&pullReqs)
	if err != nil {
		log.Fatalf("ERROR Decode - %v\n", err)
	}

	return pullReqs, parseLinkHeader(resp.Header.Get("Link"))
}

func (chkr *checker) filter() {
	delete(chkr.coreContributors, "matthew-andrews")
	delete(chkr.coreContributors, "dgem")
	delete(chkr.coreContributors, "duffj")
}
