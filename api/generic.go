package api

import (
	"log"
	"net/http"
	"html/template"
	"encoding/json"

	"github.com/andyxning/shortme/conf"
)

func version(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	versionInfo, _ := json.Marshal(Version{Version: conf.Version})
	w.Write(versionInfo)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	tpl := template.New("health.html")
	var err error
	tpl, err = tpl.ParseFiles("static/template/health.html")
	if err != nil {
		log.Printf("parse template error. %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		return
	}
	err = tpl.Execute(w, nil)
	if err != nil {
		log.Printf("execute template error. %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
	}
}
