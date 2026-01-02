package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type JobStatus string

const (
	StatusPending    JobStatus = "pending"
	StatusProcessing JobStatus = "processing"
	StatusCompleted  JobStatus = "completed"
	StatusFailed     JobStatus = "failed"
)

type ActionType string

const (
	ActionNavigate    ActionType = "navigate"
	ActionClick       ActionType = "click"
	ActionType_       ActionType = "type"
	ActionWait        ActionType = "wait"
	ActionScreenshot  ActionType = "screenshot"
	ActionExtract     ActionType = "extract"
	ActionScroll      ActionType = "scroll"
)

type ScrapingJob struct {
	ID          uuid.UUID   `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	URL         string      `gorm:"not null" json:"url"`
	Actions     []JobAction `gorm:"foreignKey:JobID;constraint:OnDelete:CASCADE" json:"actions"`
	Status      JobStatus   `gorm:"default:'pending'" json:"status"`
	UserAgent   string      `json:"user_agent,omitempty"`
	Timeout     int         `gorm:"default:30" json:"timeout"` // in seconds
	Retries     int         `gorm:"default:3" json:"retries"`
	MaxRetries  int         `gorm:"default:3" json:"max_retries"`
	Error       string      `json:"error,omitempty"`
	Result      *ScrapingResult `gorm:"foreignKey:JobID;constraint:OnDelete:CASCADE" json:"result,omitempty"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	CompletedAt *time.Time  `json:"completed_at,omitempty"`
}

type JobAction struct {
	ID        uuid.UUID               `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	JobID     uuid.UUID               `gorm:"type:uuid;not null" json:"job_id"`
	Type      ActionType              `gorm:"not null" json:"type"`
	Target    string                  `json:"target,omitempty"`    // CSS selector, URL, or text
	Value     string                  `json:"value,omitempty"`     // Text to type, wait duration, etc.
	Options   map[string]interface{}  `gorm:"type:jsonb" json:"options,omitempty"`
	Order     int                     `gorm:"not null" json:"order"`
	CreatedAt time.Time               `json:"created_at"`
}

type ScrapingResult struct {
	ID          uuid.UUID               `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	JobID       uuid.UUID               `gorm:"type:uuid;not null;unique" json:"job_id"`
	Data        map[string]interface{}  `gorm:"type:jsonb" json:"data"`
	Screenshots []string                `gorm:"type:text[]" json:"screenshots,omitempty"`
	Metadata    map[string]interface{}  `gorm:"type:jsonb" json:"metadata"`
	CreatedAt   time.Time               `json:"created_at"`
}

func (s *ScrapingJob) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}