package adaptivecard

import (
	"encoding/json"
	"fmt"
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
	Type      string `json:"type"`
	Text      string `json:"text"`
	Weight    string `json:"weight,omitempty"`
	Size      string `json:"size,omitempty"`
	Wrap      bool   `json:"wrap,omitempty"`
	Separator bool   `json:"separator,omitempty"`
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

func (t *TextBlock) WithWeight(weight string) {
	t.Weight = weight
}

func (t *TextBlock) WithSize(size string) {
	t.Size = size
}

func (t *TextBlock) WithSeparator() {
	t.Separator = true
}

// ----------------------
// Container
// ----------------------
type Container struct {
	Type      string    `json:"type"`
	Separator bool      `json:"separator"`
	Items     []Element `json:"items"`
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

func (t *Container) WithSeparator() {
	t.Separator = true
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
// Table
// ----------------------
type Table struct {
	Type              string     `json:"type"`
	Columns           []TableCol `json:"columns"`
	Rows              []TableRow `json:"rows"`
	FirstRowAsHeaders bool       `json:"firstRowAsHeaders"`
}

type TableCol struct {
	Width string `json:"width"`
}

type TableRow struct {
	Type  string      `json:"type"`
	Cells []TableCell `json:"cells"`
}

type TableCell struct {
	Type  string    `json:"type"`
	Items []Element `json:"items"`
}

func NewTable() Table {
	return Table{
		Type:              "Table",
		FirstRowAsHeaders: true,
		Columns:           []TableCol{},
		Rows:              []TableRow{},
	}
}
func NewTableCell(items ...Element) TableCell {
	return TableCell{
		Type:  "TableCell",
		Items: items,
	}
}
func (Table) isElement() {}
func (t Table) toRaw() any {
	// Convert rows and cells recursively
	rows := make([]any, len(t.Rows))
	for i, r := range t.Rows {
		rows[i] = r.toRaw()
	}
	return struct {
		Type    string     `json:"type"`
		Columns []TableCol `json:"columns"`
		Rows    []any      `json:"rows"`
	}{
		Type:    t.Type,
		Columns: t.Columns,
		Rows:    rows,
	}
}

func (tr TableRow) toRaw() any {
	cells := make([]any, len(tr.Cells))
	for i, c := range tr.Cells {
		cells[i] = c.toRaw()
	}
	return struct {
		Type  string `json:"type"`
		Cells []any  `json:"cells"`
	}{
		Type:  tr.Type,
		Cells: cells,
	}
}

func (tc TableCell) toRaw() any {
	items := make([]any, len(tc.Items))
	for i, el := range tc.Items {
		items[i] = el.toRaw()
	}
	return struct {
		Type  string `json:"type"`
		Items []any  `json:"items"`
	}{
		Type:  tc.Type,
		Items: items,
	}
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
	Type      string  `json:"type"`
	Text      string  `json:"text"`
	Mentioned Mention `json:"mentioned"`
}

type Mention struct {
	ID   string `json:"id"`
	Name string `json:"name"`
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

func (t *Table) AddColumn(width string) {
	t.Columns = append(t.Columns, TableCol{Width: width})
}

func (t *Table) AddRow(cells ...TableCell) {
	t.Rows = append(t.Rows, TableRow{Type: "TableRow", Cells: cells})
}

func (c *AdaptiveCard) AddMentionsMap(textPrefix string, mentions []string) {
	if c.MSTeams == nil {
		c.MSTeams = &MSTeamsInfo{
			Entities: []MSTeamsEntity{},
		}
	}

	// Build the text with all <at> placeholders
	text := textPrefix
	for _, displayName := range mentions {
		text += fmt.Sprintf(" <at>%s</at>", displayName)
	}

	c.AddBody(NewTextBlock(text))

	// Add entities
	entity := MSTeamsEntity{
		Type: "mention",
		Text: "@Team",
		Mentioned: Mention{
			ID:   "19:general@thread.tacv2",
			Name: "Team",
		},
	}
	c.MSTeams.Entities = append(c.MSTeams.Entities, entity)
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
