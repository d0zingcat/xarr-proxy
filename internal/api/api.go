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

	"xarr-proxy/internal/auth"
	"xarr-proxy/internal/config"
	"xarr-proxy/internal/consts"
	"xarr-proxy/internal/model"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
)

var (
	r      *chi.Mux
	server *http.Server
)

func Init(cfg *config.Config) *chi.Mux {
	r = chi.NewRouter()
	InitMiddleware()
	InitErrorHandlers()
	InitRoutes()
	return r
}

func InitRoutes() {
	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, consts.STATIC_FILE_DIR))
	FileServer(r, "/", filesDir)

	// legacy static file server, just for a backup
	// r.Handle("/*",
	// 	http.StripPrefix("", http.FileServer(http.Dir(consts.STATIC_FILE_DIR))))
	r.Route("/api", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Post("/system/user/login", userLogin)
		})
		r.Group(func(r chi.Router) {
			// jwt auth
			r.Use(jwtauth.Verifier(auth.GetVerifier()))
			r.Use(jwtauth.Authenticator)
			r.Use(MiddlewareUserInfoInjection)

			r.Get("/system/user/info", userInfo)
		})
	})
}

func InitMiddleware() {
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	// r.Use(middleware.URLFormat)
	// r.Use(render.SetContentType(render.ContentTypeJSON))
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

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}

type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

// middleware extract user info from token and format to a valid system user model
func MiddlewareUserInfoInjection(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id, username, role, validStatus := auth.GetUserInfo(r)
		userInfo := model.SystemUser{
			Id:          id,
			Username:    username,
			Role:        role,
			ValidStatus: validStatus,
		}
		ctx = context.WithValue(ctx, consts.USER_INFO_CTX_KEY, userInfo)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
