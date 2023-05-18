package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go_stream_api/api"
	env "go_stream_api/environment"
	db "go_stream_api/repository/database"
	ws "go_stream_api/repository/webscraper"
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

type workflowRunsResponse struct {
	WorkflowRuns []workflowRunDetail `json:"workflow_runs"`
}

type workflowRunDetail struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type webhookData struct {
	actionKey, githubUsername, repositoryName string
}

func (wd *webhookData) recreateWorkflow() error {
	workflowName, fileName := wd.getWorkflowNameAndFileNameBasedOnVersion()

	workflowID, err := wd.getWorkflowIDByWorkflowName(workflowName)
	if err != nil {
		return err
	}

	err = wd.deleteWorkflowByID(workflowID)
	if err != nil {
		return err
	}

	err = wd.runWorkflowByFileName(fileName)
	if err != nil {
		return err
	}

	return nil
}

func (wd *webhookData) getWorkflowNameAndFileNameBasedOnVersion() (string, string) {
	var workflowName string
	var fileName string

	switch version {
	case firstHalf:
		workflowName = "2nd"
		fileName = "second_half.yml"
	case secondHalf:
		workflowName = "1st"
		fileName = "first_half.yml"
	}

	return workflowName, fileName
}

func (wd *webhookData) getWorkflowIDByWorkflowName(name string) (int64, error) {
	url := fmt.Sprintf(
		`https://api.github.com/repos/%s/%s/actions/runs`,
		wd.githubUsername, wd.repositoryName,
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}

	authHeader := fmt.Sprintf("Bearer %s", wd.actionKey)
	req.Header.Set("Authorization", authHeader)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	jsonBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var responseBody workflowRunsResponse
	err = json.Unmarshal(jsonBytes, &responseBody)
	if err != nil {
		return 0, err
	}

	for _, run := range responseBody.WorkflowRuns {
		if strings.Contains(run.Name, name) {
			return run.ID, nil
		}
	}

	return 0, fmt.Errorf("cannot find workflow with name: %s", name)
}

func (wd *webhookData) deleteWorkflowByID(id int64) error {
	url := fmt.Sprintf(
		`https://api.github.com/repos/%s/%s/actions/runs/%d`,
		wd.githubUsername, wd.repositoryName, id,
	)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	authHeader := fmt.Sprintf("Bearer %s", wd.actionKey)
	req.Header.Set("Authorization", authHeader)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 204 {
		return fmt.Errorf(
			"expecting status code: 204 for func deleteWorkflowByID, but got %d",
			resp.StatusCode,
		)
	}

	return nil
}

// FileName is the name of .yml file inside .github/workflows
func (wd *webhookData) runWorkflowByFileName(fileName string) error {
	url := fmt.Sprintf(
		`https://api.github.com/repos/%s/%s/actions/workflows/%s/dispatches`,
		wd.githubUsername, wd.repositoryName, fileName,
	)

	// "main" is the branch of repository that you'd like to execute github actions
	reqBody := []byte(`{"ref": "main"}`)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}

	authHeader := fmt.Sprintf("Bearer %s", wd.actionKey)
	req.Header.Set("Authorization", authHeader)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 204 {
		return fmt.Errorf(
			"expecting status code: 204 for func runWorkflowByFileName, but got %d",
			resp.StatusCode,
		)
	}

	return nil
}

var (
	version currentVersion
	sched   *gocron.Scheduler
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
		cronExpression = "30 0 1 * *"
	}

	sched = gocron.NewScheduler(time.Local)
	_, err = sched.Cron(cronExpression).Do(webhookTask)
	if err != nil {
		log.Fatal(err)
	}
	sched.StartAsync()
}

func webhookTask() {
	webhook := initWebhookData(
		env.GHAuthToken,
		"flyingfrisbee",
		"go_stream_api",
	)

	for {
		err := webhook.recreateWorkflow()
		if err == nil {
			break
		}

		log.Println(err)
	}

	sched.Clear()
	api.Stop()
	ws.Stop()
	db.TerminateConnectionToDB()
}

func initWebhookData(actionKey, githubUsername, repositoryName string) *webhookData {
	return &webhookData{
		actionKey:      actionKey,
		githubUsername: githubUsername,
		repositoryName: repositoryName,
	}
}
