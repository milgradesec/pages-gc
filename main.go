package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/cloudflare/cloudflare-go"
)

var (
	email, apikey, account, project string
)

func main() {
	flag.StringVar(&email, "email", "", "Cloudflare account email")
	flag.StringVar(&apikey, "key", "", "Cloudflare API Key")
	flag.StringVar(&account, "account", "", "Cloudflare account ID")
	flag.StringVar(&project, "project", "", "Pages project name")
	flag.Parse()

	if email == "" || apikey == "" || account == "" {
		email = os.Getenv("CLOUDFLARE_EMAIl")
		apikey = os.Getenv("CLOUDFLARE_API_KEY")
		account = os.Getenv("CLOUDFLARE_ACCOUNT_ID")
	}

	api, err := cloudflare.New(apikey, email)
	if err != nil {
		log.Fatal(err)
	}

	if project == "" {
		projects, _, err := api.ListPagesProjects(context.TODO(), account, cloudflare.PaginationOptions{})
		if err != nil {
			log.Fatal(err)
		}

		var pagesProjects []string
		for _, v := range projects {
			pagesProjects = append(pagesProjects, v.Name)
		}

		var project string
		prompt := &survey.Select{
			Message: "Select a project:",
			Options: pagesProjects,
		}
		err = survey.AskOne(prompt, &project)
		if err != nil {
			log.Fatal(err)
		}
	}

	opts := cloudflare.ListPagesDeploymentsParams{
		ProjectName: project,
	}
	deployments, _, err := api.ListPagesDeployments(context.TODO(), cloudflare.AccountIdentifier(account), opts)
	if err != nil {
		log.Fatal(err)
	}

	for _, d := range deployments {
		err = api.DeletePagesDeployment(context.TODO(), cloudflare.AccountIdentifier(account), project, d.ID)
		if err != nil {
			fmt.Printf("❌ Failed to delete deployment id=%s\n", d.ID)
			continue
		}
		fmt.Printf("🧹 Deleted deployment id=%s\n", d.ID)
	}
}
