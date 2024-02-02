package routes

import (
	"errors"
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/uptrace/bun"
	"github.com/utilyre/session-auth/models"
	"github.com/utilyre/session-auth/utils"
	"github.com/utilyre/xmate"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserDuplicated = xmate.NewHTTPError(http.StatusConflict, "user already exists")
)

type UserRoute struct {
	Handler    xmate.ErrorHandler
	SignupView *template.Template
	LoginView  *template.Template
	DB         *bun.DB
}

func (ur UserRoute) Router() chi.Router {
	r := chi.NewRouter()

	r.Get("/signup", ur.Handler.HandleFunc(ur.signup))
	// r.Get("/login", ur.login)

	r.Post("/signup", ur.Handler.HandleFunc(ur.handleSignup))

	return r
}

func (ur UserRoute) signup(w http.ResponseWriter, r *http.Request) error {
	return xmate.WriteHTML(w, ur.SignupView, http.StatusOK, nil)
}

func (ur UserRoute) handleSignup(w http.ResponseWriter, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	email := r.PostForm["email"][0]
	password := r.PostForm["password"][0]
	name := r.PostForm["name"][0]

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &models.User{
		Email:    email,
		Password: hash,
		Name:     name,
	}

	if _, err := ur.DB.NewInsert().Model(user).Exec(r.Context()); err != nil {
		err = utils.WrapDBErr(err)
		if errors.Is(err, utils.ErrDuplicatedKey) {
			return ErrUserDuplicated
		}

		return err
	}

	w.Header().Set("HX-Redirect", "/")
	return xmate.WriteText(w, http.StatusCreated, "user created successfully")
}
