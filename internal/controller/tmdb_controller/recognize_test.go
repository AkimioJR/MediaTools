package tmdb_controller_test

import (
	"MediaTools/internal/config"
	"MediaTools/internal/controller/tmdb_controller"
	"MediaTools/internal/pkg/meta"
	"testing"

	_ "MediaTools/internal/logging"

	"github.com/stretchr/testify/require"
)

const TMDBApiKey = "db55323b8d3e4154498498a75642b381"

func init() {
	config.TMDB.ApiKey = TMDBApiKey
	tmdb_controller.Init()
}

func TestGetInfo(t *testing.T) {
	m := meta.MediaTypeTV
	info, err := tmdb_controller.GetInfo(271607, &m)
	require.NoError(t, err)
	require.Equal(t, "薰香花朵凛然绽放", info.TMDBInfo.TVInfo.SerieInfo.Name)
	m = meta.MediaTypeMovie
	info, err = tmdb_controller.GetInfo(874745, &m)
	require.NoError(t, err)
	require.Equal(t, "致深爱你的那个我", info.TMDBInfo.MovieInfo.Title)
	require.Equal(t, "2022-10-07", info.TMDBInfo.MovieInfo.ReleaseDate)
	require.Equal(t, "人们会经常地变动于仅有些许差别的平行世界，这点在这个世界已被证明——  因为两亲离婚而和父亲住在一起的日高历，在父亲任职的虚质科学研究所里与名为佐籐栞相遇了。  互相抱有淡淡恋心的两人，在某天因为父母的再婚问题而大大改变。  认为已经无法和对方结婚的历和栞，决定跳跃到不会成为兄妹的世界……  她不在的世界毫无意义。", info.TMDBInfo.MovieInfo.Overview)
}

func TestRecognizeMedia(t *testing.T) {
	mv := meta.ParseVideoMeta("【爪爪字幕组】★7月新番[欢迎来到实力至上主义的教室 第二季/Youkoso Jitsuryoku Shijou Shugi no Kyoushitsu e S2][11][1080p][HEVC][GB][MP4][招募翻译校对]")
	info, err := tmdb_controller.RecognizeMedia(mv, nil, nil)
	require.NoError(t, err)
	require.Equal(t, meta.MediaTypeTV, info.MediaType)
	require.Equal(t, 72517, info.TMDBID)

	mv = meta.ParseVideoMeta("[黒ネズミたち] 帝乃三姊妹意外地容易相处。 / Mikadono Sanshimai wa Angai, Choroi. - 02 (ABEMA 1280x720 AVC AAC MP4) [367.7 MB]")
	info, err = tmdb_controller.RecognizeMedia(mv, nil, nil)
	require.NoError(t, err)
	require.Equal(t, meta.MediaTypeTV, info.MediaType)
	require.Equal(t, 272556, info.TMDBID)

	mv = meta.ParseVideoMeta("【推しの子】 致深爱你的那个我 / Kimi wo Aishita Hitori no Boku e - Movie (CR 1920x1080 AVC AAC MKV)")
	info, err = tmdb_controller.RecognizeMedia(mv, nil, nil)
	require.NoError(t, err)
	require.Equal(t, meta.MediaTypeMovie, info.MediaType)
	require.Equal(t, 874745, info.TMDBID)

	mv = meta.ParseVideoMeta("【喵萌奶茶屋】★剧场版★[青春猪头少年不会梦到娇怜外出妹/Seishun Buta Yarou wa Odekake Sister no Yume o Minai][BDRip][1080p][繁日双语][招募翻译时轴]")
	info, err = tmdb_controller.RecognizeMedia(mv, nil, nil)
	require.NoError(t, err)
	require.Equal(t, meta.MediaTypeMovie, info.MediaType)
	require.Equal(t, 1056803, info.TMDBID)

	mv = meta.ParseVideoMeta("Thor Love and Thunder (2022) [1080p] [WEBRip] [5.1]")
	m := meta.MediaTypeMovie
	info, err = tmdb_controller.RecognizeMedia(mv, &m, nil)
	require.NoError(t, err)
	require.Equal(t, meta.MediaTypeMovie, info.MediaType)
	require.Equal(t, 616037, info.TMDBID)
}
