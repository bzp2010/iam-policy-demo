package database

import (
	"github.com/ory/ladon"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	Username string

	Roles    []*Role   `gorm:"many2many:user_roles"`
	Policies []*Policy `gorm:"many2many:user_policies"`
}

// AKA User Group
type Role struct {
	gorm.Model

	Name string

	Users    []*User   `gorm:"many2many:user_roles"`
	Policies []*Policy `gorm:"many2many:role_policies"`
}

type Policy struct {
	gorm.Model

	Policy datatypes.JSONType[*ladon.DefaultPolicy]

	Users []*User `gorm:"many2many:user_policies"`
	Roles []*Role `gorm:"many2many:role_policies"`
}
