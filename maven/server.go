package maven

import (
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
)

type PutHandler struct {
	rootpath string
	auth     *AuthManager
}

type AuthManager interface {
	HasAccess(path string, username string, password string) bool
}

func StartServer(conf ServerConf, auth *AuthManager) {
	handler := mux.NewRouter()
	fileserver := http.FileServer(http.Dir(conf.MavenPath))
	handler.PathPrefix(conf.BasePath).Methods("GET").Handler(http.StripPrefix(conf.BasePath, fileserver))
	handler.PathPrefix(conf.BasePath).Methods("HEAD").Handler(http.StripPrefix(conf.BasePath, fileserver))
	handler.PathPrefix(conf.BasePath).Methods("PUT").Handler(http.StripPrefix(conf.BasePath, PutHandler{rootpath: conf.MavenPath, auth: auth}))
	err := http.ListenAndServe("0.0.0.0:"+strconv.Itoa(conf.Port), handler)
	if err != nil {
		panic(err)
	}
}

func (h PutHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	filepath := h.rootpath + req.URL.Path
	auth := *h.auth
	username, password, _ := req.BasicAuth()
	if auth.HasAccess(req.URL.Path, username, password) {
		err := os.MkdirAll(path.Dir(filepath), 0777)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
		f, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		defer f.Close()
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
		_, err = io.Copy(f, req.Body)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
	} else {
		http.Error(res, "No Permission", http.StatusForbidden)
	}
}
