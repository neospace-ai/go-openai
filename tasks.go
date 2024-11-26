package openai

const DEFAULT_EXPERTISE = "discovering"

const (
	TASK_TYPE_GUARD             = "GUARD"
	TASK_TYPE_SELECT_EXPERTISES = "SELECT_EXPERTISES"
)

type TaskGuard struct {
	Safe            bool     `json:"safe"`
	GuardReasoning  string   `json:"guard_reasoning"`
	GuardCategories []string `json:"guard_categories"`
}

func (t *TaskGuard) ToGeneric() GenericTask {
	return GenericTask{
		TaskType: TASK_TYPE_GUARD,
		Task:     t,
	}
}

type PotentialExpertise struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type TaskSelectExpertises struct {
	SearchQuery           []string             `json:"search_query"`
	PotentialExpertises   []PotentialExpertise `json:"potential_expertises"`
	ChosenExpertises      []string             `json:"chosen_expertises"`
	SelectExpertiseAnswer string               `json:"select_expertise_answer"`
}

func (t *TaskSelectExpertises) ToGeneric() GenericTask {
	return GenericTask{
		TaskType: TASK_TYPE_SELECT_EXPERTISES,
		Task:     t,
	}
}

type TaskResultCollection struct {
	RawResponse         string                `json:"raw_response"`
	TaskGuard           *TaskGuard            `json:"task_guard,omitempty"`
	TaskSelectExpertise *TaskSelectExpertises `json:"task_select_expertise,omitempty"`
}

func (t *TaskResultCollection) ToGeneric() GenericTask {
	if t.TaskGuard != nil {
		return t.TaskGuard.ToGeneric()
	}
	if t.TaskSelectExpertise != nil {
		return t.TaskSelectExpertise.ToGeneric()
	}
	panic("TaskResultCollection must have a task")
}

// GenericTask is a generic task structure that can be used to represent any task.
type GenericTask struct {
	TaskType string `json:"task_type"`
	Task     any    `json:"task"`
}

type SupervisorTaskScore struct {
	Token       int    `json:"token"`
	Description string `json:"description"`
}

type SupervisorTaskCategory struct {
	Name            string            `json:"name"`
	Decription      string            `json:"description"`
	AvailableScores []SupervisorScore `json:"available_scores"`
	Chosen          *int              `json:"chosen"`
}

type TaskSupervisor struct {
	RawResponse string                   `json:"raw_response"`
	Categories  []SupervisorTaskCategory `json:"categories"`
	Reasoning   string                   `json:"reasoning"`
	Feedback    string                   `json:"feedback"`
}
