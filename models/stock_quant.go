package models

import (
	"time"
	"github.com/astaxie/beego/orm"
)

// StockQuant  	(Inventory analysis)
type StockQuant struct {
	ID                   int64               `orm:"column(id);pk;auto" json:"id"`                 
	CreateUser           *User               `orm:"rel(fk);null" json:"-"`                         
	UpdateUser           *User               `orm:"rel(fk);null" json:"-"`                         
	CreateDate           time.Time           `orm:"auto_now_add;type(datetime)" json:"-"`         
	UpdateDate           time.Time           `orm:"auto_now;type(datetime)" json:"-"`              
	Name                 string              `json:"Name"`                                        
	Product              *ProductProduct     `orm:"rel(fk)"`                                       
	Location             *StockLocation      `orm:"rel(fk);null"`                                  
	FirstUomQty          float64             `orm:"default(0)"`                                   
	SecondUomQty         float64             `orm:"default(0)"`                                    
	FirstUom             *ProductUom         `orm:"rel(fk)"`                                       
	SecondUom            *ProductUom         `orm:"rel(fk);null"`                                  
	Package              *StockQuantPackage  `orm:"rel(fk);null"`                                
	PackagingType        *ProductPackaging   `orm:"rel(fk)"`                                    
	Reservation          *StockMove          `orm:"rel(fk);null"`                                 
	Lot                  *StockProductionLot `orm:"rel(fk)"`                                       
	Cost                 float64             `orm:"default(0)"`                                    
	InDate               time.Time           `orm:"auto_now_add;type(datetime)"`                 
	Historys             []*StockMove        `orm:"reverse(many);rel_table(stock_quant_move_rel)"`
	Company              *Company            `orm:"rel(fk)"`                                      
	PropagatedFrom       *StockQuant         `orm:"rel(fk);null"`                                 
	NegativeDestLocation *StockLocation      `orm:"rel(fk);null"`                                  
	NegativeMove         *StockMove          `orm:"rel(fk);null"`                                  

	FormAction   string   `orm:"-" json:"FormAction"`   
	ActionFields []string `orm:"-" json:"ActionFields"` 
	CompanyID    int64    `orm:"-" json:"Company"`
}

func init() {
	orm.RegisterModel(new(StockQuant))
}