package internal

import (
	"TL-ID-Service/config"
	"TL-ID-Service/log"
	"TL-ID-Service/proto/id"
	"context"
	"database/sql"
	"fmt"

	uuid "github.com/satori/go.uuid"

	// for postgreq sql
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

// Generate32Bit implements id.Service
func (s *Service) Generate32Bit(ctx context.Context, request *id.Generate32BitRequest) (*id.Generate32BitReply, error) {
	// generate the sequence id
	next, source, err := s.generate(ctx)
	if err != nil {
		log.Infof("generate user id and source: %v", err)
		return &id.Generate32BitReply{Status: id.Status_Error}, nil
	}

	return &id.Generate32BitReply{Status: id.Status_Success, Id: next, Source: source}, nil
}

// GetSource implements id.Service
func (s *Service) GetSource(ctx context.Context, request *id.GetSourceRequest) (*id.GetSourceReply, error) {
	var source string

	// query the record with user id
	row := s.db.QueryRowContext(ctx, "select source from tl_user where id=$1", request.Id)
	if err := row.Scan(&source); err != nil {
		log.Infof("query source by id %d: %v", request.Id, err)
		return &id.GetSourceReply{Status: id.Status_Error}, nil
	}

	return &id.GetSourceReply{Status: id.Status_Success, Source: source}, nil
}

// generate a 32bit id
func (s *Service) generate(ctx context.Context) (id int32, source string, err error) {
	// generate the user's source
	source = uuid.NewV4().String()

	// generate the id frist
	row := s.db.QueryRowContext(ctx, "select nextval('tl_user_id_seq');")
	if err = row.Scan(&id); err != nil {
		return
	}

	log.Infof("id: %d, source: %s", id, source)

	var result sql.Result
	// insert the id to table tl_user
	result, err = s.db.ExecContext(ctx, "insert into tl_user values($1, $2)", id, source)
	if err != nil {
		return
	}
	if affected, _ := result.RowsAffected(); affected == 0 {
		err = fmt.Errorf("no rows affected")
		return
	}

	return
}
