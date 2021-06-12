package db

import "github.com/victor-nach/time-tracker/models"

//Datastore defines the required store methods
type Datastore interface {
	createUser(rate *models.User) (*models.User, error)
	GetUser(id string) (*models.User, error)

	GetSession(id, owner string) (*models.Session, error)
	GetSessions(owner string, filter sessionFilter) ([]*models.Session, error)

	createSession(rate *models.Session) (*models.Session, error)
	updateSession(id string, info models.SessionInfo) error
	DeleteSession(id string) error
}

type sessionFilter string

const (
	Day   sessionFilter = "day"
	Week  sessionFilter = "week"
	Month sessionFilter = "month"
)