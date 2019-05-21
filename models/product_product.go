package models

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"time"
	"projectERP/utils"
	"strconv"
	"github.com/astaxie/beego/orm"
)

// ProductProduct (product specification)
type ProductProduct struct {
	ID                    int64                    `orm:"column(id);pk;auto" json:"id"`         
	CreateUser            *User                    `orm:"rel(fk);null" json:"-"`                
	UpdateUser            *User                    `orm:"rel(fk);null" json:"-"`                
	CreateDate            time.Time                `orm:"auto_now_add;type(datetime)" json:"-"`
	UpdateDate            time.Time                `orm:"auto_now;type(datetime)" json:"-"`     
	Name                  string                   `orm:"index"`                                
	Company               *Company                 `orm:"rel(fk);null"`                         
	Category              *ProductCategory         `orm:"rel(fk)"`                             
	IsProductVariant      bool                     `orm:"default(true)"`                        
	ProductTags           []*ProductTag            `orm:"reverse(many)"`                        
	SaleOk                bool                     `orm:"default(true)" json:"SaleOk"`          
	Active                bool                     `orm:"default(true)"`                        
	Barcode               string                   `orm:"null" json:"Barcode"`                  
	StandardPrice         float64                  `json:"StandardPrice"`                      
	DefaultCode           string                   `orm:"unique"`                               
	ProductTemplate       *ProductTemplate         `orm:"rel(fk)"`                              
	AttributeValues       []*ProductAttributeValue `orm:"reverse(many)"`                       
	ProductType           string                   `orm:"default(stock)"`                       
	AttributeValuesString string                   `orm:"index;default()"`                      
	FirstSaleUom          *ProductUom              `orm:"rel(fk)"`                            
	SecondSaleUom         *ProductUom              `orm:"rel(fk);null"`                         
	FirstPurchaseUom      *ProductUom              `orm:"rel(fk)"`                              
	SecondPurchaseUom     *ProductUom              `orm:"rel(fk);null"`                         
	ProductPackagings     []*ProductPackaging      `orm:"reverse(many)"`                      
	PackagingDependTemp   bool                     `orm:"default(true)"`                      
	BigImages             []*ProductImage          `orm:"reverse(many)"`                     
	MidImages             []*ProductImage          `orm:"reverse(many)"`                  
	SmallImages           []*ProductImage          `orm:"reverse(many)"`                     
	PurchaseDependTemp    bool                     `orm:"default(true)"`                        

	FormAction            string                 `orm:"-" json:"FormAction"`        
	ActionFields          []string               `orm:"-" json:"ActionFields"`      
	CategoryID            int64                  `orm:"-" json:"Category"`          
	FirstSaleUomID        int64                  `orm:"-" json:"FirstSaleUom"`     
	SecondSaleUomID       int64                  `orm:"-" json:"SecondSaleUom"`    
	FirstPurchaseUomID    int64                  `orm:"-" json:"FirstPurchaseUom"`  
	SecondPurchaseUomID   int64                  `orm:"-" json:"SecondPurchaseUom"` 
	ProductCounterID      int64                  `orm:"-" json:"ProductCounter"`    
	ProductAttributeLines []ProductAttributeLine `orm:"-" json:"ProductAttributes"` 
	ProductTemplateID     int64                  `orm:"-" json:"ProductTemplateID"` 
	AttributeValueIDs     map[string][]int64     `orm:"-" json:"AttributeValueIds"` 

}

func init() {
	orm.RegisterModel(new(ProductProduct))
}

// AddProductProduct insert a new ProductProduct into database and returns
// last inserted ID on success.
func AddProductProduct(obj *ProductProduct, addUser *User) (id int64, err error) {
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
	if obj.ProductTemplateID > 0 {
		if template, err := GetProductTemplateByID(obj.ProductTemplateID); err == nil {
			obj.ProductTemplate = template
			sequence := GetVariantCount(template)
			b := bytes.Buffer{}
			b.WriteString(template.DefaultCode)
			b.WriteString("-")
			b.WriteString(strconv.FormatInt(sequence+1, 10))
			obj.DefaultCode = b.String()
		} else {
			return 0, err
		}
	}
	// if len(obj.AttributeValueIDs) > 0 {
	// 	strArr := make([]string, 0, 0)
	// 	for _, item := range obj.AttributeValueIDs {
	// 		strArr = append(strArr, strconv.FormatInt(item, 10))

	// 	}
	// 	sort.Strings(strArr)
	// 	obj.AttributeValuesString = strings.Join(strArr, "-")
	// }
	if obj.CategoryID > 0 {
		obj.Category, _ = GetProductCategoryByID(obj.CategoryID)
	}
	if obj.FirstSaleUomID > 0 {
		obj.FirstSaleUom, _ = GetProductUomByID(obj.FirstSaleUomID)
	}
	if obj.SecondSaleUomID > 0 {
		obj.SecondSaleUom, _ = GetProductUomByID(obj.SecondSaleUomID)
	}
	if obj.FirstPurchaseUomID > 0 {
		obj.FirstPurchaseUom, _ = GetProductUomByID(obj.FirstPurchaseUomID)
	}
	if obj.SecondPurchaseUomID > 0 {
		obj.SecondPurchaseUom, _ = GetProductUomByID(obj.SecondPurchaseUomID)
	}
	if id, err = o.Insert(obj); err != nil {
		return 0, err
	}
	obj.ID = id
	if createAttributeValuesRecords, ok := obj.AttributeValueIDs["create"]; ok {
		m2m := o.QueryM2M(obj, "AttributeValues")
		Oattr := orm.NewOrm()
		for _, attrValID := range createAttributeValuesRecords {
			if attributeValue, err := GetProductAttributeValueByID(attrValID); err == nil {
				m2m.Add(attributeValue)
				m2mAttr := Oattr.QueryM2M(attributeValue.Attribute, "Products")
				m2mAttr.Add(obj)
				UpdateProductAttributeValueProductsCount(attributeValue, addUser)
				UpdateProductAttributeProductsCount(attributeValue.Attribute, addUser)
			}
		}
	}
	// Oattr := orm.NewOrm()
	// for _, item := range obj.AttributeValueIDs {
	// 	m2m := o.QueryM2M(obj, "AttributeValues")
	// 	if attributeValue, err := GetProductAttributeValueByID(item); err == nil {
	// 		m2m.Add(attributeValue)
	// 		m2mAttr := Oattr.QueryM2M(attributeValue.Attribute, "Products")
	// 		m2mAttr.Add(obj)
	// 		UpdateProductAttributeValueProductsCount(attributeValue, addUser)
	// 		UpdateProductAttributeProductsCount(attributeValue.Attribute, addUser)
	// 	}
	// }

	if err != nil {
		return 0, err
	}
	errCommit := o.Commit()
	if errCommit != nil {
		return 0, errCommit
	}

	return id, err
}

// GetProductProductByID retrieves ProductProduct by ID. Returns error if
// ID doesn't exist
func GetProductProductByID(id int64) (obj *ProductProduct, err error) {
	o := orm.NewOrm()
	obj = &ProductProduct{ID: id}
	if err = o.Read(obj); err == nil {
		if obj.ProductTemplate != nil {
			o.Read(obj.ProductTemplate)
		}
		if obj.AttributeValues != nil {
			o.LoadRelated(obj.AttributeValues, "AttributeValues")
		}
		if obj.Category != nil {
			o.Read(obj.Category)
		}
		if obj.FirstSaleUom != nil {
			o.Read(obj.FirstSaleUom)
		}
		if obj.FirstPurchaseUom != nil {
			o.Read(obj.FirstPurchaseUom)
		}
		if obj.SecondSaleUom != nil {
			o.Read(obj.SecondSaleUom)
		}
		if obj.SecondPurchaseUom != nil {
			o.Read(obj.SecondPurchaseUom)
		}
		o.LoadRelated(obj, "AttributeValues")

		return obj, nil
	}
	return nil, err
}

// GetProductProductByName retrieves ProductProduct by Name. Returns error if
// Name doesn't exist
func GetProductProductByName(name string) (obj *ProductProduct, err error) {
	o := orm.NewOrm()
	obj = &ProductProduct{Name: name}
	if err = o.Read(obj); err == nil {
		return obj, nil
	}
	return nil, err
}

// GetAllProductProduct retrieves all ProductProduct matches certain condition. Returns empty list if
// no records exist
func GetAllProductProduct(query map[string]interface{}, exclude map[string]interface{}, condMap map[string]map[string]interface{}, fields []string, sortby []string, order []string, offset int64, limit int64) (utils.Paginator, []ProductProduct, error) {
	var (
		objArrs   []ProductProduct
		paginator utils.Paginator
		num       int64
		err       error
	)
	if limit == 0 {
		limit = 20
	}
	o := orm.NewOrm()
	qs := o.QueryTable(new(ProductProduct))
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
	
	for obj := range objArrs {
		o.LoadRelated(&obj, "AttributeValues")
	}
	return paginator, objArrs, err
}

// UpdateProductProductByID updates ProductProduct by ID and returns error if
// the record to be updated doesn't exist
func UpdateProductProductByID(m *ProductProduct) (err error) {
	o := orm.NewOrm()
	v := ProductProduct{ID: m.ID}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteProductProduct deletes ProductProduct by ID and returns error if
// the record to be deleted doesn't exist
func DeleteProductProduct(id int64) (err error) {
	o := orm.NewOrm()
	v := ProductProduct{ID: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&ProductProduct{ID: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

// BatchUpdateProductProduct 
func BatchUpdateProductProduct(query map[string]interface{}, fields map[string]interface{}) (err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(ProductProduct))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		qs = qs.Filter(k, v)
	}
	_, err = qs.Update(fields)
	return err
}