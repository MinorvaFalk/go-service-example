package main

import (
	"sync"

	"github.com/MinorvaFalk/go-service-example/config"
	"github.com/MinorvaFalk/go-service-example/datasource"
	"github.com/MinorvaFalk/go-service-example/job"
	"github.com/MinorvaFalk/go-service-example/service"
	"github.com/MinorvaFalk/go-service-example/utils/logger"
)

func main() {
	// Init logger
	l := logger.NewLogger()

	// Read env configs
	c := config.InitConfig()

	// Init datasources
	db := datasource.NewDB(c.Dsn)

	// Init service
	s := service.NewService(l, db)

	// Init Job
	job.InitJob(s, l)

	// Blocking operation for cron job
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
