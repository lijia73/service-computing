package service

import (
	"net/http"

	"github.com/unrolled/render"
)

func apiTableHandler(formatter *render.Render) http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		formatter.HTML(w, http.StatusOK, "table", struct {
			Username string 
			Password string 
		}{Username: req.Form["username"][0], Password: req.Form["password"][0]})
	}
}