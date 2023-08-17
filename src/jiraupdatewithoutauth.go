package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type JiraComment struct {
	Body string `json:"body"`
}

func AddCommentToJiraTicket(client *http.Client, jiraURL, issueID, commentBody string) error {
	// Prepare payload for the comment
	payload := JiraComment{
		Body: commentBody,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// Construct the POST request to add a comment
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/rest/api/2/issue/%s/comment", jiraURL, issueID), bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	// Send the request using the provided client
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to add comment to Jira ticket: %s", resp.Status)
	}

	return nil
}

// Sample usage
// You would use an authenticated http.Client from another function or setup
// client := getAuthenticatedClient() // Example function that returns an authenticated client
// err := AddCommentToJiraTicket(client, "https://your.jira.domain", "ISSUE-123", "Your comment here")
// if err != nil {
//     fmt.Println("Error:", err)
// }
