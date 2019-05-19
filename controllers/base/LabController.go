package base

import (
	"bytes"
	"encoding/json"
	md "projectERP/models"
	"strconv"
	"strings"
)

// LabController 
type LabController struct {
	BaseController
}

func (ctl *DepartmentController) Post() {
	ctl.URL = "/lab/"
	ctl.Data["URL"] = ctl.URL
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
func (ctl *LabController) Get() {
	ctl.URL = "/lab/"
	ctl.PageName = "Lab management"
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
	ctl.Data["URL"] = ctl.URL
	ctl.Data["MenuLabActive"] = "active"
}

// Put 
func (ctl *LabController) Put() {
	result := make(map[string]interface{})
	postData := ctl.GetString("postData")
	lab := new(md.Lab)
	var (
		err error
		id  int64
	)
	if err = json.Unmarshal([]byte(postData), lab); err == nil {

		// structName := reflect.Indirect(reflect.ValueOf(lab)).Type().Name()
		if id, err = md.AddLab(lab, &ctl.User); err == nil {
			result["code"] = "success"
			result["location"] = ctl.URL + strconv.FormatInt(id, 10) + "?action=detail"
		} else {
			result["code"] = "failed"
			result["message"] = "Data creation failed"
			result["debug"] = err.Error()
		}
	}
	if err != nil {
		result["code"] = "failed"
		result["debug"] = err.Error()
	}
	ctl.Data["json"] = result
	ctl.ServeJSON()
}
func (ctl *LabController) PostCreate() {
	result := make(map[string]interface{})
	postData := ctl.GetString("postData")
	lab := new(md.Lab)
	var (
		err error
		id  int64
	)
	if err = json.Unmarshal([]byte(postData), lab); err == nil {

		// structName := reflect.Indirect(reflect.ValueOf(lab)).Type().Name()

		if id, err = md.AddLab(lab, &ctl.User); err == nil {
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
func (ctl *LabController) Edit() {
	id := ctl.Ctx.Input.Param(":id")
	if id != "" {
		if idInt64, e := strconv.ParseInt(id, 10, 64); e == nil {
			if lab, err := md.GetLabByID(idInt64); err == nil {
				ctl.PageAction = lab.Name
				ctl.Data["Lab"] = lab
			}
		}
	}
	ctl.Data["Action"] = "edit"
	ctl.Data["RecordID"] = id
	ctl.Data["FormField"] = "form-edit"
	ctl.Layout = "base/base.html"
	ctl.TplName = "user/lab_form.html"
}
func (ctl *LabController) Detail() {
	ctl.Edit()
	ctl.Data["Readonly"] = true
	ctl.Data["FormTreeField"] = "form-tree-edit"
	ctl.Data["Action"] = "detail"
}
func (ctl *LabController) Create() {
	ctl.Data["Action"] = "create"
	ctl.Data["Readonly"] = false
	ctl.PageAction = "create"
	ctl.Layout = "base/base.html"
	ctl.Data["FormField"] = "form-create"
	ctl.Data["FormTreeField"] = "form-tree-create"
	ctl.TplName = "user/lab_form.html"
}

func (ctl *LabController) Validator() {
	name := strings.TrimSpace(ctl.GetString("Name"))
	recordID, _ := ctl.GetInt64("recordID")
	result := make(map[string]bool)
	obj, err := md.GetLabByName(name)
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

// Get style data that meets the requirements
func (ctl *LabController) addressTemplateList(query map[string]interface{}, exclude map[string]interface{}, cond map[string]map[string]interface{}, fields []string, sortby []string, order []string, offset int64, limit int64) (map[string]interface{}, error) {

	var arrs []md.Lab
	paginator, arrs, err := md.GetAllLab(query, exclude, cond, fields, sortby, order, offset, limit)
	result := make(map[string]interface{})
	if err == nil {


		tableLines := make([]interface{}, 0, 4)
		for _, line := range arrs {
			oneLine := make(map[string]interface{})
			oneLine["Name"] = line.Name
			if line.Leader != nil {
				leader := make(map[string]interface{})
				leader["id"] = line.Leader.ID
				leader["name"] = line.Leader.Name
				oneLine["Leader"] = leader
			}
			oneLine["ID"] = line.ID
			oneLine["id"] = line.ID
			if line.Parent != nil {
				oneLine["Parent"] = line.Parent.Name
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
func (ctl *LabController) PostList() {
	query := make(map[string]interface{})
	exclude := make(map[string]interface{})
	cond := make(map[string]map[string]interface{})
	condAnd := make(map[string]interface{})
	condOr := make(map[string]interface{})
	filterMap := make(map[string]interface{})
	fields := make([]string, 0, 0)
	sortby := make([]string, 0, 1)
	order := make([]string, 0, 1)
	if ID, err := ctl.GetInt64("Id"); err == nil {
		query["Id"] = ID
	}
	if name := strings.TrimSpace(ctl.GetString("Name")); name != "" {
		condAnd["Name.icontains"] = name
	}
	filter := ctl.GetString("filter")
	if filter != "" {
		json.Unmarshal([]byte(filter), &filterMap)
	}

	if filterName, ok := filterMap["Name"]; ok {
		filterName = strings.TrimSpace(filterName.(string))
		if filterName != "" {
			condAnd["Name.icontains"] = filterName
		}
	}
	if len(condAnd) > 0 {
		cond["and"] = condAnd
	}
	if len(condOr) > 0 {
		cond["or"] = condOr
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
	if result, err := ctl.addressTemplateList(query, exclude, cond, fields, sortby, order, offset, limit); err == nil {
		ctl.Data["json"] = result
	}
	ctl.ServeJSON()

}

func (ctl *LabController) GetList() {
	viewType := ctl.Input().Get("view")
	if viewType == "" || viewType == "table" {
		ctl.Data["ViewType"] = "table"
	}
	ctl.PageAction = "List"
	ctl.Data["tableId"] = "table-lab"
	ctl.Layout = "base/base_list_view.html"
	ctl.TplName = "user/lab_list_search.html"
}