package models

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"projectERP/utils"
	"github.com/astaxie/beego/orm"
)

// SaleOrderLine 
type SaleOrderLine struct {
	ID            int64           `orm:"column(id);pk;auto" json:"id"`         
	CreateUser    *User           `orm:"rel(fk);null" json:"-"`                
	UpdateUser    *User           `orm:"rel(fk);null" json:"-"`                
	CreateDate    time.Time       `orm:"auto_now_add;type(datetime)" json:"-"` 
	UpdateDate    time.Time       `orm:"auto_now;type(datetime)" json:"-"`     
	Name          string          `orm:"default()" json:"Name"`            
	Company       *Company        `orm:"rel(fk)"`                              
	SaleOrder     *SaleOrder      `orm:"rel(fk)"`                             
	Partner       *Partner        `orm:"rel(fk)"`                              
	Product       *ProductProduct `orm:"rel(fk)"`                              
	ProductName   string          `json:"ProductName"`                         
	ProductCode   string          `json:"ProductCode"`                        
	FirstSaleUom  *ProductUom     `orm:"rel(fk)"`                              
	SecondSaleUom *ProductUom     `orm:"rel(fk);null"`                         
	FirstSaleQty  float32         `orm:"default(1)" json:"FirstSaleQty"`       
	SecondSaleQty float32         `orm:"default(0)" json:"SecondSaleQty"`      
	State         string          `orm:"default(draft)"`                       
	PriceUnit     float32         `orm:"default(0)" json:"PriceUnit"`          
	Total         float32         `orm:"default(0)" json:"Total"`              

	FormAction   string   `orm:"-" json:"FormAction"`   
	ActionFields []string `orm:"-" json:"ActionFields"` 
	CompanyID    int64    `orm:"-" json:"Company"`
	SaleOrderID  int64    `orm:"-" json:"SaleOrder"`
	ProductID    int64    `orm:"-" json:"Product"`
}

func init() {
	orm.RegisterModel(new(SaleOrderLine))
}

// AddSaleOrderLine insert a new SaleOrderLine into database and returns
// last inserted ID on success.
func AddSaleOrderLine(obj *SaleOrderLine) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(obj)
	return id, err
}

// GetSaleOrderLineByID retrieves SaleOrderLine by ID. Returns error if
// ID doesn't exist
func GetSaleOrderLineByID(id int64) (obj *SaleOrderLine, err error) {
	o := orm.NewOrm()
	obj = &SaleOrderLine{ID: id}
	if err = o.Read(obj); err == nil {
		return obj, nil
	}
	return nil, err
}

// GetAllSaleOrderLine retrieves all SaleOrderLine matches certain condition. Returns empty list if
// no records exist
func GetAllSaleOrderLine(query map[string]interface{}, exclude map[string]interface{}, condMap map[string]map[string]interface{}, fields []string, sortby []string, order []string, offset int64, limit int64) (utils.Paginator, []SaleOrderLine, error) {
	var (
		objArrs   []SaleOrderLine
		paginator utils.Paginator
		num       int64
		err       error
	)
	if limit == 0 {
		limit = 20
	}
	o := orm.NewOrm()
	qs := o.QueryTable(new(SaleOrderLine))
	qs = qs.RelatedSel()


	cond := orm.NewCondition()
	if _, ok := condMap["and"]; ok {
		andMap := condMap["and"]
		for k, v := range andMap {
			k = strings.Replace(k, ".", "__", -1)
			cond = cond.And(k, v)
		}
	}
	if _, ok := condMap["or"]; ok {
		orMap := condMap["or"]
		for k, v := range orMap {
			k = strings.Replace(k, ".", "__", -1)
			cond = cond.Or(k, v)
		}
	}
	qs = qs.SetCond(cond)
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		qs = qs.Filter(k, v)
	}
	//exclude k=v
	for k, v := range exclude {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		qs = qs.Exclude(k, v)
	}

	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + strings.Replace(v, ".", "__", -1)
				} else if order[i] == "asc" {
					orderby = strings.Replace(v, ".", "__", -1)
				} else {
					return paginator, nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + strings.Replace(v, ".", "__", -1)
				} else if order[0] == "asc" {
					orderby = strings.Replace(v, ".", "__", -1)
				} else {
					return paginator, nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return paginator, nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return paginator, nil, errors.New("Error: unused 'order' fields")
		}
	}

	qs = qs.OrderBy(sortFields...)
	if cnt, err := qs.Count(); err == nil {
		if cnt > 0 {
			paginator = utils.GenPaginator(limit, offset, cnt)
			if num, err = qs.Limit(limit, offset).All(&objArrs, fields...); err == nil {
				paginator.CurrentPageSize = num
			}
		}
	}
	return paginator, objArrs, err
}

// UpdateSaleOrderLineByID updates SaleOrderLine by ID and returns error if
// the record to be updated doesn't exist
func UpdateSaleOrderLineByID(m *SaleOrderLine) (err error) {
	o := orm.NewOrm()
	v := SaleOrderLine{ID: m.ID}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// GetSaleOrderLineByName retrieves SaleOrderLine by Name. Returns error if
// Name doesn't exist
func GetSaleOrderLineByName(name string) (obj *SaleOrderLine, err error) {
	o := orm.NewOrm()
	obj = &SaleOrderLine{Name: name}
	if err = o.Read(obj); err == nil {
		return obj, nil
	}
	return nil, err
}

// DeleteSaleOrderLine deletes SaleOrderLine by ID and returns error if
// the record to be deleted doesn't exist
func DeleteSaleOrderLine(id int64) (err error) {
	o := orm.NewOrm()
	v := SaleOrderLine{ID: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&SaleOrderLine{ID: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}