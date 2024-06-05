package lib_wip

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
	ctx        context.Context
}

func NewClient(apiKey string) *Client {
	return &Client{
		baseURL:    "https://api.wip.co/v1",
		apiKey:     apiKey,
		httpClient: &http.Client{},
		ctx:        context.Background(),
	}
}

func (c *Client) do(req *http.Request) (*http.Response, error) {
	return c.httpClient.Do(req.WithContext(c.ctx))
}

type Todo struct {
	ID          string       `json:"id"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	Body        string       `json:"body"`
	URL         string       `json:"url"`
	Attachments []Attachment `json:"attachments"`
	UserID      string       `json:"user_id"`
}

type Attachment struct {
	URL string `json:"url"`
}

type Project struct {
	ID          string      `json:"id"`
	Slug        string      `json:"slug"`
	Name        string      `json:"name"`
	Pitch       string      `json:"pitch"`
	Description string      `json:"description"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	Hashtag     string      `json:"hashtag"`
	WebsiteURL  string      `json:"website_url"`
	Protected   bool        `json:"protected"`
	Archived    bool        `json:"archived"`
	URL         string      `json:"url"`
	Logo        interface{} `json:"logo"`
	Owner       User        `json:"owner"`
	Makers      []User      `json:"makers"`
}

type User struct {
	ID         string      `json:"id"`
	Username   string      `json:"username"`
	Streak     int         `json:"streak"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
	Protected  bool        `json:"protected"`
	FirstName  string      `json:"first_name"`
	LastName   string      `json:"last_name"`
	TodosCount int         `json:"todos_count"`
	TimeZone   string      `json:"time_zone"`
	URL        string      `json:"url"`
	Avatar     interface{} `json:"avatar"`
	BestStreak int         `json:"best_streak"`
	Streaking  bool        `json:"streaking"`
}

type PaginatedTodos struct {
	Data       []Todo `json:"data"`
	HasMore    bool   `json:"has_more"`
	TotalCount int    `json:"total_count"`
}

type PaginatedProjects struct {
	Data       []Project `json:"data"`
	HasMore    bool      `json:"has_more"`
	TotalCount int       `json:"total_count"`
}

type Upload struct {
	Filename    string            `json:"filename"`
	ByteSize    int               `json:"byte_size"`
	Checksum    string            `json:"checksum"`
	ContentType string            `json:"content_type"`
	URL         string            `json:"url"`
	Key         string            `json:"key"`
	SignedID    string            `json:"signed_id"`
	Method      string            `json:"method"`
	Headers     map[string]string `json:"headers"`
}

func (c *Client) get(path string, limit *int, startingAfter *string) ([]byte, error) {
	req, err := http.NewRequest("GET", c.baseURL+path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	q := req.URL.Query()
	q.Add("api_key", c.apiKey)

	if limit != nil {
		q.Add("limit", strconv.Itoa(*limit))
	}

	if startingAfter != nil {
		q.Add("starting_after", *startingAfter)
	}

	req.URL.RawQuery = q.Encode()

	resp, err := c.do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return body, nil
}

func (c *Client) GetMyProjects(limit *int, startingAfter *string) (*PaginatedProjects, error) {
	respBytes, err := c.get("/users/me/projects", limit, startingAfter)
	if err != nil {
		return nil, err
	}

	var projects *PaginatedProjects
	if err := json.Unmarshal(respBytes, &projects); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return projects, nil
}

func (c *Client) GetProjectTodos(projectID string, limit *int, startingAfter *string) (*PaginatedTodos, error) {
	body, err := c.get(fmt.Sprintf("/projects/%s/todos", projectID), limit, startingAfter)
	if err != nil {
		return nil, err
	}

	var todos *PaginatedTodos
	if err := json.Unmarshal(body, &todos); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return todos, nil
}
