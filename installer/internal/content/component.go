package content

import "strings"

type Kind string

const (
	Commands Kind = "commands"
	Rules    Kind = "rules"
	Skills   Kind = "skills"
	Agents   Kind = "agents"
	Hooks    Kind = "hooks"
)

type Spec struct {
	ID          Kind
	DisplayName string
	SourcePath  string
	DestPath    string
	FileType    string
}

var Registry = []Spec{
	{ID: Commands, DisplayName: "Commands", SourcePath: "commands", DestPath: "commands", FileType: "markdown"},
	{ID: Rules, DisplayName: "Rules", SourcePath: "rules", DestPath: "rules", FileType: "markdown"},
	{ID: Skills, DisplayName: "Skills", SourcePath: "skills", DestPath: "skills", FileType: "any"},
	{ID: Agents, DisplayName: "Agents", SourcePath: "agents", DestPath: "agents", FileType: "any"},
	{ID: Hooks, DisplayName: "Hooks", SourcePath: "hooks", DestPath: "hooks", FileType: "any"},
}

var byID = func() map[Kind]Spec {
	m := make(map[Kind]Spec, len(Registry))
	for _, s := range Registry {
		m[s.ID] = s
	}
	return m
}()

func Lookup(id string) (Spec, bool) {
	s, ok := byID[Kind(id)]
	return s, ok
}

func Known(id string) bool {
	_, ok := Lookup(id)
	return ok
}

func DisplayName(id string) string {
	if s, ok := Lookup(id); ok {
		return s.DisplayName
	}
	if id == "" {
		return id
	}
	return strings.ToUpper(id[:1]) + id[1:]
}

func ComponentIDs() []string {
	out := make([]string, len(Registry))
	for i, s := range Registry {
		out[i] = string(s.ID)
	}
	return out
}
