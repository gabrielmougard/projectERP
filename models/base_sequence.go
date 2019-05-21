package models

import (
	"bytes"
	"errors"
	"fmt"
	"projectERP/utils"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

//Sequence  
// Table serial number management, used for sales order number, purchase order number, etc., temporarily allows one table to adopt multiple prefixes, the serial numbers are incremented
type Sequence struct {
	ID         int64     `orm:"column(id);pk;auto" json:"id"`         
	CreateUser *User     `orm:"rel(fk);null" json:"-"`               
	UpdateUser *User     `orm:"rel(fk);null" json:"-"`            
	CreateDate time.Time `orm:"auto_now_add;type(datetime)" json:"-"` 
	UpdateDate time.Time `orm:"auto_now;type(datetime)" json:"-"`     
	Company    *Company  `orm:"rel(fk)"`                             
	Name       string    `orm:"unique" json:"Name"`             
	Prefix     string    `orm:"unique" json:"Prefix"`                 
	Current    int64     `json:"Current"`                     
	Padding    int64     `orm:"default(8)" json:"Padding"`            
	StructName string    `json:"StructName"`               
	Active     bool      `orm:"default(true)" json:"Active"`          
	IsDefault  bool      `orm:"default(true)" json:"IsDefault"`      

	FormAction   string   `orm:"-" json:"FormAction"`   
	ActionFields []string `orm:"-" json:"ActionFields"` 
	CompanyID    int64    `orm:"-" json:"Company"`
}

func init() {
	orm.RegisterModel(new(Sequence))
}

// TableName 
func (u *Sequence) TableName() string {
	return "base_sequence"
}

//GetNextSequece 
func GetNextSequece(structName string, companyID int64) (stStr string, err error) {
	o := orm.NewOrm()
	var (
		sequence Sequence
	)
	errBegin := o.Begin()
	defer func() {
		if err != nil {
			if errRollback := o.Rollback(); errRollback != nil {
				err = errRollback
			}
		}
	}()
	if errBegin != nil {
		return "", errBegin
	}
	cond := orm.NewCondition()
	cond = cond.And("StructName", structName).And("active", true).And("IsDefault", true).And("company__id", companyID)
	qs := o.QueryTable(&sequence)
	qs = qs.SetCond(cond)
	if err = qs.One(&sequence); err == nil {
		b := bytes.Buffer{}
		b.WriteString(sequence.Prefix)
		b.WriteString("%0")
		b.WriteString(strconv.Itoa(int(sequence.Padding)))
		b.WriteString("s")
		fmtStr := b.String()
		sequence.Current++
		stStr = fmt.Sprintf(fmtStr, strconv.Itoa(int(sequence.Current)))
		_, err = o.Update(&sequence)
	}
	if err == nil {
		errCommit := o.Commit()
		if errCommit != nil {
			return "", errCommit
		}
	}
	return stStr, err
}

// AddSequence insert a new Sequence into database and returns
// last inserted ID on success.
func AddSequence(obj *Sequence, addUser *User) (id int64, err error) {
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
	id, err = o.Insert(obj)
	if err == nil {
		errCommit := o.Commit()
		if errCommit != nil {
			return 0, errCommit
		}
	}
	return id, err
}

// GetSequenceByID retrieves Sequence by ID. Returns error if
// ID doesn't exist
func GetSequenceByID(id int64) (obj *Sequence, err error) {
	o := orm.NewOrm()
	obj = &Sequence{ID: id}
	if err = o.Read(obj); err == nil {
		return obj, nil
	}
	return nil, err
}

// GetAllSequence retrieves all Sequence matches certain condition. Returns empty list if
// no records exist
func GetAllSequence(query map[string]interface{}, exclude map[string]interface{}, condMap map[string]map[string]interface{}, fields []string, sortby []string, order []string, offset int64, limit int64) (utils.Paginator, []Sequence, error) {
	var (
		objArrs   []Sequence
		paginator utils.Paginator
		num       int64
		err       error
	)
	if limit == 0 {
		limit = 20
	}
	o := orm.NewOrm()
	qs := o.QueryTable(new(Sequence))
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

// UpdateSequenceByID updates Sequence by ID and returns error if
// the record to be updated doesn't exist
func UpdateSequenceByID(m *Sequence) (err error) {
	o := orm.NewOrm()
	v := Sequence{ID: m.ID}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// GetSequenceByName retrieves Sequence by Name. Returns error if
// Name doesn't exist
func GetSequenceByName(name string) (obj *Sequence, err error) {
	o := orm.NewOrm()
	obj = &Sequence{Name: name}
	if err = o.Read(obj); err == nil {
		return obj, nil
	}
	return nil, err
}

// DeleteSequence deletes Sequence by ID and returns error if
// the record to be deleted doesn't exist
func DeleteSequence(id int64) (err error) {
	o := orm.NewOrm()
	v := Sequence{ID: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Sequence{ID: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}