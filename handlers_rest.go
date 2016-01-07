// handlers_rest
package gorest2

import (
	"encoding/json"
	"fmt"
	"github.com/dvsekhvalnov/jose2go"
	"github.com/elgs/gojq"
	"github.com/elgs/gosplitargs"
	"gopkg.in/redis.v3"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var RedisMaster *redis.Client
var RedisLocal *redis.Client

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
	context["token"] = r.Header.Get("token")

	projectId := r.Header.Get("app_id")
	authorization := r.Header.Get("Authorization")
	if authorization != "" {
		authTokenArray := strings.SplitN(authorization, " ", 2)
		authToken := authTokenArray[len(authTokenArray)-1]

		payload, _, err := jose.Decode(authToken, []byte{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		jq, err := gojq.NewStringQuery(payload)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		appId, err := jq.Query("app_id")
		if err == nil && appId != nil && projectId == "" {
			projectId, _ = appId.(string)
		}

		userId, err := jq.Query("id")
		if err == nil && userId != nil {
			context["user_id"] = userId
		}

		email, err := jq.Query("email")
		if err == nil && email != nil {
			context["email"] = email
		}

		userKey := fmt.Sprint("user:", appId, ":", userId)
		if authToken != RedisLocal.HGet(userKey, "authToken").Val() {
			http.Error(w, "Authentication failed.", http.StatusInternalServerError)
			return
		}
	}

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

	clientIp := strings.Split(r.RemoteAddr, ":")[0]
	context["client_ip"] = clientIp

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
				if ts, ok := m["data"].([]map[string]string); ok {
					if len(ts) == 1 {
						t := ts[0]
						userId := ""
						payload := make(map[string]interface{})
						for k, v := range t {
							uk := strings.ToUpper(k)
							if uk == "CREATE_TIME" || uk == "CREATETIME" ||
								uk == "CREATOR_ID" || uk == "CREATORID" ||
								uk == "CREATOR_CODE" || uk == "CREATORCODE" ||
								uk == "UPDATE_TIME" || uk == "UPDATETIME" ||
								uk == "UPDATER_ID" || uk == "UPDATERID" ||
								uk == "UPDATER_CODE" || uk == "UPDATERCODE" {
								continue
							}
							if uk == "ID" {
								userId = v
							}
							payload[k] = v
						}
						payload["app_id"] = projectId
						payload["exp"] = time.Now().Add(time.Hour * 72).Unix()

						payloadBytes, err := json.Marshal(&payload)
						tokenString, err := jose.Sign(string(payloadBytes), jose.HS256, []byte{})
						if err != nil {
							http.Error(w, err.Error(), http.StatusInternalServerError)
							return
						}
						jsonString := fmt.Sprintf(`{"token":"%v"}`, tokenString)
						userKey := strings.Join([]string{"user", projectId, userId}, ":")
						err = RedisMaster.HMSet(userKey, "authToken", tokenString).Err()
						if err != nil {
							http.Error(w, err.Error(), http.StatusInternalServerError)
							return
						}
						w.Header().Set("Content-Type", "application/json; charset=utf-8")
						fmt.Fprint(w, jsonString)
						return
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
		exec := false
		if execValues != nil && execValues[0] == "1" {
			exec = true
		}

		m := make(map[string]interface{})
		if exec {
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
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			data, err := dbo.Exec(tableId, parameters, queryParams, context)
			m = map[string]interface{}{
				"data": data,
			}
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			decoder := json.NewDecoder(r.Body)
			err := decoder.Decode(&m)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
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
		dataId := urlPathData[2]
		data, err := dbo.Duplicate(tableId, dataId, context)

		m := map[string]interface{}{
			"data": data,
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
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
		if len(urlPathData) >= 3 {
			dataId = urlPathData[2]
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
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
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
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
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
