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

func AddCommentToJiraTicket(client *http.Client, jiraURL, issueID, comment string) error {
	// Construct the comment payload
	payload := JiraComment{
		Body: comment,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// Construct the request to add a comment
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/rest/api/2/issue/%s/comment", jiraURL, issueID), bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	// Send the request
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

// Sample usage:
// Assuming you have an authenticated http.Client from another function
// client := authenticateJira()  // This function should authenticate and return an http.Client
// err := AddCommentToJiraTicket(client, "https://your.jira.domain", "ISSUE-123", "This is a new comment")
// if err != nil {
//     fmt.Println("Error adding comment:", err)
// }
