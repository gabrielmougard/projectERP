package models

import (
	"errors"
	"fmt"
	"projectERP/utils"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

//Record 
type Record struct {
	ID         int64     `orm:"column(id);pk;auto" json:"id"`         
	CreateUser *User     `orm:"rel(fk);null" json:"-"`               
	UpdateUser *User     `orm:"rel(fk);null" json:"-"`                
	CreateDate time.Time `orm:"auto_now_add;type(datetime)" json:"-"` 
	UpdateDate time.Time `orm:"auto_now;type(datetime)" json:"-"`     
	FormAction string    `orm:"-" form:"FormAction"`                  
	User       *User     `orm:"rel(fk)"`
	Logout     time.Time `orm:"type(datetime);null"` 
	UserAgent  string    `orm:"null"`                
	IP         string    `orm:"column(ip) "`         
}

// User login record
func init() {
	orm.RegisterModel(new(Record))
}

// TableName 
func (u *Record) TableName() string {
	return "base_record"
}

// AddRecord insert a new Record into database and returns
// last inserted ID on success.
func AddRecord(obj *Record) (id int64, err error) {
	o := orm.NewOrm()

	id, err = o.Insert(obj)
	return id, err
}

// GetRecordByID retrieves Record by ID. Returns error if
// ID doesn't exist
func GetRecordByID(id int64) (obj *Record, err error) {
	o := orm.NewOrm()
	obj = &Record{ID: id}
	if err = o.Read(obj); err == nil {
		return obj, nil
	}
	return nil, err
}

// GetLastRecordByUserID recoed
func GetLastRecordByUserID(userID int64) (Record, error) {
	o := orm.NewOrm()
	var (
		record Record
		err    error
	)

	o.Using("default")
	err = o.QueryTable(&record).Filter("User", userID).RelatedSel().OrderBy("-id").Limit(1).One(&record)
	return record, err
}

// GetAllRecord retrieves all Record matches certain condition. Returns empty list if
// no records exist
func GetAllRecord(query map[string]interface{}, exclude map[string]interface{}, condMap map[string]map[string]interface{}, fields []string, sortby []string, order []string, offset int64, limit int64) (utils.Paginator, []Record, error) {
	var (
		objArrs   []Record
		paginator utils.Paginator
		num       int64
		err       error
	)
	if limit == 0 {
		limit = 20
	}
	o := orm.NewOrm()
	qs := o.QueryTable(new(Record))
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

// UpdateRecordByID updates Record by ID and returns error if
// the record to be updated doesn't exist
func UpdateRecordByID(m *Record) error {
	o := orm.NewOrm()
	v := Record{ID: m.ID}
	var err error
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		_, err = o.Update(m)
	}
	return err
}

// DeleteRecord deletes Record by ID and returns error if
// the record to be deleted doesn't exist
func DeleteRecord(id int64) (err error) {
	o := orm.NewOrm()
	v := Record{ID: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Record{ID: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}