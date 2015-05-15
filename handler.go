// handler
package gorest2

import (
	"net/http"
)

var handlerRegistry = make(map[string]func(dbo DataOperator) func(w http.ResponseWriter, r *http.Request))

func RegisterHandler(id string, handler func(dbo DataOperator) func(w http.ResponseWriter, r *http.Request)) {
	handlerRegistry[id] = handler
}

func GetHandler(id string) func(dbo DataOperator) func(w http.ResponseWriter, r *http.Request) {
	return handlerRegistry[id]
}

var dboRegistry = make(map[string]DataOperator)

func RegisterDbo(id string, dbo DataOperator) {
	dboRegistry[id] = dbo
}

func GetDbo(id string) DataOperator {
	return dboRegistry[id]
}
