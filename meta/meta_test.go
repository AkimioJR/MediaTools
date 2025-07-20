package meta_test

import (
	"MediaTools/meta"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseMetaVideo(t *testing.T) {
	type expectedMeta struct {
		mediaType      meta.MediaType
		cnName         string
		enName         string
		year           uint
		part           string
		season         string
		episode        string
		resourcePix    meta.ResourcePix
		resourceType   meta.ResourceType
		resourceEffect map[meta.ResourceEffect]struct{}
		videoEncode    meta.VideoEncode
		audioEncode    meta.AudioEncode
	}
	testCases := []struct {
		input    string
		expected expectedMeta
	}{
		{
			input: "The Long Season 2017 2160p WEB-DL H265 AAC-XXX",
			expected: expectedMeta{
				mediaType:      meta.MediaTypeUnknown,
				cnName:         "",
				enName:         "The Long Season",
				year:           2017,
				part:           "",
				season:         "",
				episode:        "",
				resourcePix:    meta.ResourcePix2160p,
				resourceType:   meta.ResourceTypeWebDL,
				resourceEffect: make(map[meta.ResourceEffect]struct{}),
				videoEncode:    meta.VideoEncodeH265,
				audioEncode:    meta.AudioEncodeAAC,
			},
		},
		{
			input: "Cherry Season S01 2014 2160p WEB-DL H265 AAC-XXX",
			expected: expectedMeta{
				mediaType:      meta.MediaTypeTV,
				cnName:         "",
				enName:         "Cherry Season",
				year:           2014,
				part:           "",
				season:         "S01",
				episode:        "",
				resourcePix:    meta.ResourcePix2160p,
				resourceType:   meta.ResourceTypeWebDL,
				resourceEffect: make(map[meta.ResourceEffect]struct{}),
				videoEncode:    meta.VideoEncodeH265,
				audioEncode:    meta.AudioEncodeAAC,
			},
		},
		{
			input: "【爪爪字幕组】★7月新番[欢迎来到实力至上主义的教室 第二季/Youkoso Jitsuryoku Shijou Shugi no Kyoushitsu e S2][11][1080p][HEVC][GB][MP4][招募翻译校对]",
			expected: expectedMeta{
				mediaType:      meta.MediaTypeTV,
				cnName:         "欢迎来到实力至上主义的教室",
				enName:         "Youkoso Jitsuryoku Shijou Shugi no Kyoushitsu e",
				year:           0,
				part:           "",
				season:         "S02",
				episode:        "E11",
				resourcePix:    meta.ResourcePix1080p,
				resourceType:   meta.ResourceTypeUnknown,
				resourceEffect: make(map[meta.ResourceEffect]struct{}),
				videoEncode:    meta.VideoEncodeH265,
				audioEncode:    meta.AudioEncodeUnknown,
			},
		},
		{
			input: "National.Parks.Adventure.AKA.America.Wild:.National.Parks.Adventure.3D.2016.1080p.Blu-ray.AVC.TrueHD.7.1",
			expected: expectedMeta{
				mediaType:    meta.MediaTypeUnknown,
				cnName:       "",
				enName:       "National Parks Adventure AKA America Wild: National Parks Adventure",
				year:         2016,
				part:         "",
				season:       "",
				episode:      "",
				resourcePix:  meta.ResourcePix1080p,
				resourceType: meta.ResourceTypeBluRay,
				resourceEffect: map[meta.ResourceEffect]struct{}{
					meta.ResourceEffect3D: {},
				},
				videoEncode: meta.VideoEncodeH264,
				audioEncode: meta.AudioEncodeTrueHD,
			},
		},
		{
			input: "哆啦A梦：大雄的宇宙小战争 2021 (2022) - 1080p.mp4",
			expected: expectedMeta{
				mediaType:      meta.MediaTypeUnknown,
				cnName:         "哆啦A梦：大雄的宇宙小战争 2021",
				enName:         "",
				year:           2022,
				part:           "",
				season:         "",
				episode:        "",
				resourcePix:    meta.ResourcePix1080p,
				resourceType:   meta.ResourceTypeUnknown,
				resourceEffect: make(map[meta.ResourceEffect]struct{}),
				videoEncode:    meta.VideoEncodeUnknown,
				audioEncode:    meta.AudioEncodeUnknown,
			},
		},
		{
			input: "新精武门1991 (1991).mkv",
			expected: expectedMeta{
				mediaType:      meta.MediaTypeUnknown,
				cnName:         "新精武门1991",
				enName:         "",
				year:           1991,
				part:           "",
				season:         "",
				episode:        "",
				resourcePix:    meta.ResourcePixUnknown,
				resourceType:   meta.ResourceTypeUnknown,
				resourceEffect: make(map[meta.ResourceEffect]struct{}),
				videoEncode:    meta.VideoEncodeUnknown,
				audioEncode:    meta.AudioEncodeUnknown,
			},
		},
		{
			input: "24 S01 1080p WEB-DL AAC2.0 H.264-BTN",
			expected: expectedMeta{
				mediaType:      meta.MediaTypeTV,
				cnName:         "",
				enName:         "24",
				year:           0,
				part:           "",
				season:         "S01",
				episode:        "",
				resourcePix:    meta.ResourcePix1080p,
				resourceType:   meta.ResourceTypeWebDL,
				resourceEffect: make(map[meta.ResourceEffect]struct{}),
				videoEncode:    meta.VideoEncodeUnknown,
				audioEncode:    meta.AudioEncodeUnknown,
			},
		},
		{
			input: "Qi Refining for 3000 Years S01E06 2022 1080p B-Blobal WEB-DL X264 AAC-AnimeS@AdWeb",
			expected: expectedMeta{
				mediaType:      meta.MediaTypeTV,
				cnName:         "",
				enName:         "Qi Refining for 3000 Years",
				year:           2022,
				part:           "",
				season:         "S01",
				episode:        "E06",
				resourcePix:    meta.ResourcePix1080p,
				resourceType:   meta.ResourceTypeWebDL,
				resourceEffect: make(map[meta.ResourceEffect]struct{}),
				videoEncode:    meta.VideoEncodeH264,
				audioEncode:    meta.AudioEncodeAAC,
			},
		},
		{
			input: "Noumin Kanren no Skill Bakka Agetetara Naze ka Tsuyoku Natta S01E02 2022 1080p B-Global WEB-DL X264 AAC-AnimeS@ADWeb[2022年10月新番]",
			expected: expectedMeta{
				mediaType:      meta.MediaTypeTV,
				cnName:         "",
				enName:         "Noumin Kanren no Skill Bakka Agetetara Naze ka Tsuyoku Natta",
				year:           2022,
				part:           "",
				season:         "S01",
				episode:        "E02",
				resourcePix:    meta.ResourcePix1080p,
				resourceType:   meta.ResourceTypeWebDL,
				resourceEffect: make(map[meta.ResourceEffect]struct{}),
				videoEncode:    meta.VideoEncodeH264,
				audioEncode:    meta.AudioEncodeAAC,
			},
		},
		{
			input: "Thor Love and Thunder (2022) [1080p] [WEBRip] [5.1]",
			expected: expectedMeta{
				mediaType:      meta.MediaTypeUnknown,
				cnName:         "",
				enName:         "Thor Love and Thunder",
				year:           2022,
				part:           "",
				season:         "",
				episode:        "",
				resourcePix:    meta.ResourcePix1080p,
				resourceType:   meta.ResourceTypeWebRip,
				resourceEffect: make(map[meta.ResourceEffect]struct{}),
				videoEncode:    meta.VideoEncodeUnknown,
				audioEncode:    meta.AudioEncodeUnknown,
			},
		},
		{
			input: "钢铁侠2 (2010) 1080p AC3.mp4",
			expected: expectedMeta{
				mediaType:      meta.MediaTypeUnknown,
				cnName:         "钢铁侠2",
				enName:         "",
				year:           2010,
				part:           "",
				season:         "",
				episode:        "",
				resourcePix:    meta.ResourcePix1080p,
				resourceType:   meta.ResourceTypeUnknown,
				resourceEffect: make(map[meta.ResourceEffect]struct{}),
				videoEncode:    meta.VideoEncodeUnknown,
				audioEncode:    meta.AudioEncodeAC3,
			},
		},
		{
			input: "Wonder Woman 1984 2020 BluRay 1080p Atmos TrueHD 7.1 X264-EPiC",
			expected: expectedMeta{
				mediaType:      meta.MediaTypeUnknown,
				cnName:         "",
				enName:         "Wonder Woman 1984",
				year:           2020,
				part:           "",
				season:         "",
				episode:        "",
				resourcePix:    meta.ResourcePix1080p,
				resourceType:   meta.ResourceTypeBluRay,
				resourceEffect: make(map[meta.ResourceEffect]struct{}),
				videoEncode:    meta.VideoEncodeH264,
				audioEncode:    meta.AudioEncodeAtmos,
			},
		},
		{
			input: "9-1-1 - S04E03 - Future Tense WEBDL-1080p.mp4",
			expected: expectedMeta{
				mediaType:      meta.MediaTypeTV,
				cnName:         "",
				enName:         "9 1 1",
				year:           0,
				part:           "",
				season:         "S04",
				episode:        "E03",
				resourcePix:    meta.ResourcePix1080p,
				resourceType:   meta.ResourceTypeWebDL,
				resourceEffect: make(map[meta.ResourceEffect]struct{}),
				videoEncode:    meta.VideoEncodeUnknown,
				audioEncode:    meta.AudioEncodeUnknown,
			},
		},
		{
			input: "【幻月字幕组】【22年日剧】【据幸存的六人所说】【04】【1080P】【中日双语】",
			expected: expectedMeta{
				mediaType:      meta.MediaTypeTV,
				cnName:         "据幸存的六人所说",
				enName:         "",
				year:           0,
				part:           "",
				season:         "S01",
				episode:        "E04",
				resourcePix:    meta.ResourcePix1080p,
				resourceType:   meta.ResourceTypeUnknown,
				resourceEffect: make(map[meta.ResourceEffect]struct{}),
				videoEncode:    meta.VideoEncodeUnknown,
				audioEncode:    meta.AudioEncodeUnknown,
			},
		},
		{
			input: "【爪爪字幕组】★7月新番[即使如此依旧步步进逼/Soredemo Ayumu wa Yosetekuru][09][1080p][HEVC][GB][MP4][招募翻译校对]",
			expected: expectedMeta{
				mediaType:      meta.MediaTypeTV,
				cnName:         "即使如此依旧步步进逼",
				enName:         "Soredemo Ayumu wa Yosetekuru",
				year:           0,
				part:           "",
				season:         "S01",
				episode:        "E09",
				resourcePix:    meta.ResourcePix1080p,
				resourceType:   meta.ResourceTypeUnknown,
				resourceEffect: make(map[meta.ResourceEffect]struct{}),
				videoEncode:    meta.VideoEncodeH265,
				audioEncode:    meta.AudioEncodeUnknown,
			},
		},
		{
			input: "[GM-Team][国漫][寻剑 第1季][Sword Quest Season 1][2022][02][AVC][GB][1080P]",
			expected: expectedMeta{
				mediaType:      meta.MediaTypeTV,
				cnName:         "寻剑",
				enName:         "Sword Quest Season 1",
				year:           2022,
				part:           "",
				season:         "S01",
				episode:        "E02",
				resourcePix:    meta.ResourcePix1080p,
				resourceType:   meta.ResourceTypeUnknown,
				resourceEffect: make(map[meta.ResourceEffect]struct{}),
				videoEncode:    meta.VideoEncodeH264,
				audioEncode:    meta.AudioEncodeUnknown,
			},
		},
		{
			input: "[猎户不鸽发布组] 组长女儿与照料专员 / 组长女儿与保姆 Kumichou Musume to Sewagakari [09] [1080p+] [简中内嵌] [2022年7月番]",
			expected: expectedMeta{
				mediaType:      meta.MediaTypeTV,
				cnName:         "组长女儿与照料专员 组长女儿与保姆",
				enName:         "Kumichou Musume to Sewagakari",
				year:           0,
				part:           "",
				season:         "S01",
				episode:        "E09",
				resourcePix:    meta.ResourcePix1080p,
				resourceType:   meta.ResourceTypeUnknown,
				resourceEffect: make(map[meta.ResourceEffect]struct{}),
				videoEncode:    meta.VideoEncodeUnknown,
				audioEncode:    meta.AudioEncodeUnknown,
			},
		},
		{
			input: "Nande Koko ni Sensei ga!? 2019 Blu-ray Remux 1080p AVC LPCM-7³ ACG",
			expected: expectedMeta{
				mediaType:      meta.MediaTypeUnknown,
				cnName:         "",
				enName:         "Nande Koko ni Sensei ga!?",
				year:           2019,
				part:           "",
				season:         "",
				episode:        "",
				resourcePix:    meta.ResourcePix1080p,
				resourceType:   meta.ResourceTypeBluRayRemux,
				resourceEffect: make(map[meta.ResourceEffect]struct{}),
				videoEncode:    meta.VideoEncodeH264,
				audioEncode:    meta.AudioEncodeLPCM,
			},
		},
		{
			input: "30.Rock.S02E01.1080p.UHD.BluRay.X264-BORDURE.mkv",
			expected: expectedMeta{
				mediaType:      meta.MediaTypeTV,
				cnName:         "",
				enName:         "30 Rock",
				year:           0,
				part:           "",
				season:         "S02",
				episode:        "E01",
				resourcePix:    meta.ResourcePix1080p,
				resourceType:   meta.ResourceTypeUHDBluRay,
				resourceEffect: make(map[meta.ResourceEffect]struct{}),
				videoEncode:    meta.VideoEncodeH264,
				audioEncode:    meta.AudioEncodeUnknown,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			meta := meta.ParseMetaVideo(tc.input, true)
			require.Equal(t, tc.expected.mediaType, meta.GetType(), "媒体类型不匹配")
			require.Equal(t, tc.expected.cnName, meta.GetCNName(), "中文名不匹配")
			require.Equal(t, tc.expected.enName, meta.GetENName(), "英文名不匹配")
			require.Equal(t, tc.expected.year, meta.GetYear(), "年份不匹配")
			require.Equal(t, tc.expected.part, meta.GetPart(), "分部信息不匹配")
			require.Equal(t, tc.expected.season, meta.GetSeasonStr(), "季信息不匹配")
			require.Equal(t, tc.expected.episode, meta.GetEpisodeStr(), "集信息不匹配")
			require.Equal(t, tc.expected.resourcePix, meta.GetResourcePix(), "资源像素不匹配")
			require.Equal(t, tc.expected.resourceType, meta.GetResourceType(), "资源类型不匹配")
			require.Equal(t, tc.expected.resourceEffect, meta.GetResourceEffect(), "资源效果不匹配")
			require.Equal(t, tc.expected.videoEncode, meta.GetVideoEncode(), "视频编码不匹配")
			require.Equal(t, tc.expected.audioEncode, meta.GetAudioEncode(), "音频编码不匹配")
		})

	}
}
