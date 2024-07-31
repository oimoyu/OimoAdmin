package OimoAdmin

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/oimoyu/OimoAdmin/src/db"
	"github.com/oimoyu/OimoAdmin/src/http/router"
	"github.com/oimoyu/OimoAdmin/src/utils/_log"
	"github.com/oimoyu/OimoAdmin/src/utils/_type"
	"github.com/oimoyu/OimoAdmin/src/utils/config"
	_const "github.com/oimoyu/OimoAdmin/src/utils/const"
	"github.com/oimoyu/OimoAdmin/src/utils/front"
	"github.com/oimoyu/OimoAdmin/src/utils/functions"
	"github.com/oimoyu/OimoAdmin/src/utils/g"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"path"
)

func Init(r *gin.Engine, d *gorm.DB, options ...interface{}) {
	var adminConfig _type.AdminConfig
	var redisClient *redis.Client
	for _, option := range options {
		switch opt := option.(type) {
		case _type.AdminConfig:
			adminConfig = opt
		case *redis.Client:
			redisClient = opt
		default:
			panic(fmt.Sprintf("Unknown input type：%T", option))
		}
	}

	// validate gorm db
	if _, err := d.DB(); err != nil {
		panic(fmt.Sprintf("Database connection error: %s", err))
	}

	// validate redis client
	if redisClient != nil {
		if _, err := redisClient.Ping(context.Background()).Result(); err != nil {
			panic(fmt.Sprintf("Redis connection error: %s", err))
		}
	}

	initDir()
	config.LoadAllConfig()

	fileLogPath := path.Join(functions.GetEnvDir(), ".log")

	g.OimoAdmin = &_type.OimoAdminStruct{
		Logger:        _log.NewLogger(_const.OimoAdminString, fileLogPath),
		DB:            d,
		RDB:           redisClient,
		Router:        r,
		Tables:        db.InitTables(d, adminConfig),
		DBName:        getDBName(d),
		RawSqlRecords: make([]string, 0),
	}

	initSDK()
	initFront()
	router.SetupGinRouter()

}

func initDir() {
	var err error
	runtimeDir := functions.GetRuntimeDir()
	err = functions.EnsureDirExists(runtimeDir)
	if err != nil {
		panic(err)
	}

	envDir := functions.GetEnvDir()
	err = functions.EnsureDirExists(envDir)
	if err != nil {
		panic(err)
	}

	wwwRootDir := functions.GetWWWRootDir()
	err = functions.EnsureDirExists(wwwRootDir)
	if err != nil {
		panic(err)
	}

}
func initSDK() {
	sdkUrl := "https://github.com/baidu/amis/releases/download/6.6.0/sdk.tar.gz"
	wwwRootDir := functions.GetWWWRootDir()
	runtimeDir := functions.GetRuntimeDir()
	distDir := path.Join(wwwRootDir, "sdk")
	tempPath := path.Join(runtimeDir, "download_temp")

	// skip if sdk exist
	if functions.IsPathExists(distDir) {
		return
	}

	g.OimoAdmin.Logger.Info("Downloading amis SDK, url: [%s], dist dir: [%s]", sdkUrl, distDir)
	err := functions.DownloadAndExtract(sdkUrl, distDir, tempPath)
	if err != nil {
		g.OimoAdmin.Logger.Error(err.Error())
		g.OimoAdmin.Logger.Error("Failed to download sdk, you can manually download and unzip sdk to dist dir")
		panic(err)
	}
	g.OimoAdmin.Logger.Info("SDK Download Complete.")

}

func initFront() {
	if err := front.GenerateAllFront(); err != nil {
		panic(err)
	}
}

func getDBName(db *gorm.DB) string {
	var dbName string
	var query string

	switch db.Dialector.Name() {
	case "mysql":
		query = "SELECT DATABASE()"
	case "postgres":
		query = "SELECT current_database()"
	case "sqlserver":
		query = "SELECT DB_NAME()"
	case "sqlite":
		var result []struct {
			Name string
		}
		err := db.Raw("PRAGMA database_list").Scan(&result).Error
		if err != nil {
			return ""
		}
		if len(result) > 0 {
			dbName = result[0].Name
		}
		return dbName
	default:
		return ""
	}

	err := db.Raw(query).Scan(&dbName).Error
	if err != nil {
		return ""
	}
	return dbName
}
