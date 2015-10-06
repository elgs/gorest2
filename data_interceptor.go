package gorest2

import (
	"database/sql"
	"net/http"
	"strings"
)

var DataInterceptorRegistry = map[string]DataInterceptor{}
var GlobalDataInterceptorRegistry = []DataInterceptor{}

var GlobalHandlerInterceptorRegistry = []HandlerInterceptor{}
var HandlerInterceptorRegistry = map[string]HandlerInterceptor{}

func RegisterDataInterceptor(id string, dataInterceptor DataInterceptor) {
	DataInterceptorRegistry[strings.ToUpper(id)] = dataInterceptor
}

func GetDataInterceptor(id string) DataInterceptor {
	return DataInterceptorRegistry[strings.ToUpper(strings.Replace(id, "`", "", -1))]
}

func RegisterGlobalDataInterceptor(globalDataInterceptor DataInterceptor) {
	GlobalDataInterceptorRegistry = append(GlobalDataInterceptorRegistry, globalDataInterceptor)
}

type DataInterceptor interface {
	BeforeLoad(resourceId string, db *sql.DB, fields string, context map[string]interface{}, id string) (bool, error)
	AfterLoad(resourceId string, db *sql.DB, fields string, context map[string]interface{}, data map[string]string) error
	BeforeCreate(resourceId string, db *sql.DB, context map[string]interface{}, data map[string]interface{}) (bool, error)
	AfterCreate(resourceId string, db *sql.DB, context map[string]interface{}, data map[string]interface{}) error
	BeforeUpdate(resourceId string, db *sql.DB, context map[string]interface{}, data map[string]interface{}) (bool, error)
	AfterUpdate(resourceId string, db *sql.DB, context map[string]interface{}, data map[string]interface{}) error
	BeforeDuplicate(resourceId string, db *sql.DB, context map[string]interface{}, id string) (bool, error)
	AfterDuplicate(resourceId string, db *sql.DB, context map[string]interface{}, oldId string, newId string) error
	BeforeDelete(resourceId string, db *sql.DB, context map[string]interface{}, id string) (bool, error)
	AfterDelete(resourceId string, db *sql.DB, context map[string]interface{}, id string) error
	BeforeListMap(resourceId string, db *sql.DB, fields string, context map[string]interface{}, filter *string, sort *string, group *string, start int64, limit int64) (bool, error)
	AfterListMap(resourceId string, db *sql.DB, fields string, context map[string]interface{}, data []map[string]string, total int64) error
	BeforeListArray(resourceId string, db *sql.DB, fields string, context map[string]interface{}, filter *string, sort *string, group *string, start int64, limit int64) (bool, error)
	AfterListArray(resourceId string, db *sql.DB, fields string, context map[string]interface{}, headers []string, data [][]string, total int64) error
	BeforeQueryMap(resourceId string, script string, params []interface{}, db *sql.DB, context map[string]interface{}) (bool, error)
	AfterQueryMap(resourceId string, script string, params []interface{}, db *sql.DB, context map[string]interface{}, data []map[string]string) error
	BeforeQueryArray(resourceId string, script string, params []interface{}, db *sql.DB, context map[string]interface{}) (bool, error)
	AfterQueryArray(resourceId string, script string, params []interface{}, db *sql.DB, context map[string]interface{}, headers []string, data [][]string) error
	BeforeExec(resourceId string, scripts string, params []interface{}, tx *sql.Tx, context map[string]interface{}) (bool, error)
	AfterExec(resourceId string, scripts string, params []interface{}, tx *sql.Tx, context map[string]interface{}) error
}

type HandlerInterceptor interface {
	BeforeHandle(w http.ResponseWriter, r *http.Request) (bool, error)
	AfterHandle(w http.ResponseWriter, r *http.Request) error
}

type DefaultDataInterceptor struct{}
type DefaultHandlerInterceptor struct{}

func (this *DefaultDataInterceptor) BeforeLoad(resourceId string, db *sql.DB, fields string, context map[string]interface{}, id string) (bool, error) {
	return true, nil
}
func (this *DefaultDataInterceptor) AfterLoad(resourceId string, db *sql.DB, fields string, context map[string]interface{}, data map[string]string) error {
	return nil
}
func (this *DefaultDataInterceptor) BeforeCreate(resourceId string, db *sql.DB, context map[string]interface{}, data map[string]interface{}) (bool, error) {
	return true, nil
}
func (this *DefaultDataInterceptor) AfterCreate(resourceId string, db *sql.DB, context map[string]interface{}, data map[string]interface{}) error {
	return nil
}
func (this *DefaultDataInterceptor) BeforeUpdate(resourceId string, db *sql.DB, context map[string]interface{}, data map[string]interface{}) (bool, error) {
	return true, nil
}
func (this *DefaultDataInterceptor) AfterUpdate(resourceId string, db *sql.DB, context map[string]interface{}, data map[string]interface{}) error {
	return nil
}
func (this *DefaultDataInterceptor) BeforeDuplicate(resourceId string, db *sql.DB, context map[string]interface{}, id string) (bool, error) {
	return true, nil
}
func (this *DefaultDataInterceptor) AfterDuplicate(resourceId string, db *sql.DB, context map[string]interface{}, oldId string, newId string) error {
	return nil
}
func (this *DefaultDataInterceptor) BeforeDelete(resourceId string, db *sql.DB, context map[string]interface{}, id string) (bool, error) {
	return true, nil
}
func (this *DefaultDataInterceptor) AfterDelete(resourceId string, db *sql.DB, context map[string]interface{}, id string) error {
	return nil
}
func (this *DefaultDataInterceptor) BeforeListMap(resourceId string, db *sql.DB, fields string, context map[string]interface{}, filter *string, sort *string, group *string, start int64, limit int64) (bool, error) {
	return true, nil
}
func (this *DefaultDataInterceptor) AfterListMap(resourceId string, db *sql.DB, fields string, context map[string]interface{}, data []map[string]string, total int64) error {
	return nil
}
func (this *DefaultDataInterceptor) BeforeListArray(resourceId string, db *sql.DB, fields string, context map[string]interface{}, filter *string, sort *string, group *string, start int64, limit int64) (bool, error) {
	return true, nil
}
func (this *DefaultDataInterceptor) AfterListArray(resourceId string, db *sql.DB, fields string, context map[string]interface{}, headers []string, data [][]string, total int64) error {
	return nil
}
func (this *DefaultDataInterceptor) BeforeQueryMap(resourceId string, script string, params []interface{}, db *sql.DB, context map[string]interface{}) (bool, error) {
	return true, nil
}
func (this *DefaultDataInterceptor) AfterQueryMap(resourceId string, script string, params []interface{}, db *sql.DB, context map[string]interface{}, data []map[string]string) error {
	return nil
}
func (this *DefaultDataInterceptor) BeforeQueryArray(resourceId string, script string, params []interface{}, db *sql.DB, context map[string]interface{}) (bool, error) {
	return true, nil
}
func (this *DefaultDataInterceptor) AfterQueryArray(resourceId string, script string, params []interface{}, db *sql.DB, context map[string]interface{}, headers []string, data [][]string) error {
	return nil
}
func (this *DefaultDataInterceptor) BeforeExec(resourceId string, scripts string, params []interface{}, tx *sql.Tx, context map[string]interface{}) (bool, error) {
	return true, nil
}
func (this *DefaultDataInterceptor) AfterExec(resourceId string, scripts string, params []interface{}, tx *sql.Tx, context map[string]interface{}) error {
	return nil
}

func (this *DefaultHandlerInterceptor) BeforeHandle(w http.ResponseWriter, r *http.Request) (bool, error) {
	return true, nil
}
func (this *DefaultHandlerInterceptor) AfterHandle(w http.ResponseWriter, r *http.Request) error {
	return nil
}
