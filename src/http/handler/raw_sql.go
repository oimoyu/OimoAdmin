package handler

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/oimoyu/OimoAdmin/src/utils/g"
	"github.com/oimoyu/OimoAdmin/src/utils/restful"
	"gorm.io/gorm"
)

func addSqlRecord(str string, strSlice *[]string) {
	// not consider sync lock
	maxLength := 10
	*strSlice = append(*strSlice, str)
	if len(*strSlice) > maxLength {
		*strSlice = (*strSlice)[1:]
	}
}

func ExecuteRawSql(c *gin.Context) {
	var requestData struct {
		RawSql   string `json:"raw_sql" binding:"required"`
		DryRun   bool   `json:"dry_run"`
		IsSelect bool   `json:"is_select"` // don't use required,false is default zero value, seen as err
	}
	if err := c.ShouldBindBodyWith(&requestData, binding.JSON); err != nil {
		restful.ParamErr(c, fmt.Sprintf("parse input params failed: %v", err))
		return
	}
	if requestData.DryRun {
		var rowsAffected int64
		var ErrDryRunRollback = errors.New("dry run, manually rolled back")

		err := g.OimoAdmin.DB.Transaction(func(tx *gorm.DB) error {
			result := tx.Exec(requestData.RawSql)
			if result.Error != nil {
				return result.Error
			}
			rowsAffected = result.RowsAffected

			tx.Rollback() // manually rollback
			return ErrDryRunRollback
		})

		if err != nil && !errors.Is(err, ErrDryRunRollback) {
			errMsg := fmt.Sprintf("failed to dry run: %v", err)
			restful.ParamErr(c, errMsg)
			return
		}

		returnData := map[string]interface{}{
			"affect_rows": rowsAffected,
		}
		restful.Ok(c, returnData)
	} else {
		var sqlReturn []map[string]interface{}
		var result *gorm.DB
		if requestData.IsSelect {
			result = g.OimoAdmin.DB.Raw(requestData.RawSql).Scan(&sqlReturn)
		} else {
			result = g.OimoAdmin.DB.Exec(requestData.RawSql)
		}

		if result.Error != nil {
			errMsg := fmt.Sprintf("failed to execute: %v", result.Error)
			restful.ParamErr(c, errMsg)
			return
		}
		returnData := map[string]interface{}{
			"affect_rows": result.RowsAffected,
			"sql_return":  sqlReturn,
		}
		g.OimoAdmin.Logger.FileLog(fmt.Sprintf("%s: Execute raw sql: %s", c.ClientIP(), requestData.RawSql))

		addSqlRecord(requestData.RawSql, &g.OimoAdmin.RawSqlRecords)

		restful.Ok(c, fmt.Sprintf("Execute Success, Affected Rows: %d", result.RowsAffected), returnData)
	}

}

func RawSqlRecords(c *gin.Context) {
	returnData := map[string]interface{}{
		"raw_sql_records": g.OimoAdmin.RawSqlRecords,
	}
	restful.Ok(c, returnData)
}
