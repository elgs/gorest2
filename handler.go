// handler
package gorest2

import (
	"fmt"
	"github.com/elgs/gosqljson"
	"net/http"
)

var handlerRegistry = make(map[string]func(w http.ResponseWriter, r *http.Request))

func RegisterHandler(id string, handler func(w http.ResponseWriter, r *http.Request)) {
	handlerRegistry[id] = handler
}

func GetHandler(id string) func(w http.ResponseWriter, r *http.Request) {
	return handlerRegistry[id]
}

var dboRegistry = make(map[string]DataOperator)

func RegisterDbo(id string, dbo DataOperator) {
	dboRegistry[id] = dbo
}

func GetDbo(id string) DataOperator {
	return LoadDbo(id)
}

func LoadDbo(id string) DataOperator {
	ret := dboRegistry[id]
	if ret != nil {
		return ret
	}
	defaultDbo := GetDbo("default")
	db, err := defaultDbo.GetConn()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	query := `SELECT data_store.* FROM project
		INNER JOIN data_store ON project.DATA_STORE_NAME = data_store.DATA_STORE_KEY
		WHERE project.PROJECT_KEY=?`
	data, err := gosqljson.QueryDbToMap(db, query, id)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	if data == nil || len(data) == 0 {
		return nil
	}
	dboData := data[0]
	ret = &MySqlDataOperator{
		Ds: fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", dboData["USERNAME"], dboData["PASSWORD"],
			dboData["HOST"], dboData["PORT"], dboData["DB"]),
		DbType: "mysql",
	}
	RegisterDbo(id, ret)
	return ret
}
