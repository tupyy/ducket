package reader

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/tupyy/finance/internal/entity"
	"go.uber.org/zap"
	"sigs.k8s.io/yaml"
)

type rules struct {
	Include []string      `yaml:"include"`
	Rules   []entity.Rule `yaml:"rules"`
}

func ReadRules(filepath string) ([]entity.Rule, error) {
	r, err := os.Open(filepath)
	if err != nil {
		return []entity.Rule{}, err
	}

	content, err := ioutil.ReadAll(r)
	if err != nil {
		return []entity.Rule{}, err
	}

	var rrules rules
	if err := yaml.Unmarshal(content, &rrules); err != nil {
		return []entity.Rule{}, err
	}

	rules := rrules.Rules

	for _, i := range rrules.Include {
		filepath := path.Join(path.Dir(filepath), fmt.Sprintf("%s.yaml", i))
		r, err := ReadRules(filepath)
		if err != nil {
			zap.S().Errorf("unable to read rule file %q: %w", filepath, err)
			continue
		}
		rules = append(rules, r...)
	}

	return rules, nil
}
