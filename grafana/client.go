package grafana

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

var (
	ErrBadResp = errors.New("bad resp from grafana")
)

type Client struct {
	cl     *http.Client
	token  string
	url    string
	folder int
}

func New(token, url string, folder int) *Client {
	return &Client{
		cl:     &http.Client{Timeout: time.Second * 10},
		token:  token,
		url:    url,
		folder: folder,
	}
}

type CreateDashboardReq struct {
	Dashboard json.RawMessage `json:"dashboard"`
	FolderID  int             `json:"folderID"`
	Overwrite bool            `json:"overwrite"`
}

func (cl *Client) CreateDashboards(dashboards map[string][]byte) error {
	for name, content := range dashboards {
		if err := cl.CreateDashboard(name, content); err != nil {
			return err
		}
	}
	return nil
}

func (cl *Client) CreateDashboard(name string, content []byte) error {
	url := fmt.Sprintf("%s/api/dashboards/db", cl.url)

	payload := &CreateDashboardReq{
		Dashboard: content,
		FolderID:  cl.folder,
		Overwrite: true,
	}

	b, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal json: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(b))
	if err != nil {
		return fmt.Errorf("failed to create http request for %s: %w", name, err)
	}
	if _, err := cl.call(req); err != nil {
		return fmt.Errorf("failed to make http request to %s: %w", url, err)
	}

	return nil
}

func (cl *Client) call(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", cl.token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	resp, err := cl.cl.Do(req)
	if err != nil {
		return nil, err
	}

	statusOK := resp.StatusCode >= 200 && resp.StatusCode < 300
	if !statusOK {
		return nil, fmt.Errorf("%w: %s", ErrBadResp, resp.Status)
	}
	return resp, nil
}
