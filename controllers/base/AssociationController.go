package base

import (
	"bytes"
	"encoding/json"
	md "goERP/models"
	"strconv"
	"strings"
)

// CompanyController 公司
type AssociationController struct {
	BaseController
}

func (ctl *AssociationController) Post() {
	ctl.URL = "/association/"
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
func (ctl *AssociationController) Get() {
	ctl.URL = "/association/"
	ctl.PageName = "Association management"
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
	ctl.Data["MenuCompanyActive"] = "active"
}

// Put 
func (ctl *AssociationController) Put() {
	result := make(map[string]interface{})
	postData := ctl.GetString("postData")
	association := new(md.Association)
	
	var (
		err error
		id  int64
	)

	if err = json.Unmarshal([]byte(postData), association); err == nil {

		// structName := reflect.Indirect(reflect.ValueOf(association)).Type().Name()
		if id, err = md.UpdateAssociation(association, &ctl.User); err == nil {
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
func (ctl *AssociationController) PostCreate() {
	result := make(map[string]interface{})
	postData := ctl.GetString("postData")
	association := new(md.Association)
	var (
		err error
		id  int64
	)
	if err = json.Unmarshal([]byte(postData), association); err == nil {
	
		// structName := reflect.Indirect(reflect.ValueOf(company)).Type().Name()

		if id, err = md.AddAssociation(association, &ctl.User); err == nil {
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
func (ctl *AssociationController) Edit() {
	id := ctl.Ctx.Input.Param(":id")
	if id != "" {
		if idInt64, e := strconv.ParseInt(id, 10, 64); e == nil {
			if association, err := md.GetAssociationByID(idInt64); err == nil {
				ctl.PageAction = association.Name
				ctl.Data["Association"] = association
			}
		}
	}
	ctl.Data["Action"] = "edit"
	ctl.Data["RecordID"] = id
	ctl.Data["FormField"] = "form-edit"
	ctl.Layout = "base/base.html"
	ctl.TplName = "user/association_form.html"
}
func (ctl *AssociationController) Detail() {
	ctl.Edit()
	ctl.Data["Readonly"] = true
	ctl.Data["FormTreeField"] = "form-tree-edit"
	ctl.Data["Action"] = "detail"
}
func (ctl *AssociationController) Create() {
	ctl.Data["Action"] = "create"
	ctl.Data["Readonly"] = false
	ctl.PageAction = "create"
	ctl.Layout = "base/base.html"
	ctl.Data["FormField"] = "form-create"
	ctl.Data["FormTreeField"] = "form-tree-create"
	ctl.TplName = "user/association_form.html"
}

func (ctl *AssociationController) Validator() {
	name := strings.TrimSpace(ctl.GetString("Name"))
	recordID, _ := ctl.GetInt64("recordID")
	result := make(map[string]bool)
	obj, err := md.GetAssociationByName(name)
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


func (ctl *AssociationController) addressTemplateList(query map[string]interface{}, exclude map[string]interface{}, cond map[string]map[string]interface{}, fields []string, sortby []string, order []string, offset int64, limit int64) (map[string]interface{}, error) {

	var arrs []md.Association
	paginator, arrs, err := md.GetAllAssociation(query, exclude, cond, fields, sortby, order, offset, limit)
	result := make(map[string]interface{})
	if err == nil {

		
		tableLines := make([]interface{}, 0, 4)
		for _, line := range arrs {
			oneLine := make(map[string]interface{})
			oneLine["Name"] = line.Name
			oneLine["Code"] = line.Code
			oneLine["ID"] = line.ID
			oneLine["id"] = line.ID
			if line.Parent != nil {
				oneLine["Parent"] = line.Parent.Name
			}
			b := bytes.Buffer{}
			if line.Country != nil {
				b.WriteString(line.Country.Name)
			}
			if line.Province != nil {
				b.WriteString(line.Province.Name)
			}
			if line.City != nil {
				b.WriteString(line.City.Name)
			}
			if line.District != nil {
				b.WriteString(line.District.Name)
			}
			b.WriteString(line.Street)
			oneLine["Address"] = b.String()

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
func (ctl *AssociationController) PostList() {
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
	if len(condAnd) > 0 {
		cond["and"] = condAnd
	}
	if len(condOr) > 0 {
		cond["or"] = condOr
	}
	if result, err := ctl.addressTemplateList(query, exclude, cond, fields, sortby, order, offset, limit); err == nil {
		ctl.Data["json"] = result
	}
	ctl.ServeJSON()

}

func (ctl *AssociationController) GetList() {
	viewType := ctl.Input().Get("view")
	if viewType == "" || viewType == "table" {
		ctl.Data["ViewType"] = "table"
	}
	ctl.PageAction = "List"
	ctl.Data["tableId"] = "table-association"
	ctl.Layout = "base/base_list_view.html"
	ctl.TplName = "user/association_list_search.html"
}