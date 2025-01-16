package dtos

import "media-agent/enums"

type Facility struct {
	ID             string
	PlayBackServer []PlayBackServer
}

type PlayBackServer struct {
	ID           string
	AuditoriumID string
	Type         enums.ServerType
	IsActive     bool
	IP           string
	ThumbPrint   string
}
