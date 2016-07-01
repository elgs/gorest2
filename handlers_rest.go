// handlers_rest
package gorest2

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	//	"time"

	//	"github.com/dvsekhvalnov/jose2go"
	"github.com/elgs/gosplitargs"
)

var RequestWrites = map[string]int{}
var RequestReads = map[string]int{}

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

	appId := r.Header.Get("app")
	token := r.Header.Get("token")

	context["token"] = token

	if appId == "" {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		fmt.Fprint(w, `{"err":"Invalid app."}`)
		return
	}
	context["app_id"] = appId

	if r.Method == "GET" {
		RequestReads[appId] += 1
	} else {
		RequestWrites[appId] += 1
	}

	dbo, err := GetDbo(appId)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		fmt.Fprint(w, fmt.Sprintf(`{"err":"%v"}`, err))
		return
	}
	if dbo == nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		fmt.Fprint(w, `{"err":"Invalid project."}`)
		return
	}

	sepIndex := strings.LastIndex(r.RemoteAddr, ":")
	clientIp := r.RemoteAddr[0:sepIndex]
	context["client_ip"] = strings.Replace(strings.Replace(clientIp, "[", "", -1), "]", "", -1)

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
			qp := r.FormValue("query_params")
			context["case"] = c
			filter := r.Form["filter"]
			array := translateBoolParam(r.FormValue("array"), false)
			query := translateBoolParam(r.FormValue("query"), false)
			jwtToken := translateBoolParam(r.FormValue("jwt"), false)
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
			params, err := gosplitargs.SplitArgs(p, ",", false)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			parameters := make([]interface{}, len(params))
			for i, v := range params {
				parameters[i] = v
			}

			queryParams, err := gosplitargs.SplitArgs(qp, ",", false)
			_ = queryParams
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			var data interface{}
			var total int64 = -1
			m := map[string]interface{}{}
			if array {
				var headers []string
				var dataArray [][]string
				if query {
					headers, dataArray, err = dbo.QueryArray(tableId, parameters, queryParams, context)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					} else {
						m["headers"] = headers
						m["data"] = dataArray
					}

				} else {
					headers, dataArray, total, err = dbo.ListArray(tableId, fields, filter, sort, group, start, limit, context)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					} else {
						m["headers"] = headers
						m["data"] = dataArray
						m["total"] = total
					}
				}
			} else {
				if query {
					data, err = dbo.QueryMap(tableId, parameters, queryParams, context)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					} else {
						m["data"] = data
					}
				} else {
					data, total, err = dbo.ListMap(tableId, fields, filter, sort, group, start, limit, context)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					} else {
						m["data"] = data
						m["total"] = total
					}
				}
			}
			if jwtToken && !array {
				//				if ts, ok := m["data"].([]map[string]string); ok {
				//					if len(ts) == 1 {
				//						t := ts[0]
				//						userId := ""
				//						payload := make(map[string]interface{})
				//						for k, v := range t {
				//							uk := strings.ToUpper(k)
				//							if uk == "CREATE_TIME" || uk == "CREATETIME" ||
				//								uk == "CREATOR_ID" || uk == "CREATORID" ||
				//								uk == "CREATOR_CODE" || uk == "CREATORCODE" ||
				//								uk == "UPDATE_TIME" || uk == "UPDATETIME" ||
				//								uk == "UPDATER_ID" || uk == "UPDATERID" ||
				//								uk == "UPDATER_CODE" || uk == "UPDATERCODE" {
				//								continue
				//							}
				//							if uk == "ID" {
				//								userId = v
				//							}
				//							payload[k] = v
				//						}
				//						payload["app_id"] = appId
				//						payload["exp"] = time.Now().Add(time.Hour * 72).Unix()

				//						payloadBytes, err := json.Marshal(&payload)
				//						tokenString, err := jose.Sign(string(payloadBytes), jose.HS256, []byte{})
				//						if err != nil {
				//							http.Error(w, err.Error(), http.StatusInternalServerError)
				//							return
				//						}
				//						jsonString := fmt.Sprintf(`{"token":"%v"}`, tokenString)
				//						userKey := strings.Join([]string{"user", appId, userId}, ":")
				//						err = RedisMaster.HMSet(userKey, "authToken", tokenString).Err()
				//						if err != nil {
				//							http.Error(w, err.Error(), http.StatusInternalServerError)
				//							return
				//						}
				//						w.Header().Set("Content-Type", "application/json; charset=utf-8")
				//						fmt.Fprint(w, jsonString)
				//						return
				//					}
				//				}
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
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
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
		exec := 0
		if execValues != nil && execValues[0] == "1" {
			exec = 1
		} else if execValues != nil && execValues[0] == "2" {
			exec = 2
		}

		m := map[string]interface{}{}
		if exec == 1 {
			parameters := make([]interface{}, 0, 10)
			p := r.FormValue("params")
			qp := r.FormValue("query_params")
			params, err := gosplitargs.SplitArgs(p, ",", false)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			for _, v := range params {
				parameters = append(parameters, v)
			}
			queryParams, err := gosplitargs.SplitArgs(qp, ",", false)
			_ = queryParams
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			data, err := dbo.Exec(tableId, [][]interface{}{parameters}, queryParams, context)
			if data != nil && len(data) == 1 {
				m = map[string]interface{}{
					"data": data[0],
				}
			} else {
				m = map[string]interface{}{
					"data": data,
				}
			}
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else if exec == 2 {
			qp := ""
			qpArray := r.URL.Query()["query_params"]
			if qpArray != nil && len(qpArray) > 0 {
				qp = qpArray[0]
			}
			queryParams, err := gosplitargs.SplitArgs(qp, ",", false)
			_ = queryParams
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			var postData []interface{}
			decoder := json.NewDecoder(r.Body)
			err = decoder.Decode(&postData)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			parametersArray := [][]interface{}{}

			for _, postData1 := range postData {
				parameters := []interface{}{}

				if m1, ok := postData1.([]interface{}); ok {
					for _, v := range m1 {
						parameters = append(parameters, v)
					}
				}
				parametersArray = append(parametersArray, parameters)
			}
			data, err := dbo.Exec(tableId, parametersArray, queryParams, context)
			if data != nil && len(data) == 1 {
				m = map[string]interface{}{
					"data": data[0],
				}
			} else {
				m = map[string]interface{}{
					"data": data,
				}
			}
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			var postData interface{}
			decoder := json.NewDecoder(r.Body)
			err := decoder.Decode(&postData)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			inputMode := 1
			postDataArray := []interface{}{}
			switch v := postData.(type) {
			case []interface{}:
				inputMode = 2
				postDataArray = v
			case map[string]interface{}:
				postDataArray = append(postDataArray, v)
			default:
				http.Error(w, "Error parsing post data.", http.StatusInternalServerError)
				return
			}

			upperCasePostDataArray := []map[string]interface{}{}
			for _, m := range postDataArray {
				mUpper := map[string]interface{}{}
				if m1, ok := m.(map[string]interface{}); ok {
					for k, v := range m1 {
						if !strings.HasPrefix(k, "_") {
							mUpper[strings.ToUpper(k)] = v
						}
					}
					upperCasePostDataArray = append(upperCasePostDataArray, mUpper)
				}
			}
			data, err := dbo.Create(tableId, upperCasePostDataArray, context)
			if inputMode == 1 && data != nil && len(data) == 1 {
				m["data"] = data[0]
			} else {
				m["data"] = data
			}
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		jsonData, err := json.Marshal(m)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		jsonString := string(jsonData)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		fmt.Fprint(w, jsonString)
	case "COPY":
		// Duplicate a new record.
		dataIds := []string{}
		if len(urlPathData) >= 3 && len(urlPathData[2]) > 0 {
			dataIds = append(dataIds, urlPathData[2])
		} else {
			var postData interface{}
			decoder := json.NewDecoder(r.Body)
			err := decoder.Decode(&postData)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			postDataArray := []interface{}{}
			switch v := postData.(type) {
			case []interface{}:
				postDataArray = v
			default:
				http.Error(w, "Error parsing post data.", http.StatusInternalServerError)
				return
			}

			for _, postData := range postDataArray {
				if dataId, ok := postData.(string); ok {
					dataIds = append(dataIds, dataId)
				}
			}
		}
		data, err := dbo.Duplicate(tableId, dataIds, context)

		m := map[string]interface{}{}
		if data != nil && len(data) == 1 {
			m["data"] = data[0]
		} else {
			m["data"] = data
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
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
		var postData interface{}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&postData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		inputMode := 1
		postDataArray := []interface{}{}
		switch v := postData.(type) {
		case []interface{}:
			inputMode = 2
			postDataArray = v
		case map[string]interface{}:
			if len(urlPathData) >= 3 && len(urlPathData[2]) > 0 {
				dataId = urlPathData[2]
			}
			postDataArray = append(postDataArray, v)
		default:
			http.Error(w, "Error parsing post data.", http.StatusInternalServerError)
			return
		}

		upperCasePostDataArray := []map[string]interface{}{}
		for _, m := range postDataArray {
			mUpper := map[string]interface{}{}
			if m1, ok := m.(map[string]interface{}); ok {
				for k, v := range m1 {
					if !strings.HasPrefix(k, "_") {
						mUpper[strings.ToUpper(k)] = v
					}
				}
				if inputMode == 1 && dataId != "" {
					mUpper["ID"] = dataId
				}
				upperCasePostDataArray = append(upperCasePostDataArray, mUpper)
			}
		}
		data, err := dbo.Update(tableId, upperCasePostDataArray, context)
		m := map[string]interface{}{}
		if inputMode == 1 && data != nil && len(data) == 1 {
			m["data"] = data[0]
		} else {
			m["data"] = data
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		jsonData, err := json.Marshal(m)
		jsonString := string(jsonData)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		fmt.Fprint(w, jsonString)
	case "DELETE":
		// Remove the record.
		dataIds := []string{}
		if len(urlPathData) >= 3 && len(urlPathData[2]) > 0 {
			dataIds = append(dataIds, urlPathData[2])
		} else {
			var postData interface{}
			decoder := json.NewDecoder(r.Body)
			err := decoder.Decode(&postData)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			postDataArray := []interface{}{}
			switch v := postData.(type) {
			case []interface{}:
				postDataArray = v
			default:
				http.Error(w, "Error parsing post data.", http.StatusInternalServerError)
				return
			}

			for _, postData := range postDataArray {
				if dataId, ok := postData.(string); ok {
					dataIds = append(dataIds, dataId)
				}
			}
		}
		data, err := dbo.Delete(tableId, dataIds, context)

		m := map[string]interface{}{}
		if data != nil && len(data) == 1 {
			m["data"] = data[0]
		} else {
			m["data"] = data
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
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
