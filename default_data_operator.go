// default_data_operator
package gorest2

import (
	"database/sql"
	"strings"
)

var DataInterceptorRegistry = map[string]DataInterceptor{}
var GlobalDataInterceptorRegistry = []DataInterceptor{}

func RegisterDataInterceptor(id string, dataInterceptor DataInterceptor) {
	DataInterceptorRegistry[strings.ToUpper(id)] = dataInterceptor
}

func GetDataInterceptor(id string) DataInterceptor {
	return DataInterceptorRegistry[strings.ToUpper(strings.Replace(id, "`", "", -1))]
}

func RegisterGlobalDataInterceptor(globalDataInterceptor DataInterceptor) {
	GlobalDataInterceptorRegistry = append(GlobalDataInterceptorRegistry, globalDataInterceptor)
}

type DefaultDataOperator struct {
}

func (this *DefaultDataOperator) Load(resourceId string, id string, fields string, context map[string]interface{}) (map[string]string, error) {
	return nil, nil
}
func (this *DefaultDataOperator) ListMap(resourceId string, fields string, filter []string, sort string, group string, start int64, limit int64, context map[string]interface{}) ([]map[string]string, int64, error) {
	return nil, -1, nil
}
func (this *DefaultDataOperator) ListArray(resourceId string, fields string, filter []string, sort string, group string, start int64, limit int64, context map[string]interface{}) ([]string, [][]string, int64, error) {
	return nil, nil, -1, nil
}
func (this *DefaultDataOperator) Create(resourceId string, data map[string]interface{}, context map[string]interface{}) (interface{}, error) {
	return nil, nil
}
func (this *DefaultDataOperator) Update(resourceId string, data map[string]interface{}, context map[string]interface{}) (int64, error) {
	return -1, nil
}
func (this *DefaultDataOperator) Duplicate(resourceId string, id string, context map[string]interface{}) (interface{}, error) {
	return nil, nil
}
func (this *DefaultDataOperator) Delete(resourceId string, id string, context map[string]interface{}) (int64, error) {
	return -1, nil
}
func (this *DefaultDataOperator) MetaData(resourceId string) ([]map[string]string, error) {
	return nil, nil
}
func (this *DefaultDataOperator) QueryMap(resourceId string, params []interface{}, context map[string]interface{}) ([]map[string]string, error) {
	return nil, nil
}
func (this *DefaultDataOperator) QueryArray(resourceId string, params []interface{}, context map[string]interface{}) ([]string, [][]string, error) {
	return nil, nil, nil
}
func (this *DefaultDataOperator) Exec(resourceId string, params []interface{}, context map[string]interface{}) (int64, error) {
	return -1, nil
}
func (this *DefaultDataOperator) GetConn() (*sql.DB, error) {
	return nil, nil
}
