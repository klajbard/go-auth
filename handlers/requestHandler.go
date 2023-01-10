package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type RequestHandler struct {
	R http.Request
	W http.ResponseWriter
}

func (handler *RequestHandler) RemoveCookie(value string) {
	cSession := &http.Cookie{
		Name:     value,
		MaxAge:   -1,
		Value:    "",
		Secure:   false,
		HttpOnly: true,
		Path:     "/",
	}
	http.SetCookie(handler.W, cSession)
}

func (handler *RequestHandler) CreateSessionToken(username string) *http.Cookie {
	newSessiontoken := uuid.NewString()
	expiresAt := time.Now().Add(2 * time.Hour)
	cSession := &http.Cookie{
		Name:     "session_token",
		Expires:  expiresAt,
		Value:    newSessiontoken,
		Secure:   false,
		HttpOnly: true,
		Path:     "/",
	}
	Sessions[newSessiontoken] = Session{
		Username: username,
		expiry:   expiresAt,
	}
	return cSession
}

func (handler *RequestHandler) GetSessionToken() (*http.Cookie, bool) {
	c, err := handler.R.Cookie("session_token")

	if err != nil {
		if err == http.ErrNoCookie {
			fmt.Println("no session token", err)
			return &http.Cookie{}, true
		} else {
			fmt.Println("unable to get session token", err)
			return &http.Cookie{}, false
		}
	}

	return c, true
}

func (handler *RequestHandler) CreateRefreshToken() *http.Cookie {
	refreshSessionTokenString := uuid.NewString()
	expiresAt := time.Now().AddDate(1, 0, 0)
	cRefresh := &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshSessionTokenString,
		Expires:  expiresAt,
		Secure:   false,
		HttpOnly: true,
		Path:     "/",
	}
	RefreshSessions[refreshSessionTokenString] = RefreshSession{
		expiry: expiresAt,
	}
	return cRefresh
}

func (handler *RequestHandler) GetRefreshToken() (*http.Cookie, bool) {
	c, err := handler.R.Cookie("refresh_token")

	if err != nil {
		if err == http.ErrNoCookie {
			fmt.Println("no refresh token", err)
			return &http.Cookie{}, true
		} else {
			fmt.Println("unable to get refresh token", err)
			return &http.Cookie{}, false
		}
	}

	return c, true
}

func (handler *RequestHandler) GetUserSession() (int, Session) {
	c, err := handler.R.Cookie("session_token")
	var userSession Session

	if err != nil {
		if err != http.ErrNoCookie {
			return http.StatusInternalServerError, Session{}
		}

		fmt.Println("unauthorized", err)

		cRefresh, ok := handler.GetRefreshToken()

		if !ok {
			fmt.Println("unable to create refresh token")
			return http.StatusInternalServerError, Session{}
		}

		if cRefresh.Value == "" {
			cRefresh = handler.CreateRefreshToken()
		}

		refreshSession, ok := RefreshSessions[cRefresh.Value]

		if !ok || refreshSession.isExpired() {
			delete(RefreshSessions, cRefresh.Value)
			fmt.Println("unauthorized - user refresh session not found or expired")
			cRefresh = handler.CreateRefreshToken()
			http.SetCookie(handler.W, cRefresh)
			return http.StatusUnauthorized, Session{}
		}
		http.SetCookie(handler.W, cRefresh)

		userSession, ok := Sessions[refreshSession.SessionToken]

		if !ok {
			fmt.Println("unauthorized - user session not found for refresh token")
			return http.StatusUnauthorized, Session{}
		}

		if userSession.isExpired() {
			fmt.Println("unauthorized - user session expired for refresh token")
		}
		fmt.Println("creating new session token from refresh token")
		cSession := handler.CreateSessionToken(userSession.Username)
		http.SetCookie(handler.W, cSession)
	} else {
		userSession, ok := Sessions[c.Value]

		if !ok {
			fmt.Println("unauthorized - user session not found")
			return http.StatusUnauthorized, Session{}
		}

		if userSession.isExpired() {
			fmt.Println("unauthorized - user session expired")
			matchingSessionToken := false
			delete(Sessions, c.Value)

			for _, rSession := range RefreshSessions {
				if rSession.SessionToken == c.Value {
					matchingSessionToken = true
				}
			}

			if !matchingSessionToken {
				return http.StatusUnauthorized, Session{}
			}

			cSession := handler.CreateSessionToken(userSession.Username)

			http.SetCookie(handler.W, cSession)
		}
	}

	return http.StatusOK, userSession
}
