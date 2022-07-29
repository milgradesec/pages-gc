package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/cloudflare/cloudflare-go"
	"github.com/spf13/viper"
)

var (
	email, apikey, account, project string
)

func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok { //nolint
			fmt.Println("config.yml not found, reading configuration from environment.")
		} else {
			fmt.Printf("Error reading configuration file: %v\n", err)
			os.Exit(1)
		}
	}

	email = viper.GetString("cloudflare.email")
	apikey = viper.GetString("cloudflare.api_key")
	account = viper.GetString("cloudflare.account_id")

	api, err := cloudflare.New(apikey, email)
	if err != nil {
		log.Fatal(err)
	}

	if !viper.IsSet("cloudflare.pages_project") {
		project, err = promptUserToSelectProject(api)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		project = viper.GetString("cloudflare.pages_project")
	}

	if project == "all" {
		err := purgeAllProjects(api)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err = purgeProject(api, project)
		if err != nil {
			log.Fatal(err)
		}
	}
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

func promptUserToSelectProject(api *cloudflare.API) (string, error) {
	resp, _, err := api.ListPagesProjects(context.TODO(), account, cloudflare.PaginationOptions{})
	if err != nil {
		return "", nil
	}

	var projects []string
	for _, v := range resp {
		projects = append(projects, v.Name)
	}

	var project string
	prompt := &survey.Select{
		Message: "Select a project:",
		Options: projects,
	}
	err = survey.AskOne(prompt, &project)
	if err != nil {
		return "", err
	}
	return project, nil
}
