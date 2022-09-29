package test

import (
	"context"
	"testing"
	"time"

	"github.com/MinorvaFalk/go-service-example/config"
	"github.com/MinorvaFalk/go-service-example/datasource"
	servicea "github.com/MinorvaFalk/go-service-example/service/service-a"
	"github.com/MinorvaFalk/go-service-example/utils/logger"
)

func TestServiceA(t *testing.T) {
	l := logger.NewLogger()
	c := config.InitConfig()
	db := datasource.NewDB(c.Dsn)

	serviceA := servicea.NewServiceA(l, db)

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	serviceA.CopyToCSV(ctx)
}
