package gorest2

import (
	"fmt"
	"net/http"
	"strings"
)

type Gorest map[string]interface{}

func (this Gorest) Serve() {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", r.Header.Get("Access-Control-Request-Method"))
		w.Header().Set("Access-Control-Allow-Headers", r.Header.Get("Access-Control-Request-Headers"))

		if r.Method == "OPTIONS" {
			return
		}

		urlPath := r.URL.Path
		var dataHandler func(w http.ResponseWriter, r *http.Request)
		if strings.HasPrefix(urlPath, "/api/") {
			dataHandler = GetHandler("/api")
			dataHandler(w, r)
		} else {
			dataHandler = GetHandler(urlPath)
			if dataHandler == nil {
				http.Error(w, "Not found.", http.StatusNotFound)
				return
			}
			for _, globalHandlerInterceptor := range GlobalHandlerInterceptorRegistry {
				ctn, err := globalHandlerInterceptor.BeforeHandle(w, r)
				if !ctn || err != nil {
					fmt.Fprint(w, err.Error())
					return
				}
			}
			handlerInterceptor := HandlerInterceptorRegistry[urlPath]
			if handlerInterceptor != nil {
				ctn, err := handlerInterceptor.BeforeHandle(w, r)
				if !ctn || err != nil {
					fmt.Fprint(w, err.Error())
					return
				}
			}
			dataHandler(w, r)
		}
	}
	http.HandleFunc("/", handler)

	enableHttp := this["enable_http"].(bool)
	enableHttps := this["enable_https"].(bool)

	if enableHttp {
		go func() {
			hostHttp := this["host_http"].(string)
			portHttp := uint16(this["port_http"].(float64))
			fmt.Println(fmt.Sprint("Listening on http://", hostHttp, ":", portHttp, "/"))
			err := http.ListenAndServe(fmt.Sprint(hostHttp, ":", portHttp), nil)
			if err != nil {
				fmt.Println(err)
			}
		}()
	}
	if enableHttps {
		go func() {
			hostHttps := this["host_https"].(string)
			portHttps := uint16(this["port_https"].(float64))
			certFileHttps := this["cert_file_https"].(string)
			keyFileHttps := this["key_file_https"].(string)
			fmt.Println(fmt.Sprint("Listening on https://", hostHttps, ":", portHttps, "/"))
			err := http.ListenAndServeTLS(fmt.Sprint(hostHttps, ":", portHttps), certFileHttps, keyFileHttps, nil)
			if err != nil {
				fmt.Println(err)
			}
		}()
	}
	if enableHttp || enableHttps {
		select {}
	} else {
		fmt.Println("Neither http nor https is listening, therefore I am quiting.")
	}
}
