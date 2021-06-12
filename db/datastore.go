package db

import "github.com/victor-nach/time-tracker/models"

//Datastore defines the required store methods
type Datastore interface {
	CreateUser(user *models.User) (*models.User, error)
	GetUser(id string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)

	GetSession(id, owner string) (*models.Session, error)
	GetSessions(owner string, filter sessionFilter) ([]*models.Session, error)

	CreateSession(session *models.Session) (*models.Session, error)
	UpdateSession(id string, info models.SessionInfo) error
	DeleteSession(id string) error
}

type sessionFilter string

const (
	Day   sessionFilter = "day"
	Week  sessionFilter = "week"
	Month sessionFilter = "month"
)