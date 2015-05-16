package gorest2

import (
	"fmt"
	"github.com/elgs/gosqljson"
	"net/http"
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

	FileBasePath string
}

func (this *Gorest) Serve() {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", r.Header.Get("Access-Control-Request-Method"))
		w.Header().Set("Access-Control-Allow-Headers", r.Header.Get("Access-Control-Request-Headers"))

		if r.Method == "OPTIONS" {
			return
		}

		urlPath := r.URL.Path
		var dataHandler func(dbo DataOperator) func(w http.ResponseWriter, r *http.Request)
		if strings.HasPrefix(urlPath, "/api/") {
			dataHandler = GetHandler("/api")
		} else {
			dataHandler = GetHandler(urlPath)
		}

		urlPathData := strings.Split(urlPath[1:], "/")
		dboId := urlPathData[0]
		dbo := GetDbo(dboId)
		if dbo == nil {
			dbo = LoadDbo(dboId)
			if dbo == nil {
				dbo = GetDbo("api")
			}
		}
		dataHandler(dbo)(w, r)
	}
	http.HandleFunc("/", handler)

	if this.EnableHttp {
		go func() {
			fmt.Println(fmt.Sprint("Listening on http://", this.HostHttp, ":", this.PortHttp, "/"))
			err := http.ListenAndServe(fmt.Sprint(this.HostHttp, ":", this.PortHttp), nil)
			if err != nil {
				fmt.Println(err)
			}
		}()
	}
	if this.EnableHttps {
		go func() {
			fmt.Println(fmt.Sprint("Listening on https://", this.HostHttps, ":", this.PortHttps, "/"))
			err := http.ListenAndServeTLS(fmt.Sprint(this.HostHttps, ":", this.PortHttps), this.CertFileHttps, this.KeyFileHttps, nil)
			if err != nil {
				fmt.Println(err)
			}
		}()
	}
	if this.EnableHttp || this.EnableHttps {
		select {}
	} else {
		fmt.Println("Neither http nor https is listening, therefore I am quiting.")
	}
}

func LoadDbo(id string) DataOperator {
	defaultDbo := GetDbo("api")
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
	dbo := &MySqlDataOperator{
		Ds: fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", dboData["USERNAME"], dboData["PASSWORD"],
			dboData["HOST"], dboData["PORT"], dboData["DB"]),
		DbType: "mysql",
	}
	RegisterDbo(id, dbo)
	return dbo
}
