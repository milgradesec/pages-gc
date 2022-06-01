package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/cloudflare/cloudflare-go"
)

func main() {
	email := os.Getenv("CLOUDFLARE_EMAIl")
	apikey := os.Getenv("CLOUDFLARE_API_KEY")
	accountID := os.Getenv("CLOUDFLARE_ACCOUNT_ID")

	api, err := cloudflare.New(apikey, email)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	projects, _, err := api.ListPagesProjects(ctx, accountID, cloudflare.PaginationOptions{})
	if err != nil {
		log.Fatal(err)
	}

	var projectNames []string
	for _, v := range projects {
		projectNames = append(projectNames, v.Name)
	}

	var projectName string
	prompt := &survey.Select{
		Message: "Select a project:",
		Options: projectNames,
	}
	err = survey.AskOne(prompt, &projectName)
	if err != nil {
		log.Fatal(err)
	}

	opts := cloudflare.ListPagesDeploymentsParams{
		AccountID:   accountID,
		ProjectName: projectName,
	}
	deployments, _, err := api.ListPagesDeployments(ctx, opts)
	if err != nil {
		log.Fatal(err)
	}

	for _, d := range deployments {
		opts := cloudflare.DeletePagesDeploymentParams{
			AccountID:    accountID,
			ProjectName:  projectName,
			DeploymentID: d.ID}
		err = api.DeletePagesDeployment(ctx, opts)
		if err != nil {
			fmt.Printf("❌ Failed to delete deployment id=%s\n", d.ID)
			continue
		}
		fmt.Printf("🧹 Deleted deployment id=%s\n", d.ID)
	}
}
