package main

import (
	. "projectERP/init"
	_ "projectERP/models"
	_ "projectERP/routers"

	"projectERP/utils"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/beego/i18n"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"

)

const (
	APP_VER = "1.0"
)

func init() {
	dbType := beego.AppConfig.String("db_type")
	dbAlias := beego.AppConfig.String(dbType + "::db_alias")
	dbName := beego.AppConfig.String(dbType + "::db_name")
	dbUser := beego.AppConfig.String(dbType + "::db_user")
	dbPwd := beego.AppConfig.String(dbType + "::db_pwd")
	dbPort := beego.AppConfig.String(dbType + "::db_port")
	dbHost := beego.AppConfig.String(dbType + "::db_host")

	orm.RegisterDriver(dbType,orm.DRPostgres)

	switch dbType {
	case "postgres":

		dbSslmode := "disable"
		dataSource := "user=" + dbUser + " password=" + dbPwd + " dbname=" + dbName + " host=" + dbHost + " port=" + dbPort + " sslmode=" + dbSslmode
		orm.RegisterDataBase(dbAlias, dbType, dataSource)

	case "mysql":

		dbCharset := beego.AppConfig.String(dbType + "db_charset")
		dataSource := dbUser + ":" + dbPwd + "@/" + dbName + "?charset=" + dbCharset
		orm.RegisterDataBase(dbAlias, dbType, dataSource)

	case "sqlite3":

		orm.RegisterDataBase(dbAlias, "sqlite3", dbName)
	
	default :
		utils.LogOut("info","the database type hasn't been chosed !")
	
	}
	
	utils.LogOut("info","Use the database as : "+dbType)

	//If we rerun the program it will not delete the original table
	coverDb := true
	if coverDb {
		utils.LogOut("info","coverDb is true")
	}
	

	//automatic table construction
	orm.RunSyncdb(dbAlias, coverDb, true)
	InitApp()
	InitDb()

	//cache initialization
	utils.InitCache()
	beego.AddFuncMap("i18n",i18n.Tr) //internationalization


}

func main() {

	
	utils.LogOut("info", "start server")
	beego.Run()
}