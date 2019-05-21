package models

import (
	"errors"
	"fmt"
	"projectERP/utils"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

//Team 
type Team struct {
	ID          int64       `orm:"column(id);pk;auto" json:"id"`       
	CreateUser  *User       `orm:"rel(fk);null" json:"-"`             
	UpdateUser  *User       `orm:"rel(fk);null" json:"-"`                
	CreateDate  time.Time   `orm:"auto_now_add;type(datetime)" json:"-"` 
	UpdateDate  time.Time   `orm:"auto_now;type(datetime)" json:"-"`    
	Name        string      `orm:"unique" json:"Name"`                   
	Leader      *User       `orm:"rel(fk);null"`                 
	Company     *Company    `orm:"rel(fk);null"`                     
	Lab  *Lab `orm:"rel(fk);null"`                         
	Members     []*User     `orm:"reverse(many)"`                        
	Active      bool        `orm:"default(true)"`               
	Description string      `orm:"default()" json:"Description"`         

	FormAction   string   `orm:"-" json:"FormAction"`   
	ActionFields []string `orm:"-" json:"ActionFields"` 
	LeaderID     int64    `orm:"-" json:"Leader"`      
	CompanyID    int64    `orm:"-" json:"Company"`      
	LabID int64    `orm:"-" json:"Lab"`   

}

func init() {
	orm.RegisterModel(new(Team))
}

// TableName 
func (u *Team) TableName() string {
	return "base_team"
}

// AddTeam insert a new Team into database and returns
// last inserted ID on success.
func AddTeam(obj *Team, addUser *User) (id int64, err error) {
	o := orm.NewOrm()
	obj.CreateUser = addUser
	obj.UpdateUser = addUser
	errBegin := o.Begin()
	defer func() {
		if err != nil {
			utils.LogOut("error", err.Error())
			if errRollback := o.Rollback(); errRollback != nil {
				err = errRollback
			}

		}
	}()
	if obj.CompanyID > 0 {
		obj.Company, _ = GetCompanyByID(obj.CompanyID)
	}
	if obj.LeaderID > 0 {
		obj.Leader, _ = GetUserByID(obj.LeaderID)
	}
	if obj.LabID > 0 {
		obj.Lab, _ = GetLabByID(obj.LabID)
	}
	if errBegin != nil {
		return 0, errBegin
	}
	id, err = o.Insert(obj)
	if err != nil {
		return 0, err
	}
	errCommit := o.Commit()
	if errCommit != nil {
		return 0, errCommit
	}

	return id, err
}

// GetTeamByID retrieves Team by ID. Returns error if
// ID doesn't exist
func GetTeamByID(id int64) (obj *Team, err error) {
	o := orm.NewOrm()
	obj = &Team{ID: id}
	if err = o.Read(obj); err == nil {
		o.LoadRelated(obj, "Members")
		if obj.Leader != nil {
			o.Read(obj.Leader)
		}
		if obj.Company != nil {
			o.Read(obj.Company)
		}
		if obj.Lab != nil {
			o.Read(obj.Lab)
		}
		return obj, err
	}
	return nil, err
}

// GetTeamByName retrieves Team by Name. Returns error if
// Name doesn't exist
func GetTeamByName(name string) (*Team, error) {
	o := orm.NewOrm()
	var obj Team
	cond := orm.NewCondition()
	cond = cond.And("Name", name)
	qs := o.QueryTable(&obj)
	qs = qs.SetCond(cond)
	err := qs.One(&obj)
	return &obj, err
}

// GetAllTeam retrieves all Team matches certain condition. Returns empty list if
// no records exist
func GetAllTeam(query map[string]interface{}, exclude map[string]interface{}, condMap map[string]map[string]interface{}, fields []string, sortby []string, order []string, offset int64, limit int64) (utils.Paginator, []Team, error) {
	var (
		objArrs   []Team
		paginator utils.Paginator
		num       int64
		err       error
	)
	if limit == 0 {
		limit = 20
	}

	o := orm.NewOrm()
	qs := o.QueryTable(new(Team))
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
				for obj := range objArrs {
					o.LoadRelated(&obj, "Members")
				}
			}
		}
	}

	return paginator, objArrs, err
}

// UpdateTeam updates Team by ID and returns error if
// the record to be updated doesn't exist
func UpdateTeam(obj *Team, updateUser *User) (id int64, err error) {
	o := orm.NewOrm()
	obj.UpdateUser = updateUser
	var num int64
	if num, err = o.Update(obj); err == nil {
		fmt.Println("Number of records updated in database:", num)
	}
	return obj.ID, err
}

// DeleteTeam deletes Team by ID and returns error if
// the record to be deleted doesn't exist
func DeleteTeam(id int64) (err error) {
	o := orm.NewOrm()
	v := Team{ID: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Team{ID: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}