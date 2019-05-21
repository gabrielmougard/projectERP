package models

import (
	"errors"
	"fmt"
	"projectERP/utils"
	"strings"
	"time"
	"github.com/astaxie/beego/orm"
)

//ProductPricelistItem 
type ProductPricelistItem struct {
	ID         int64     `orm:"column(id);pk;auto" json:"id"`        
	CreateUser *User     `orm:"rel(fk);null" json:"-"`               
	UpdateUser *User     `orm:"rel(fk);null" json:"-"`                
	CreateDate time.Time `orm:"auto_now_add;type(datetime)" json:"-"` 
	UpdateDate time.Time `orm:"auto_now;type(datetime)" json:"-"`     

	FormAction   string   `orm:"-" json:"FormAction"`   
	ActionFields []string `orm:"-" json:"ActionFields"` 
}

func init() {
	orm.RegisterModel(new(ProductPricelistItem))
}

// AddProductPricelistItem insert a new ProductPricelistItem into database and returns
// last inserted ID on success.
func AddProductPricelistItem(obj *ProductPricelistItem) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(obj)
	return id, err
}

// GetProductPricelistItemByID retrieves ProductPricelistItem by ID. Returns error if
// ID doesn't exist
func GetProductPricelistItemByID(id int64) (obj *ProductPricelistItem, err error) {
	o := orm.NewOrm()
	obj = &ProductPricelistItem{ID: id}
	if err = o.Read(obj); err == nil {
		return obj, nil
	}
	return nil, err
}

// GetAllProductPricelistItem retrieves all ProductPricelistItem matches certain condition. Returns empty list if
// no records exist
func GetAllProductPricelistItem(query map[string]interface{}, exclude map[string]interface{}, condMap map[string]map[string]interface{}, fields []string, sortby []string, order []string, offset int64, limit int64) (utils.Paginator, []ProductPricelistItem, error) {
	var (
		objArrs   []ProductPricelistItem
		paginator utils.Paginator
		num       int64
		err       error
	)
	if limit == 0 {
		limit = 20
	}
	o := orm.NewOrm()
	qs := o.QueryTable(new(ProductPricelistItem))
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

// UpdateProductPricelistItemByID updates ProductPricelistItem by ID and returns error if
// the record to be updated doesn't exist
func UpdateProductPricelistItemByID(m *ProductPricelistItem) (err error) {
	o := orm.NewOrm()
	v := ProductPricelistItem{ID: m.ID}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteProductPricelistItem deletes ProductPricelistItem by ID and returns error if
// the record to be deleted doesn't exist
func DeleteProductPricelistItem(id int64) (err error) {
	o := orm.NewOrm()
	v := ProductPricelistItem{ID: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&ProductPricelistItem{ID: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}