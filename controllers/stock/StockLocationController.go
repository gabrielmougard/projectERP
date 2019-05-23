package stock

import (
	"bytes"
	"encoding/json"
	"projectERP/controllers/base"
	md "projectERP/models"
	"strconv"
	"strings"
)

// StockLocationController 
type StockLocationController struct {
	base.BaseController
}

// Post 
func (ctl *StockLocationController) Post() {
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
func (ctl *StockLocationController) Put() {
	id := ctl.Ctx.Input.Param(":id")
	ctl.URL = "/stock/location/"
	if idInt64, e := strconv.ParseInt(id, 10, 64); e == nil {
		if location, err := md.GetStockLocationByID(idInt64); err == nil {
			if err := ctl.ParseForm(&location); err == nil {

				if err := md.UpdateStockLocationByID(location); err == nil {
					ctl.Redirect(ctl.URL+id+"?action=detail", 302)
				}
			}
		}
	}
	ctl.Redirect(ctl.URL+id+"?action=edit", 302)

}

// Get 
func (ctl *StockLocationController) Get() {
	ctl.PageName = "Location management"
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
	ctl.URL = "/stock/location/"
	ctl.Data["URL"] = ctl.URL

	ctl.Data["MenuStockLocationActive"] = "active"
}

// Edit 
func (ctl *StockLocationController) Edit() {
	id := ctl.Ctx.Input.Param(":id")
	if id != "" {
		if idInt64, e := strconv.ParseInt(id, 10, 64); e == nil {
			if location, err := md.GetStockLocationByID(idInt64); err == nil {
				ctl.PageAction = location.Name
				ctl.Data["Location"] = location

			}
		}
	}
	ctl.Data["FormField"] = "form-edit"
	ctl.Data["Action"] = "edit"
	ctl.Data["RecordID"] = id
	ctl.Layout = "base/base.html"
	ctl.TplName = "stock/stock_location_form.html"
}

// Create 
func (ctl *StockLocationController) Create() {
	ctl.Data["Action"] = "create"
	ctl.Data["Readonly"] = false
	ctl.Data["FormField"] = "form-create"
	ctl.PageAction = "创建"
	ctl.Layout = "base/base.html"
	ctl.TplName = "stock/stock_location_form.html"
}

// Detail (Product attribute information shows get request, information cannot be modified)
func (ctl *StockLocationController) Detail() {

	ctl.Edit()
	ctl.Data["Readonly"] = true
	ctl.Data["Action"] = "detail"
}

// PostCreate 
func (ctl *StockLocationController) PostCreate() {
	result := make(map[string]interface{})
	postData := ctl.GetString("postData")
	location := new(md.StockLocation)
	var (
		err error
		id  int64
	)
	if err = json.Unmarshal([]byte(postData), location); err == nil {
	
		// structName := reflect.Indirect(reflect.ValueOf(category)).Type().Name()
		if id, err = md.AddStockLocation(location, &ctl.User); err == nil {
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
func (ctl *StockLocationController) Validator() {
	name := ctl.GetString("name")
	name = strings.TrimSpace(name)
	recordID, _ := ctl.GetInt64("recordID")
	result := make(map[string]bool)
	obj, err := md.GetStockLocationByName(name)
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
func (ctl *StockLocationController) stockLocationList(query map[string]interface{}, exclude map[string]interface{}, condMap map[string]map[string]interface{}, fields []string, sortby []string, order []string, offset int64, limit int64) (map[string]interface{}, error) {

	var arrs []md.StockLocation
	paginator, arrs, err := md.GetAllStockLocation(query, exclude, condMap, fields, sortby, order, offset, limit)
	result := make(map[string]interface{})
	if err == nil {


		tableLines := make([]interface{}, 0, 4)
		for _, line := range arrs {
			oneLine := make(map[string]interface{})
			oneLine["Name"] = line.Name
			oneLine["ID"] = line.ID
			oneLine["id"] = line.ID
			oneLine["Usage"] = line.Usage
			oneLine["Active"] = line.Active
			oneLine["Barcode"] = line.Barcode
			oneLine["ReturnLocation"] = line.ReturnLocation
			oneLine["ScrapLocation"] = line.ScrapLocation
			oneLine["Posx"] = line.Posx
			oneLine["Posy"] = line.Posy
			oneLine["Posz"] = line.Posz

			if line.Company != nil {
				company := make(map[string]interface{})
				company["id"] = line.Company.ID
				company["name"] = line.Company.Name
				oneLine["Company"] = company
			}
			if line.Parent != nil {
				parent := make(map[string]interface{})
				parent["id"] = line.Parent.ID
				parent["name"] = line.Parent.Name
				oneLine["Parent"] = parent
			}
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
func (ctl *StockLocationController) PostList() {
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
	if result, err := ctl.stockLocationList(query, exclude, cond, fields, sortby, order, offset, limit); err == nil {
		ctl.Data["json"] = result
	}
	ctl.ServeJSON()

}

// GetList 
func (ctl *StockLocationController) GetList() {
	viewType := ctl.Input().Get("view")
	if viewType == "" || viewType == "table" {
		ctl.Data["ViewType"] = "table"
	}
	ctl.PageAction = "list"
	ctl.Data["tableId"] = "table-stock-location"
	ctl.Layout = "base/base_list_view.html"
	ctl.TplName = "stock/stock_location_list_search.html"
}