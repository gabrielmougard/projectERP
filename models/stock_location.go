package models

import (
	"errors"
	"fmt"
	"projectERP/utils"
	"strings"
	"time"
	"github.com/astaxie/beego/orm"
)

// StockLocation 
type StockLocation struct {
	ID             int64            `orm:"column(id);pk;auto" json:"id"`         
	CreateUser     *User            `orm:"rel(fk);null" json:"-"`                
	UpdateUser     *User            `orm:"rel(fk);null" json:"-"`                
	CreateDate     time.Time        `orm:"auto_now_add;type(datetime)" json:"-"` 
	UpdateDate     time.Time        `orm:"auto_now;type(datetime)" json:"-"`    
	Name           string           `orm:"unique"`                             
	Company        *Company         `orm:"rel(fk);null"`                         
	Usage          string           `json:"Usage"`                               
	Active         bool             `orm:"default(true)"`                        
	Barcode        string           `json:"Barcode"`                             
	Parent         *StockLocation   `orm:"rel(fk);null"`                        
	Childs         []*StockLocation `orm:"reverse(many)"`                        
	ReturnLocation bool             `orm:"default(false)"`                       
	ScrapLocation  bool             `orm:"default(false)"`                       
	Posx           int64            `json:"Posx"`                                
	Posy           int64            `json:"Posy"`                                
	Posz           int64            `json:"Posz"`                                
	// PutawayStrategy string           ``                                           
	// RemovalStrategy string           ``                                           

	FormAction   string   `orm:"-" json:"FormAction"`   
	ActionFields []string `orm:"-" json:"ActionFields"` 
	CompanyID    int64    `orm:"-" json:"Company"`
	ParentID     int64    `orm:"-" json:"Parent"`
}

func init() {
	orm.RegisterModel(new(StockLocation))
}

// AddStockLocation insert a new StockLocation into database and returns
// last inserted ID on success.
func AddStockLocation(obj *StockLocation, addUser *User) (id int64, err error) {

	o := orm.NewOrm()
	obj.CreateUser = addUser
	obj.UpdateUser = addUser
	errBegin := o.Begin()
	defer func() {
		if err != nil {
			if errRollback := o.Rollback(); errRollback != nil {
				err = errRollback
			}
		}
	}()
	if errBegin != nil {
		return 0, errBegin
	}
	if obj.CompanyID > 0 {
		obj.Company, _ = GetCompanyByID(obj.CompanyID)
	}
	if obj.ParentID > 0 {
		obj.Parent, _ = GetStockLocationByID(obj.ParentID)
	}
	id, err = o.Insert(obj)
	if err == nil {
		errCommit := o.Commit()
		if errCommit != nil {
			return 0, errCommit
		}
	}
	return id, err
}

// GetStockLocationByID retrieves StockLocation by ID. Returns error if
// ID doesn't exist
func GetStockLocationByID(id int64) (obj *StockLocation, err error) {
	o := orm.NewOrm()
	obj = &StockLocation{ID: id}
	if err = o.Read(obj); err == nil {
		return obj, nil
	}
	return nil, err
}

// GetStockLocationByName retrieves StockLocation by Name. Returns error if
// Name doesn't exist
func GetStockLocationByName(name string) (obj *StockLocation, err error) {
	o := orm.NewOrm()
	obj = &StockLocation{Name: name}
	if err = o.Read(obj); err == nil {

		return obj, nil
	}
	return nil, err
}

// GetAllStockLocation retrieves all StockLocation matches certain condition. Returns empty list if
// no records exist
func GetAllStockLocation(query map[string]interface{}, exclude map[string]interface{}, condMap map[string]map[string]interface{}, fields []string, sortby []string, order []string, offset int64, limit int64) (utils.Paginator, []StockLocation, error) {
	var (
		objArrs   []StockLocation
		paginator utils.Paginator
		num       int64
		err       error
	)
	if limit == 0 {
		limit = 20
	}
	o := orm.NewOrm()
	qs := o.QueryTable(new(StockLocation))
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

// UpdateStockLocationByID updates StockLocation by ID and returns error if
// the record to be updated doesn't exist
func UpdateStockLocationByID(m *StockLocation) (err error) {
	o := orm.NewOrm()
	v := StockLocation{ID: m.ID}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteStockLocation deletes StockLocation by ID and returns error if
// the record to be deleted doesn't exist
func DeleteStockLocation(id int64) (err error) {
	o := orm.NewOrm()
	v := StockLocation{ID: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&StockLocation{ID: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}