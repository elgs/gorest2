// handlers_rest
package gorest2

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

var RestFunc = func(dbo DataOperator) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		context := make(map[string]interface{})
		context["api_token_id"] = r.Header.Get("api_token_id")
		context["api_token_key"] = r.Header.Get("api_token_key")

		urlPath := r.URL.Path
		//fmt.Println(urlPath)
		restData := strings.Split(urlPath, "/")
		tableId := restData[0]

		switch r.Method {
		case "GET":
			if len(restData) == 1 ||
				strings.HasPrefix(restData[1], "?") ||
				len(restData[1]) == 0 {
				//List records.
				t := r.FormValue("total")
				a := r.FormValue("array")
				filter := r.Form["filter"]
				fields := strings.ToUpper(r.FormValue("fields"))
				sort := r.FormValue("sort")
				group := r.FormValue("group")
				s := r.FormValue("start")
				l := r.FormValue("limit")
				c := r.FormValue("case")
				context["case"] = c
				includeTotal := false
				array := false
				if t == "1" {
					includeTotal = true
				}
				if a == "1" {
					array = true
				}
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
				var data interface{}
				var total int64 = -1
				if array {
					var headers []string
					var dataArray [][]string
					headers, dataArray, total, err = dbo.ListArray(tableId, fields, filter, sort, group, start, limit, includeTotal, context)
					data = map[string]interface{}{
						"headers": headers,
						"data":    dataArray,
					}
				} else {
					data, total, err = dbo.ListMap(tableId, fields, filter, sort, group, start, limit, includeTotal, context)
				}
				m := map[string]interface{}{
					"data":  data,
					"total": total,
				}
				if err != nil {
					m["err"] = err.Error()
				}
				jsonData, err := json.Marshal(m)
				jsonString := string(jsonData)
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				fmt.Fprint(w, jsonString)
			} else {
				// Load record by id.
				dataId := restData[1]
				c := r.FormValue("case")
				context["case"] = c

				fields := strings.ToUpper(r.FormValue("fields"))

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
			meta := false
			if metaValues != nil && metaValues[0] == "1" {
				meta = true
			}
			context["meta"] = meta

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
				mUpper[strings.ToUpper(k)] = v
			}
			data, info, err := dbo.Create(tableId, mUpper, context)
			m = map[string]interface{}{
				"data": data,
				"info": info,
			}
			if err != nil {
				m["err"] = err.Error()
			}
			jsonData, err := json.Marshal(m)
			jsonString := string(jsonData)
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			fmt.Fprint(w, jsonString)
		case "COPY":
			// Duplicate a new record.
			dataId := restData[1]
			data, info, err := dbo.Duplicate(tableId, dataId, context)

			m := map[string]interface{}{
				"data": data,
				"info": info,
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
			meta := false
			if metaValues != nil && metaValues[0] == "1" {
				meta = true
			}
			context["meta"] = meta
			dataId := ""
			if len(restData) >= 2 {
				dataId = restData[1]
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
				mUpper[strings.ToUpper(k)] = v
			}
			if dataId != "" {
				mUpper["ID"] = dataId
			}
			data, info, err := dbo.Update(tableId, mUpper, context)
			m = map[string]interface{}{
				"data": data,
				"info": info,
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
			dataId := restData[1]

			load := false
			l := r.FormValue("load")
			if l == "1" {
				load = true
			}
			context["load"] = load

			data, info, err := dbo.Delete(tableId, dataId, context)

			m := map[string]interface{}{
				"data": data,
				"info": info,
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
}
