package webhook

import (
	"encoding/json"
	"fmt"
	env "go_stream_api/environment"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-co-op/gocron"
)

type currentVersion int

const (
	firstHalf currentVersion = iota
	secondHalf
)

type workflows struct {
	WorkflowRuns []jobData `json:"workflow_runs"`
}

type jobData struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var (
	githubOwner = "flyingfrisbee"
	repoName    = "go_stream_api"
	version     currentVersion
	sched       *gocron.Scheduler
)

func StartWebhookService() {
	ver, err := strconv.Atoi(env.CurrentVersion)
	if err != nil {
		log.Fatal(err)
	}

	cronExpression := ""
	version = currentVersion(ver)
	switch version {
	case firstHalf:
		cronExpression = "0 0 18 * *"
	default:
		cronExpression = "0 0 1 * *"
	}

	sched = gocron.NewScheduler(time.Local)
	_, err = sched.Cron(cronExpression).Do(webhookTask)
	if err != nil {
		log.Fatal(err)
	}
	sched.StartAsync()
}

func webhookTask() {
	workflowID := getID(version)
	if workflowID == -1 {
		return
	}

	successRerun := false
	for !successRerun {
		successRerun = rerunWorkflowByID(workflowID)
	}

	// cancel context webscraper + notif + database connection
	sched.Clear()
}

func getID(version currentVersion) int {
	url := fmt.Sprintf(
		"https://api.github.com/repos/%s/%s/actions/runs",
		githubOwner,
		repoName,
	)
	authBearer := fmt.Sprintf("Bearer %s", env.GHAuthToken)

	req, err := http.NewRequest(
		"GET",
		url,
		nil,
	)
	if err != nil {
		log.Println(err)
		return -1
	}

	req.Header.Set("Authorization", authBearer)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return -1
	}
	defer resp.Body.Close()

	jsonBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return -1
	}

	result := workflows{}

	err = json.Unmarshal(jsonBytes, &result)
	if err != nil {
		log.Println(err)
		return -1
	}

	return findJobIDBasedOnVersion(&result, version)
}

func findJobIDBasedOnVersion(result *workflows, version currentVersion) int {
	keywordToFind := "1st"
	if version == firstHalf {
		keywordToFind = "2nd"
	}

	for _, job := range result.WorkflowRuns {
		if strings.Contains(job.Name, keywordToFind) {
			return job.ID
		}
	}

	log.Printf("Cannot find workflow with name consisting of %s", keywordToFind)
	return -1
}

func rerunWorkflowByID(id int) bool {
	url := fmt.Sprintf(
		"https://api.github.com/repos/%s/%s/actions/runs/%d/rerun",
		githubOwner,
		repoName,
		id,
	)
	authValue := fmt.Sprintf("Bearer %s", env.GHAuthToken)

	req, err := http.NewRequest(
		"POST",
		url,
		nil,
	)
	if err != nil {
		log.Println(err)
		return false
	}

	req.Header.Set("Authorization", authValue)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == 201
}
