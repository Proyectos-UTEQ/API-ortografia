package handlers

import (
	"context"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"github.com/spf13/viper"
	"net/http"
)

const (
	key    = "random-key"
	MaxAge = 86400 * 30
	IsProd = false
)

type AuthHandler struct {
	config *viper.Viper
}

func NewAuthHandler(config *viper.Viper) *AuthHandler {
	return &AuthHandler{
		config: config,
	}
}

func (h *AuthHandler) ConfigProvider() {

	googleClientID := h.config.GetString("GOOGLE_CLIENT_ID")
	googleClientSecret := h.config.GetString("GOOGLE_CLIENT_SECRET")

	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(MaxAge)

	store.Options.Path = "http://localhost:5173/"
	store.Options.HttpOnly = true
	store.Options.Secure = IsProd

	gothic.Store = store
	goth.UseProviders(
		google.New(googleClientID, googleClientSecret, "http://localhost:3000/api/auth/google/callback"),
	)
}

func (h *AuthHandler) BeginAuthGoogle(w http.ResponseWriter, r *http.Request) {
	provider := "google"
	r = r.WithContext(context.WithValue(context.Background(), "provider", provider))
	// Probamos la autenticación
	if gothUser, err := gothic.CompleteUserAuth(w, r); err == nil {
		fmt.Println(gothUser)

	} else {
		gothic.BeginAuthHandler(w, r)
	}
	return
}

func (h *AuthHandler) GetAuthCallbackFunction(w http.ResponseWriter, r *http.Request) {
	provider := "google"
	r = r.WithContext(context.WithValue(context.Background(), "provider", provider))
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		return
	}

	fmt.Println(user)

	// redirigir a la página de inicio
	http.Redirect(w, r, "http://localhost:5173/onboard", http.StatusFound)
}

func (h *AuthHandler) GetAuthSuccessFunction(w http.ResponseWriter, r *http.Request) {
	provider := "google"
	r = r.WithContext(context.WithValue(context.Background(), "provider", provider))
	gothUser, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	fmt.Fprintln(w, gothUser)
}
