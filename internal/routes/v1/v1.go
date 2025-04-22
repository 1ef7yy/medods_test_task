package v1

import "github.com/1ef7yy/medods_test_task/internal/view"

type Router struct {
	View view.View
}

func NewRouter(view view.View) *Router {
	return &Router{
		View: view,
	}
}
