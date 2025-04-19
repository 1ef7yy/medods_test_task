package view

import (
	"encoding/json"
	"net/http"

	"github.com/1ef7yy/medods_test_task/models"
	"github.com/1ef7yy/medods_test_task/pkg/utils"
)

func (v *view) Login(w http.ResponseWriter, r *http.Request) {
	guid := r.URL.Query().Get("guid")

	ip := r.RemoteAddr

	if guid == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("guid not in request"))
		if err != nil {
			v.log.Errorf("error writing to client: %s", err.Error())
		}
		return
	}

	req := models.GenerateTokenRequest{
		IP:   ip,
		Guid: guid,
	}

	tokens, err := v.domain.Login(r.Context(), req)

	if err != nil {
		v.log.Errorf("error logging in by guid: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if tokens.RefreshToken == "" || tokens.AccessToken == "" {
		v.log.Warnf("one or both of tokens are empty: %s %s", tokens.AccessToken, tokens.RefreshToken)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	refreshCookie := http.Cookie{
		Name:     "refresh_token",
		Value:    utils.StringToBase64(tokens.RefreshToken),
		MaxAge:   604800, // 7 days
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, &refreshCookie)

	resp, err := json.Marshal(struct{
		AccessToken string `json:"access_token"`
	} {
		AccessToken: tokens.AccessToken,
	})

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
