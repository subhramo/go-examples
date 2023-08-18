package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	jira "github.com/andygrunwald/go-jira"
)

type CommentRequest struct {
	IssueID     string `json:"issue_id"`
	CommentBody string `json:"comment_body"`
}

func getJiraClient() *jira.Client {
	tp := jira.BasicAuthTransport{
		Username: "yourUsername",
		Password: "yourPassword",
	}

	client, err := jira.NewClient(tp.Client(), "https://your.jira.instance.url/")
	if err != nil {
		log.Fatalf("Error initializing Jira client: %v", err)
	}
	return client
}

func AddCommentToJiraTicket(jiraClient *jira.Client, issueID, commentBody string) error {
	comment := &jira.Comment{
		Body: commentBody,
	}

	_, _, err := jiraClient.Issue.AddComment(issueID, comment)
	return err
}

func handleComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is supported", http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	var req CommentRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, "Error unmarshaling JSON", http.StatusBadRequest)
		return
	}

	jiraClient := getJiraClient()
	err = AddCommentToJiraTicket(jiraClient, req.IssueID, req.CommentBody)
	if err != nil {
		http.Error(w, "Error adding comment to Jira", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Successfully added comment to issue: %s", req.IssueID)
}

func main() {
	http.HandleFunc("/add-comment", handleComment)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
