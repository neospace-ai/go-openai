package openai

const DEFAULT_EXPERTISE = "discovering"

type TaskGuard struct {
	Safe            bool     `json:"safe"`
	GuardReasoning  string   `json:"guard_reasoning"`
	GuardCategories []string `json:"guard_categories"`
}

type PotentialExpertise struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type TaskSelectExpertises struct {
	SearchQuery           string               `json:"search_query"`
	PotentialExpertises   []PotentialExpertise `json:"potential_expertises"`
	ChosenExpertises      []string             `json:"chosen_expertises"`
	SelectExpertiseAnswer string               `json:"select_expertise_answer"`
}

type TaskResultCollection struct {
	RawResponse         string                `json:"raw_response"`
	TaskGuard           *TaskGuard            `json:"task_guard,omitempty"`
	TaskSelectExpertise *TaskSelectExpertises `json:"task_select_expertise,omitempty"`
}

type QuerySupervisorTask struct {
	RawResponse      string   `json:"raw_response"`
	SearchQuery      []string `json:"search_query"`
	Category         string   `json:"category"`
	Description      string   `json:"description"`
	AvailableScores  string   `json:"available_scores"`
	Reasoning        string   `json:"reasoning"`
	Score            string   `json:"score"`
	Feedback         string   `json:"feedback"`
	SuggestedQueries []string `json:"suggested_queries"`
}
