package v1

import "net/http"

func (v *Router) Auth() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("POST /login", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v.View.Login(w, r)
	}))
	mux.Handle("POST /refresh", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v.View.Refresh(w, r)
	}))

	return http.StripPrefix("/api/v1/auth", mux)
}
