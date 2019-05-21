package models

import (
	"errors"
	"fmt"
	"projectERP/utils"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

//Partners, including suppliers, will automatically create a login account for each partner later
type Partner struct {
	ID         int64            `orm:"column(id);pk;auto" json:"id"`         
	CreateUser *User            `orm:"rel(fk);null" json:"-"`            
	UpdateUser *User            `orm:"rel(fk);null" json:"-"`                
	CreateDate time.Time        `orm:"auto_now_add;type(datetime)" json:"-"` 
	UpdateDate time.Time        `orm:"auto_now;type(datetime)" json:"-"`     
	Name       string           `orm:"unique" json:"Name"`                  
	IsCompany  bool             `orm:"default(true)" json:"IsCompany"`       
	IsSupplier bool             `orm:"default(false)" json:"IsSupplier"`    
	IsCustomer bool             `orm:"default(true)" json:"IsCustomer"`   
	Active     bool             `orm:"default(true)" json:"Active"`          
	Country    *AddressCountry  `orm:"rel(fk);null"`                        
	Province   *AddressProvince `orm:"rel(fk);null"`                      
	City       *AddressCity     `orm:"rel(fk);null"`                   
	District   *AddressDistrict `orm:"rel(fk);null"`                    
	Street     string           `orm:"default(\"\")" json:"Street"`         
	Parent     *Partner         `orm:"rel(fk);null"`                         
	Childs     []*Partner       `orm:"reverse(many)"`                        
	Mobile     string           `orm:"default(\"\")" json:"Mobile"`          
	Tel        string           `orm:"default(\"\")" json:"Tel"`           
	Email      string           `orm:"default(\"\")" json:"Email"`   
	Qq         string           `orm:"default(\"\")" json:"Qq"`              
	WeChat     string           `orm:"default(\"\")" json:"WeChat"`       
	Comment    string           `orm:"type(text)" json:"Comment"`            

	FormAction string `orm:"-" json:"FormAction"` 
	ParentID   int64  `orm:"-" json:"Parent"`    
	CountryID  int64  `orm:"-" json:"Country"`   
	ProvinceID int64  `orm:"-" json:"Province"`
	CityID     int64  `orm:"-" json:"City"` 
	DistrictID int64  `orm:"-" json:"District"`  

}

func init() {
	orm.RegisterModel(new(Partner))
}

// TableName 
func (u *Partner) TableName() string {
	return "base_partner"
}

// AddPartner insert a new Partner into database and returns
// last inserted ID on success.
func AddPartner(obj *Partner, addUser *User) (id int64, err error) {
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
	if obj.ParentID > 0 {
		obj.Parent, _ = GetPartnerByID(obj.ParentID)
	}
	if obj.CountryID > 0 {
		obj.Country, _ = GetAddressCountryByID(obj.CountryID)
	}
	if obj.ProvinceID > 0 {
		obj.Province, _ = GetAddressProvinceByID(obj.ProvinceID)
	}
	if obj.CityID > 0 {
		obj.City, _ = GetAddressCityByID(obj.CityID)
	}
	if obj.DistrictID > 0 {
		obj.District, _ = GetAddressDistrictByID(obj.DistrictID)
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

// GetPartnerByID retrieves Partner by ID. Returns error if
// ID doesn't exist
func GetPartnerByID(id int64) (obj *Partner, err error) {
	o := orm.NewOrm()
	obj = &Partner{ID: id}
	if err = o.Read(obj); err == nil {
		o.LoadRelated(obj, "Parent")
		o.LoadRelated(obj, "Country")
		o.LoadRelated(obj, "Province")
		o.LoadRelated(obj, "City")
		o.LoadRelated(obj, "District")
		return obj, nil
	}
	return nil, err
}

// GetPartnerByName retrieves Partner by Name. Returns error if
// Name doesn't exist
func GetPartnerByName(name string) (*Partner, error) {
	o := orm.NewOrm()
	var obj Partner
	cond := orm.NewCondition()
	cond = cond.And("Name", name)
	qs := o.QueryTable(&obj)
	qs = qs.SetCond(cond)
	err := qs.One(&obj)
	return &obj, err
}

// GetAllPartner retrieves all Partner matches certain condition. Returns empty list if
// no records exist
func GetAllPartner(query map[string]interface{}, exclude map[string]interface{}, condMap map[string]map[string]interface{}, fields []string, sortby []string, order []string, offset int64, limit int64) (utils.Paginator, []Partner, error) {
	var (
		objArrs   []Partner
		paginator utils.Paginator
		num       int64
		err       error
	)
	if limit == 0 {
		limit = 20
	}

	o := orm.NewOrm()
	qs := o.QueryTable(new(Partner))
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

// UpdatePartner updates Partner by ID and returns error if
// the record to be updated doesn't exist
func UpdatePartner(obj *Partner, updateUser *User) (id int64, err error) {
	o := orm.NewOrm()
	obj.UpdateUser = updateUser
	var num int64
	if num, err = o.Update(obj); err == nil {
		fmt.Println("Number of records updated in database:", num)
	}
	return obj.ID, err
}

// DeletePartner deletes Partner by ID and returns error if
// the record to be deleted doesn't exist
func DeletePartner(id int64) (err error) {
	o := orm.NewOrm()
	v := Partner{ID: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Partner{ID: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}