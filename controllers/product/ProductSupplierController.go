package product

import (
	"bytes"
	"encoding/json"
	"projectERP/controllers/base"
	md "projectERP/models"
	"strconv"
	"strings"
)

type ProductSupplierController struct {
	base.BaseController
}

func (ctl *ProductSupplierController) Post() {
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
func (ctl *ProductSupplierController) Put() {
	id := ctl.Ctx.Input.Param(":id")
	ctl.URL = "/product/supplier/"
	if idInt64, e := strconv.ParseInt(id, 10, 64); e == nil {
		if supplier, err := md.GetProductSupplierByID(idInt64); err == nil {
			if err := ctl.ParseForm(&supplier); err == nil {
				if err := md.UpdateProductSupplierByID(supplier); err == nil {
					ctl.Redirect(ctl.URL+id+"?action=detail", 302)
				}
			}
		}
	}
	ctl.Redirect(ctl.URL+id+"?action=edit", 302)

}
func (ctl *ProductSupplierController) Get() {
	ctl.PageName = "Product supplier management"
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
	ctl.URL = "/product/supplier/"
	ctl.Data["URL"] = ctl.URL
	ctl.Data["MenuProductSupplierActive"] = "active"
}
func (ctl *ProductSupplierController) Edit() {
	id := ctl.Ctx.Input.Param(":id")
	if id != "" {
		if idInt64, e := strconv.ParseInt(id, 10, 64); e == nil {

			if supplier, err := md.GetProductSupplierByID(idInt64); err == nil {
				ctl.PageAction = supplier.Supplier.Name
				ctl.Data["Supplier"] = supplier
			}
		}
	}
	ctl.Data["FormField"] = "form-edit"
	ctl.Data["Action"] = "edit"
	ctl.Data["RecordID"] = id
	ctl.Layout = "base/base.html"

	ctl.TplName = "product/product_supplier_form.html"
}

func (ctl *ProductSupplierController) Detail() {
	//Just like getting information, call Edit directly
	ctl.Edit()
	ctl.Data["Readonly"] = true
	ctl.Data["Action"] = "detail"
}

//post
func (ctl *ProductSupplierController) PostCreate() {
	result := make(map[string]interface{})
	postData := ctl.GetString("postData")
	supplier := new(md.ProductSupplier)
	var (
		err error
		id  int64
	)
	if err = json.Unmarshal([]byte(postData), supplier); err == nil {

		// structName := reflect.Indirect(reflect.ValueOf(supplier)).Type().Name()
		if id, err = md.AddProductSupplier(supplier, &ctl.User); err == nil {
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
func (ctl *ProductSupplierController) Create() {
	ctl.Data["Action"] = "create"
	ctl.Data["Readonly"] = false
	ctl.PageAction = "create"
	ctl.Data["FormField"] = "form-create"
	ctl.Layout = "base/base.html"
	ctl.TplName = "product/product_supplier_form.html"
}
func (ctl *ProductSupplierController) Validator() {
	name := ctl.GetString("name")
	// recordID, _ := ctl.GetInt64("recordID")
	name = strings.TrimSpace(name)
	result := make(map[string]bool)
	// obj, err := md.GetProductSupplierByName(name)
	// if err != nil {
	// 	result["valid"] = true
	// } else {
	// 	if obj.Name == name {
	// 		if recordID == obj.ID {
	// 			result["valid"] = true
	// 		} else {
	// 			result["valid"] = false
	// 		}

	// 	} else {
	// 		result["valid"] = true
	// 	}

	// }
	ctl.Data["json"] = result
	ctl.ServeJSON()
}

// Get qualified city data
func (ctl *ProductSupplierController) productSupplierList(query map[string]interface{}, exclude map[string]interface{}, condMap map[string]map[string]interface{}, fields []string, sortby []string, order []string, offset int64, limit int64) (map[string]interface{}, error) {

	var arrs []md.ProductSupplier
	paginator, arrs, err := md.GetAllProductSupplier(query, exclude, condMap, fields, sortby, order, offset, limit)
	result := make(map[string]interface{})
	if err == nil {

		// result["recordsFiltered"] = paginator.TotalCount
		tableLines := make([]interface{}, 0, 4)
		for _, line := range arrs {
			oneLine := make(map[string]interface{})
			oneLine["name"] = line.Supplier.Name
			oneLine["ID"] = line.ID
			oneLine["id"] = line.ID
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
func (ctl *ProductSupplierController) PostList() {
	query := make(map[string]interface{})
	exclude := make(map[string]interface{})
	cond := make(map[string]map[string]interface{})

	fields := make([]string, 0, 0)
	sortby := make([]string, 0, 1)
	order := make([]string, 0, 1)
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
	if result, err := ctl.productSupplierList(query, exclude, cond, fields, sortby, order, offset, limit); err == nil {
		ctl.Data["json"] = result
	}
	ctl.ServeJSON()

}

func (ctl *ProductSupplierController) GetList() {
	viewType := ctl.Input().Get("view")
	if viewType == "" || viewType == "table" {
		ctl.Data["ViewType"] = "table"
	}
	ctl.PageAction = "list"
	ctl.Data["tableId"] = "table-product-supplier"
	ctl.Layout = "base/base_list_view.html"
	ctl.TplName = "product/product_supplier_list_search.html"
}