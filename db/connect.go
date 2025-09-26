package db

import (
	"context"
	"log"
	"log/slog"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Conn *gorm.DB
var Roles gorm.Interface[Role]
var Guilds gorm.Interface[Guild]
var Users gorm.Interface[User]
var Projects gorm.Interface[Project]
var Issues gorm.Interface[Issue]
var Relationships gorm.Interface[Relationship]
var ProjectViewStates gorm.Interface[ProjectViewState]
var Tags gorm.Interface[Tag]

var Ctx = context.Background()

func Connect(path string) (*gorm.DB, error) {
	slog.Info("Connecting to db at", "path", path)
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Warn, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,        // Don't include params in the SQL log
			Colorful:                  false,       // Disable color
		},
	)
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{Logger: newLogger})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&Role{}, &Guild{}, &User{}, &Project{}, &Issue{}, &Relationship{}, &ProjectViewState{}, &Tag{})
	Conn = db

	Roles = gorm.G[Role](Conn)
	Guilds = gorm.G[Guild](Conn)
	Users = gorm.G[User](Conn)
	Projects = gorm.G[Project](Conn)
	Issues = gorm.G[Issue](Conn)
	Relationships = gorm.G[Relationship](Conn)
	ProjectViewStates = gorm.G[ProjectViewState](Conn)
	Tags = gorm.G[Tag](Conn)

	slog.Info("Connected to db")
	return db, nil
}
