package scrape_controller

type InfoType uint8

const (
	InfoTypeMovie InfoType = iota
	InfoTypeTV
	InfoTypeTVSeason
	InfoTypeTVEpisode
)

type InfoData interface {
	XML() ([]byte, error) // 返回 XML 格式的数据
}
