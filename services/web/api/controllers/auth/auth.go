package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/golang/protobuf/ptypes"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	usersGrpc "github.com/ravilushqa/highload/services/users/api/grpc"
	apiLib "github.com/ravilushqa/highload/services/web/api/lib"
	"github.com/ravilushqa/highload/services/web/lib"
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
		c.renderSigninWithError(w, r, "Please provide both email and password to login", "", nil)
		return
	}

	userResponse, err := c.usersClient.GetByEmail(r.Context(), &usersGrpc.GetByEmailRequest{Email: email})
	if err != nil {
		c.logger.Error("failed get email", zap.Error(err))
		c.renderSigninWithError(w, r, "Account not found. Please check your email or register a new account.", email, nil)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(userResponse.User.Password), []byte(password))

	if err == nil {
		token, err := c.auth.EncodeToken(int(userResponse.User.Id))
		if err != nil {
			c.renderSigninWithError(w, r, "Error generating authentication token. Please try again later.", email, nil)
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
		c.renderSigninWithError(w, r, "Incorrect password. Please try again.", email, nil)
		return
	}
}

// renderSigninWithError renders the signin page with an error message
func (c *Controller) renderSigninWithError(w http.ResponseWriter, r *http.Request, loginError string, email string, regData map[string]string) {
	w.WriteHeader(http.StatusUnauthorized)

	templateFiles := []string{
		"resources/views/base.html",
		"resources/views/auth/signin.html",
	}

	data := struct {
		LoginError    string
		Email         string
		RegisterError string
		RegEmail      string
		FirstName     string
		LastName      string
		Birthday      string
		Sex           string
		Interests     string
		City          string
	}{
		LoginError: loginError,
		Email:      email,
	}

	// If we have registration data
	if regData != nil {
		data.RegisterError = regData["error"]
		data.RegEmail = regData["email"]
		data.FirstName = regData["firstName"]
		data.LastName = regData["lastName"]
		data.Birthday = regData["birthday"]
		data.Sex = regData["sex"]
		data.Interests = regData["interests"]
		data.City = regData["city"]
	}

	err := apiLib.RenderTemplate(w, r, templateFiles, data)
	if err != nil {
		c.logger.Error("failed to render signin template with error", zap.NamedError("error", err))
		_, _ = w.Write([]byte("An error occurred. Please try again."))
	}
}

func (c *Controller) signin(w http.ResponseWriter, r *http.Request) {
	if lib.IsAuth(r) {
		uid, _ := lib.GetAuthUserID(r.Context())
		http.Redirect(w, r, fmt.Sprintf("/users/%d", uid), http.StatusTemporaryRedirect)
		return
	}

	templateFiles := []string{
		"resources/views/base.html",
		"resources/views/auth/signin.html",
	}

	// Empty data structure with no errors
	data := struct {
		LoginError    string
		Email         string
		RegisterError string
		RegEmail      string
		FirstName     string
		LastName      string
		Birthday      string
		Sex           string
		Interests     string
		City          string
	}{}

	err := apiLib.RenderTemplate(w, r, templateFiles, data)
	if err != nil {
		c.logger.Error("failed to render signin template", zap.NamedError("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (c *Controller) register(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		c.renderSigninWithRegisterError(w, r, "Error processing form data. Please try again.", r)
		return
	}

	// Validate required fields
	firstName := r.FormValue("register-form-first-name")
	lastName := r.FormValue("register-form-last-name")
	email := r.FormValue("register-form-email")
	password := r.FormValue("register-form-password")
	birthday := r.FormValue("register-form-birthday")
	sex := r.FormValue("register-form-sex")
	interests := r.FormValue("register-form-interests")
	city := r.FormValue("register-form-city")

	if firstName == "" || lastName == "" || email == "" || password == "" || birthday == "" {
		c.renderSigninWithRegisterError(w, r, "All required fields must be filled out.", r)
		return
	}

	bd, err := time.Parse("2006-01-02", birthday)
	if err != nil {
		c.renderSigninWithRegisterError(w, r, "Invalid birthday format. Please use YYYY-MM-DD format.", r)
		return
	}
	bdProto, err := ptypes.TimestampProto(bd)
	if err != nil {
		c.renderSigninWithRegisterError(w, r, "Error processing birthday. Please try again.", r)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.renderSigninWithRegisterError(w, r, "Error processing password. Please try again.", r)
		return
	}

	storeResponse, err := c.usersClient.Store(r.Context(), &usersGrpc.StoreRequest{
		Email:     email,
		Password:  string(hashedPassword),
		FirstName: firstName,
		LastName:  lastName,
		Birthday:  bdProto,
		Interests: interests,
		Sex:       usersGrpc.Sex(usersGrpc.Sex_value[strings.Title(sex)]),
		City:      city,
	})

	if err != nil {
		c.renderSigninWithRegisterError(w, r, "Error creating account. The email might already be in use.", r)
		return
	}

	token, err := c.auth.EncodeToken(int(storeResponse.Id))
	if err != nil {
		c.renderSigninWithRegisterError(w, r, "Error generating authentication token. Please try again later.", r)
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

// renderSigninWithRegisterError renders the signin page with a registration error
func (c *Controller) renderSigninWithRegisterError(w http.ResponseWriter, r *http.Request, errorMessage string, formData *http.Request) {
	w.WriteHeader(http.StatusUnprocessableEntity)

	regData := map[string]string{
		"error":     errorMessage,
		"email":     formData.FormValue("register-form-email"),
		"firstName": formData.FormValue("register-form-first-name"),
		"lastName":  formData.FormValue("register-form-last-name"),
		"birthday":  formData.FormValue("register-form-birthday"),
		"sex":       formData.FormValue("register-form-sex"),
		"interests": formData.FormValue("register-form-interests"),
		"city":      formData.FormValue("register-form-city"),
	}

	c.renderSigninWithError(w, r, "", "", regData)
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
