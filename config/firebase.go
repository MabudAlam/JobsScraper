package config

import (
	"context"
	"fmt"
	"goscraper/models"
	"goscraper/services"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

// Defined a list of scraper functions
var jobScrapers = []struct {
	name    string
	scraper func() ([]models.Job, error)
}{
	{
		name:    "Amazon",
		scraper: services.AmazonScrapper,
	},
	{
		name:    "Atlassian",
		scraper: services.AtlassianScrapper,
	},
	{
		name:    "Coursera",
		scraper: services.CourseraScraper,
	},
	{
		name:    "CRED",
		scraper: services.CredScraper,
	}, {
		name:    "FreshWorks",
		scraper: services.FreshWorksScraper,
	},
	{
		name:    "Gojek",
		scraper: services.GojekScraper,
	},
	{
		name:    "MPL",
		scraper: services.MplScrapper,
	},
}

var app *firebase.App

func init() {
	opt := option.WithCredentialsFile("cred.json")
	var err error
	app, err = firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		fmt.Printf("error initializing app: %v\n", err)
		panic(err)
	}
}
