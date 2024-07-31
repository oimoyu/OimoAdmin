package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/oimoyu/OimoAdmin/src/utils/_type"
	"github.com/oimoyu/OimoAdmin/src/utils/g"
	"github.com/oimoyu/OimoAdmin/src/utils/restful"
)

type redisItemStruct struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
	Type  string      `json:"type"`
	TTL   int64       `json:"ttl"`
}

func CreateRedisRow(c *gin.Context) {
	restful.ParamErr(c, "This function is not completed yet")
}
func DeleteRedisRows(c *gin.Context) {
	var fetchRequestData struct {
		Pattern string `json:"pattern"`
	}
	if err := c.ShouldBindBodyWith(&fetchRequestData, binding.JSON); err != nil {
		restful.ParamErr(c, fmt.Sprintf("parse input params failed: %v", err))
		return
	}
	if fetchRequestData.Pattern == "" {
		restful.ParamErr(c, fmt.Sprintf("Pattern is empty"))
		return
	}

	// TODO: bad performance for large data
	keys, err := g.OimoAdmin.RDB.Keys(context.Background(), fetchRequestData.Pattern).Result()
	if err != nil {
		restful.ParamErr(c, err.Error())
		return
	}

	num := 0
	for _, key := range keys {
		err := g.OimoAdmin.RDB.Del(context.Background(), key).Err()
		if err != nil {
			restful.ParamErr(c, err.Error())
			return
		}
		num += 1
	}

	msg := fmt.Sprintf("Number of keys deleted: %d", num)
	restful.Ok(c, msg)
}
func UpdateRedisRows(c *gin.Context) {
	restful.ParamErr(c, "This function is not completed yet")
}
func FetchRedisList(c *gin.Context) {
	// ShouldBindBodyWith will not clear data after read
	var paginationRequest _type.PaginationRequestStruct
	if err := c.ShouldBindBodyWith(&paginationRequest, binding.JSON); err != nil {
		restful.ParamErr(c, fmt.Sprintf("input pagination error: %v", err))
		return
	}
	offset := (paginationRequest.Page - 1) * paginationRequest.PerPage

	var fetchRequestData struct {
		Pattern string `json:"pattern"`
	}
	if err := c.ShouldBindBodyWith(&fetchRequestData, binding.JSON); err != nil {
		restful.ParamErr(c, fmt.Sprintf("parse input params failed: %v", err))
		return
	}

	redisKeyPattern := "*"
	if fetchRequestData.Pattern != "" {
		redisKeyPattern = fetchRequestData.Pattern
	}

	rdb := g.OimoAdmin.RDB
	keys, _, err := rdb.Scan(context.Background(), offset, redisKeyPattern, int64(paginationRequest.Page)).Result()
	if err != nil {
		fmt.Println("get key failed: ", err)
	}

	items := make([]redisItemStruct, 0)
	// 遍历每个键名
	for _, key := range keys {
		keyType, err := rdb.Type(context.Background(), key).Result()
		if err != nil {
			g.OimoAdmin.Logger.Error("Failed to get type for key %s: %v\n", key, err)
			continue
		}

		// 获取过期时间（TTL）
		ttl, err := rdb.TTL(context.Background(), key).Result()
		if err != nil {
			g.OimoAdmin.Logger.Error("Failed to get TTL for key %s: %v\n", key, err)
			continue
		}

		var value interface{}
		redisVarType := keyType
		failedParseMsg := "Failed to parse value"

		// 根据键的类型执行不同的操作
		switch keyType {
		case "string":
			val, err := rdb.Get(context.Background(), key).Result()
			if err != nil {
				g.OimoAdmin.Logger.Error("Failed to get value for string key %s: %v\n", key, err)
				value = failedParseMsg
			} else {
				value = val
			}
		case "list":
			vals, err := rdb.LRange(context.Background(), key, 0, -1).Result()
			if err != nil {
				g.OimoAdmin.Logger.Error("Failed to get list values for key %s: %v\n", key, err)
				value = failedParseMsg
			} else {
				value = vals
			}
		case "set":
			vals, err := rdb.SMembers(context.Background(), key).Result()
			if err != nil {
				g.OimoAdmin.Logger.Error("Failed to get set members for key %s: %v\n", key, err)
				value = failedParseMsg
			} else {
				value = vals
			}
		case "zset":
			vals, err := rdb.ZRangeWithScores(context.Background(), key, 0, -1).Result()
			if err != nil {
				g.OimoAdmin.Logger.Error("Failed to get zset values for key %s: %v\n", key, err)
				value = failedParseMsg
			} else {
				value = vals
			}
		case "hash":
			fields, err := rdb.HGetAll(context.Background(), key).Result()
			if err != nil {
				g.OimoAdmin.Logger.Error("Failed to get hash fields for key %s: %v\n", key, err)
				value = failedParseMsg
			} else {
				value = fields
			}
		default:
			//fmt.Printf("Unsupported key type %s for key %s\n", keyType, key)
		}
		items = append(items, redisItemStruct{Key: key, Value: value, Type: redisVarType, TTL: int64(ttl.Seconds())})
	}

	var itemMap interface{}
	data, err := json.Marshal(items) // Convert to a json string
	if err != nil {
		restful.ParamErr(c, "解析json错误")
		return
	}
	err = json.Unmarshal(data, &itemMap) // Convert to a map

	returnData := map[string]interface{}{
		"items": itemMap,
	}

	restful.Ok(c, returnData)
}

//func getRedisKeys(requestStruct _type.PaginationRequestStruct) {}
