package view

import "net/http"

func (v *view) Login(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusContinue)
}

func (v *view) Refresh(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:  "testcookie",
		Value: "testvalue",
	}
	http.SetCookie(w, &cookie)
}
