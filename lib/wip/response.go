package lib_wip

import "time"

type WIPAPIResponse struct {
	Data Data `json:"data"`
}
type Attachment struct {
	URL string `json:"url"`
}
type Todo struct {
	ID          string       `json:"id"`
	Body        string       `json:"body"`
	CompletedAt time.Time    `json:"completed_at"`
	Attachments []Attachment `json:"attachments"`
}
type Project struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Pitch      string `json:"pitch"`
	WebsiteURL string `json:"website_url"`
	Todos      []Todo `json:"todos"`
}
type Viewer struct {
	Projects []Project `json:"projects"`
}
type Data struct {
	Viewer Viewer `json:"viewer"`
}
