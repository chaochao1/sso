package models

import (
	"errors"
	"github.com/cicdi-go/sso/src/utils"
	"github.com/xormplus/xorm"
)

type ActiveRecod interface {
	GetDb() (e *xorm.Engine, err error)
	TableName() string
}

type Base struct {
}

func (u *Base) TableName() string {
	return utils.Config.TablePrefix + "user"
}

func (u *Base) GetDb() (e *xorm.Engine, err error) {
	var found bool
	if e, found = utils.Engin.GetXormEngin("default"); !found {
		err = errors.New("Database default is not found!")
	}
	return
}