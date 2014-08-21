package gorest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"strconv"
	"strings"
)

type Gorest struct {
	EnableHttp bool
	PortHttp   uint16
	HostHttp   string

	EnableHttps   bool
	PortHttps     uint16
	HostHttps     string
	CertFileHttps string
	KeyFileHttps  string

	UrlPrefix string
	Dbo       DataOperator
}

func (this *Gorest) Serve() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		urlPath := r.URL.Path
		urlPrefix := fmt.Sprint("/", this.UrlPrefix, "/")
		if !strings.HasPrefix(urlPath, urlPrefix) {
			return
		}
		restUrl := urlPath[len(urlPrefix):]
		restData := strings.Split(restUrl, "/")
		tableId := restData[0]
		switch r.Method {
		case "GET":
			if len(restData) == 1 ||
				strings.HasPrefix(restData[1], "?") ||
				len(restData[1]) == 0 {
				//List records.
				t := r.FormValue("total")
				a := r.FormValue("array")
				where := r.FormValue("where")
				order := r.FormValue("order")
				s := r.FormValue("start")
				l := r.FormValue("limit")
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
				}
				limit, err := strconv.ParseInt(l, 10, 0)
				if err != nil {
					limit = 25
				}
				var data interface{}
				var total int64 = -1
				if array {
					data, total = this.Dbo.ListArray(tableId, where, order, start, limit, includeTotal)
				} else {
					data, total = this.Dbo.ListMap(tableId, where, order, start, limit, includeTotal)
				}
				m := map[string]interface{}{
					"data":  data,
					"total": total,
				}
				json, err := json.Marshal(m)
				if err != nil {
					http.Error(w, err.Error(), 500)
					return
				}
				jsonString := string(json)
				fmt.Fprint(w, jsonString)
			} else {
				// Load record by id.
				dataId := restData[1]

				data := this.Dbo.Load(tableId, dataId)

				json, err := json.Marshal(data)
				if err != nil {
					http.Error(w, err.Error(), 500)
					return
				}
				jsonString := string(json)
				fmt.Fprint(w, jsonString)
			}
		case "POST":
			if tableId == "_query" {
				sqlSelect := r.FormValue("sql_select")
				sqlSelectCount := r.FormValue("sql_select_count")
				t := r.FormValue("total")
				a := r.FormValue("array")
				s := r.FormValue("start")
				l := r.FormValue("limit")
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
				}
				limit, err := strconv.ParseInt(l, 10, 0)
				if err != nil {
					limit = 25
				}
				var data interface{}
				var total int64 = -1
				if array {
					data, total = this.Dbo.QueryArray(tableId, sqlSelect, sqlSelectCount, start, limit, includeTotal)
				} else {
					data, total = this.Dbo.QueryMap(tableId, sqlSelect, sqlSelectCount, start, limit, includeTotal)
				}
				m := map[string]interface{}{
					"data":  data,
					"total": total,
				}
				json, err := json.Marshal(m)
				if err != nil {
					http.Error(w, err.Error(), 500)
					return
				}
				jsonString := string(json)
				fmt.Fprint(w, jsonString)
			} else {
				// Create the record.
				decoder := json.NewDecoder(r.Body)
				var m map[string]interface{}
				err := decoder.Decode(&m)
				if err != nil {
					fmt.Println(err)
					return
				}
				mUpper := make(map[string]interface{})
				for k, v := range m {
					mUpper[strings.ToUpper(k)] = v
				}
				data := this.Dbo.Create(tableId, mUpper)
				json, err := json.Marshal(data)
				if err != nil {
					http.Error(w, err.Error(), 500)
					return
				}
				jsonString := string(json)
				fmt.Fprint(w, jsonString)
			}
		case "COPY":
			// Duplicate a new record.
			dataId := restData[1]

			data := this.Dbo.Duplicate(tableId, dataId)

			json, err := json.Marshal(data)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			jsonString := string(json)
			fmt.Fprint(w, jsonString)
		case "PUT":
			// Update an existing record.
			decoder := json.NewDecoder(r.Body)
			var m map[string]interface{}
			err := decoder.Decode(&m)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			mUpper := make(map[string]interface{})
			for k, v := range m {
				mUpper[strings.ToUpper(k)] = v
			}
			data := this.Dbo.Update(tableId, mUpper)
			json, err := json.Marshal(data)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			jsonString := string(json)
			fmt.Fprint(w, jsonString)
		case "DELETE":
			// Remove the record.
			dataId := restData[1]

			data := this.Dbo.Delete(tableId, dataId)

			json, err := json.Marshal(data)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			jsonString := string(json)
			fmt.Fprint(w, jsonString)
		default:
			// Give an error message.
		}
	}
	http.HandleFunc("/", handler)

	if this.EnableHttp {
		go func() {
			fmt.Println(fmt.Sprint("Listening on http://", this.HostHttp, ":", this.PortHttp, "/", this.UrlPrefix))
			http.ListenAndServe(fmt.Sprint(this.HostHttp, ":", this.PortHttp), nil)
		}()
	}
	if this.EnableHttps {
		go func() {
			fmt.Println(fmt.Sprint("Listening on https://", this.HostHttps, ":", this.PortHttps, "/", this.UrlPrefix))
			http.ListenAndServeTLS(fmt.Sprint(this.HostHttps, ":", this.PortHttps), this.CertFileHttps, this.KeyFileHttps, nil)
		}()
	}
	if this.EnableHttp || this.EnableHttps {
		select {}
	} else {
		fmt.Println("Neither http nor https is listening, therefore I am quiting.")
	}
}
