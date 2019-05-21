package models

import (
	"errors"
	"fmt"
	"projectERP/utils"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

// Permission  
type Permission struct {
	ID         int64     `orm:"column(id);pk;auto" json:"id"`         
	CreateUser *User     `orm:"rel(fk);null" json:"-"`           
	UpdateUser *User     `orm:"rel(fk);null" json:"-"`                
	CreateDate time.Time `orm:"auto_now_add;type(datetime)" json:"-"` 
	UpdateDate time.Time `orm:"auto_now;type(datetime)" json:"-"`     
	Name       string    `orm:"unique" json:"Name"`              
	Source     *Source   `orm:"rel(fk)"`                              
	PermCreate bool      `orm:"default(true)" json:"PermCreate"`      
	PermRead   bool      `orm:"default(true)" json:"PermRead"`        
	PermWrite  bool      `orm:"default(true)" json:"PermWrite"`       
	PermDelete bool      `orm:"default(false)" json:"PermDelete"`    
	Relation   string    `json:"Relation"`                           
	Roles      []*Role   `orm:"rel(m2m)"`                            

	FormAction   string   `orm:"-" json:"FormAction"`   
	ActionFields []string `orm:"-" json:"ActionFields"` 
}

func init() {
	orm.RegisterModel(new(Permission))
}

// TableName 
func (u *Permission) TableName() string {
	return "base_permission"
}

// AddPermission insert a new Permission into database and returns
// last inserted ID on success.
func AddPermission(obj *Permission, addUser *User) (id int64, err error) {
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

// GetPermissionByID retrieves Permission by ID. Returns error if
// ID doesn't exist
func GetPermissionByID(id int64) (obj *Permission, err error) {
	o := orm.NewOrm()
	obj = &Permission{ID: id}
	if err = o.Read(obj); err == nil {
		return obj, err
	}
	return nil, err
}

// GetPermissionByName retrieves Permission by Name. Returns error if
// Name doesn't exist
func GetPermissionByName(name string) (*Permission, error) {
	o := orm.NewOrm()
	var obj Permission
	cond := orm.NewCondition()
	cond = cond.And("Name", name)
	qs := o.QueryTable(&obj)
	qs = qs.SetCond(cond)
	err := qs.One(&obj)
	return &obj, err
}

// GetAllPermission retrieves all Permission matches certain condition. Returns empty list if
// no records exist
func GetAllPermission(query map[string]interface{}, exclude map[string]interface{}, condMap map[string]map[string]interface{}, fields []string, sortby []string, order []string, offset int64, limit int64) (utils.Paginator, []Permission, error) {
	var (
		objArrs   []Permission
		paginator utils.Paginator
		num       int64
		err       error
	)
	if limit == 0 {
		limit = 20
	}

	o := orm.NewOrm()
	qs := o.QueryTable(new(Permission))
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

// UpdatePermission updates Permission by ID and returns error if
// the record to be updated doesn't exist
func UpdatePermission(obj *Permission, updateUser *User) (id int64, err error) {
	o := orm.NewOrm()
	obj.UpdateUser = updateUser
	var num int64
	if num, err = o.Update(obj); err == nil {
		fmt.Println("Number of records updated in database:", num)
	}
	return obj.ID, err
}

// DeletePermission deletes Permission by ID and returns error if
// the record to be deleted doesn't exist
func DeletePermission(id int64) (err error) {
	o := orm.NewOrm()
	v := Permission{ID: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Permission{ID: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}