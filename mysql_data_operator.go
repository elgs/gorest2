package gorest

import (
	"bytes"
	"code.google.com/p/go-uuid/uuid"
	"database/sql"
	"errors"
	"fmt"
	"github.com/elgs/exparser"
	"github.com/elgs/gosqljson"
	"strconv"
	"strings"
)

type MySqlDataOperator struct {
	*DefaultDataOperator
	Ds         string
	DbType     string
	TokenTable string
}

func (this *MySqlDataOperator) Load(tableId string, id string, context map[string]interface{}) (map[string]string, error) {
	ret := make(map[string]string, 0)
	tableId = normalizeTableId(tableId, this.DbType, this.Ds)
	context["token_table"] = this.TokenTable
	db, err := getConn(this.DbType, this.Ds)
	defer db.Close()
	if err != nil {
		fmt.Println(err)
		return ret, err
	}

	for _, globalDataInterceptor := range GlobalDataInterceptorRegistry {
		ctn, err := globalDataInterceptor.BeforeLoad(tableId, db, context, id)
		if !ctn {
			return ret, err
		}
	}
	dataInterceptor := GetDataInterceptor(tableId)
	if dataInterceptor != nil {
		ctn, err := dataInterceptor.BeforeLoad(tableId, db, context, id)
		if !ctn {
			return ret, err
		}
	}

	// Load the record
	extraFilter := context["extra_filter"]
	if extraFilter == nil {
		extraFilter = ""
	}
	m, err := gosqljson.QueryDbToMap(db, true,
		fmt.Sprint("SELECT * FROM ", tableId, " WHERE ID=? ", extraFilter), id)
	if err != nil {
		fmt.Println(err)
		return ret, err
	}

	if len(m) == 0 {
		m = []map[string]string{
			make(map[string]string, 0),
		}
	}

	if dataInterceptor != nil {
		dataInterceptor.AfterLoad(tableId, db, context, m[0])
	}
	for _, globalDataInterceptor := range GlobalDataInterceptorRegistry {
		globalDataInterceptor.AfterLoad(tableId, db, context, m[0])
	}

	if m != nil && len(m) == 1 {
		return m[0], err
	} else {
		return ret, err
	}

}
func (this *MySqlDataOperator) ListMap(tableId string, filter []string, sort string,
	start int64, limit int64, includeTotal bool, context map[string]interface{}) ([]map[string]string, int64, error) {
	ret := make([]map[string]string, 0)
	tableId = normalizeTableId(tableId, this.DbType, this.Ds)
	context["token_table"] = this.TokenTable
	db, err := getConn(this.DbType, this.Ds)
	defer db.Close()
	if err != nil {
		fmt.Println(err)
		return ret, -1, err
	}

	sort = parseSort(sort)
	where := parseFilters(filter)
	for _, globalDataInterceptor := range GlobalDataInterceptorRegistry {
		ctn, err := globalDataInterceptor.BeforeListMap(tableId, db, context, &where, &sort, start, limit, includeTotal)
		if !ctn {
			return ret, -1, err
		}
	}
	dataInterceptor := GetDataInterceptor(tableId)
	if dataInterceptor != nil {
		ctn, err := dataInterceptor.BeforeListMap(tableId, db, context, &where, &sort, start, limit, includeTotal)
		if !ctn {
			return ret, -1, err
		}
	}

	m, err := gosqljson.QueryDbToMap(db, true,
		fmt.Sprint("SELECT * FROM ", tableId, where, sort, " LIMIT ?,?"), start, limit)
	if err != nil {
		fmt.Println(err)
		return ret, -1, err
	}
	cnt := -1
	if includeTotal {
		c, err := gosqljson.QueryDbToMap(db, false,
			fmt.Sprint("SELECT COUNT(*) AS CNT FROM ", tableId, where))
		if err != nil {
			fmt.Println(err)
			return ret, -1, err
		}
		cnt, err = strconv.Atoi(c[0]["CNT"])
		if err != nil {
			fmt.Println(err)
			return ret, -1, err
		}
	}

	if dataInterceptor != nil {
		dataInterceptor.AfterListMap(tableId, db, context, m, int64(cnt))
	}
	for _, globalDataInterceptor := range GlobalDataInterceptorRegistry {
		globalDataInterceptor.AfterListMap(tableId, db, context, m, int64(cnt))
	}

	return m, int64(cnt), err
}
func (this *MySqlDataOperator) ListArray(tableId string, filter []string, sort string,
	start int64, limit int64, includeTotal bool, context map[string]interface{}) ([][]string, int64, error) {
	ret := make([][]string, 0)
	tableId = normalizeTableId(tableId, this.DbType, this.Ds)
	context["token_table"] = this.TokenTable
	db, err := getConn(this.DbType, this.Ds)
	defer db.Close()
	if err != nil {
		fmt.Println(err)
		return ret, -1, err
	}

	sort = parseSort(sort)
	where := parseFilters(filter)
	for _, globalDataInterceptor := range GlobalDataInterceptorRegistry {
		ctn, err := globalDataInterceptor.BeforeListArray(tableId, db, context, &where, &sort, start, limit, includeTotal)
		if !ctn {
			return ret, -1, err
		}
	}
	dataInterceptor := GetDataInterceptor(tableId)
	if dataInterceptor != nil {
		ctn, err := dataInterceptor.BeforeListArray(tableId, db, context, &where, &sort, start, limit, includeTotal)
		if !ctn {
			return ret, -1, err
		}
	}

	a, err := gosqljson.QueryDbToArray(db, true,
		fmt.Sprint("SELECT * FROM ", tableId, where, sort, " LIMIT ?,?"), start, limit)
	if err != nil {
		fmt.Println(err)
		return ret, -1, err
	}
	cnt := -1
	if includeTotal {
		c, err := gosqljson.QueryDbToMap(db, false,
			fmt.Sprint("SELECT COUNT(*) AS CNT FROM ", tableId, where))
		if err != nil {
			fmt.Println(err)
			return ret, -1, err
		}
		cnt, err = strconv.Atoi(c[0]["CNT"])
		if err != nil {
			fmt.Println(err)
			return ret, -1, err
		}
	}

	if dataInterceptor != nil {
		dataInterceptor.AfterListArray(tableId, db, context, a, int64(cnt))
	}
	for _, globalDataInterceptor := range GlobalDataInterceptorRegistry {
		globalDataInterceptor.AfterListArray(tableId, db, context, a, int64(cnt))
	}

	return a, int64(cnt), err
}
func (this *MySqlDataOperator) Create(tableId string, data map[string]interface{}, context map[string]interface{}) (interface{}, error) {
	tableId = normalizeTableId(tableId, this.DbType, this.Ds)
	context["token_table"] = this.TokenTable
	db, err := getConn(this.DbType, this.Ds)
	defer db.Close()
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	for _, globalDataInterceptor := range GlobalDataInterceptorRegistry {
		ctn, err := globalDataInterceptor.BeforeCreate(tableId, db, context, data)
		if !ctn {
			return "", err
		}
	}
	dataInterceptor := GetDataInterceptor(tableId)
	if dataInterceptor != nil {
		ctn, err := dataInterceptor.BeforeCreate(tableId, db, context, data)
		if !ctn {
			return "", err
		}
	}

	// Create the record
	if data["ID"] == nil {
		data["ID"] = uuid.New()
	}
	dataLen := len(data)
	values := make([]interface{}, 0, dataLen)
	var fieldBuffer bytes.Buffer
	var qmBuffer bytes.Buffer
	count := 0
	for k, v := range data {
		count++
		if count == dataLen {
			fieldBuffer.WriteString(k)
			qmBuffer.WriteString("?")
		} else {
			fieldBuffer.WriteString(fmt.Sprint(k, ","))
			qmBuffer.WriteString("?,")
		}
		values = append(values, v)
	}
	fields := fieldBuffer.String()
	qms := qmBuffer.String()
	gosqljson.ExecDb(db, fmt.Sprint("INSERT INTO ", tableId, " (", fields, ") VALUES (", qms, ")"), values...)

	if dataInterceptor != nil {
		dataInterceptor.AfterCreate(tableId, db, context, data)
	}
	for _, globalDataInterceptor := range GlobalDataInterceptorRegistry {
		globalDataInterceptor.AfterCreate(tableId, db, context, data)
	}

	return data["ID"], err
}
func (this *MySqlDataOperator) Update(tableId string, data map[string]interface{}, context map[string]interface{}) (int64, error) {
	tableId = normalizeTableId(tableId, this.DbType, this.Ds)
	context["token_table"] = this.TokenTable
	db, err := getConn(this.DbType, this.Ds)
	defer db.Close()
	if err != nil {
		fmt.Println(err)
		return -1, err
	}

	for _, globalDataInterceptor := range GlobalDataInterceptorRegistry {
		ctn, err := globalDataInterceptor.BeforeUpdate(tableId, db, context, data)
		if !ctn {
			return 0, err
		}
	}
	dataInterceptor := GetDataInterceptor(tableId)
	if dataInterceptor != nil {
		ctn, err := dataInterceptor.BeforeUpdate(tableId, db, context, data)
		if !ctn {
			return 0, err
		}
	}
	// Update the record
	id := data["ID"]
	if id == nil {
		fmt.Println("ID is not found.")
		return 0, err
	}
	delete(data, "ID")
	dataLen := len(data)
	values := make([]interface{}, 0, dataLen)
	var buffer bytes.Buffer
	for k, v := range data {
		buffer.WriteString(fmt.Sprint(k, "=?,"))
		values = append(values, v)
	}
	values = append(values, id)
	sets := buffer.String()
	sets = sets[0 : len(sets)-1]
	rowsAffected, err := gosqljson.ExecDb(db, fmt.Sprint("UPDATE ", tableId, " SET ", sets, " WHERE ID=?"), values...)
	if err != nil {
		fmt.Println(err)
		return -1, err
	}

	if dataInterceptor != nil {
		dataInterceptor.AfterUpdate(tableId, db, context, data)
	}
	for _, globalDataInterceptor := range GlobalDataInterceptorRegistry {
		globalDataInterceptor.AfterUpdate(tableId, db, context, data)
	}

	return rowsAffected, err
}
func (this *MySqlDataOperator) Duplicate(tableId string, id string, context map[string]interface{}) (interface{}, error) {
	tableId = normalizeTableId(tableId, this.DbType, this.Ds)
	context["token_table"] = this.TokenTable
	db, err := getConn(this.DbType, this.Ds)
	defer db.Close()
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	for _, globalDataInterceptor := range GlobalDataInterceptorRegistry {
		ctn, err := globalDataInterceptor.BeforeDuplicate(tableId, db, context, id)
		if !ctn {
			return "", err
		}
	}
	dataInterceptor := GetDataInterceptor(tableId)
	if dataInterceptor != nil {
		ctn, err := dataInterceptor.BeforeDuplicate(tableId, db, context, id)
		if !ctn {
			return "", err
		}
	}

	// Duplicate the record
	data, err := gosqljson.QueryDbToMap(db, false,
		fmt.Sprint("SELECT * FROM ", tableId, " WHERE ID=?"), id)
	if data == nil || len(data) != 1 {
		return "", err
	}
	newData := make(map[string]interface{}, len(data[0]))
	for k, v := range data[0] {
		newData[k] = v
	}
	newId := uuid.New()
	newData["ID"] = newId

	newDataLen := len(newData)
	newValues := make([]interface{}, 0, newDataLen)
	var fieldBuffer bytes.Buffer
	var qmBuffer bytes.Buffer
	count := 0
	for k, v := range newData {
		count++
		if count == newDataLen {
			fieldBuffer.WriteString(k)
			qmBuffer.WriteString("?")
		} else {
			fieldBuffer.WriteString(fmt.Sprint(k, ","))
			qmBuffer.WriteString("?,")
		}
		newValues = append(newValues, v)
	}
	fields := fieldBuffer.String()
	qms := qmBuffer.String()
	gosqljson.ExecDb(db, fmt.Sprint("INSERT INTO ", tableId, " (", fields, ") VALUES (", qms, ")"), newValues...)

	if dataInterceptor != nil {
		dataInterceptor.AfterDuplicate(tableId, db, context, id, newId)
	}
	for _, globalDataInterceptor := range GlobalDataInterceptorRegistry {
		globalDataInterceptor.AfterDuplicate(tableId, db, context, id, newId)
	}

	return newId, err
}
func (this *MySqlDataOperator) Delete(tableId string, id string, context map[string]interface{}) (int64, error) {
	tableId = normalizeTableId(tableId, this.DbType, this.Ds)
	context["token_table"] = this.TokenTable
	db, err := getConn(this.DbType, this.Ds)
	defer db.Close()
	if err != nil {
		fmt.Println(err)
		return -1, err
	}

	for _, globalDataInterceptor := range GlobalDataInterceptorRegistry {
		ctn, err := globalDataInterceptor.BeforeDelete(tableId, db, context, tableId)
		if !ctn {
			return 0, err
		}
	}
	dataInterceptor := GetDataInterceptor(tableId)
	if dataInterceptor != nil {
		ctn, err := dataInterceptor.BeforeDelete(tableId, db, context, tableId)
		if !ctn {
			return 0, err
		}
	}
	// Delete the record
	rowsAffected, err := gosqljson.ExecDb(db, fmt.Sprint("DELETE FROM ", tableId, " WHERE ID=?"), id)
	if err != nil {
		fmt.Println(err)
		return -1, err
	}

	if dataInterceptor != nil {
		dataInterceptor.AfterDelete(tableId, db, context, tableId)
	}
	for _, globalDataInterceptor := range GlobalDataInterceptorRegistry {
		globalDataInterceptor.AfterDelete(tableId, db, context, tableId)
	}

	return rowsAffected, err
}
func (this *MySqlDataOperator) QueryMap(tableId string, sqlSelect string, sqlSelectCount string,
	start int64, limit int64, includeTotal bool, context map[string]interface{}) ([]map[string]string, int64, error) {
	ret := make([]map[string]string, 0)
	tableId = normalizeTableId(tableId, this.DbType, this.Ds)
	context["token_table"] = this.TokenTable
	if !isSelect(sqlSelect) {
		return ret, -1, nil
	}
	if includeTotal && !isSelect(sqlSelectCount) {
		return ret, -1, nil
	}
	db, err := getConn(this.DbType, this.Ds)
	defer db.Close()
	if err != nil {
		fmt.Println(err)
		return ret, -1, err
	}

	for _, globalDataInterceptor := range GlobalDataInterceptorRegistry {
		ctn, err := globalDataInterceptor.BeforeQueryMap(tableId, db, context, &sqlSelect, &sqlSelectCount, start, limit, includeTotal)
		if !ctn {
			return ret, -1, err
		}
	}
	dataInterceptor := GetDataInterceptor(tableId)
	if dataInterceptor != nil {
		ctn, err := dataInterceptor.BeforeQueryMap(tableId, db, context, &sqlSelect, &sqlSelectCount, start, limit, includeTotal)
		if !ctn {
			return ret, -1, err
		}
	}

	m, err := gosqljson.QueryDbToMap(db, true,
		fmt.Sprint(sqlSelect, " LIMIT ?,?"), start, limit)
	cnt := -1
	if err != nil {
		fmt.Println(err)
		return ret, -1, err
	}
	if includeTotal {
		c, err := gosqljson.QueryDbToMap(db, false, sqlSelectCount)
		if err != nil {
			fmt.Println(err)
			return ret, -1, err
		}
		for _, v := range c[0] {
			cnt, err = strconv.Atoi(v)
		}
		if err != nil {
			fmt.Println(err)
			return ret, -1, err
		}
	}

	if dataInterceptor != nil {
		dataInterceptor.AfterQueryMap(tableId, db, context, m, int64(cnt))
	}
	for _, globalDataInterceptor := range GlobalDataInterceptorRegistry {
		globalDataInterceptor.AfterQueryMap(tableId, db, context, m, int64(cnt))
	}

	return m, int64(cnt), err
}
func (this *MySqlDataOperator) QueryArray(tableId string, sqlSelect string, sqlSelectCount string,
	start int64, limit int64, includeTotal bool, context map[string]interface{}) ([][]string, int64, error) {
	ret := make([][]string, 0)
	tableId = normalizeTableId(tableId, this.DbType, this.Ds)
	context["token_table"] = this.TokenTable
	if !isSelect(sqlSelect) {
		return ret, -1, errors.New("Invalid query.")
	}
	if includeTotal && !isSelect(sqlSelectCount) {
		return ret, -1, errors.New("Invalid query.")
	}
	db, err := getConn(this.DbType, this.Ds)
	defer db.Close()
	if err != nil {
		fmt.Println(err)
		return ret, -1, err
	}

	for _, globalDataInterceptor := range GlobalDataInterceptorRegistry {
		ctn, err := globalDataInterceptor.BeforeQueryArray(tableId, db, context, &sqlSelect, &sqlSelectCount, start, limit, includeTotal)
		if !ctn {
			return ret, -1, err
		}
	}
	dataInterceptor := GetDataInterceptor(tableId)
	if dataInterceptor != nil {
		ctn, err := dataInterceptor.BeforeQueryArray(tableId, db, context, &sqlSelect, &sqlSelectCount, start, limit, includeTotal)
		if !ctn {
			return ret, -1, err
		}
	}

	a, err := gosqljson.QueryDbToArray(db, true,
		fmt.Sprint(sqlSelect, " LIMIT ?,?"), start, limit)
	if err != nil {
		fmt.Println(err)
		return ret, -1, err
	}
	cnt := -1
	if includeTotal {
		c, err := gosqljson.QueryDbToMap(db, false, sqlSelectCount)
		if err != nil {
			fmt.Println(err)
			return ret, -1, err
		}
		for _, v := range c[0] {
			cnt, err = strconv.Atoi(v)
		}
		if err != nil {
			fmt.Println(err)
			return ret, -1, err
		}
	}

	if dataInterceptor != nil {
		dataInterceptor.AfterQueryArray(tableId, db, context, a, int64(cnt))
	}
	for _, globalDataInterceptor := range GlobalDataInterceptorRegistry {
		globalDataInterceptor.AfterQueryArray(tableId, db, context, a, int64(cnt))
	}

	return a, int64(cnt), err
}

func isSelect(sqlSelect string) bool {
	return strings.HasPrefix(strings.ToUpper(sqlSelect), "SELECT ")
}

func getConn(dbType string, ds string) (*sql.DB, error) {
	if len(strings.TrimSpace(dbType)) == 0 {
		dbType = "mysql"
	}
	db, err := sql.Open(dbType, ds)
	return db, err
}

func extractDbNameFromDs(dbType string, ds string) string {
	switch dbType {
	case "sqlite3":
		return ""
	default:
		a := strings.LastIndex(ds, "/")
		b := ds[a+1:]
		c := strings.Index(b, "?")
		if c < 0 {
			return b
		}
		return b[:c]
	}
}

func normalizeTableId(tableId string, dbType string, ds string) string {
	if strings.Contains(tableId, ".") {
		a := strings.Split(tableId, ".")
		return fmt.Sprint(a[0], ".", a[1])
	}
	db := extractDbNameFromDs(dbType, ds)

	MysqlSafe(&tableId)
	if len(strings.TrimSpace(db)) == 0 {
		return tableId
	} else {
		MysqlSafe(&db)
		return fmt.Sprint(db, ".", tableId)
	}
}

func MysqlSafe(s *string) {
	*s = strings.Replace(*s, "'", "''", -1)
	*s = strings.Replace(*s, "--", "", -1)
}

func parseSort(sort string) string {
	if len(strings.TrimSpace(sort)) == 0 {
		return ""
	}
	return fmt.Sprint(" ORDER BY ", strings.ToUpper(strings.Replace(sort, ":", " ", -1)), " ")
}

func parseFilter(filter string) string {
	r, err := parser.Calculate(filter)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return r
}

func parseFilters(filters []string) (r string) {
	for _, v := range filters {
		r += fmt.Sprint("AND ", parseFilter(v))
	}
	r = fmt.Sprint(" WHERE 1=1 ", r)
	//fmt.Println(r)
	return
}

var parser = &exparser.Parser{
	Operators: exparser.MysqlOperators,
}
