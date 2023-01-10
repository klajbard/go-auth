package handlers

import "time"

var Sessions = map[string]Session{}
var RefreshSessions = map[string]RefreshSession{}

type Session struct {
	Username string
	expiry   time.Time
}

type RefreshSession struct {
	SessionToken string
	expiry       time.Time
}

func (s Session) isExpired() bool {
	return s.expiry.Before(time.Now())
}

func (s RefreshSession) isExpired() bool {
	return s.expiry.Before(time.Now())
}
