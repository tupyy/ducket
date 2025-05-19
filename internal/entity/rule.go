package entity

type Rule struct {
	Name    string   `yaml:"name"`
	Pattern string   `yaml:"pattern"`
	Tags    []string `yaml:"tags"`
}
