package commands

import "time"

// UserStory represents a user story for project planning
type UserStory struct {
	Role       string `json:"role"`
	Goal       string `json:"goal"`
	Reason     string `json:"reason"`
	Acceptance string `json:"acceptance"`
	Priority   string `json:"priority"`
	Complexity string `json:"complexity"`
}

// ProjectIntent represents the intent and requirements for a project
type ProjectIntent struct {
	Description       string            `json:"description"`
	Framework         string            `json:"framework"`
	Language          string            `json:"language"`
	ProjectType       string            `json:"project_type"`
	CoreFeatures      []string          `json:"core_features"`
	TechnicalNeeds    []string          `json:"technical_needs"`
	UserStories       []UserStory       `json:"user_stories"`
	ArchitectureHints map[string]string `json:"architecture_hints"`
	CreatedAt         time.Time         `json:"created_at"`
}
