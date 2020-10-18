package main

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	"github.com/ravilushqa/highload/controllers/auth"
	"github.com/ravilushqa/highload/controllers/users"
	"github.com/ravilushqa/highload/lib"
)

type API struct {
	serv    *http.Server
	mux     *chi.Mux
	config  *config
	logger  *zap.Logger
	auth    *auth.Controller
	users   *users.Controller
	libAuth *lib.Auth
}

func NewAPI(config *config, logger *zap.Logger, auth *auth.Controller, users *users.Controller, libAuth *lib.Auth) *API {
	return &API{config: config, logger: logger, auth: auth, users: users, libAuth: libAuth}
}

func (a *API) Run(ctx context.Context) error {
	a.mux = chi.NewRouter()
	a.mux.Use(
		middleware.Logger,
		middleware.RequestID,
		middleware.RealIP,
		middleware.Recoverer,
		middleware.StripSlashes,
		middleware.Timeout(a.config.APITimeout),
		jwtauth.Verifier(a.libAuth.GetToken()),
	)
	a.registerRoutes()

	a.serv = &http.Server{
		Addr:    a.config.Addr,
		Handler: a.mux,
	}

	go func() {
		<-ctx.Done()
		_ = a.Shutdown(ctx)
	}()

	a.logger.Info("service started", zap.String("listen", a.config.Addr))
	if err := a.serv.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (a *API) Shutdown(ctx context.Context) error {
	a.logger.Info("api shutdown")
	if err := a.serv.Shutdown(ctx); err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (a *API) registerRoutes() {
	a.mux.Get("/health-check", func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, map[string]string{"status": "ok"})
	})
	a.mux.Handle("/metrics", promhttp.Handler())

	// public group
	a.mux.Group(func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/signin", http.StatusTemporaryRedirect)
		})
		// static files
		workDir, _ := os.Getwd()
		filesDir := http.Dir(filepath.Join(workDir, "public"))
		a.FileServer(r, "/public", filesDir)

		a.auth.Router(r)
	})

	// auth group
	a.mux.Group(func(r chi.Router) {
		r.Use(jwtauth.Authenticator)

		a.users.Router(r)
	})
}

func (a *API) FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
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
