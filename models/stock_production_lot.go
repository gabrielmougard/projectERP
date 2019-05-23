package models

import (
	"time"
	"github.com/astaxie/beego/orm"
)

// StockProductionLot 
type StockProductionLot struct {
	ID         int64     `orm:"column(id);pk;auto" json:"id"`         
	CreateUser *User     `orm:"rel(fk);null" json:"-"`                
	UpdateUser *User     `orm:"rel(fk);null" json:"-"`               
	CreateDate time.Time `orm:"auto_now_add;type(datetime)" json:"-"` 
	UpdateDate time.Time `orm:"auto_now;type(datetime)" json:"-"`     
}

func init() {
	orm.RegisterModel(new(StockProductionLot))
}