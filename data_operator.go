package gorest2

import (
	"database/sql"
	"strings"
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

type DefaultDataOperator struct {
	DataInterceptorRegistry       map[string]DataInterceptor
	GlobalDataInterceptorRegistry []DataInterceptor
}

func NewDbo(dataInterceptorRegistry map[string]DataInterceptor, globalDataInterceptorRegistry []DataInterceptor) DataOperator {
	ret := &DefaultDataOperator{
		DataInterceptorRegistry:       dataInterceptorRegistry,
		GlobalDataInterceptorRegistry: globalDataInterceptorRegistry,
	}
	return ret
}

func (this *DefaultDataOperator) RegisterDataInterceptor(id string, dataInterceptor DataInterceptor) {
	this.DataInterceptorRegistry[strings.ToUpper(id)] = dataInterceptor
}

func (this *DefaultDataOperator) GetDataInterceptor(id string) DataInterceptor {
	return this.DataInterceptorRegistry[strings.ToUpper(strings.Replace(id, "`", "", -1))]
}

func (this *DefaultDataOperator) RegisterGlobalDataInterceptor(globalDataInterceptor DataInterceptor) {
	this.GlobalDataInterceptorRegistry = append(this.GlobalDataInterceptorRegistry, globalDataInterceptor)
}

func (this *DefaultDataOperator) Load(resourceId string, id string, fields string, context map[string]interface{}) (map[string]string, error) {
	return nil, nil
}
func (this *DefaultDataOperator) ListMap(resourceId string, fields string, filter []string, sort string, group string, start int64, limit int64, includeTotal bool, context map[string]interface{}) ([]map[string]string, int64, error) {
	return nil, -1, nil
}
func (this *DefaultDataOperator) ListArray(resourceId string, fields string, filter []string, sort string, group string, start int64, limit int64, includeTotal bool, context map[string]interface{}) ([]string, [][]string, int64, error) {
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
func (this *DefaultDataOperator) QueryMap(resourceId string, start int64, limit int64, includeTotal bool, context map[string]interface{}) ([]map[string]string, int64, error) {
	return nil, -1, nil
}
func (this *DefaultDataOperator) QueryArray(resourceId string, start int64, limit int64, includeTotal bool, context map[string]interface{}) ([]string, [][]string, int64, error) {
	return nil, nil, -1, nil
}
func (this *DefaultDataOperator) Exec(resourceId string, context map[string]interface{}) (int64, error) {
	return -1, nil
}
func (this *DefaultDataOperator) GetConn() (*sql.DB, error) {
	return nil, nil
}
