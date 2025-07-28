package scrape_controller

type InfoType uint8

const (
	InfoTypeMovie InfoType = iota
	InfoTypeTV
	InfoTypeTVSeason
	InfoTypeTVEpisode
)
