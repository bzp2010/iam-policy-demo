package ladonconditions

import (
	"context"

	"github.com/ory/ladon"
)

type LabelsContainCondition struct {
	Labels map[string]string `json:labels`
}

func (c *LabelsContainCondition) Fulfills(ctx context.Context, v interface{}, _ *ladon.Request) bool {
	pairs, ok := v.(map[string]string)
	if !ok {
		return false
	}

	for key, value := range pairs {
		if c.Labels[key] != value {
			return false
		}
	}

	return true
}

func (c *LabelsContainCondition) GetName() string {
	return "LabelsContainCondition"
}
