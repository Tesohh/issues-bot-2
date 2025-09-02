package db

import (
	"context"
	"log/slog"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var Conn *gorm.DB

var Roles gorm.Interface[Role]
var Guilds gorm.Interface[Guild]
var Users gorm.Interface[User]
var Projects gorm.Interface[Project]
var Issues gorm.Interface[Issue]
var Relationships gorm.Interface[Relationship]
var ProjectViewStates gorm.Interface[ProjectViewState]

var Ctx = context.Background()

func Connect(path string) (*gorm.DB, error) {
	slog.Info("Connecting to db at", "path", path)
	db, err := gorm.Open(sqlite.Open(path))
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&Role{}, &Guild{}, &User{}, &Project{}, &Issue{}, &Relationship{}, &ProjectViewState{})
	Conn = db

	Roles = gorm.G[Role](Conn)
	Guilds = gorm.G[Guild](Conn)
	Users = gorm.G[User](Conn)
	Projects = gorm.G[Project](Conn)
	Issues = gorm.G[Issue](Conn)
	Relationships = gorm.G[Relationship](Conn)
	ProjectViewStates = gorm.G[ProjectViewState](Conn)

	slog.Info("Connected to db")
	return db, nil
}
