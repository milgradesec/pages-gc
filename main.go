package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/cloudflare/cloudflare-go"
	"github.com/spf13/viper"
)

var (
	email, apikey, account, project string
)

func main() {
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok { //nolint
			fmt.Println("Config file not found.")
		} else {
			fmt.Printf("Error reading config file: %v\n", err)
			os.Exit(1)
		}
	}

	if !viper.IsSet("CLOUDFLARE_EMAIL") {
		fmt.Println("CLOUDFLARE_EMAIL not set")
	}
	if !viper.IsSet("CLOUDFLARE_API_KEY") {
		fmt.Println("CLOUDFLARE_API_KEY not set")
	}
	if !viper.IsSet("CLOUDFLARE_ACCOUNT_ID") {
		fmt.Println("CLOUDFLARE_ACCOUNT_ID not set")
	}

	email = viper.GetString("CLOUDFLARE_EMAIl")
	apikey = viper.GetString("CLOUDFLARE_API_KEY")
	account = viper.GetString("CLOUDFLARE_ACCOUNT_ID")

	if viper.IsSet("CLOUDFLARE_PAGES_PROJECT") {
		project = viper.GetString("CLOUDFLARE_PAGES_PROJECT")
	} else {
		// TODO: ask user
	}

	api, err := cloudflare.New(apikey, email)
	if err != nil {
		log.Fatal(err)
	}

	if project == "all" {
		err := purgeAllProjects(api)
		if err != nil {
			log.Fatal(err)
		}
	}
	/*if project == "" {
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
	}*/
}

func purgeAllProjects(api *cloudflare.API) error {
	projects, _, err := api.ListPagesProjects(context.TODO(), account, cloudflare.PaginationOptions{})
	if err != nil {
		return err
	}

	for _, p := range projects {
		err = purgeProject(api, p.Name)
		if err != nil {
			fmt.Printf("❌ Failed to cleanup project '%s': %v", p.Name, err)
		}
	}
	return nil
}

func purgeProject(api *cloudflare.API, name string) error {
	opts := cloudflare.ListPagesDeploymentsParams{
		ProjectName: name,
	}
	deployments, _, err := api.ListPagesDeployments(context.TODO(), cloudflare.AccountIdentifier(account), opts)
	if err != nil {
		return err
	}

	for _, d := range deployments {
		err = api.DeletePagesDeployment(context.TODO(), cloudflare.AccountIdentifier(account), project, d.ID)
		if err != nil {
			fmt.Printf("❌ Failed to delete deployment id=%s project=%s\n", d.ID, d.ProjectName)
			continue
		}
		fmt.Printf("🧹 Deleted deployment id=%s project=%s\n", d.ID, d.ProjectName)
	}
	return nil
}
