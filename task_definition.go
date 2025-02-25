package openai

type Task struct {
	Name              string            `json:"name" bson:"name"`
	Description       string            `json:"description" bson:"description"`
	Type              string            `json:"type" bson:"type"`
	Components        []TaskComponent   `json:"components" bson:"components"` // ComponentName: ComponentDetails
	SupervisorProfile SupervisorProfile `json:"supervisor_profile" bson:"supervisor_profile"`
}

type TaskComponent struct {
	Name        string   `json:"name" bson:"name"`
	Description string   `json:"description" bson:"description"`
	Examples    []string `json:"examples" bson:"examples"`
	Type        string   `json:"type" bson:"type"`
}

type SupervisorProfile struct {
	Description string               `json:"description" bson:"description"`
	Components  SupervisorComponents `json:"components" bson:"components"`
}

type SupervisorComponents []SupervisorComponent // ComponentName: ComponentDetails

type SupervisorComponent struct {
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Scores      SupervisorScores `json:"scores"`
	Type        string           `json:"type"`
}

type SupervisorScores map[string]SupervisorScore // SpecialToken: ScoreDetails

type SupervisorScore struct {
	Description string  `json:"description" bson:"description"`
	Label       string  `json:"label" bson:"label"`
	Perfect     bool    `json:"perfect" bson:"perfect"`
	Weight      float64 `json:"weight" bson:"weight"`
}
