package routes

import (
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/uptrace/bun"
	"github.com/utilyre/session-auth/models"
	"github.com/utilyre/xmate"
	"golang.org/x/crypto/bcrypt"
)

type Users struct {
	Handler    xmate.ErrorHandler
	SignupView *template.Template
	LoginView  *template.Template
	DB         *bun.DB
}

func (u Users) Router() chi.Router {
	r := chi.NewRouter()

	r.Get("/signup", u.Handler.HandleFunc(u.signup))
	// r.Get("/login", ur.login)

	r.Post("/signup", u.Handler.HandleFunc(u.handleSignup))

	return r
}

func (u Users) signup(w http.ResponseWriter, r *http.Request) error {
	return xmate.WriteHTML(w, u.SignupView, http.StatusOK, nil)
}

func (u Users) handleSignup(w http.ResponseWriter, r *http.Request) error {
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

	if _, err := u.DB.NewInsert().Model(user).Exec(r.Context()); err != nil {
		return err
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}
