package router

import (
	"fmt"
	"github.com/oimoyu/OimoAdmin/src/http/handler"
	"github.com/oimoyu/OimoAdmin/src/http/middleware"
	_const "github.com/oimoyu/OimoAdmin/src/utils/const"
	"github.com/oimoyu/OimoAdmin/src/utils/front"
	"github.com/oimoyu/OimoAdmin/src/utils/functions"
	"github.com/oimoyu/OimoAdmin/src/utils/g"
	"path"
)

func SetupGinRouter() {
	router := g.OimoAdmin.Router

	//apiPrefix := fmt.Sprintf("%s/%s", _const.OimoAdminApiPrefix, g.AdminPathSecret)
	secretPathPrefix := fmt.Sprintf("%s/%s", _const.OimoAdminStaticPrefix, g.AdminPathSecret)

	router.POST(secretPathPrefix+"/login", handler.Login)
	router.POST(secretPathPrefix+"/logout", middleware.AdminAuthMiddleware(), handler.Logout)

	router.POST(secretPathPrefix+"/sys_info", middleware.AdminAuthMiddleware(), handler.SysInfo)
	router.POST(secretPathPrefix+"/dashboard", middleware.AdminAuthMiddleware(), handler.Dashboard)
	router.POST(secretPathPrefix+"/raw_sql_records", middleware.AdminAuthMiddleware(), handler.RawSqlRecords)
	router.POST(secretPathPrefix+"/execute_raw_sql", middleware.AdminAuthMiddleware(), handler.ExecuteRawSql)
	router.POST(secretPathPrefix+"/create_row", middleware.AdminAuthMiddleware(), handler.CreateRow)
	router.POST(secretPathPrefix+"/delete_rows", middleware.AdminAuthMiddleware(), handler.DeleteRows)
	router.POST(secretPathPrefix+"/update_rows", middleware.AdminAuthMiddleware(), handler.UpdateRows)
	router.POST(secretPathPrefix+"/fetch_list", middleware.AdminAuthMiddleware(), handler.FetchList)
	router.POST(secretPathPrefix+"/log", middleware.AdminAuthMiddleware(), handler.Log)

	wwwRootDir := functions.GetWWWRootDir()
	router.StaticFile(fmt.Sprintf("%s/", secretPathPrefix), path.Join(wwwRootDir, "index.html"))
	router.StaticFile(fmt.Sprintf("%s/site.json", secretPathPrefix), path.Join(wwwRootDir, "site.json"))
	router.Static(fmt.Sprintf("%s/sdk", secretPathPrefix), path.Join(wwwRootDir, "sdk"))

	router.POST(secretPathPrefix+"/page", middleware.AdminAuthMiddleware(), front.PageHandler)
}
