package library_controller_test

import (
	"MediaTools/internal/controller/library_controller"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseEpisodeFormat(t *testing.T) {
	tests := []struct {
		name        string
		filename    string
		format      string
		expected    int
		expectError bool
	}{
		{
			name:        "十二国記 第45話",
			filename:    "(BD)十二国記 第45話「東の海神 西の滄海 五章」(1440x1080 x264-10bpp flac).mkv",
			format:      "(BD)十二国記 第{ep}話{a}(1440x1080 x264-10bpp flac).mkv",
			expected:    45,
			expectError: false,
		},
		{
			name:        "十二国記 第32話",
			filename:    "(BD)十二国記 第32話「風の万里 黎明の空　九章」(1440x1080 x264-10bpp flac).mkv",
			format:      "(BD)十二国記 第{ep}話{a}(1440x1080 x264-10bpp flac).mkv",
			expected:    32,
			expectError: false,
		},
		{
			name:        "简单集数格式",
			filename:    "Episode 12.mp4",
			format:      "Episode {ep}.mp4",
			expected:    12,
			expectError: false,
		},
		{
			name:        "多个通配符",
			filename:    "Series S01E05 Title Here 1080p.mkv",
			format:      "Series {a}E{ep} {a}.mkv",
			expected:    5,
			expectError: false,
		},
		{
			name:        "集数在中间",
			filename:    "Anime_Episode_08_Final.mp4",
			format:      "Anime_Episode_{ep}_{a}.mp4",
			expected:    8,
			expectError: false,
		},
		{
			name:        "三位数集数",
			filename:    "Show EP123 HD.mkv",
			format:      "Show EP{ep} HD.mkv",
			expected:    123,
			expectError: false,
		},
		{
			name:        "集数带空格",
			filename:    "Test  15  End.mp4",
			format:      "Test  {ep}  End.mp4",
			expected:    15,
			expectError: false,
		},
		{
			name:        "不匹配的格式",
			filename:    "Different File Name.mp4",
			format:      "Expected Format {ep}.mp4",
			expected:    -1,
			expectError: true,
		},
		{
			name:        "无效的集数格式",
			filename:    "Episode ABC.mp4",
			format:      "Episode {ep}.mp4",
			expected:    -1,
			expectError: true,
		},
		{
			name:        "无效的正则表达式格式",
			filename:    "test.mp4",
			format:      "test[.mp4",
			expected:    -1,
			expectError: true,
		},
		{
			name:        "没有{ep}占位符",
			filename:    "test.mp4",
			format:      "test.mp4",
			expected:    -1,
			expectError: true,
		},
		{
			name:        "特殊字符转义",
			filename:    "Test.Episode.01.720p.mkv",
			format:      "Test.Episode.{ep}.720p.mkv",
			expected:    1,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := library_controller.ParseEpisodeFormat(tt.filename, tt.format)

			if tt.expectError {
				require.Error(t, err, "Expected error for test case: %s", tt.name)
				require.Equal(t, tt.expected, result, "Expected result to match even on error")
			} else {
				require.NoError(t, err, "Expected no error for test case: %s", tt.name)
				require.Equal(t, tt.expected, result, "Expected episode number to match for test case: %s", tt.name)
			}
		})
	}
}
