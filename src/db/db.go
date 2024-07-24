package db

import (
	"fmt"
	"github.com/oimoyu/OimoAdmin/src/utils/_type"
	"gorm.io/gorm"
	"reflect"
	"strings"
)

func InitTables(db *gorm.DB, adminConfig _type.AdminConfig) []_type.TableStruct {
	var tableNames []string
	if len(adminConfig.TableNames) != 0 {
		tableNames = adminConfig.TableNames
	} else {
		var err error
		tableNames, err = db.Migrator().GetTables()
		if err != nil {
			panic(err)
		}
	}

	tables := make([]_type.TableStruct, 0)
	for _, tableName := range tableNames {
		gormColTypes, err := db.Migrator().ColumnTypes(tableName)
		if err != nil {
			panic(err)
		}
		columns := make([]_type.ColumnStruct, 0)
		for _, gormColType := range gormColTypes {
			colName := gormColType.Name()
			colTypeString := strings.ToLower(gormColType.DatabaseTypeName())
			//colType, err := TypeFromString(colTypeString)
			//if err != nil {
			//	panic(err)
			//}

			columns = append(columns, _type.ColumnStruct{
				Name: colName,
				//Type: colType,
				DBType: colTypeString,
			})
		}
		tables = append(tables, _type.TableStruct{
			Name:    tableName,
			Columns: columns,
		})

	}

	return tables
}

func TypeFromString(typeStr string) (reflect.Type, error) {
	typeStr = strings.ToLower(typeStr)

	switch typeStr {
	case "text", "char", "varchar", "tinytext", "mediumtext", "longtext":
		return reflect.TypeOf(""), nil
	case "integer", "int", "smallint", "mediumint", "bigint":
		return reflect.TypeOf(int64(0)), nil
	case "tinyint":
		return reflect.TypeOf(int8(0)), nil
	case "real", "float":
		return reflect.TypeOf(float32(0)), nil
	case "double", "decimal", "numeric":
		return reflect.TypeOf(float64(0)), nil
	case "boolean":
		return reflect.TypeOf(true), nil
	case "binary", "varbinary", "blob", "tinyblob", "mediumblob", "longblob", "bytea":
		return reflect.TypeOf([]byte{}), nil
	case "date", "datetime", "timestamp", "time":
		return reflect.TypeOf(""), nil
	case "int8":
		return reflect.TypeOf(int8(0)), nil
	case "int16":
		return reflect.TypeOf(int16(0)), nil
	case "int32":
		return reflect.TypeOf(int32(0)), nil
	case "int64":
		return reflect.TypeOf(int64(0)), nil
	case "float32":
		return reflect.TypeOf(float32(0)), nil
	case "float64":
		return reflect.TypeOf(float64(0)), nil
	default:
		return nil, fmt.Errorf("can't handle this type: %s", typeStr)
	}
}
