package entity

type Rule struct {
	Name    string            `yaml:"name"`
	Pattern string            `yaml:"pattern"`
	Labels  map[string]string `yaml:"labels"`
}
