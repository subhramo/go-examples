package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	Url   = "https://api.github.com/repos/subhramo/go-examples/pulls"
	Token = "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"
)

func makeRequest(url, token, email, jiraID string) (*http.Response, error) {
	jsonData := fmt.Sprintf(`{"title":"Amazing new feature","body":"Pull request for Jira ID: %s by %s","head":"feature","base":"main"}`, jiraID, email)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer([]byte(jsonData)))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}

	return resp, nil
}

func main() {
	email := flag.String("email", "", "Email to use in the pull request")
	jiraID := flag.String("jiraID", "", "Jira ID to use in the pull request")

	flag.Parse()

	if *email == "" || *jiraID == "" {
		log.Fatalf("You must provide an email and Jira ID.")
	}

	resp, err := makeRequest(Url, Token, *email, *jiraID)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	fmt.Println("Response body: ", string(body))
}
