package gorest2

import (
	"fmt"
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

		//		apiTokenId := r.Header.Get("api_token_id")
		//		apiTokenKey := r.Header.Get("api_token_key")
		//		dataStoreKey := r.Header.Get("data_store_key")

		if r.Method == "OPTIONS" {
			return
		}

		urlPath := r.URL.Path
		var dataHandler func(w http.ResponseWriter, r *http.Request)
		if strings.HasPrefix(urlPath, "/api/") {
			dataHandler = GetHandler("/api")
		} else {
			dataHandler = GetHandler(urlPath)
		}
		dataHandler(w, r)
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
