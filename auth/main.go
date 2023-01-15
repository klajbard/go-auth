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

func getUser(u string) (types.User, bool) {
	user := types.User{}

	err := config.Users.Find(bson.M{"username": u}).One(&user)
	if err != nil {
		return types.User{}, false
	}

	return user, true
}

func main() {
	whiteListUrls := []string{os.Getenv("CLIENT_DEV_URL")}
	config.Conf.GetConf()
	mux := httprouter.New()

	mux.GET("/services", services)
	mux.GET("/auth", auth)
	mux.POST("/auth", login)
	mux.POST("/logout", logout)

	for _, v := range config.Conf.Services {
		whiteListUrls = append(whiteListUrls, v.URL)
	}

	handler := cors.New(cors.Options{
		AllowCredentials: true,
		AllowedOrigins:   whiteListUrls,
	}).Handler(mux)
	http.ListenAndServe(":8080", handler)
}

func services(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	request := handlers.RequestHandler{R: *r, W: w}
	isValidSession := request.ValidateSession()

	if isValidSession {
		json.NewEncoder(w).Encode(&config.Conf.Services)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(&types.Status{Message: "Unauthorized"})
	}
}

func auth(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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
	var u types.User
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

		if refreshSession, ok := handlers.GetRefreshSession(c.Value); ok && u.RememberMe {
			refreshSession.Update(cSession.Value)
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
		if refreshSession, ok := handlers.GetRefreshSession(cRefresh.Value); ok {
			refreshSession.Remove()
		}
	}
	if cSession.Value != "" {
		request.RemoveCookie("session_token")
		if Session, ok := handlers.GetSession(cRefresh.Value); ok {
			Session.Remove()
		}
	}
	w.WriteHeader(http.StatusAccepted)
}
