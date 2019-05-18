package routers

import (
	"projectERP/controllers/address"
	"projectERP/controllers/base"
	"projectERP/controllers/product"
	"projectERP/controllers/purchase"
	"projectERP/controllers/sale"
	"projectERP/controllers/stock"

	"github.com/astaxie/beego"

)

func init() {

	beego.Router("/", &base.IndexController{})
	//======================================Basic operations===========================================
	//Log In
	beego.Router("/login/:action([A-Za-z]+)/", &base.LoginController{})
	//User
	beego.Router("/user/?:id", &base.UserController{})
	//the company
	beego.Router("/company/?:id", &base.CompanyController{})
	//the departement
	beego.Router("/department/?:id", &base.DepartmentController{})
	//the position
	beego.Router("/position/?:id", &base.PositionController{})
	//the team
	beego.Router("/team/?:id", &base.TeamController{})
	//the record
	beego.Router("/record/", &base.RecordController{})
	//Sequence generator
	beego.Router("/sequence/?:id", &base.SequenceController{})
	//File template
	beego.Router("/templatefile/?:id", &base.TemplateFileController{})
	// ===============================Access control===========================================
	// System ressource
	beego.Router("/source/?:id", &base.SourceController{})
	// Character
	beego.Router("/role/?:id", &base.RoleController{})
	// Access control
	beego.Router("/permission/?:id", &base.PermissionController{})
	// Menu control
	beego.Router("/menu/?:id", &base.MenuController{})
	// ===============================Address===========================================
	//country
	beego.Router("/address/country/?:id", &address.AddressCountryController{})
	//province
	beego.Router("/address/province/?:id", &address.AddressProvinceController{})
	//city
	beego.Router("/address/city/?:id", &address.AddressCityController{})
	//district
	beego.Router("/address/district/?:id", &address.AddressDistrictController{})
	//=======================================Product management===========================================

	//Product category
	beego.Router("/product/category/?:id", &product.ProductCategoryController{})
	//Attributes
	beego.Router("/product/attribute/?:id", &product.ProductAttributeController{})
	//Attribute details
	beego.Router("/product/attributevalue/?:id", &product.ProductAttributeValueController{})
	//Attribute value details
	beego.Router("/product/attributeline/?:id", &product.ProductAttributeLineController{})
	//Product template
	beego.Router("/product/template/?:id", &product.ProductTemplateController{})
	//Product specifications
	beego.Router("/product/product/?:id", &product.ProductProductController{})

	//Product label
	beego.Router("/product/tag/:action([A-Za-z]+)/?:id", &product.ProductTagController{})
	//Product packaging
	beego.Router("/product/packaging/:action([A-Za-z]+)/?:id", &product.ProductPackagingController{})
	//Product attribute price
	beego.Router("/product/attributeprice/:action([A-Za-z]+)/?:id", &product.ProductAttributePriceController{})
	//Product unit of measure
	beego.Router("/product/uom/?:id", &product.ProductUomController{})
	//category product unit of measure
	beego.Router("/product/uomcateg/?:id", &product.ProductUomCategController{})
	//========================================Partner management===============================
	//Partner management
	beego.Router("/partner/?:id", &base.PartnerController{})
	//=======================================Sales order management===========================================
	//Sales setup
	beego.Router("/sale/config/?:id", &sale.SaleConfigController{})
	//counter
	beego.Router("/sale/counter/?:id", &sale.SaleCounterController{})
	//Counter product
	beego.Router("/sale/counter/product/?:id", &sale.SaleCounterProductController{})
	//Sales order
	beego.Router("/sale/order/?:id", &sale.SaleOrderController{})
	//Sales order details
	beego.Router("/sale/order/line/?:id", &sale.SaleOrderLineController{})
	//Sales order
	beego.Router("/sale/order/state/?:id", &sale.SaleOrderStateController{})
	//========================================Purchase order management=====================================
	//Purchase settigns
	beego.Router("/purchase/config/?:id", &purchase.PurchaseConfigController{})
	//Purchase order
	beego.Router("/purchase/order/?:id", &purchase.PurchaseOrderController{})
	//Purchase order details
	beego.Router("/purchase/order/line/?:id", &purchase.PurchaseOrderLineController{})
	//Purchase order 
	beego.Router("/purchase/order/state/?:id", &purchase.PurchaseOrderStateController{})
	//========================================Warehouse management=====================================
	//  Warehouse management
	beego.Router("/stock/warehouse/?:id", &stock.StockWarehouseController{})
	//  Warehouse document management
	beego.Router("/stock/picking/type/?:id", &stock.StockPickingTypeController{})
	//  Warehouse document management
	beego.Router("/stock/picking/?:id", &stock.StockPickingController{})
	// Location management
	beego.Router("/stock/location/?:id", &stock.StockLocationController{})
	// Inventory management
	beego.Router("/stock/inventory/?:id", &stock.StockInventoryController{})
}