package api

import (
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
)

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}

// [port from](https://future-architect.github.io/articles/20210408/)
func tryRead(filesDir, requestedPath string, w http.ResponseWriter) error {
	fmt.Println(filesDir, requestedPath)
	f, err := os.Open(filesDir + requestedPath)
	if err != nil {
		return err
	}
	defer f.Close()

	// Goのfs.Openはディレクトリを読みこもとうしてもエラーにはならないがここでは邪魔なのでエラー扱いにする
	stat, _ := f.Stat()
	if stat.IsDir() {
		return errors.New("path is directory")
	}

	contentType := mime.TypeByExtension(filepath.Ext(requestedPath))
	w.Header().Set("Content-Type", contentType)
	_, err = io.Copy(w, f)
	return err
}
