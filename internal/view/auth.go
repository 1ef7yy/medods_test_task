package view

import (
	"encoding/json"
	"net/http"
)

func (v *view) Login(w http.ResponseWriter, r *http.Request) {
	guid := r.URL.Query().Get("guid")

	if guid == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tokens, err := v.domain.Login(guid)

	if err != nil {
		v.log.Errorf("error logging in by guid: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	refreshCookie := http.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken,
		MaxAge:   604800, // 7 days
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, &refreshCookie)

	resp, err := json.Marshal(tokens.AccessToken)

	if err != nil {
		v.log.Errorf("error marshalling access token: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write(resp)

	if err != nil {
		v.log.Errorf("error writing to client: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (v *view) Refresh(w http.ResponseWriter, r *http.Request) {
}
