package dtos

import "time"

type MediaInfo struct {
	LastUpdatedTime time.Time `json:"lastUpdatedTime"`
	CPLs            []CPL     `json:"cpls"`
}

type CPL struct {
	ID                string  `json:"id"`
	Name              string  `json:"name"`
	ContentKind       string  `json:"contentKind"`
	SizeInBytes       int     `json:"sizeInBytes"`
	DurationInSeconds float64 `json:"durationInSeconds"`
	KeepUntil         string  `json:"keepUntil"`
	Resolution        string  `json:"resolution"`
	PictureType       string  `json:"pictureType"`
	HasAtmos          bool    `json:"hasAtmos"`
	IsSmpte           bool    `json:"isSmpte"`
	IsEncrypted       bool    `json:"isEncrypted"`
	Assets            []Asset `json:"assets"`
	KDMS              []KDM   `json:"kdms"`
}

type KDM struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	NotValidAfter  string `json:"notValidAfter"`
	NotValidBefore string `json:"notValidBefore"`
}

type Asset struct {
	ID                 string `json:"id"`
	Name               string `json:"name"`
	Type               string `json:"type"`
	SizeInBytes        int    `json:"sizeInBytes"`
	VerificationStatus string `json:"verificationStatus"`
	ChannelCount       int    `json:"channelCount"`
	SampleRate         int    `json:"sampleRate"`
	SampleSize         int    `json:"sampleSize"`
	Language           string `json:"language"`
}
