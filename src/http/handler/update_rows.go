package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/oimoyu/OimoAdmin/src/utils/g"
	"github.com/oimoyu/OimoAdmin/src/utils/restful"
	"gorm.io/gorm"
)

func UpdateRows(c *gin.Context) {
	var requestData struct {
		RowsDiff  []map[string]interface{} `json:"rowsDiff" binding:"required"`
		TableName string                   `json:"table_name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&requestData); err != nil {
		restful.ParamErr(c, fmt.Sprintf("parse input params failed: %v", err))
		return
	}

	updateData := make(map[interface{}]map[string]interface{})
	for _, rowDiff := range requestData.RowsDiff {
		if _, exists := rowDiff["id"]; !exists {
			restful.ParamErr(c, fmt.Sprintf("no id field in request data, please check table setting: table should have id field and make sure it is primary key"))
			return
		}

		rowCopy := make(map[string]interface{})
		for k, v := range rowDiff {
			if k != "id" {
				rowCopy[k] = v
			}
		}
		updateData[rowDiff["id"]] = rowCopy
	}

	err := g.OimoAdmin.DB.Transaction(func(tx *gorm.DB) error {
		for id, updateData := range updateData {
			if err := tx.Table(requestData.TableName).Where("id = ?", id).Updates(updateData).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		restful.ParamErr(c, fmt.Sprintf("failed to update rows: %v", err))
		return
	}

	g.OimoAdmin.Logger.FileLog(fmt.Sprintf("%s: update rows, table name: %s, update data: %v", c.ClientIP(), requestData.TableName, updateData))

	restful.Ok(c, "Update Success")
}
