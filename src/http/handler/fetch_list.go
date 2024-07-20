package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/oimoyu/OimoAdmin/src/utils/_type"
	"github.com/oimoyu/OimoAdmin/src/utils/g"
	"github.com/oimoyu/OimoAdmin/src/utils/restful"
	"strings"
)

func FetchList(c *gin.Context) {
	// ShouldBindBodyWith will not clear data after read
	var paginationRequest _type.PaginationRequestStruct
	if err := c.ShouldBindBodyWith(&paginationRequest, binding.JSON); err != nil {
		restful.ParamErr(c, fmt.Sprintf("input pagination error: %v", err))
		return
	}

	var fetchRequestData struct {
		TableName string `json:"table_name" binding:"required"`
		OrderBy   string `json:"orderBy"`
		OrderDir  string `json:"orderDir"`
		Keyword   string `json:"keyword"`
	}
	if err := c.ShouldBindBodyWith(&fetchRequestData, binding.JSON); err != nil {
		restful.ParamErr(c, fmt.Sprintf("parse input params failed: %v", err))
		return
	}

	items := make([]map[string]interface{}, 0)
	db := g.OimoAdmin.DB
	query := db.Table(fetchRequestData.TableName)

	// pagination
	offset := (paginationRequest.Page - 1) * paginationRequest.PerPage
	query = query.Limit(int(paginationRequest.PerPage)).Offset(int(offset))

	// order
	if fetchRequestData.OrderBy != "" {
		orderDir := "desc"
		if fetchRequestData.OrderDir == "asc" || fetchRequestData.OrderDir == "desc" {
			orderDir = fetchRequestData.OrderDir
		}
		query = query.Order(fmt.Sprintf("%s %s", fetchRequestData.OrderBy, orderDir))
	}

	// keyword
	if fetchRequestData.Keyword != "" {
		// get table object
		var currentTable *_type.TableStruct
		for i := range g.OimoAdmin.Tables {
			if g.OimoAdmin.Tables[i].Name == fetchRequestData.TableName {
				currentTable = &g.OimoAdmin.Tables[i] // when using a loop var, attention to pointer
			}
		}
		if currentTable == nil {
			restful.ParamErr(c, fmt.Sprintf("current table not exist: %s", fetchRequestData.TableName))
			return
		}

		// get column name
		columnNames := make([]string, 0)
		for _, column := range currentTable.Columns {
			columnNames = append(columnNames, column.Name)
		}

		keyword := "%" + fetchRequestData.Keyword + "%"
		var conditions []string
		var args []interface{}
		for _, column := range columnNames {
			conditions = append(conditions, fmt.Sprintf("%s LIKE ?", column)) // sql concat, but the columnNames is trusted
			args = append(args, keyword)
		}
		query = query.Where(strings.Join(conditions, " OR "), args...)

	}

	// copy a query for count
	countQuery := *query

	result := query.Find(&items)
	if result.Error != nil {
		errMsg := fmt.Sprintf("failed to get items: %v", result.Error)
		g.OimoAdmin.Logger.Error(errMsg)
		restful.ParamErr(c, errMsg)
		return
	}

	var total int64
	result = countQuery.Count(&total)
	if result.Error != nil {
		errMsg := fmt.Sprintf("failed to count items: %v", result.Error)
		g.OimoAdmin.Logger.Error(errMsg)
		restful.ParamErr(c, errMsg)
		return
	}

	returnData := map[string]interface{}{
		"items": items,
		"total": total,
	}
	restful.Ok(c, returnData)
}
