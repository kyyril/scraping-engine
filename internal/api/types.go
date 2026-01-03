package api

type CreateJobRequest struct {
	URL        string          `json:"url" example:"https://example.com"`
	Actions    []ActionRequest `json:"actions,omitempty"`
	UserAgent  string          `json:"user_agent,omitempty" example:"Mozilla/5.0..."`
	Timeout    int             `json:"timeout,omitempty" example:"30"`
	MaxRetries int             `json:"max_retries,omitempty" example:"3"`
}

type ActionRequest struct {
	Type    string                 `json:"type" example:"navigate" enums:"navigate,click,type,wait,screenshot,extract,scroll"`
	Target  string                 `json:"target,omitempty" example:"#search-input"`
	Value   string                 `json:"value,omitempty" example:"search query"`
	Options map[string]interface{} `json:"options,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type QueueStatusResponse struct {
	QueueSize int  `json:"queue_size"`
	IsRunning bool `json:"is_running"`
}