package meta

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseVersion(t *testing.T) {
	tests := []struct {
		input    string
		expected uint8
	}{
		{"[SweetSub&LoliHouse] Made in Abyss S2 - 03v2 [WebRip 1080p HEVC-10bit AAC ASSx2].mkv", 2},
		{"[HKGS]昨日之歌[12A][1080P][WEB-DL].mp4", 1},
		{"[ANi] 聖劍傳說 Legend of Mana - The Teardrop Crystal - - 12B [1080P][Baha][WEB-DL][AAC AVC][CHT].mkv", 2},
		{"[Haruhana] Kaoru Hana wa Rin to Saku - 02v2 [WebRip][HEVC-10bit 1080p][CHI_JPN].mkv", 2},
		{"[SweetSub&VCB-Studio] Boogiepop Phantom [02v2][Ma10p_720p][x265_flac].mkv", 2},
		{"[UHA-WINGS][Super Cub][02v2][x264 1080p][sc_jp].mkv", 2},
		{"[Sakurato] 86 Eiti Shikkusu [02v2][HEVC-10bit 1080p AAC][CHS&CHT].mkv", 2},
		{"[ANi]異世界魔王與召喚少女的奴隸魔術 Ω[02v2][1080P][Baha][WEB-DL].mp4", 2},
		{"[ANi]總之就是非常可愛[12v2][1080P][WEB-DL-B].mp4", 2},
		{"[ANi] D4DJ All Mix - 02v2 [1080P][Baha][WEB-DL][AAC AVC][CHT].mkv", 2},
	}

	for _, test := range tests {
		require.Equal(t, test.expected, ParseVersion(test.input), "Input: %s", test.input)
	}
}
