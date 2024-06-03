package ladonhelper

import (
	"context"
	"fmt"
	"policydemo/pkg/database"
	ladonconditions "policydemo/pkg/ladon_conditions"
	"slices"

	"github.com/ory/ladon"
	manager "github.com/ory/ladon/manager/memory"
	"gorm.io/gorm"
)

type CachedDBManager struct {
	*manager.MemoryManager
}

func NewCachedDBManager(db *gorm.DB) (*CachedDBManager, error) {
	ladon.ConditionFactories[new(ladonconditions.LabelsContainCondition).GetName()] = func() ladon.Condition {
		return new(ladonconditions.LabelsContainCondition)
	}

	m := &CachedDBManager{
		MemoryManager: manager.NewMemoryManager(),
	}

	m.RebuildCache(db)

	return m, nil
}

func (m *CachedDBManager) RebuildCache(db *gorm.DB) error {
	fmt.Println("Rebuilding cache")
	m.Policies = make(map[string]ladon.Policy)

	var dbPolicise []database.Policy
	err := db.
		Model(&database.Policy{}).
		Preload("Users").
		Preload("Roles").
		Preload("Roles.Users").
		Find(&dbPolicise).Error
	if err != nil {
		return err
	}

	ctx := context.Background()
	for _, dbPolicy := range dbPolicise {
		var subjects []string
		for _, user := range dbPolicy.Users {
			subjects = append(subjects, user.Username)
		}
		for _, role := range dbPolicy.Roles {
			for _, user := range role.Users {
				subjects = append(subjects, user.Username)
			}
		}
		slices.Sort(subjects) // remove duplicates

		policy := dbPolicy.Policy.Data()
		policy.Subjects = subjects
		m.MemoryManager.Create(ctx, policy)
	}

	return nil
}
