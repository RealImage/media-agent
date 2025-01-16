package services

import (
	"context"
	"media-agent/dtos"
	"time"

	sms "github.com/RealImage/libsms/v2"
)

type PlayBackServer struct {
	client sms.Client
}

func NewPlayBackServer(client sms.Client) *PlayBackServer {
	return &PlayBackServer{client: client}
}

func (p *PlayBackServer) GetMediaInfo(device *dtos.Facility) (*dtos.MediaInfo, error) {
	cpls, err := p.client.ListCPLs(context.Background())
	if err != nil {
		return nil, err
	}

	kdms, err := p.client.ListKDMS(context.Background())
	if err != nil {
		return nil, err
	}

	return buildMediaInfoHelper(cpls, kdms), nil
}

func buildMediaInfoHelper(cpls []sms.CPLInfo, kdms []sms.KDMInfo) *dtos.MediaInfo {
	mediaInfo := &dtos.MediaInfo{
		LastUpdatedTime: time.Now().UTC(),
		CPLs:            make([]dtos.CPL, 0),
	}

	for _, cpl := range cpls {
		cplDTO := dtos.CPL{
			ID:                cpl.ID.String(),
			Name:              cpl.ContentTitleText,
			KeepUntil:         cpl.KeepUntil,
			ContentKind:       cpl.ContentKind,
			SizeInBytes:       int(cpl.SizeInBytes),
			DurationInSeconds: cpl.Duration.Seconds(),
			Resolution:        cpl.Resolution,
			PictureType:       cpl.PictureType,
			HasAtmos:          cpl.HasAtmos,
			IsSmpte:           cpl.IsSmpte,
			IsEncrypted:       cpl.IsEncrypted,
		}

		for _, asset := range cpl.Assets {
			cplDTO.Assets = append(cplDTO.Assets, dtos.Asset{
				ID:                 asset.ID.String(),
				Name:               asset.Name,
				Type:               asset.Type,
				SizeInBytes:        int(asset.SizeInBytes),
				VerificationStatus: asset.VerificationStatus,
				ChannelCount:       cpl.Audio.ChannelCount,
				SampleRate:         cpl.Audio.SampleRate,
				SampleSize:         cpl.Audio.SampleSize,
				Language:           cpl.Audio.Language,
			})

		}

		for _, kdm := range kdms {
			if kdm.CPLID == cpl.ID {
				cplDTO.KDMS = append(cplDTO.KDMS, dtos.KDM{
					ID:             kdm.ID.String(),
					Name:           kdm.Name,
					NotValidAfter:  kdm.NotValidAfter,
					NotValidBefore: kdm.NotValidBefore,
				})
			}
		}

		mediaInfo.CPLs = append(mediaInfo.CPLs, cplDTO)
	}

	return mediaInfo
}
