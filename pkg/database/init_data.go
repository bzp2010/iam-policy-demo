package database

import (
	ladonconditions "policydemo/pkg/ladon_conditions"
	"strconv"

	"github.com/ory/ladon"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func InitDatabase(db *gorm.DB) {
	db.Exec("DROP TABLE users; DROP TABLE user_roles; DROP TABLE user_policies; DROP TABLE roles; DROP TABLE role_policies; DROP TABLE policies;")

	// Migrate the schema
	db.AutoMigrate(&User{}, &Role{}, &Policy{})

	policies := []*Policy{
		newDBPolicy(&ladon.DefaultPolicy{
			Description: "Allow all service on all actions",
			Resources:   []string{"api:gateway:service:<.*>"},
			Actions:     []string{"gateway:<.*>"},
			Effect:      ladon.AllowAccess,
		}, 1),
		newDBPolicy(&ladon.DefaultPolicy{
			Description: "Allow all service update basic",
			Resources:   []string{"api:gateway:service:<.*>"},
			Actions:     []string{"gateway:UpdateServiceBasic"},
			Effect:      ladon.AllowAccess,
		}, 2),
		newDBPolicy(&ladon.DefaultPolicy{
			Description: "Allow service A&B&3 on all actions",
			Resources:   []string{"api:gateway:service:a", "api:gateway:service:b", "api:gateway:service:3"},
			Actions:     []string{"gateway:<.*>"},
			Effect:      ladon.AllowAccess,
		}, 3),
		newDBPolicy(&ladon.DefaultPolicy{
			Description: "Allow service A on update basic",
			Resources:   []string{"api:gateway:service:<.*>"},
			Actions:     []string{"gateway:UpdateServiceBasic"},
			Effect:      ladon.AllowAccess,
		}, 4),
		newDBPolicy(&ladon.DefaultPolicy{
			Description: "Allow all resources and all actions",
			Resources:   []string{"api:gateway:<.*>"},
			Actions:     []string{"gateway:<.*>"},
			Effect:      ladon.AllowAccess,
		}, 5),
		newDBPolicy(&ladon.DefaultPolicy{
			Description: "Allow all resources and get actions",
			Resources:   []string{"api:gateway:<.*>"},
			Actions:     []string{"gateway:Get<.*>"},
			Effect:      ladon.AllowAccess,
		}, 6),
		newDBPolicy(&ladon.DefaultPolicy{
			Description: "Allow all resources and get actions BUT need token auth",
			Resources:   []string{"api:gateway:<.*>"},
			Actions:     []string{"gateway:Get<.*>"},
			Effect:      ladon.AllowAccess,
			Conditions: ladon.Conditions{
				"needToken": &ladon.BooleanCondition{
					BooleanValue: true,
				},
			},
		}, 7),
		newDBPolicy(&ladon.DefaultPolicy{
			Description: "Allow service 8 and routes belong to service 8",
			Resources:   []string{"api:gateway:service:8", "api:gateway:route:<.*>"},
			Actions:     []string{"gateway:<.*>"},
			Effect:      ladon.AllowAccess,
			Conditions: ladon.Conditions{
				"belongTo": &ladon.StringEqualCondition{
					Equals: "api:gateway:service:8",
				},
			},
		}, 8),
		newDBPolicy(&ladon.DefaultPolicy{
			Description: "Allow services by multiple labels",
			Resources:   []string{"api:gateway:service:<.*>"},
			Actions:     []string{"gateway:<.*>"},
			Effect:      ladon.AllowAccess,
			Conditions: ladon.Conditions{
				"api:labels": &ladonconditions.LabelsContainCondition{
					Labels: map[string]string{
						"key:a": "value:a",
						"key:b": "value:b",
					},
				},
			},
		}, 9),
		newDBPolicy(&ladon.DefaultPolicy{
			Description: "Allow services by multiple labels",
			Resources:   []string{"api:gateway:service:<.*>"},
			Actions:     []string{"gateway:<.*>"},
			Effect:      ladon.AllowAccess,
			Conditions: ladon.Conditions{
				"api:labels": &ladonconditions.LabelsContainCondition{
					Labels: map[string]string{
						"key:c": "value:c",
					},
				},
			},
		}, 10),
	}

	db.Create(policies)

	roles := []*Role{
		{
			Model: gorm.Model{
				ID: 1,
			},
			Name: "SuperAdmin",
			Policies: []*Policy{
				{
					Model: gorm.Model{
						ID: 5,
					},
				},
			},
		},
		{
			Model: gorm.Model{
				ID: 2,
			},
			Name: "Viewer",
			Policies: []*Policy{
				{
					Model: gorm.Model{
						ID: 6,
					},
				},
			},
		},
	}
	db.Create(roles)

	users := []*User{
		{Username: "superadmin", Roles: []*Role{
			{
				Model: gorm.Model{
					ID: 1, // role 1
				},
			},
		}},
		{Username: "alice", Policies: []*Policy{
			{
				Model: gorm.Model{
					ID: 3, // policy 3
				},
			},
		}},
		{Username: "forcetoken", Policies: []*Policy{
			{
				Model: gorm.Model{
					ID: 7, // policy 7
				},
			},
		}},
		{Username: "limitroutes", Policies: []*Policy{
			{
				Model: gorm.Model{
					ID: 8, // policy 8
				},
			},
		}},
		{Username: "multiplelabels", Policies: []*Policy{
			{
				Model: gorm.Model{
					ID: 9, // policy 9
				},
			},
		}},
		{Username: "singlelabel", Policies: []*Policy{
			{
				Model: gorm.Model{
					ID: 10, // policy 10
				},
			},
		}},
	}

	db.Create(users)
}

func newDBPolicy(p *ladon.DefaultPolicy, id uint) *Policy {
	p.ID = strconv.Itoa(int(id))
	return &Policy{
		Model:  gorm.Model{ID: id},
		Policy: datatypes.NewJSONType(p),
	}
}
