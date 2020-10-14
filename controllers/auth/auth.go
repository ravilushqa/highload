package auth

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"github.com/ravilushqa/highload/lib"
	"github.com/ravilushqa/highload/lib/user"
)

type Controller struct {
	logger      *zap.Logger
	auth        *lib.Auth
	userManager *user.Manager
}

func NewController(logger *zap.Logger, auth *lib.Auth, userManager *user.Manager) *Controller {
	return &Controller{logger: logger, auth: auth, userManager: userManager}
}

func (c *Controller) Router(r chi.Router) chi.Router {
	return r.Route("/", func(r chi.Router) {
		r.Get("/signin", c.signin)
		r.Post("/login", c.login)
		r.Post("/register", c.register)
		r.Get("/logout", c.logout)
	})
}

func (c *Controller) login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	email := r.FormValue("login-form-email")
	password := r.FormValue("login-form-password")

	if len(email) == 0 || len(password) == 0 {
		_, _ = w.Write([]byte("Please provide email and password to obtain the token"))
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	u, err := c.userManager.GetByEmail(r.Context(), email)
	if err != nil {
		c.logger.Error("failed get email", zap.Error(err))
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("wrong email"))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))

	if err == nil {
		token, err := c.auth.EncodeToken(u.ID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("Error generating JWT token: " + err.Error()))
		} else {
			http.SetCookie(w, &http.Cookie{
				Name:     "jwt",
				Value:    token,
				Expires:  time.Now().AddDate(0, 0, 14),
				HttpOnly: true,
			})

			http.Redirect(w, r, fmt.Sprintf("/users/%d", u.ID), http.StatusTemporaryRedirect)
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("Name and password do not match"))
		return
	}
}

func (c *Controller) signin(w http.ResponseWriter, r *http.Request) {
	if lib.IsAuth(r) {
		uid, _ := lib.GetUsedIDFromCtx(r.Context())
		http.Redirect(w, r, fmt.Sprintf("/users/%d", uid), http.StatusTemporaryRedirect)
	}
	tmpl, err := template.ParseFiles("resources/views/base.html", "resources/views/auth/signin.html")
	if err != nil {
		c.logger.Error("failed parse templates", zap.NamedError("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tmpl.ExecuteTemplate(w, "layout", nil)
}

func (c *Controller) register(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	bd, err := time.Parse("2006-01-02", r.FormValue("register-form-birthday"))
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(r.FormValue("register-form-password")), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	u := &user.User{
		Email:     r.FormValue("register-form-email"),
		Password:  string(hashedPassword),
		FirstName: r.FormValue("register-form-first-name"),
		LastName:  r.FormValue("register-form-last-name"),
		Birthday:  bd,
		Sex:       user.Sex(r.FormValue("register-form-sex")),
		Interests: r.FormValue("register-form-interests"),
		City:      r.FormValue("register-form-city"),
	}

	u.ID, err = c.userManager.Store(r.Context(), u)

	if err != nil {
		http.Redirect(w, r, r.Header.Get("Referer"), 302)
		return
	}

	token, err := c.auth.EncodeToken(u.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Error generating JWT token: " + err.Error()))
	} else {
		http.SetCookie(w, &http.Cookie{
			Name:     "jwt",
			Value:    token,
			Expires:  time.Now().AddDate(0, 0, 14),
			HttpOnly: true,
		})

		http.Redirect(w, r, fmt.Sprintf("/users/%d", u.ID), http.StatusTemporaryRedirect)
	}

}

func (c *Controller) logout(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:    "jwt",
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),

		HttpOnly: true,
	}

	http.SetCookie(w, cookie)
	http.Redirect(w, r, fmt.Sprintf("/signin"), http.StatusTemporaryRedirect)
}
