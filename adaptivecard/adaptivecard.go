package adaptivecard

import (
	"encoding/json"
)

// AdaptiveCard root
type AdaptiveCard struct {
	Type    string       `json:"type"`
	Version string       `json:"version"`
	Body    []Element    `json:"body"`
	Schema  string       `json:"$schema"`
	Actions []Action     `json:"actions,omitempty"`
	MSTeams *MSTeamsInfo `json:"msteams,omitempty"`
}

// --- ELEMENT INTERFACE ---
type Element interface {
	isElement()
	toRaw() any
}

// ----------------------
// TextBlock
// ----------------------
type TextBlock struct {
	Type   string `json:"type"`
	Text   string `json:"text"`
	Weight string `json:"weight,omitempty"`
	Size   string `json:"size,omitempty"`
	Wrap   bool   `json:"wrap,omitempty"`
}

func NewTextBlock(text string) TextBlock {
	return TextBlock{
		Type: "TextBlock",
		Text: text,
		Wrap: true,
	}
}
func (TextBlock) isElement() {}
func (t TextBlock) toRaw() any {
	return t
}

// ----------------------
// Container
// ----------------------
type Container struct {
	Type  string    `json:"type"`
	Items []Element `json:"items"`
}

func NewContainer(items ...Element) Container {
	return Container{
		Type:  "Container",
		Items: items,
	}
}
func (Container) isElement() {}
func (c Container) toRaw() any {
	// recursively flatten inner elements
	items := make([]any, len(c.Items))
	for i, el := range c.Items {
		items[i] = el.toRaw()
	}
	return struct {
		Type  string `json:"type"`
		Items []any  `json:"items"`
	}{
		Type:  "Container",
		Items: items,
	}
}

// ----------------------
// FactSet
// ----------------------
type FactSet struct {
	Type  string `json:"type"`
	Facts []Fact `json:"facts"`
}
type Fact struct {
	Title string `json:"title"`
	Value string `json:"value"`
}

func NewFactSet(facts ...Fact) FactSet {
	return FactSet{
		Type:  "FactSet",
		Facts: facts,
	}
}
func (FactSet) isElement() {}
func (fs FactSet) toRaw() any {
	return fs
}

// ----------------------
// Action
// ----------------------
type Action struct {
	Type  string `json:"type"`
	Title string `json:"title"`
	Url   string `json:"url,omitempty"`
}

// ----------------------
// MSTeams
// ----------------------
type MSTeamsInfo struct {
	Entities []MSTeamsEntity `json:"entities"`
}
type MSTeamsEntity struct {
	Type string `json:"type"`
	Text string `json:"text"`
	ID   string `json:"id"`
}

// ----------------------
// Helpers (Receiver Functions)
// ----------------------
func (c *AdaptiveCard) AddBody(el Element) {
	c.Body = append(c.Body, el)
}

func (c *AdaptiveCard) AddAction(action Action) {
	c.Actions = append(c.Actions, action)
}

func (c *Container) AddItem(el Element) {
	c.Items = append(c.Items, el)
}

// ----------------------
// MarshalJSON for AdaptiveCard
// ----------------------
func (c AdaptiveCard) MarshalJSON() ([]byte, error) {
	body := make([]any, len(c.Body))
	for i, el := range c.Body {
		body[i] = el.toRaw()
	}

	// build a raw struct to marshal
	raw := struct {
		Type    string       `json:"type"`
		Version string       `json:"version"`
		Body    []any        `json:"body"`
		Schema  string       `json:"$schema"`
		Actions []Action     `json:"actions,omitempty"`
		MSTeams *MSTeamsInfo `json:"msteams,omitempty"`
	}{
		Type:    c.Type,
		Version: c.Version,
		Body:    body,
		Schema:  c.Schema,
		Actions: c.Actions,
		MSTeams: c.MSTeams,
	}

	return json.Marshal(raw)
}
