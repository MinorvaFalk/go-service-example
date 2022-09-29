package servicea

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/MinorvaFalk/go-service-example/datasource"
	"github.com/MinorvaFalk/go-service-example/utils/logger"
)

type ServiceA struct {
	l  *logger.Logger
	db *datasource.DB
}

func NewServiceA(l *logger.Logger, db *datasource.DB) *ServiceA {
	return &ServiceA{
		l:  l,
		db: db,
	}
}

func (s *ServiceA) DoSomething() {
	fmt.Println("Hello world from service-a")
}

func (s *ServiceA) CopyToCSV(ctx context.Context) {
	conn := s.db.Conn
	defer conn.Close(context.Background())

	query := `
		COPY (SELECT * 
			FROM users 
		) TO STDOUT DELIMITER ',' CSV HEADER
	`

	file := s.openFile()
	defer file.Close()

	res, err := conn.PgConn().CopyTo(ctx, file, query)
	if err != nil {
		panic(err)
	}

	s.l.Sugar.Infof("row affected: %v", res.RowsAffected())
}

func (s *ServiceA) openFile() *os.File {
	fileName := "../files/user.csv"

	file, err := os.OpenFile(
		filepath.Join(rootDir(), fileName),
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
		0666,
	)
	if err != nil {
		s.l.Sugar.Panic(err)
	}

	return file
}

func rootDir() string {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	return filepath.Dir(d)
}
