package reader

import (
	"io"
	"io/ioutil"

	"github.com/tupyy/finance/internal/entity"
	"sigs.k8s.io/yaml"
)

type rules struct {
	Rules []entity.Rule `yaml:"rules"`
}

func ReadRules(r io.Reader) ([]entity.Rule, error) {
	content, err := ioutil.ReadAll(r)
	if err != nil {
		return []entity.Rule{}, err
	}

	var rrules rules
	if err := yaml.Unmarshal(content, &rrules); err != nil {
		return []entity.Rule{}, err
	}
	return rrules.Rules, nil
}
