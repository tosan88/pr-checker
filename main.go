package main

import (
	"github.com/jawher/mow.cli"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"
)

type checker struct {
	conf             config
	coreContributors map[string]User
	client           *http.Client
}

type config struct {
	token        string
	contributors string
	minDays      int
}

func main() {
	log.Printf("Application starting with args %s", os.Args)
	app := cli.App("pr-checker", "Checks FT UPP's PRs which are too long time open")

	token := app.String(cli.StringOpt{
		Name:   "token",
		Value:  "",
		Desc:   "The GitHub Api's OAuth Token.",
		EnvVar: "TOKEN",
	})

	contributors := app.String(cli.StringOpt{
		Name:   "contributors",
		Value:  "",
		Desc:   "The list of contributors. Only those repos will be considered where these people contributed. Optional",
		EnvVar: "CONTRIBUTORS",
	})

	minDays := app.String(cli.StringOpt{
		Name:   "min-days",
		Value:  "14",
		Desc:   "The number of minimum days which an open PR could stay open. Only PRs which are opened more than that number of days are retrieved. Optional",
		EnvVar: "MIN_DAYS",
	})

	app.Before = func() {
		if *token == "" {
			app.PrintHelp()
			os.Exit(1)
		}
	}

	app.Action = func() {
		defer func(start time.Time) {
			elapsed := time.Since(start)
			log.Printf("Application finished. It was active %v seconds", elapsed.Seconds())
		}(time.Now())

		tr := &http.Transport{
			MaxIdleConnsPerHost: 32,
			Dial: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial,
		}
		c := &http.Client{
			Transport: tr,
			Timeout:   30 * time.Second,
		}

		days, err := strconv.Atoi(*minDays)
		if err != nil {
			log.Printf("INFO - Using default min days: 14\n")
			days = 14
		}

		conf := config{
			token:        *token,
			contributors: *contributors,
			minDays:      days,
		}

		chkr := checker{
			conf:             conf,
			client:           c,
			coreContributors: make(map[string]User),
		}

		runChecker(&chkr)

	}
	if err := app.Run(os.Args); err != nil {
		log.Fatalln(err)
	}
}

func runChecker(chkr *checker) {
	chkr.decideCoreContributors()
	repos := chkr.collectRepos()
	var count int
	defer func() {
		log.Printf("There are %v pull requests open across FT repositories relevant for UPP\n", count)
	}()
	for _, repo := range repos {
		contributors := chkr.collectContributors(repo.ContributorsURL)

		if chkr.isAnyCoreContributor(contributors) {
			prs := chkr.collectPullRequests(repo.PullsURL.Url)

			for _, pr := range prs {
				minDateTime := time.Now().AddDate(0, 0, -chkr.conf.minDays)

				parsedCreatedAt, err := time.Parse("2006-01-02T15:04:05Z", pr.CreatedAt)
				if err != nil {
					log.Fatalf("ERROR time.Parse - %v\n", err)
				}
				if parsedCreatedAt.Before(minDateTime) {
					realUserName := chkr.getUser(pr.User.UserURL)
					log.Printf("PR %v (%v) open by %v(%v) since %v, updated at %v\n", pr.HTMLURL, pr.Title, pr.User.User, realUserName, pr.CreatedAt, pr.UpdatedAt)
					count++
				}

			}
		}
	}
}
