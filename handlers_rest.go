// handlers_rest
package gorest2

import (
	"encoding/json"
	"fmt"
	"github.com/elgs/gosplitargs"
	"net/http"
	"strconv"
	"strings"
)

var translateBoolParam = func(field string, defaultValue bool) bool {
	if field == "1" {
		return true
	} else if field == "0" {
		return false
	} else {
		return defaultValue
	}
}

var RestFunc = func(w http.ResponseWriter, r *http.Request) {
	context := make(map[string]interface{})
	//	context["api_token_id"] = r.Header.Get("api_token_id")
	context["token"] = r.Header.Get("token")

	projectId := r.Header.Get("app_id")
	if projectId == "" {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		fmt.Fprint(w, `{"err":"Invalid project."}`)
		return
	}
	context["app_id"] = projectId
	dbo := GetDbo(projectId)
	if dbo == nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		fmt.Fprint(w, `{"err":"Invalid project."}`)
		return
	}

	urlPath := r.URL.Path
	urlPathData := strings.Split(urlPath[1:], "/")
	tableId := urlPathData[1]

	switch r.Method {
	case "GET":
		if len(urlPathData) == 2 || len(urlPathData[2]) == 0 {
			//List records.
			fields := strings.ToUpper(r.FormValue("fields"))
			sort := r.FormValue("sort")
			group := r.FormValue("group")
			s := r.FormValue("start")
			l := r.FormValue("limit")
			c := r.FormValue("case")
			p := r.FormValue("params")
			context["case"] = c
			filter := r.Form["filter"]
			array := translateBoolParam(r.FormValue("array"), false)
			query := translateBoolParam(r.FormValue("query"), false)
			start, err := strconv.ParseInt(s, 10, 0)
			if err != nil {
				start = 0
				err = nil
			}
			limit, err := strconv.ParseInt(l, 10, 0)
			if err != nil {
				limit = 25
				err = nil
			}
			if fields == "" {
				fields = "*"
			}
			params, err := gosplitargs.SplitArgs(p, ",", true)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Email already used."))
				return
			}
			parameters := make([]interface{}, len(params))
			for i, v := range params {
				parameters[i] = v
			}
			var data interface{}
			var total int64 = -1
			m := map[string]interface{}{}
			if array {
				var headers []string
				var dataArray [][]string
				if query {
					headers, dataArray, err = dbo.QueryArray(tableId, parameters, context)
					if err != nil {
						m["err"] = err.Error()
					} else {
						m["headers"] = headers
						m["data"] = dataArray
					}

				} else {
					headers, dataArray, total, err = dbo.ListArray(tableId, fields, filter, sort, group, start, limit, context)
					if err != nil {
						m["err"] = err.Error()
					} else {
						m["headers"] = headers
						m["data"] = dataArray
						m["total"] = total
					}
				}
			} else {
				if query {
					data, err = dbo.QueryMap(tableId, parameters, context)
					if err != nil {
						m["err"] = err.Error()
					} else {
						m["data"] = data
					}
				} else {
					data, total, err = dbo.ListMap(tableId, fields, filter, sort, group, start, limit, context)
					if err != nil {
						m["err"] = err.Error()
					} else {
						m["data"] = data
						m["total"] = total
					}
				}
			}
			jsonData, err := json.Marshal(m)
			jsonString := string(jsonData)
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			fmt.Fprint(w, jsonString)
		} else {
			// Load record by id.
			dataId := urlPathData[2]
			c := r.FormValue("case")
			context["case"] = c

			fields := strings.ToUpper(r.FormValue("fields"))
			if fields == "" {
				fields = "*"
			}

			data, err := dbo.Load(tableId, dataId, fields, context)

			m := map[string]interface{}{
				"data": data,
			}
			if err != nil {
				m["err"] = err.Error()
			}
			jsonData, _ := json.Marshal(m)
			jsonString := string(jsonData)
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			fmt.Fprint(w, jsonString)
		}
	case "POST":
		// Create the record.
		metaValues := r.URL.Query()["meta"]
		meta := true
		if metaValues != nil && metaValues[0] == "0" {
			meta = false
		}
		context["meta"] = meta

		execValues := r.URL.Query()["exec"]
		exec := false
		if execValues != nil && execValues[0] == "1" {
			exec = true
		}

		noParamValues := r.URL.Query()["noparam"]
		noParam := false
		if noParamValues != nil && noParamValues[0] == "1" {
			noParam = true
		}
		parameters := make([]interface{}, 0)
		if !noParam {
			p := ""
			if len(r.URL.Query()["params"]) > 0 {
				p = r.URL.Query()["params"][0]
			}
			params, err := gosplitargs.SplitArgs(p, ",", true)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
			for i, v := range params {
				parameters[i] = v
			}
		}

		m := make(map[string]interface{})
		if exec {
			data, err := dbo.Exec(tableId, parameters, context)
			m = map[string]interface{}{
				"data": data,
			}
			if err != nil {
				m["err"] = err.Error()
			}
		} else {
			decoder := json.NewDecoder(r.Body)
			err := decoder.Decode(&m)
			if err != nil {
				m["err"] = err.Error()
				jsonData, _ := json.Marshal(m)
				jsonString := string(jsonData)
				fmt.Fprint(w, jsonString)
				return
			}
			mUpper := make(map[string]interface{})
			for k, v := range m {
				if !strings.HasPrefix(k, "_") {
					mUpper[strings.ToUpper(k)] = v
				}
			}
			data, err := dbo.Create(tableId, mUpper, context)
			m = map[string]interface{}{
				"data": data,
			}
			if err != nil {
				m["err"] = err.Error()
			}
		}
		jsonData, err := json.Marshal(m)
		if err != nil {
			fmt.Println(err)
		}
		jsonString := string(jsonData)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		fmt.Fprint(w, jsonString)
	case "COPY":
		// Duplicate a new record.
		dataId := urlPathData[2]
		data, err := dbo.Duplicate(tableId, dataId, context)

		m := map[string]interface{}{
			"data": data,
		}
		if err != nil {
			m["err"] = err.Error()
		}
		jsonData, err := json.Marshal(m)
		jsonString := string(jsonData)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		fmt.Fprint(w, jsonString)
	case "PUT":
		// Update an existing record.
		metaValues := r.URL.Query()["meta"]
		meta := true
		if metaValues != nil && metaValues[0] == "0" {
			meta = false
		}
		context["meta"] = meta
		dataId := ""
		if len(urlPathData) >= 3 {
			dataId = urlPathData[1]
		}
		decoder := json.NewDecoder(r.Body)
		m := make(map[string]interface{})
		err := decoder.Decode(&m)
		if err != nil {
			m["err"] = err.Error()
			jsonData, _ := json.Marshal(m)
			jsonString := string(jsonData)
			fmt.Fprint(w, jsonString)
			return
		}
		mUpper := make(map[string]interface{})
		for k, v := range m {
			if !strings.HasPrefix(k, "_") {
				mUpper[strings.ToUpper(k)] = v
			}
		}
		if dataId != "" {
			mUpper["ID"] = dataId
		}
		data, err := dbo.Update(tableId, mUpper, context)
		m = map[string]interface{}{
			"data": data,
		}
		if err != nil {
			m["err"] = err.Error()
		}
		jsonData, err := json.Marshal(m)
		jsonString := string(jsonData)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		fmt.Fprint(w, jsonString)
	case "DELETE":
		// Remove the record.
		dataId := urlPathData[2]

		//		load := false
		//		l := r.FormValue("load")
		//		if l == "1" {
		//			load = true
		//		}
		//		context["load"] = load

		data, err := dbo.Delete(tableId, dataId, context)

		m := map[string]interface{}{
			"data": data,
		}
		if err != nil {
			m["err"] = err.Error()
		}
		jsonData, err := json.Marshal(m)
		jsonString := string(jsonData)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		fmt.Fprint(w, jsonString)
	case "OPTIONS":
	default:
		// Give an error message.
	}
}
