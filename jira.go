package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type DumbWorklogAppender struct {
	username string
	token    string
	endpoint string

	debug bool
}

type JiraDate time.Time

func (j JiraDate) MarshalJSON() ([]byte, error) {
	return []byte(time.Time(j).Format("\"2006-01-02T15:04:05.000-0700\"")), nil
}

type Worklog struct {
	Started   JiraDate `json:"started"`
	TimeSpent string   `json:"timeSpent"`
	Comment   string   `json:"comment"`
}

func NewDumpWorklogAppender(jiraEndpont, username, token string) (*DumbWorklogAppender, error) {
	return &DumbWorklogAppender{
		username: username,
		token:    token,
		endpoint: jiraEndpont,
	}, nil
}

func (d *DumbWorklogAppender) Append(ticketId string, wlog Worklog) error {
	worklogUrl := fmt.Sprintf("%s/rest/api/2/issue/%s/worklog", d.endpoint, ticketId)
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(wlog); err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, worklogUrl, &buf)
	if err != nil {
		return err
	}
	req.Header.Add("content-type", "application/json")
	req.SetBasicAuth(d.username, d.token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		errBody, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		return fmt.Errorf("failed to append worklog. Received %d: %s", resp.StatusCode, string(errBody))
	}
	return nil
}
