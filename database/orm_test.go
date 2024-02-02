package database

import (
	"context"
	"errors"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	contractsorm "github.com/goravel/framework/contracts/database/orm"
	"github.com/goravel/framework/database/gorm"
	"github.com/goravel/framework/database/orm"
	"github.com/goravel/framework/support/env"
	"github.com/goravel/framework/support/file"
)

var connections = []contractsorm.Driver{
	contractsorm.DriverMysql,
	contractsorm.DriverPostgresql,
	contractsorm.DriverSqlite,
	contractsorm.DriverSqlserver,
}

type User struct {
	orm.Model
	orm.SoftDeletes
	Name   string
	Avatar string
}

type OrmSuite struct {
	suite.Suite
	orm             *OrmImpl
	mysqlQuery      contractsorm.Query
	postgresqlQuery contractsorm.Query
	sqliteQuery     contractsorm.Query
	sqlserverDB     contractsorm.Query
}

func TestOrmSuite(t *testing.T) {
	if env.IsWindows() {
		t.Skip("Skipping tests of using docker")
	}

	if err := testDatabaseDocker.Fresh(); err != nil {
		t.Fatal(err)
	}

	mysqlDocker := gorm.NewMysqlDocker(testDatabaseDocker)
	mysqlQuery, err := mysqlDocker.New()
	if err != nil {
		log.Fatalf("Init mysql docker error: %v", err)
	}

	postgresqlDocker := gorm.NewPostgresqlDocker(testDatabaseDocker)
	postgresqlQuery, err := postgresqlDocker.New()
	if err != nil {
		log.Fatalf("Init postgresql docker error: %v", err)
	}

	sqliteDocker := gorm.NewSqliteDocker("goravel")
	sqliteQuery, err := sqliteDocker.New()
	if err != nil {
		log.Fatalf("Get sqlite error: %s", err)
	}

	sqlserverDocker := gorm.NewSqlserverDocker(testDatabaseDocker)
	sqlserverQuery, err := sqlserverDocker.New()
	if err != nil {
		log.Fatalf("Init sqlserver docker error: %v", err)
	}

	suite.Run(t, &OrmSuite{
		mysqlQuery:      mysqlQuery,
		postgresqlQuery: postgresqlQuery,
		sqliteQuery:     sqliteQuery,
		sqlserverDB:     sqlserverQuery,
	})

	assert.Nil(t, file.Remove("goravel"))
}

func (s *OrmSuite) SetupTest() {
	s.orm = &OrmImpl{
		ctx:   context.Background(),
		query: s.mysqlQuery,
		queries: map[string]contractsorm.Query{
			contractsorm.DriverMysql.String():      s.mysqlQuery,
			contractsorm.DriverPostgresql.String(): s.postgresqlQuery,
			contractsorm.DriverSqlite.String():     s.sqliteQuery,
			contractsorm.DriverSqlserver.String():  s.sqlserverDB,
		},
	}
}

func (s *OrmSuite) TestConnection() {
	for _, connection := range connections {
		s.NotNil(s.orm.Connection(connection.String()))
	}
}

func (s *OrmSuite) TestDB() {
	db, err := s.orm.DB()
	s.NotNil(db)
	s.Nil(err)

	for _, connection := range connections {
		db, err := s.orm.Connection(connection.String()).DB()
		s.NotNil(db)
		s.Nil(err)
	}
}

func (s *OrmSuite) TestQuery() {
	s.NotNil(s.orm.Query())

	s.NotPanics(func() {
		for i := 0; i < 5; i++ {
			go func() {
				var user User
				_ = s.orm.Query().Find(&user, 1)
			}()
		}
	})

	for _, connection := range connections {
		s.NotNil(s.orm.Connection(connection.String()).Query())
	}
}

func (s *OrmSuite) TestFactory() {
	s.NotNil(s.orm.Factory())

	for _, connection := range connections {
		s.NotNil(s.orm.Connection(connection.String()).Factory())
	}
}

func (s *OrmSuite) TestObserve() {
	s.orm.Observe(User{}, &UserObserver{})

	s.Equal([]orm.Observer{
		{Model: User{}, Observer: &UserObserver{}},
	}, orm.Observers)

	for _, connection := range connections {
		user := User{Name: "observer_name"}
		s.EqualError(s.orm.Connection(connection.String()).Query().Create(&user), "error")
	}
}

func (s *OrmSuite) TestTransactionSuccess() {
	for _, connection := range connections {
		user := User{Name: "transaction_success_user", Avatar: "transaction_success_avatar"}
		user1 := User{Name: "transaction_success_user1", Avatar: "transaction_success_avatar1"}
		s.Nil(s.orm.Connection(connection.String()).Transaction(func(tx contractsorm.Transaction) error {
			s.Nil(tx.Create(&user))
			s.Nil(tx.Create(&user1))

			return nil
		}))

		var user2, user3 User
		s.Nil(s.orm.Connection(connection.String()).Query().Find(&user2, user.ID))
		s.Nil(s.orm.Connection(connection.String()).Query().Find(&user3, user1.ID))
	}
}

func (s *OrmSuite) TestTransactionError() {
	for _, connection := range connections {
		s.NotNil(s.orm.Connection(connection.String()).Transaction(func(tx contractsorm.Transaction) error {
			user := User{Name: "transaction_error_user", Avatar: "transaction_error_avatar"}
			s.Nil(tx.Create(&user))

			user1 := User{Name: "transaction_error_user1", Avatar: "transaction_error_avatar1"}
			s.Nil(tx.Create(&user1))

			return errors.New("error")
		}))

		var users []User
		s.Nil(s.orm.Connection(connection.String()).Query().Find(&users))
		s.Equal(0, len(users))
	}
}

type UserObserver struct{}

func (u *UserObserver) Retrieved(event contractsorm.Event) error {
	return nil
}

func (u *UserObserver) Creating(event contractsorm.Event) error {
	name := event.GetAttribute("name")
	if name != nil && name.(string) == "observer_name" {
		return errors.New("error")
	}

	return nil
}

func (u *UserObserver) Created(event contractsorm.Event) error {
	return nil
}

func (u *UserObserver) Updating(event contractsorm.Event) error {
	return nil
}

func (u *UserObserver) Updated(event contractsorm.Event) error {
	return nil
}

func (u *UserObserver) Saving(event contractsorm.Event) error {
	return nil
}

func (u *UserObserver) Saved(event contractsorm.Event) error {
	return nil
}

func (u *UserObserver) Deleting(event contractsorm.Event) error {
	return nil
}

func (u *UserObserver) Deleted(event contractsorm.Event) error {
	return nil
}

func (u *UserObserver) ForceDeleting(event contractsorm.Event) error {
	return nil
}

func (u *UserObserver) ForceDeleted(event contractsorm.Event) error {
	return nil
}
