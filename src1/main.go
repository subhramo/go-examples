package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	Url   = "https://api.github.com/repos/subhramo/go-examples/pulls"
	Token = "xxxxxx"
)

type PullRequest struct {
	Email    string `json:"email" binding:"required"`
	JiraID   string `json:"jiraID" binding:"required"`
	SnsTopic string `json:"snsTopic" binding:"required"`
}

func makeRequest(method, url, token, email, jiraID string) (*http.Response, error) {
	var req *http.Request
	var err error

	if method == http.MethodPost {
		jsonData := fmt.Sprintf(`{"title":"Amazing new feature","body":"Pull request for Jira ID: %s by %s","head":"feature","base":"main"}`, jiraID, email)
		req, err = http.NewRequest(method, url, bytes.NewBuffer([]byte(jsonData)))
	} else {
		req, err = http.NewRequest(method, url, nil)
	}

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

func createPR(c *gin.Context) {
	var pr PullRequest
	if err := c.ShouldBindJSON(&pr); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := makeRequest(http.MethodPost, Url, Token, pr.Email, pr.JiraID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"response": string(body)})
}

func getPR(c *gin.Context) {
	resp, err := makeRequest(http.MethodGet, Url, Token, "", "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	var pullRequests []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&pullRequests); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pullRequests)
}

func main() {
	router := gin.Default()

	router.POST("/createPR", createPR)
	router.GET("/getPR", getPR)

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
