package product

import (
	"bytes"
	"projectERP/controllers/base"
)

// ProductAttributePriceController 
type ProductAttributePriceController struct {
	base.BaseController
}

// Get get
func (ctl *ProductAttributePriceController) Get() {
	action := ctl.GetString(":action")
	viewType := ctl.Input().Get("view_type")
	switch action {
	case "list":
		switch viewType {
		case "list":
			ctl.List()
		default:
			ctl.List()
		}
	default:
		ctl.List()
	}

	
	b := bytes.Buffer{}
	b.WriteString(ctl.PageName)
	b.WriteString("\\")
	b.WriteString(ctl.PageAction)
	ctl.Data["PageName"] = b.String()
	ctl.Layout = "base/base.html"
}
func (ctl *ProductAttributePriceController) List() {

}