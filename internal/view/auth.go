package view

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/1ef7yy/medods_test_task/internal/errors"
	"github.com/1ef7yy/medods_test_task/models"
	"github.com/1ef7yy/medods_test_task/pkg/utils"
)

func (v *view) Login(w http.ResponseWriter, r *http.Request) {
	guid := r.URL.Query().Get("guid")

	ip := strings.Split(r.RemoteAddr, ":")[0]

	if guid == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("guid not in request"))
		if err != nil {
			v.log.Errorf("error writing to client: %s", err.Error())
		}
		return
	}

	req := models.GenerateTokenRequest{
		IP:         ip,
		Guid:       guid,
		Generation: 1,
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

	resp, err := json.Marshal(struct {
		AccessToken string `json:"access_token"`
	}{
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
	refreshCookie, err := r.Cookie("refresh_token")
	if err != nil {
		v.log.Errorf("error getting refresh cookie :%s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// workaround for scanning
	var accessToken struct {
		Data string `json:"access_token"`
	}

	err = json.NewDecoder(r.Body).Decode(&accessToken)
	if err != nil {
		v.log.Errorf("error decoding body: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	refreshToken, err := utils.Base64ToString(refreshCookie.Value)
	if err != nil {
		v.log.Errorf("error decoding base64 cookie: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token := models.Token{
		RefreshToken: refreshToken,
		AccessToken:  accessToken.Data,
	}

	ip := strings.Split(r.RemoteAddr, ":")[0]

	req := models.RefreshTokenRequest{
		IP:     ip,
		Tokens: token,
	}

	newToken, err := v.domain.Refresh(r.Context(), req)
	if err == errors.TokenInvalidErr {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("token is not valid"))
		if err != nil {
			v.log.Errorf("error writing to client: %s", err)
		}
		return
	}

	if err == errors.GuidIsDifferentErr {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("guid in the tokens are different"))
		if err != nil {
			v.log.Errorf("error writing to client: %s", err)
		}
		return
	}

	if err != nil {
		v.log.Errorf("error refreshing access token %s with refresh token %s: %s", refreshToken, accessToken, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	newRefreshCookie := http.Cookie{
		Name:     "refresh_token",
		Value:    utils.StringToBase64(newToken.RefreshToken),
		MaxAge:   604800,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, &newRefreshCookie)

	resp, err := json.Marshal(struct {
		AccessToken string `json:"access_token"`
	}{
		AccessToken: newToken.AccessToken,
	})

	if err != nil {
		v.log.Errorf("error marshalling: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = w.Write(resp)

	if err != nil {
		v.log.Errorf("error writing to client: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
}
