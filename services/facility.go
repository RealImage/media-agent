package services

import (
	"encoding/json"
	"fmt"
	"media-agent/config"
	"media-agent/dtos"
	"media-agent/enums"
	"media-agent/logger"
	"net/http"
	"net/url"
	"strings"
	"time"

	"go.uber.org/zap"
)

const PlaybackServerCode = "PLY"

type ScreenConfig struct {
	DeviceId   string
	DeviceIP   string
	ServerType string
}

type Facility struct {
	Auditoriums []Auditorium `json:"auditoriums"`
}
type Auditorium struct {
	Id      string   `json:"uuid"`
	Devices []Device `json:"devices"`
}

type Device struct {
	Id             string `json:"id"`
	Code           string `json:"code"`
	Type           string `json:"type"`
	IsActive       bool   `json:"is_active"`
	ThumbPrint     string `json:"thumbprint"`
	DeviceIP       string `json:"device_ip"`
	SofwareVersion string `json:"software_version"`
	Manufacturer   string `json:"manufacturer"`
}

type FacilityService struct {
	cfg *config.EnvConfig
}

func NewFacilityService(cfg *config.EnvConfig) *FacilityService {
	return &FacilityService{
		cfg: cfg,
	}
}

func (fs *FacilityService) GetFaciltiyById(facilityId string) (*dtos.Facility, error) {
	log := logger.GetLogger()

	faciity, err := fs.getFaciltiyById(facilityId)
	if err != nil {
		return nil, err
	}

	facility := dtos.Facility{
		ID:             facilityId,
		PlayBackServer: []dtos.PlayBackServer{},
	}

	for _, auditorium := range faciity.Auditoriums {
		for _, device := range auditorium.Devices {
			playBackServer, err := deviceToPlayBackServer(device, auditorium.Id)
			if err != nil {
				log.Error("failed to convert device to playback server", zap.Error(err))
				continue
			}
			facility.PlayBackServer = append(facility.PlayBackServer, *playBackServer)

		}
	}

	return &facility, nil
}

func deviceToPlayBackServer(device Device, auditoriumId string) (*dtos.PlayBackServer, error) {
	if !strings.EqualFold(device.Code, PlaybackServerCode) {
		return nil, fmt.Errorf("device is not a playback server, id: %s", device.Id)
	}

	serverType, err := deviceToServerType(device)
	if err != nil {
		return nil, fmt.Errorf("failed to get server type, id: %s, error: %s", device.Id, err.Error())
	}

	if device.DeviceIP == "" {
		return nil, fmt.Errorf("device IP is missing, id: %s", device.Id)
	}

	return &dtos.PlayBackServer{
		ID:           device.Id,
		AuditoriumID: auditoriumId,
		Type:         serverType,
		IsActive:     device.IsActive,
		IP:           device.DeviceIP,
		ThumbPrint:   device.ThumbPrint,
	}, nil

}

func deviceToServerType(device Device) (enums.ServerType, error) {
	if strings.EqualFold(device.Manufacturer, enums.DeviceManufacturerDolby.String()) {
		return enums.ServerTypeDolby, nil
	}

	if strings.EqualFold(device.Manufacturer, enums.DeviceManufacturerQube.String()) {
		if device.SofwareVersion == "" {
			return "", fmt.Errorf("software version is missing, id: %s", device.Id)
		}

		switch device.SofwareVersion[0] {
		case '4':
			return enums.ServerTypeQubeXP4, nil
		case '3', '2':
			return enums.ServerTypeQubeLegacyXP, nil
		default:
			return "", fmt.Errorf("unknown software version, id: %s", device.Id)

		}
	}

	return "", fmt.Errorf("unknown manufacturer, id: %s", device.Id)
}

func (fs *FacilityService) getFaciltiyById(facilityId string) (*Facility, error) {
	url, err := url.Parse(fs.cfg.QWCURL)
	if err != nil {
		return nil, nil
	}

	url = url.JoinPath("/facilities/", facilityId)

	finalURL := url.String() + "?token=" + fs.cfg.QWCTOKEN

	req, err := http.NewRequest("GET", finalURL, nil)
	if err != nil {
		return nil, err
	}

	client := http.Client{
		Timeout: 5 * time.Minute,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get screen config: %s", resp.Status)
	}

	var facility Facility

	if err := json.NewDecoder(resp.Body).Decode(&facility); err != nil {
		return nil, err
	}

	return &facility, nil
}
