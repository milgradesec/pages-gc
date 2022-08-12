package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/cloudflare/cloudflare-go"
)

var (
	email   = flag.String("email", "", "Cloudflare account email")
	apikey  = flag.String("apikey", "", "Cloudflare API key")
	account = flag.String("account", "", "Cloudflare account ID")
	project = flag.String("project", "all", "Pages project name")
)

func main() {
	flag.Parse()

	if *email == "" || *apikey == "" || *account == "" {
		flag.PrintDefaults()
		return
	}

	api, err := cloudflare.New(*apikey, *email)
	if err != nil {
		log.Fatal(err)
	}

	if *project == "all" {
		err := purgeAllProjects(api)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err = purgeProject(api, *project)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func purgeAllProjects(api *cloudflare.API) error {
	projects, _, err := api.ListPagesProjects(context.TODO(), *account, cloudflare.PaginationOptions{})
	if err != nil {
		return err
	}

	for _, p := range projects {
		err = purgeProject(api, p.Name)
		if err != nil {
			fmt.Printf("❌ Failed to cleanup project '%s': %v\n", p.Name, err)
		}
	}
	return nil
}

func purgeProject(api *cloudflare.API, name string) error {
	opts := cloudflare.ListPagesDeploymentsParams{
		ProjectName: name,
	}
	deployments, _, err := api.ListPagesDeployments(context.TODO(), cloudflare.AccountIdentifier(*account), opts)
	if err != nil {
		return err
	}

	for _, d := range deployments {
		err = api.DeletePagesDeployment(context.TODO(), cloudflare.AccountIdentifier(*account), *project, d.ID)
		if err != nil {
			fmt.Printf("❌ Failed to delete deployment id=%s project=%s error=%v\n", d.ID, d.ProjectName, err)
			continue
		}
		fmt.Printf("🧹 Deleted deployment id=%s project=%s\n", d.ID, d.ProjectName)
	}
	return nil
}
