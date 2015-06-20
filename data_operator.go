package gorest2

import (
	"database/sql"
)

type DataOperator interface {
	Load(resourceId string, id string, fields string, context map[string]interface{}) (map[string]string, error)
	ListMap(resourceId string, fields string, filter []string, sort string, group string, start int64, limit int64, includeTotal bool, context map[string]interface{}) ([]map[string]string, int64, error)
	ListArray(resourceId string, fields string, filter []string, sort string, group string, start int64, limit int64, includeTotal bool, context map[string]interface{}) ([]string, [][]string, int64, error)
	Create(resourceId string, data map[string]interface{}, context map[string]interface{}) (interface{}, error)
	Update(resourceId string, data map[string]interface{}, context map[string]interface{}) (int64, error)
	Duplicate(resourceId string, id string, context map[string]interface{}) (interface{}, error)
	Delete(resourceId string, id string, context map[string]interface{}) (int64, error)
	MetaData(resourceId string) ([]map[string]string, error)
	QueryMap(resourceId string, start int64, limit int64, includeTotal bool, context map[string]interface{}) ([]map[string]string, int64, error)
	QueryArray(resourceId string, start int64, limit int64, includeTotal bool, context map[string]interface{}) ([]string, [][]string, int64, error)
	Exec(resourceId string, context map[string]interface{}) (int64, error)
	GetConn() (*sql.DB, error)
}
