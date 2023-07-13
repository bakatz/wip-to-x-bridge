package lib_wip

import "time"

type WIPAPIResponse struct {
	Data Data `json:"data"`
}
type Attachments struct {
	URL string `json:"url"`
}
type Todo struct {
	ID          string        `json:"id"`
	Body        string        `json:"body"`
	CompletedAt time.Time     `json:"completed_at"`
	Attachments []Attachments `json:"attachments"`
}
type Project struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	WebsiteURL string `json:"website_url"`
	Todos      []Todo `json:"todos"`
}
type Viewer struct {
	Projects []Project `json:"projects"`
}
type Data struct {
	Viewer Viewer `json:"viewer"`
}
