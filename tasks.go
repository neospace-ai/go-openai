package openai

const DEFAULT_EXPERTISE = "discovering"

const (
	TASK_TYPE_GUARD             = "task_guard"
	TASK_TYPE_SELECT_EXPERTISES = "task_select_expertises"
)

type CategorySuggestions struct {
	MacroCategory string `json:"macro_category" bson:"macro_category"`
	SubCategory   string `json:"sub_category" bson:"sub_category"`
	Justification string `json:"justification" bson:"justification"`
}

type TaskGuard struct {
	GuardSafe           bool                `json:"guard_safe" bson:"guard_safe"`
	GuardReasoning      string              `json:"guard_reasoning" bson:"guard_reasoning"`
	GuardCategory       []string            `json:"guard_category" bson:"guard_category"`
	CategorySuggestions CategorySuggestions `json:"category_suggestions" bson:"category_suggestions"`
}

func (t *TaskGuard) ToGeneric() GenericTask {
	return GenericTask{
		TaskType: TASK_TYPE_GUARD,
		Task:     t,
	}
}

type PotentialExpertise struct {
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`
}

type TaskSelectExpertises struct {
	SearchQuery           []string             `json:"search_query" bson:"search_query"`
	PotentialExpertises   []PotentialExpertise `json:"potential_expertises" bson:"potential_expertises"`
	ChosenExpertises      []string             `json:"chosen_expertises" bson:"chosen_expertises"`
	SelectExpertiseAnswer string               `json:"select_expertise_answer" bson:"select_expertise_answer"`
}

func (t *TaskSelectExpertises) ToGeneric() GenericTask {
	return GenericTask{
		TaskType: TASK_TYPE_SELECT_EXPERTISES,
		Task:     t,
	}
}

type TaskResultCollection struct {
	RawResponse         string                `json:"raw_response" bson:"raw_response"`
	TaskGuard           *TaskGuard            `json:"task_guard,omitempty" bson:"task_guard,omitempty"`
	TaskSelectExpertise *TaskSelectExpertises `json:"task_select_expertise,omitempty" bson:"task_select_expertise,omitempty"`
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
	TaskType string `json:"task_type" bson:"task_type"`
	Task     any    `json:"task" bson:"task"`
}

type SupervisorTaskScore struct {
	Token       int    `json:"token" bson:"token"`
	TokenName   string `json:"token_name" bson:"token_name"`
	Description string `json:"description" bson:"description"`
}

type SupervisorTaskComponent struct {
	Name            string                `json:"name" bson:"name"`
	Description     string                `json:"description" bson:"description"`
	AvailableScores []SupervisorTaskScore `json:"available_scores" bson:"available_scores"`
	Chosen          *int                  `json:"chosen" bson:"chosen"`
}

type TaskSupervisor struct {
	RawResponse         string                    `json:"raw_response" bson:"raw_response"`
	Components          []SupervisorTaskComponent `json:"components" bson:"components"`
	SupervisorReasoning string                    `json:"supervisor_reasoning" bson:"supervisor_reasoning"`
	Feedback            string                    `json:"feedback" bson:"feedback"`
	Score               map[string]string         `json:"score" bson:"score"` //{ComponentName: TokenName}
}

type TaskDefinition struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	FieldTypes  map[string]string `json:"field_types"`
}

var GUARD_TASK_DEFINITION = TaskDefinition{
	Name:        "task_guard",
	Description: "A task Guard têm como objetivo analisar e assegurar a segurança de mensagens de usuários e classificá-las em categorias de risco de modo que siga princípios éticos, morais e legais.",
	FieldTypes: map[string]string{
		"guard_safe":      "bool",
		"guard_reasoning": "str",
		"guard_category":  "List[str]",
	},
}
