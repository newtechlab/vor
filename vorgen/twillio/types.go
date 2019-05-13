package twillio

type Project struct {
	Description string  `json:"description"`
	States      []State `json:"states"`
}

type State struct {
	Type        string       `json:"type"`
	Name        string       `json:"name"`
	Sid         string       `json:"sid,omitempty"`
	Properties  Properties   `json:"properties"`
	Transitions []Transition `json:"transitions"`
}

type Transition struct {
	Event      Event       `json:"event"`
	Conditions []Condition `json:"conditions"`
	Next       *string     `json:"next"`
	UUID       string      `json:"uuid"`
	WidgetID   string      `json:"widgetId"`
}

type Condition map[string]interface{}

type Event string

type Properties map[string]interface{}

func (p *Project) Add(s State) {
	if p.States == nil {
		p.States = []State{s}
	} else {
		p.States = append(p.States, s)
	}
}
