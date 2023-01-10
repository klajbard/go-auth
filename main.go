package main

import (
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	"github.com/klajbard/go-auth/config"
	"github.com/klajbard/go-auth/handlers"
	"github.com/klajbard/go-auth/types"
	"github.com/rs/cors"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	Username string // `json:"username`
	Password string // `json:"password`
}

func getUser(u string) (User, bool) {
	user := User{}

	err := config.Users.Find(bson.M{"username": u}).One(&user)
	if err != nil {
		return User{}, false
	}

	return user, true
}

func main() {
	mux := httprouter.New()

	mux.GET("/", handlers.HandleStaticPage("client/dist/index.html"))
	mux.GET("/logout", handlers.HandleStaticPage("client/dist/logout/index.html"))
	mux.GET("/services", handlers.WithAuth(handlers.HandleStaticPage("client/dist/services/index.html")))
	mux.POST("/login", login)
	mux.POST("/logout", logout)
	mux.GET("/hello", hello)

	mux.ServeFiles("/assets/*filepath", http.Dir("client/dist/assets"))
	mux.NotFound = http.HandlerFunc(handlers.NotFoundHandler)

	handler := cors.New(cors.Options{
		AllowCredentials: true,
		AllowedOrigins:   []string{os.Getenv("CLIENT_DEV_URL")},
	}).Handler(mux)
	http.ListenAndServe(":8080", handler)
}

func hello(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	request := handlers.RequestHandler{R: *r, W: w}
	status, userSession := request.GetUserSession()
	statusMessage := types.Status{}

	w.WriteHeader(status)

	switch status {
	case http.StatusOK:
		statusMessage.Message = fmt.Sprintf("Welcome %s!", userSession.Username)
	default:
		statusMessage.Message = "Authentication failed!"
	}

	json.NewEncoder(w).Encode(&status)
}

func login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var u User
	var pwMatched bool
	request := handlers.RequestHandler{R: *r, W: w}

	c, err := r.Cookie("refresh_token")

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	status := types.Status{
		Message: "Username or password is incorrect",
	}
	err = json.NewDecoder(r.Body).Decode(&u)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	defer r.Body.Close()

	user, ok := getUser(u.Username)

	if ok {
		hash := sha512.New()
		hash.Write([]byte(u.Password))

		if err != nil {
			fmt.Println(err)
		}
		pwHashed := hex.EncodeToString(hash.Sum(nil))

		pwMatched = pwHashed == user.Password
	}

	if pwMatched {
		cSession := request.CreateSessionToken(u.Username)

		http.SetCookie(w, cSession)

		if refreshSession, ok := handlers.RefreshSessions[c.Value]; ok {
			refreshSession.SessionToken = cSession.Value
			handlers.RefreshSessions[c.Value] = refreshSession
		}

		status.Message = "Successfully logged in!"
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}

	json.NewEncoder(w).Encode(&status)
}

func logout(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	request := handlers.RequestHandler{R: *r, W: w}
	cRefresh, okRefresh := request.GetRefreshToken()
	cSession, okSession := request.GetSessionToken()

	if !(okSession && okRefresh) {
		fmt.Println("something went wrong while logging out")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if cRefresh.Value != "" {
		request.RemoveCookie("refresh_token")
		delete(handlers.RefreshSessions, cRefresh.Value)
	}
	if cSession.Value != "" {
		request.RemoveCookie("session_token")
		delete(handlers.Sessions, cSession.Value)
	}
	w.WriteHeader(http.StatusAccepted)
}
