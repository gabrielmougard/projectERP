package models

import (
	"time"
	"github.com/astaxie/beego/orm"
)

// StockMoveOperationLink (Link inventory transfers to packaging operations)
type StockMoveOperationLink struct {
	ID         int64      `orm:"column(id);pk;auto" json:"id"`         
	CreateUser *User      `orm:"rel(fk);null" json:"-"`                
	UpdateUser *User      `orm:"rel(fk);null" json:"-"`                
	CreateDate time.Time  `orm:"auto_now_add;type(datetime)" json:"-"` 
	UpdateDate time.Time  `orm:"auto_now;type(datetime)" json:"-"`     
	Move       *StockMove `orm:"rel(fk)"`
}

func init() {
	orm.RegisterModel(new(StockMoveOperationLink))
}