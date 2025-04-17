package routes

import (
	"net/http"

	v1 "github.com/1ef7yy/medods_test_task/internal/routes/v1"
	"github.com/1ef7yy/medods_test_task/internal/view"
)

func InitRouter(view view.View) *http.ServeMux {
	mux := http.NewServeMux()
	v1 := v1.NewRouter(view)

	mux.Handle("/api/v1/auth/", v1.Auth())

	return mux
}
