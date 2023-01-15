package handlers

import (
	"fmt"
	"time"

	"github.com/klajbard/go-auth/config"
	"gopkg.in/mgo.v2/bson"
)

type Session struct {
	Token    string    // `json:"token"`
	Username string    // `json:"username"`
	Expiry   time.Time // `json:"expiry"`
}

type RefreshSession struct {
	Token        string    // `json:"token"`
	SessionToken string    // `json:"sessiontoken"`
	Expiry       time.Time // `json:"expiry"`
}

func GetSession(token string) (Session, bool) {
	fmt.Println(token)
	session := Session{}

	err := config.Sessions.Find(bson.M{"token": token}).One(&session)
	if err != nil {
		return Session{}, false
	}

	return session, true
}

func (s Session) isExpired() bool {
	return s.Expiry.Before(time.Now())
}

func (s Session) Create() bool {
	err := config.Sessions.Insert(bson.M{"token": s.Token, "username": s.Username, "expiry": s.Expiry})
	if err != nil {
		fmt.Println(err)
	}

	return err == nil
}

func (s Session) Remove() bool {
	err := config.Sessions.Find(bson.M{"token": s.Token})
	if err != nil {
		fmt.Println(err)
	}

	return err == nil
}

func (s Session) FindRefreshSession() bool {
	err := config.RefreshSessions.Find(bson.M{"sessiontoken": s.Token})
	if err != nil {
		fmt.Println(err)
	}

	return err == nil
}

func (s RefreshSession) isExpired() bool {
	return s.Expiry.Before(time.Now())
}

func (s RefreshSession) Create() bool {
	err := config.RefreshSessions.Insert(bson.M{"token": s.Token, "expiry": s.Expiry})
	if err != nil {
		fmt.Println(err)
	}

	return err == nil
}

func (s RefreshSession) Remove() bool {
	err := config.RefreshSessions.Remove(bson.M{"token": s.Token})
	if err != nil {
		fmt.Println(err)
	}

	return err == nil
}

func (s RefreshSession) Update(sessionToken string) bool {
	fmt.Println("updating", s.Token, "with", sessionToken)
	err := config.RefreshSessions.Update(bson.M{"token": s.Token}, bson.M{"$set": bson.M{"sessiontoken": sessionToken}})
	if err != nil {
		fmt.Println(err)
	}

	return err == nil
}

func GetRefreshSession(token string) (RefreshSession, bool) {
	refreshSession := RefreshSession{}

	err := config.RefreshSessions.Find(bson.M{"token": token}).One(&refreshSession)
	fmt.Println("asd", token, refreshSession, err != nil)
	if err != nil {
		return RefreshSession{}, false
	}

	return refreshSession, true
}
