package auth

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/golang/protobuf/ptypes"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"github.com/ravilushqa/highload/lib"
	usersGrpc "github.com/ravilushqa/highload/services/users/api/grpc"
)

type Controller struct {
	logger      *zap.Logger
	auth        *lib.Auth
	usersClient usersGrpc.UsersClient
}

func NewController(logger *zap.Logger, auth *lib.Auth, usersClient usersGrpc.UsersClient) *Controller {
	return &Controller{logger: logger, auth: auth, usersClient: usersClient}
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
	_ = r.ParseForm()
	email := r.FormValue("login-form-email")
	password := r.FormValue("login-form-password")

	if len(email) == 0 || len(password) == 0 {
		_, _ = w.Write([]byte("Please provide email and password to obtain the token"))
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	userResponse, err := c.usersClient.GetByEmail(r.Context(), &usersGrpc.GetByEmailRequest{Email: email})
	if err != nil {
		c.logger.Error("failed get email", zap.Error(err))
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("wrong email"))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(userResponse.User.Password), []byte(password))

	if err == nil {
		token, err := c.auth.EncodeToken(int(userResponse.User.Id))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("Error generating JWT token: " + err.Error()))
		} else {
			http.SetCookie(w, &http.Cookie{
				Name:    "jwt",
				Value:   token,
				Expires: time.Now().AddDate(0, 0, 14),
				//HttpOnly: true,
			})

			http.Redirect(w, r, fmt.Sprintf("/users/%d", userResponse.User.Id), http.StatusTemporaryRedirect)
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("Name and password do not match"))
		return
	}
}

func (c *Controller) signin(w http.ResponseWriter, r *http.Request) {
	if lib.IsAuth(r) {
		uid, _ := lib.GetAuthUserID(r.Context())
		http.Redirect(w, r, fmt.Sprintf("/users/%d", uid), http.StatusTemporaryRedirect)
	}
	tmpl, err := template.ParseFiles("resources/views/base.html", "resources/views/auth/signin.html")
	if err != nil {
		c.logger.Error("failed parse templates", zap.NamedError("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_ = tmpl.ExecuteTemplate(w, "layout", nil)
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
	bdProto, err := ptypes.TimestampProto(bd)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(r.FormValue("register-form-password")), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	storeResponse, err := c.usersClient.Store(r.Context(), &usersGrpc.StoreRequest{
		Email:     r.FormValue("register-form-email"),
		Password:  string(hashedPassword),
		FirstName: r.FormValue("register-form-first-name"),
		LastName:  r.FormValue("register-form-last-name"),
		Birthday:  bdProto,
		Interests: r.FormValue("register-form-interests"),
		Sex:       usersGrpc.Sex(usersGrpc.Sex_value[strings.Title(r.FormValue("register-form-sex"))]),
		City:      r.FormValue("register-form-city"),
	})

	if err != nil {
		http.Redirect(w, r, r.Header.Get("Referer"), 302)
		return
	}

	token, err := c.auth.EncodeToken(int(storeResponse.Id))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Error generating JWT token: " + err.Error()))
	} else {
		http.SetCookie(w, &http.Cookie{
			Name:    "jwt",
			Value:   token,
			Expires: time.Now().AddDate(0, 0, 14),
			//HttpOnly: true,
		})

		http.Redirect(w, r, fmt.Sprintf("/users/%d", storeResponse.Id), http.StatusTemporaryRedirect)
	}

}

func (c *Controller) logout(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:    "jwt",
		Value:   "",
		Path:    "/",
		Expires: time.Time{},
		//HttpOnly: true,
	}

	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/signin", http.StatusTemporaryRedirect)
}
