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
	email, apikey, accountID string
)

func main() {
	flag.StringVar(&email, "email", "", "Cloudflare account email")
	flag.StringVar(&apikey, "key", "", "Cloudflare API Key")
	flag.StringVar(&accountID, "account", "", "Cloudflare account ID")
	flag.Parse()

	if email == "" || apikey == "" || accountID == "" {
		email = os.Getenv("CLOUDFLARE_EMAIl")
		apikey = os.Getenv("CLOUDFLARE_API_KEY")
		accountID = os.Getenv("CLOUDFLARE_ACCOUNT_ID")
	}

	api, err := cloudflare.New(apikey, email)
	if err != nil {
		log.Fatal(err)
	}

	projects, _, err := api.ListPagesProjects(context.TODO(), accountID, cloudflare.PaginationOptions{})
	if err != nil {
		log.Fatal(err)
	}

	var pagesProjects []string //nolint
	for _, v := range projects {
		pagesProjects = append(pagesProjects, v.Name)
	}

	var projectName string
	prompt := &survey.Select{
		Message: "Select a project:",
		Options: pagesProjects,
	}
	err = survey.AskOne(prompt, &projectName)
	if err != nil {
		log.Fatal(err)
	}

	opts := cloudflare.ListPagesDeploymentsParams{
		AccountID:   accountID,
		ProjectName: projectName,
	}
	deployments, _, err := api.ListPagesDeployments(context.TODO(), opts)
	if err != nil {
		log.Fatal(err)
	}

	for _, d := range deployments {
		opts := cloudflare.DeletePagesDeploymentParams{
			AccountID:    accountID,
			ProjectName:  projectName,
			DeploymentID: d.ID}
		err = api.DeletePagesDeployment(context.TODO(), opts)
		if err != nil {
			fmt.Printf("❌ Failed to delete deployment id=%s\n", d.ID)
			continue
		}
		fmt.Printf("🧹 Deleted deployment id=%s\n", d.ID)
	}
}