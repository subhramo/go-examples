package main

import (
	"context"
	"fmt"

	"github.com/google/go-github/v53/github"
	"golang.org/x/oauth2"
)

func main() {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "xxxxx"},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	// Define pull request parameters
	newPR := &github.NewPullRequest{
		Title:               github.String("Pull Request Title"),
		Head:                github.String("BranchName"), // e.g., "feature-branch"
		Base:                github.String("main"),       // Typically "main" or "master"
		Body:                github.String("Pull request description."),
		MaintainerCanModify: github.Bool(true),
	}

	// Create the pull request
	pr, _, err := client.PullRequests.Create(ctx, "<Owner>", "<Repo>", newPR)
	if err != nil {
		fmt.Printf("Error while creating the pull request: %v\n", err)
		return
	}

	// Print the created pull request
	fmt.Printf("Created PR: %v\n", pr)
}
