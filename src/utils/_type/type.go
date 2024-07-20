package _type

import (
	"github.com/gin-gonic/gin"
	"github.com/oimoyu/OimoAdmin/src/utils/_log"
	"gorm.io/gorm"
	"reflect"
)

type AdminConfig struct {
	TableNames []string
}

type TableStruct struct {
	Name    string
	Columns []ColumnStruct
}
type ColumnStruct struct {
	Name   string
	Desc   []string
	Type   reflect.Type
	DBType string
}

type OimoAdminStruct struct {
	Logger *_log.LoggerStruct
	DB     *gorm.DB
	Router *gin.Engine

	DBName        string
	Tables        []TableStruct
	RawSqlRecords []string
}
type ConfigStruct struct {
	string `json:"tg_bot_token"`
}

var ADMIN = "ADMIN"

type PaginationRequestStruct struct {
	Page    uint `json:"page"`
	PerPage uint `json:"perPage"`
}
