package job

import (
	"context"
	"fmt"
	"time"

	"github.com/MinorvaFalk/go-service-example/service"
	"github.com/MinorvaFalk/go-service-example/utils/logger"
	"github.com/robfig/cron/v3"
)

func InitJob(s *service.Service, l *logger.Logger) {
	var cronOption []cron.Option

	// Init timezone for cron job
	tz, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		panic(fmt.Errorf("failed to load timezone\n%v", err))
	}

	cronLogger := logger.NewCronLogger(l.Logger)

	cronOption = append(
		cronOption,
		cron.WithLocation(tz),
		cron.WithLogger(cronLogger),
		cron.WithChain(
			cron.Recover(cronLogger),
		),
	)

	c := cron.New(cronOption...)
	// refer https://pkg.go.dev/github.com/robfig/cron for more info
	// about cronjob functions

	// Job every Mon-Sat at 23:50
	c.AddFunc("50 23 * * 1-6", func() {
		s.ServiceA().DoSomething()
	})

	// Job every 5s
	c.AddFunc("@every 10s", func() {
		s.ServiceB().DoSomething()
	})

	// Job every 1m
	c.AddFunc("@every 1m", func() {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, 40*time.Second)
		defer cancel()

		s.ServiceA().CopyToCSV(ctx)
	})

	c.Start()
}
