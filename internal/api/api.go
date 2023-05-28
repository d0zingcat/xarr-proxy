package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"xarr-proxy/internal/config"
	"xarr-proxy/internal/consts"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
)

var (
	r         *chi.Mux
	server    *http.Server
	tokenAuth *jwtauth.JWTAuth
)

func Init(cfg *config.Config) *chi.Mux {
	r = chi.NewRouter()
	InitMiddleware()
	InitErrorHandlers()
	InitRoutes()
	return r
}

func InitAuth(cfg *config.Config) {
	tokenAuth = jwtauth.New("HS256", []byte(cfg.JWTSecret), nil)
}

func InitErrorHandlers() {
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		w.Write([]byte("{\"code\":-1,\"message\":\"route does not exist\"}"))
	})
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(405)
		w.Write([]byte("{\"code\":-1,\"message\":\"method is not valid\"}"))
	})
}

func InitMiddleware() {
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	// r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))
}

func InitRoutes() {
	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, consts.STATIC_FILE_DIR))
	FileServer(r, "/", filesDir)

	// r.Handle("/*",
	// 	http.StripPrefix("", http.FileServer(http.Dir(consts.STATIC_FILE_DIR))))
	r.Route("/api", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Post("/system/user/login", userLogin)
			r.Get("/system/user/info", userInfo)
		})
		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Authenticator)
			r.Get("/auth", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("Authorized!"))
			})
		})
	})
	// r.Get("/", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("Welcome!"))
	// })
}

func Start(cfg *config.Config) {
	server = &http.Server{Addr: fmt.Sprintf(":%v", config.ServerPort), Handler: r}
	serverCtx, serverStopCtx := context.WithCancel(context.Background())
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, c := context.WithTimeout(serverCtx, 10*time.Second)
		defer c()

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				panic("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			panic(err)
		}
		serverStopCtx()
	}()
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(err)
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
}
