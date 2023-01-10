package handlers

import (
	"net/http"
	"text/template"

	"github.com/julienschmidt/httprouter"
	"github.com/klajbard/go-auth/types"
)

func HandleStaticPage(route string) func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		index := template.Must(template.ParseFiles(route))
		index.Execute(w, nil)
		defer r.Body.Close()
	}
}

func WithAuth(handlerFn func(w http.ResponseWriter, r *http.Request, p httprouter.Params)) func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		defer r.Body.Close()
		request := RequestHandler{R: *r, W: w}
		status, _ := request.GetUserSession()
		statusMessage := types.Status{}

		switch status {
		case http.StatusOK:
			handlerFn(w, r, p)
		default:
			statusMessage.Message = "Authentication failed!"
			http.Redirect(w, r, "/", 303)
			return
		}
	}
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	index := template.Must(template.ParseFiles("client/dist/404.html"))
	index.Execute(w, nil)
	defer r.Body.Close()
}
