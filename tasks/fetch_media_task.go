package tasks

import (
	"context"
	"encoding/json"
	"media-agent/client"
	"media-agent/config"
	"media-agent/dtos"
	"media-agent/enums"
	"media-agent/services"
	"sync"

	"github.com/RealImage/libsms/v2/clientutil"
	"go.uber.org/zap"
)

type FetchMediaTask struct {
	*services.FacilityService
	xpCredentials    *services.Credentials
	dolbyCredentials *services.Credentials
	s3client         *client.S3HttpClient
	cfg              *config.EnvConfig
}

func NewFetchMediaTask(facilityService *services.FacilityService, xpCrdentials *services.Credentials, dolbyCredentials *services.Credentials, s3Client *client.S3HttpClient, cfg *config.EnvConfig) *FetchMediaTask {
	return &FetchMediaTask{
		FacilityService:  facilityService,
		xpCredentials:    xpCrdentials,
		dolbyCredentials: dolbyCredentials,
		s3client:         s3Client,
		cfg:              cfg,
	}
}

func (f *FetchMediaTask) FetchMedia(ctx context.Context, facilityID string, log *zap.Logger) {
	log.Info("Fetching media")

	facility, err := f.GetFaciltiyById(facilityID)
	if err != nil {
		log.Error("Failed to get facility", zap.Error(err))
		return
	}

	wg := sync.WaitGroup{}

	for _, auditorium := range facility.PlayBackServer {
		if auditorium.Type == enums.ServerTypeQubeXP4 {

			wg.Add(1)
			go func() {
				defer wg.Done()

				log.Info("Fetching media info for auditorium", zap.String("auditoriumID", auditorium.ID))

				client, err := clientutil.NewClientContext(ctx, auditorium.IP, f.xpCredentials.GetCredentials())
				if err != nil {
					log.Error("Failed to create client", zap.Error(err), zap.String("auditoriumID", auditorium.ID))
					return
				}

				mediaInfo, err := services.NewPlayBackServer(client).GetMediaInfo(facility)
				if err != nil {
					log.Error("Failed to get media info", zap.Error(err), zap.String("auditoriumID", auditorium.ID))
					return
				}

				err = f.UploadMedia(ctx, "media-info", auditorium.ID, mediaInfo)
				if err != nil {
					log.Error("Failed to upload media info", zap.Error(err), zap.String("auditoriumID", auditorium.ID))
					return
				}

				log.Info("Media info uploaded successfully", zap.String("auditoriumID", auditorium.ID))
			}()
		}
	}

	wg.Wait()
}

func (f *FetchMediaTask) UploadMedia(ctx context.Context, bucketName, key string, mediaInfo *dtos.MediaInfo) error {
	presignedURL, err := f.s3client.CreatePresignedURL(ctx, f.cfg.MEDIAS3LAMBDA, bucketName, key)
	if err != nil {
		return err
	}

	jsonData, err := json.Marshal(mediaInfo)
	if err != nil {
		return err
	}

	return f.s3client.UploadMedia(ctx, presignedURL, jsonData)
}
