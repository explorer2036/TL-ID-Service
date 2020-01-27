package internal

import (
	"TL-ID-Service/config"
	"TL-ID-Service/proto/id"
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

const (
	// DefaultMaxOpenConns - max open connections for database
	DefaultMaxOpenConns = 20
	// DefaultMaxIdleConns - max idle connections for database
	DefaultMaxIdleConns = 5
)

// Service represents the interfaces for id service
type Service struct {
	id.UnimplementedServiceServer

	db       *sql.DB        // database connection pool
	settings *config.Config // settings for the id service
}

// NewService create the internal service
func NewService(settings *config.Config) *Service {
	s := &Service{settings: settings}

	// init the database connections
	source := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		settings.DB.Host,
		settings.DB.Port,
		settings.DB.User,
		settings.DB.Passwd,
		settings.DB.Name,
	)
	db, err := sql.Open("postgres", source)
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(DefaultMaxOpenConns)
	db.SetMaxIdleConns(DefaultMaxIdleConns)
	if err := db.Ping(); err != nil {
		panic(err)
	}
	s.db = db

	return s
}

// Generate32Bit implements gateway.Service
func (s *Service) Generate32Bit(ctx context.Context, request *id.Generate32BitRequest) (*id.Generate32BitReply, error) {
	// generate the sequence id
	next, err := s.generate(ctx)
	if err != nil {
		return &id.Generate32BitReply{Status: id.Status_Error}, nil
	}

	return &id.Generate32BitReply{Status: id.Status_Success, Id: next}, nil
}

// generate a 32bit id
func (s *Service) generate(ctx context.Context) (id int32, err error) {
	// generate the id frist
	row := s.db.QueryRowContext(ctx, "select nextval('user_id_sequence');")
	if err = row.Scan(&id); err != nil {
		return
	}

	var result sql.Result
	// insert the id to table tl_user
	result, err = s.db.ExecContext(ctx, "insert into tl_user values($1)", id)
	if err != nil {
		return
	}
	if affected, _ := result.RowsAffected(); affected == 0 {
		err = fmt.Errorf("no rows affected")
		return
	}

	return
}
