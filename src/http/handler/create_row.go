package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/oimoyu/OimoAdmin/src/utils/g"
	"github.com/oimoyu/OimoAdmin/src/utils/restful"
)

func CreateRow(c *gin.Context) {
	var requestData map[string]interface{}
	if err := c.ShouldBindJSON(&requestData); err != nil {
		restful.ParamErr(c, fmt.Sprintf("parse input params failed: %v", err))
		return
	}

	_, ok := requestData["id"]
	if !ok {
		restful.ParamErr(c, "id field is required")
		return
	}

	tableName, ok := requestData["table_name"]
	if !ok {
		restful.ParamErr(c, "table_name field is required")
		return
	}
	tableNameString, ok := tableName.(string)
	if !ok {
		restful.ParamErr(c, "table_name is not string")
		return
	}

	createData := make(map[string]interface{})
	for key, value := range requestData {
		if key != "table_name" {
			createData[key] = value
		}
	}

	result := g.OimoAdmin.DB.Table(tableNameString).Create(createData)
	if result.Error != nil {
		errMsg := fmt.Sprintf("failed to create row: %v", result.Error)
		restful.ParamErr(c, errMsg)
		return
	}

	g.OimoAdmin.Logger.FileLog(fmt.Sprintf("%s: Create row, table: %s, create data: %v", c.ClientIP(), tableName, createData))

	restful.Ok(c, "Create Success")
}
