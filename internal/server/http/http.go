package http

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"task-scheduler/internal/api"
	"task-scheduler/internal/users"

	"github.com/bnkamalesh/errors"
	"github.com/bnkamalesh/webgo/v6"
	"github.com/bnkamalesh/webgo/v6/middleware/accesslog"
	"github.com/dgrijalva/jwt-go"
	"go.elastic.co/apm"
	"go.elastic.co/apm/module/apmhttp"
)

type Handlers struct {
	api *api.API
}

type HTTP struct {
	server *http.Server
	cfg    *Config
}

type Config struct {
	Host              string
	Port              string
	TemplatesBasePath string
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	DialTimeout       time.Duration
}

func errResponder(w http.ResponseWriter, err error) {
	status, msg, _ := errors.HTTPStatusCodeMessage(err)
	webgo.SendError(w, msg, status)
}

func authRoute(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		if len(authHeader) != 2 {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Authorization token is missing"))
		} else {
			jwtToken := authHeader[1]
			customClaims := &users.Claims{}
			token, err := jwt.ParseWithClaims(jwtToken, customClaims, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}
				secret_key := []byte(os.Getenv("JWT_SECRET_KEY"))
				return []byte(secret_key), nil
			})

			if token.Valid {
				ctx := context.WithValue(r.Context(), "props", customClaims)
				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				fmt.Println(err)
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unauthorized"))
			}
		}
	})
}

func (h *Handlers) routes() []*webgo.Route {
	return []*webgo.Route{
		&webgo.Route{
			Name:          "health",
			Pattern:       "/-/health",
			Method:        http.MethodGet,
			Handlers:      []http.HandlerFunc{h.Health},
			TrailingSlash: true,
		},
		&webgo.Route{
			Name:          "register",
			Pattern:       "/api/auth/register",
			Method:        http.MethodPost,
			Handlers:      []http.HandlerFunc{h.Register},
			TrailingSlash: true,
		},
		&webgo.Route{
			Name:          "login",
			Pattern:       "/api/auth/login",
			Method:        http.MethodPost,
			Handlers:      []http.HandlerFunc{h.Login},
			TrailingSlash: true,
		},
		&webgo.Route{
			Name:          "add-task",
			Pattern:       "/api/tasks",
			Method:        http.MethodPost,
			Handlers:      []http.HandlerFunc{authRoute(http.HandlerFunc(h.AddTask))},
			TrailingSlash: true,
		},
		&webgo.Route{
			Name:          "edit-task",
			Pattern:       "/api/tasks/:tid",
			Method:        http.MethodPut,
			Handlers:      []http.HandlerFunc{authRoute(http.HandlerFunc(h.EditTask))},
			TrailingSlash: true,
		},
		&webgo.Route{
			Name:          "delete-task",
			Pattern:       "/api/tasks/:tid",
			Method:        http.MethodDelete,
			Handlers:      []http.HandlerFunc{authRoute(http.HandlerFunc(h.DeleteTask))},
			TrailingSlash: true,
		},
		&webgo.Route{
			Name:          "get-tasks",
			Pattern:       "/api/tasks",
			Method:        http.MethodGet,
			Handlers:      []http.HandlerFunc{authRoute(http.HandlerFunc(h.GetAllTasks))},
			TrailingSlash: true,
		},
		// this should be authorized with an admin token or whoever has access to assign tasks
		&webgo.Route{
			Name:          "assign-tasks",
			Pattern:       "/api/tasks/assign",
			Method:        http.MethodPost,
			Handlers:      []http.HandlerFunc{http.HandlerFunc(h.AssignTask)},
			TrailingSlash: true,
		},
	}
}

func (h *Handlers) Health(w http.ResponseWriter, r *http.Request) {
	out, err := h.api.Health()
	if err != nil {
		errResponder(w, err)
		return
	}
	webgo.R200(w, out)
}

func (h *HTTP) Start() {
	webgo.LOGHANDLER.Info("HTTP server, listening on", h.cfg.Host, h.cfg.Port)
	h.server.ListenAndServe()
}

func NewService(cfg *Config, a *api.API) (*HTTP, error) {
	h := &Handlers{
		api: a,
	}

	router := webgo.NewRouter(
		&webgo.Config{
			Host:            cfg.Host,
			Port:            cfg.Port,
			ReadTimeout:     cfg.ReadTimeout,
			WriteTimeout:    cfg.WriteTimeout,
			ShutdownTimeout: cfg.WriteTimeout * 2,
		},
		h.routes()...,
	)

	router.Use(accesslog.AccessLog)
	tracer, _ := apm.NewTracer("task-scheduler", "v1.0.0")

	serverHandler := apmhttp.Wrap(
		router,
		apmhttp.WithRecovery(apmhttp.NewTraceRecovery(
			tracer,
		)),
	)

	httpServer := &http.Server{
		Addr:              fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Handler:           serverHandler,
		ReadTimeout:       cfg.ReadTimeout,
		ReadHeaderTimeout: cfg.ReadTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.ReadTimeout * 2,
	}

	return &HTTP{
		server: httpServer,
		cfg:    cfg,
	}, nil
}
