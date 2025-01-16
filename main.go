package main

import (
	"context"
	"media-agent/client"
	"media-agent/config"
	"media-agent/logger"
	"media-agent/scheduler"
	"media-agent/services"
	"media-agent/tasks"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"go.uber.org/zap"
)

var (
	cfg *config.EnvConfig
)

func init() {
	cfg = config.LoadConfig()

	logger.Init()
}

func main() {
	defer func() {
		logger.Sync()
	}()

	s3Client := client.NewS3HttpClient()

	xpCredentials := services.NewCredentials(
		cfg.XPUSERNAME,
		cfg.XPPASSWORD,
	)

	dolbyCredentials := services.NewCredentials(
		cfg.DOLBYUSERNAME,
		cfg.DOLBYPASSWORD,
	)

	facilityService := services.NewFacilityService(cfg)

	fetchMediaTask := tasks.NewFetchMediaTask(facilityService, xpCredentials, dolbyCredentials, s3Client, cfg)

	log := logger.GetLogger()

	var Commit = func() string {
		if info, ok := debug.ReadBuildInfo(); ok {
			for _, setting := range info.Settings {
				if setting.Key == "vcs.revision" {
					return setting.Value
				}
			}
		}

		return ""
	}()

	sc, err := scheduler.NewScheduler(log)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	sc.ScheduleJobs(ctx, []string{"2590463e-4ab8-4b6e-aa72-2ea67985acf2"}, 1*time.Minute, fetchMediaTask.FetchMedia)

	sc.Start()

	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		log.Info("Application stopped")
		sc.Shutdown()
		cancel() // Cancel the context
	}()

	log.Info("scheduling jobs", zap.String("commitId", Commit))

	<-ctx.Done()

}
