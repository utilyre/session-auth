package routes

import (
	"context"
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/gorilla/securecookie"
	"github.com/uptrace/bun"
	"github.com/utilyre/session-auth/models"
	"github.com/utilyre/xmate"
)

type UserKey struct{}

type Me struct {
	Handler    xmate.ErrorHandler
	SecretView *template.Template
	SCookie    *securecookie.SecureCookie
	DB         *bun.DB
}

func (m Me) Router() chi.Router {
	r := chi.NewRouter()

	r.Use(func(next http.Handler) http.Handler {
		return m.Handler.HandleFunc(func(w http.ResponseWriter, r *http.Request) error {
			cookie, err := r.Cookie("token")
			if err != nil {
				return err
			}

			newUUID := uuid.UUID{}
			if err := m.SCookie.Decode("token", cookie.Value, &newUUID); err != nil {
				return err
			}

			session := new(models.Session)
			if err := m.DB.NewSelect().
				Model(session).
				Relation("User").
				Where("uuid = ?", newUUID).
				Scan(r.Context()); err != nil {
				return err
			}

			r2 := r.WithContext(context.WithValue(r.Context(), UserKey{}, session.User))
			next.ServeHTTP(w, r2)
			return nil
		})
	})
	r.Get("/", m.Handler.HandleFunc(m.secret))

	return r
}

func (m Me) secret(w http.ResponseWriter, r *http.Request) error {
	user := r.Context().Value(UserKey{}).(*models.User)
	return xmate.WriteHTML(w, m.SecretView, http.StatusOK, user)
}
