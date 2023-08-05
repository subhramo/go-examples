package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

const (
	GithubAPIURL = "https://api.github.com/repos/subhramo/go-examples/pulls"
)

type PullRequest struct {
	Title  string `json:"title"`
	Body   string `json:"body"`
	Head   string `json:"head"`
	Base   string `json:"base"`
	Email  string `json:"email"`
	JiraID string `json:"jiraID"`
	Topic  string `json:"topic"`
}

func main() {
	args := os.Args
	if len(args) < 4 {
		fmt.Println("Usage: go run main.go <email> <jiraNumber> [<snsTopic>]")
		return
	}

	//email := args[1]
	//jiraNumber := args[2]
	//var snsTopic string
	//if len(args) >= 4 {
	//		snsTopic = args[3]
	//	}

	router := gin.Default()

	router.POST("/createPR/:owner/:repo", createPullRequestHandler)
	router.GET("/getPR/:owner/:repo", getPullRequestHandler)

	router.Run(":8080")
}

func createPullRequestHandler(c *gin.Context) {
	owner := c.Param("owner")
	repo := c.Param("repo")

	var pull PullRequest
	if err := c.BindJSON(&pull); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}

	// If base branch is not provided, set it to "main" by default
	if pull.Base == "" {
		pull.Base = "main"
	}

	url := fmt.Sprintf("%s/repos/%s/%s/pulls", GithubAPIURL, owner, repo)
	resp, err := createPullRequest(url, &pull)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"response": resp})
}

func getPullRequestHandler(c *gin.Context) {
	owner := c.Param("owner")
	repo := c.Param("repo")

	url := fmt.Sprintf("%s/repos/%s/%s/pulls", GithubAPIURL, owner, repo)
	resp, err := getPullRequest(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"response": resp})
}

func createPullRequest(url string, pull *PullRequest) (string, error) {
	token := os.Getenv("GITHUB_TOKEN")

	data, err := json.Marshal(pull)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}

	req.Header.Set(http.CanonicalHeaderKey("Accept"), "application/vnd.github.v3+json")
	req.Header.Set(http.CanonicalHeaderKey("Authorization"), "token "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func getPullRequest(url string) (string, error) {
	token := os.Getenv("GITHUB_TOKEN")

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set(http.CanonicalHeaderKey("Accept"), "application/vnd.github.v3+json")
	req.Header.Set(http.CanonicalHeaderKey("Authorization"), "token "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
