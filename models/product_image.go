package models

import (
	"time"
	"github.com/astaxie/beego/orm"
)

//ProductImage 
type ProductImage struct {
	ID              int64            `orm:"column(id);pk;auto" json:"id"`      
	CreateUser      *User            `orm:"rel(fk);null" json:"-"`            
	UpdateUser      *User            `orm:"rel(fk);null" json:"-"`               
	CreateDate      time.Time        `orm:"auto_now_add;type(datetime)" json:"-"` 
	UpdateDate      time.Time        `orm:"auto_now;type(datetime)" json:"-"`     
	Name            string           `orm:"unique" form:"name"`                 
	ProductTemplate *ProductTemplate `orm:"rel(fk);null"`                        
	ProductProduct  *ProductProduct  `orm:"rel(fk);null"`                        

	FormAction   string   `orm:"-" json:"FormAction"`   
	ActionFields []string `orm:"-" json:"ActionFields"` 
}

func init() {
	orm.RegisterModel(new(ProductImage))
}