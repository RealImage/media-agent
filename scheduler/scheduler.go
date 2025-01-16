package scheduler

import (
	"context"
	"fmt"
	"media-agent/logger"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Scheduler struct {
	gocron.Scheduler
}

func NewScheduler(log *zap.Logger) (*Scheduler, error) {
	sch, err := gocron.NewScheduler(
		gocron.WithGlobalJobOptions(
			gocron.WithEventListeners(
				gocron.BeforeJobRuns(func(jobID uuid.UUID, jobName string) {
					log.Info("task started", zap.String("task_name", jobName),
						zap.String("job_id", jobID.String()))
				}),
				gocron.AfterJobRuns(func(jobID uuid.UUID, jobName string) {
					log.Info("task completed", zap.String("task_name", jobName),
						zap.String("job_id", jobID.String()))
				}),
				gocron.AfterJobRunsWithPanic(func(jobID uuid.UUID, jobName string, r interface{}) {
					log.Error("task panicked", zap.String("task_name", jobName),
						zap.Any("panic", r),
						zap.String("job_id", jobID.String()))
				}),
			),
			gocron.WithSingletonMode(gocron.LimitModeReschedule),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("error creating scheduler: %w", err)
	}
	return &Scheduler{sch}, nil
}

func (s *Scheduler) ScheduleJobs(ctx context.Context, facilityIds []string, interval time.Duration, task func(context.Context, string, *zap.Logger)) error {
	log := logger.GetLogger()

	log.Info("scheduling jobs", zap.Strings("facility_ids", facilityIds))

	for _, facilityId := range facilityIds {
		logger := log.With(zap.String("facility_id", facilityId))

		_, err := s.NewJob(
			gocron.DurationJob(
				interval,
			),
			gocron.NewTask(
				task,
				ctx,
				facilityId,
				logger,
			),
		)
		if err != nil {
			return fmt.Errorf("error scheduling job: %w", err)
		}
	}

	return nil
}
