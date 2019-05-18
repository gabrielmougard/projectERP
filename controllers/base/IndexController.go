package base

// IndexController home controller
type IndexController struct {
	BaseController
}

// Get home page
func (ctl *IndexController) Get() {

	// Basic layout page
	ctl.Layout = "base/base.html"
	ctl.TplName = "base/module_dashboard.html"

}