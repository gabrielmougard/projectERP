package product

import (
	"bytes"
	"encoding/json"
	"projectERP/controllers/base"
	md "projectERP/models"
	"strconv"
	"strings"
)

// ProductAttributeController 
type ProductAttributeController struct {
	base.BaseController
}

// Post 
func (ctl *ProductAttributeController) Post() {
	action := ctl.Input().Get("action")
	switch action {
	case "validator":
		ctl.Validator()
	case "table": 
		ctl.PostList()
	case "create":
		ctl.PostCreate()
	default:
		ctl.PostList()
	}
}

// Put 
func (ctl *ProductAttributeController) Put() {
	id := ctl.Ctx.Input.Param(":id")
	ctl.URL = "/product/attribute/"
	if idInt64, e := strconv.ParseInt(id, 10, 64); e == nil {
		if attribute, err := md.GetProductAttributeByID(idInt64); err == nil {
			if err := ctl.ParseForm(&attribute); err == nil {

				if err := md.UpdateProductAttributeByID(attribute); err == nil {
					ctl.Redirect(ctl.URL+id+"?action=detail", 302)
				}
			}
		}
	}
	ctl.Redirect(ctl.URL+id+"?action=edit", 302)

}

// Get 
func (ctl *ProductAttributeController) Get() {
	ctl.PageName = "Product attribute management"
	action := ctl.Input().Get("action")
	switch action {
	case "create":
		ctl.Create()
	case "edit":
		ctl.Edit()
	case "detail":
		ctl.Detail()
	default:
		ctl.GetList()

	}

	b := bytes.Buffer{}
	b.WriteString(ctl.PageName)
	b.WriteString("\\")
	b.WriteString(ctl.PageAction)
	ctl.Data["PageName"] = b.String()
	ctl.URL = "/product/attribute/"
	ctl.Data["URL"] = ctl.URL

	ctl.Data["MenuProductAttributeActive"] = "active"
}

// Edit 
func (ctl *ProductAttributeController) Edit() {
	id := ctl.Ctx.Input.Param(":id")
	if id != "" {
		if idInt64, e := strconv.ParseInt(id, 10, 64); e == nil {
			if attribute, err := md.GetProductAttributeByID(idInt64); err == nil {
				ctl.PageAction = attribute.Name
				ctl.Data["Attribute"] = attribute

			}
		}
	}
	ctl.Data["FormField"] = "form-edit"
	ctl.Data["Action"] = "edit"
	ctl.Data["RecordID"] = id
	ctl.Layout = "base/base.html"
	ctl.TplName = "product/product_attribute_form.html"
}

// Create 
func (ctl *ProductAttributeController) Create() {
	ctl.Data["Action"] = "create"
	ctl.Data["Readonly"] = false
	ctl.Data["FormField"] = "form-create"
	ctl.PageAction = "create"
	ctl.Layout = "base/base.html"
	ctl.TplName = "product/product_attribute_form.html"
}

// Detail (Product attribute information shows get request, information cannot be modified)
func (ctl *ProductAttributeController) Detail() {

	ctl.Edit()
	ctl.Data["Readonly"] = true
	ctl.Data["Action"] = "detail"
}

// PostCreate 
func (ctl *ProductAttributeController) PostCreate() {
	result := make(map[string]interface{})
	postData := ctl.GetString("postData")
	attribute := new(md.ProductAttribute)
	var (
		err error
		id  int64
	)
	if err = json.Unmarshal([]byte(postData), attribute); err == nil {

		// structName := reflect.Indirect(reflect.ValueOf(category)).Type().Name()
		if id, err = md.AddProductAttribute(attribute, &ctl.User); err == nil {
			result["code"] = "success"
			result["location"] = ctl.URL + strconv.FormatInt(id, 10) + "?action=detail"
		} else {
			result["code"] = "failed"
			result["message"] = "Data creation failed"
			result["debug"] = err.Error()
		}
	} else {
		result["code"] = "failed"
		result["message"] = "Request data parsing failed"
		result["debug"] = err.Error()
	}
	ctl.Data["json"] = result
	ctl.ServeJSON()
}

// Validator (Product attribute information post request for verification)
func (ctl *ProductAttributeController) Validator() {
	name := ctl.GetString("name")
	name = strings.TrimSpace(name)
	recordID, _ := ctl.GetInt64("recordID")
	result := make(map[string]bool)
	obj, err := md.GetProductAttributeByName(name)
	if err != nil {
		result["valid"] = true
	} else {
		if obj.Name == name {
			if recordID == obj.ID {
				result["valid"] = true
			} else {
				result["valid"] = false
			}

		} else {
			result["valid"] = true
		}

	}
	ctl.Data["json"] = result
	ctl.ServeJSON()
}

// Get the data that meets the requirements
func (ctl *ProductAttributeController) productAttributeList(query map[string]interface{}, exclude map[string]interface{}, condMap map[string]map[string]interface{}, fields []string, sortby []string, order []string, offset int64, limit int64) (map[string]interface{}, error) {

	var arrs []md.ProductAttribute
	paginator, arrs, err := md.GetAllProductAttribute(query, exclude, condMap, fields, sortby, order, offset, limit)
	result := make(map[string]interface{})
	if err == nil {


		tableLines := make([]interface{}, 0, 4)
		for _, line := range arrs {
			oneLine := make(map[string]interface{})
			oneLine["Name"] = line.Name
			oneLine["Code"] = line.Code
			oneLine["Sequence"] = line.Sequence
			oneLine["ProductsCount"] = line.ProductsCount
			oneLine["TemplatesCount"] = line.TemplatesCount
			oneLine["ID"] = line.ID
			oneLine["id"] = line.ID
			mapValues := make(map[int64]string)
			values := line.ValueIDs
			for _, line := range values {
				mapValues[line.ID] = line.Name
			}
			oneLine["values"] = mapValues
			tableLines = append(tableLines, oneLine)
		}
		result["data"] = tableLines
		if jsonResult, er := json.Marshal(&paginator); er == nil {
			result["paginator"] = string(jsonResult)
			result["total"] = paginator.TotalCount
		}
	}
	return result, err
}

// PostList 
func (ctl *ProductAttributeController) PostList() {
	query := make(map[string]interface{})
	exclude := make(map[string]interface{})
	fields := make([]string, 0, 0)
	sortby := make([]string, 0, 1)
	order := make([]string, 0, 1)
	cond := make(map[string]map[string]interface{})

	excludeIdsStr := ctl.GetStrings("exclude[]")
	var excludeIds []int64
	for _, v := range excludeIdsStr {
		if val, err := strconv.ParseInt(v, 10, 64); err == nil {
			excludeIds = append(excludeIds, val)
		}
	}
	if len(excludeIds) > 0 {
		exclude["Id.in"] = excludeIds
	}

	offset, _ := ctl.GetInt64("offset")
	limit, _ := ctl.GetInt64("limit")
	orderStr := ctl.GetString("order")
	sortStr := ctl.GetString("sort")
	if orderStr != "" && sortStr != "" {
		sortby = append(sortby, sortStr)
		order = append(order, orderStr)
	} else {
		sortby = append(sortby, "Id")
		order = append(order, "desc")
	}
	if result, err := ctl.productAttributeList(query, exclude, cond, fields, sortby, order, offset, limit); err == nil {
		ctl.Data["json"] = result
	}
	ctl.ServeJSON()

}

// GetList (Product attribute information get request, listing product attributes)
func (ctl *ProductAttributeController) GetList() {
	viewType := ctl.Input().Get("view")
	if viewType == "" || viewType == "table" {
		ctl.Data["ViewType"] = "table"
	}
	ctl.PageAction = "List"
	ctl.Data["tableId"] = "table-product-attribute"
	ctl.Layout = "base/base_list_view.html"
	ctl.TplName = "product/product_attribute_list_search.html"
}