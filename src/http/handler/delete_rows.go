package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/oimoyu/OimoAdmin/src/utils/g"
	"github.com/oimoyu/OimoAdmin/src/utils/restful"
	"strings"
)

func DeleteRows(c *gin.Context) {
	var requestData struct {
		Ids       string `json:"ids" binding:"required"`
		TableName string `json:"table_name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&requestData); err != nil {
		restful.ParamErr(c, fmt.Sprintf("parse input params failed: %v", err))
		return
	}

	ids := strings.Split(requestData.Ids, ",")

	err := g.OimoAdmin.DB.Table(requestData.TableName).Where("id IN ?", ids).Delete(nil).Error
	if err != nil {
		restful.ParamErr(c, fmt.Sprintf("failed to delete rows: %v", err))
		return
	}

	g.OimoAdmin.Logger.FileLog(fmt.Sprintf("%s: delete rows, table: %s, ids: %v", c.ClientIP(), requestData.TableName, ids))

	restful.Ok(c, "Delete Success")
}
