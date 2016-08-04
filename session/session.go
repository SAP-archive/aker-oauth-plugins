package session

import "github.com/gorilla/sessions"

//go:generate counterfeiter . Store

type Store interface {
	sessions.Store
}
