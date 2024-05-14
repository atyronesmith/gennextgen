package types

type KustomizeFile struct {
	Path      string                 `yaml:"path"`
	Data      Kustomize              `yaml:"data"`
	Resources map[string]interface{} `yaml:"resources"`
}

type Kustomize struct {
	APIVersion   string         `yaml:"apiVersion"`
	Kind         string         `yaml:"kind"`
	Components   []string       `yaml:"components"`
	Resources    []string       `yaml:"resources"`
	Patches      []Patches      `yaml:"patches"`
	Replacements []Replacements `yaml:"replacements"`
}
type Target struct {
	Kind          string `yaml:"kind"`
	LabelSelector string `yaml:"labelSelector"`
}
type Patches struct {
	Target Target `yaml:"target"`
	Path   string `yaml:"path"`
}
type Source struct {
	Kind      string `yaml:"kind"`
	Name      string `yaml:"name"`
	FieldPath string `yaml:"fieldPath"`
}
type Select struct {
	Kind string `yaml:"kind"`
	Name string `yaml:"name"`
}
type Targets struct {
	Select     Select   `yaml:"select"`
	FieldPaths []string `yaml:"fieldPaths"`
}
type Replacements struct {
	Source  Source    `yaml:"source"`
	Targets []Targets `yaml:"targets"`
}

func NewKustomizeFile() *KustomizeFile {
	return &KustomizeFile{
		Resources: make(map[string]interface{}),
	}
}
