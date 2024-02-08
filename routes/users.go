package routes

import (
	"database/sql"
	"errors"
	"html/template"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/gorilla/securecookie"
	"github.com/uptrace/bun"
	"github.com/utilyre/session-auth/models"
	"github.com/utilyre/xmate"
	"golang.org/x/crypto/bcrypt"
)

var ErrUserNotFound = xmate.NewHTTPError(http.StatusNotFound, "user not found")

type Users struct {
	Handler    xmate.ErrorHandler
	SignupView *template.Template
	LoginView  *template.Template
	SCookie    *securecookie.SecureCookie
	DB         *bun.DB
}

func (u Users) Router() chi.Router {
	r := chi.NewRouter()

	r.Get("/signup", u.Handler.HandleFunc(u.signup))
	r.Get("/login", u.Handler.HandleFunc(u.login))

	r.Post("/signup", u.Handler.HandleFunc(u.handleSignup))
	r.Post("/login", u.Handler.HandleFunc(u.handleLogin))

	return r
}

func (u Users) signup(w http.ResponseWriter, r *http.Request) error {
	return xmate.WriteHTML(w, u.SignupView, http.StatusOK, nil)
}

func (u Users) handleSignup(w http.ResponseWriter, r *http.Request) error {
	email := r.FormValue("email")
	password := r.FormValue("password")
	name := r.FormValue("name")

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

func (u Users) login(w http.ResponseWriter, r *http.Request) error {
	return xmate.WriteHTML(w, u.LoginView, http.StatusOK, nil)
}

func (u Users) handleLogin(w http.ResponseWriter, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	user := new(models.User)
	if err := u.DB.NewSelect().Model(user).Where("email = ?", email).Scan(r.Context()); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrUserNotFound
		}

		return err
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ErrUserNotFound
		}

		return err
	}

	newUUID := uuid.New()
	session := &models.Session{
		CreatedAt: time.Now(),
		UUID:      newUUID,
		LastIP:    r.RemoteAddr,
		UserID:    user.ID,
	}

	if _, err := u.DB.NewInsert().Model(session).Exec(r.Context()); err != nil {
		return err
	}

	cookieValue, err := u.SCookie.Encode("token", newUUID)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    cookieValue,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}
