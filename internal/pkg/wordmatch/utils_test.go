package wordmatch_test

import (
	"MediaTools/internal/pkg/wordmatch"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseOffsetExpr(t *testing.T) {
	tests := []struct {
		expr     string
		episode  int
		expected int
	}{
		{"EP+1", 5, 6},
		{"EP-2", 10, 8},
		{"2*EP", 3, 6},
		{"2*EP-4", 5, 6},
		{"2*EP+3", 4, 11},
		{"2*EP", 5, 10},
		{"EP", 7, 7},
		{"-1*EP", 5, -5},
		{"2*EP+3", 5, 13},
		{"EP+6", 5, 11},
		{"2*EP-2", 5, 8},
	}

	for _, test := range tests {
		result, err := wordmatch.ParseOffsetExpr(test.expr, test.episode)
		require.NoError(t, err, fmt.Sprintf("Failed to parse expression: %s with episode: %d", test.expr, test.episode))
		require.Equal(t, test.expected, result, fmt.Sprintf("Expected %d but got %d for expression: %s with episode: %d", test.expected, result, test.expr, test.episode))
	}
}

func TestParseLine(t *testing.T) {
	tests := []struct {
		input       string
		expected    *wordmatch.CustomWord
		shouldError bool
		description string
	}{
		{
			input: "被驅逐出勇者隊伍的白魔導師，被.S.級冒險者撿到",
			expected: &wordmatch.CustomWord{
				ReplaceFrom: "被驅逐出勇者隊伍的白魔導師，被.S.級冒險者撿到",
				ReplaceTo:   "",
				PrefixWord:  "",
				SuffixWord:  "",
				OffsetExpr:  "",
			},
			shouldError: false,
			description: "屏蔽词测试",
		},
		{
			input: "被驅逐出勇者隊伍的白魔導師，被.S.級冒險者撿到 => {[tmdbid=284771;type=tv;s=1]}被驅逐出勇者隊伍的白魔導師，被S級冒險者撿到",
			expected: &wordmatch.CustomWord{
				ReplaceFrom: "被驅逐出勇者隊伍的白魔導師，被.S.級冒險者撿到",
				ReplaceTo:   "{[tmdbid=284771;type=tv;s=1]}被驅逐出勇者隊伍的白魔導師，被S級冒險者撿到",
				PrefixWord:  "",
				SuffixWord:  "",
				OffsetExpr:  "",
			},
			shouldError: false,
			description: "仅替换词测试",
		},
		{
			input: "前缀 <> 后缀 >> 2*EP",
			expected: &wordmatch.CustomWord{
				ReplaceFrom: "",
				ReplaceTo:   "",
				PrefixWord:  "前缀",
				SuffixWord:  "后缀",
				OffsetExpr:  "2*EP",
			},
			shouldError: false,
			description: "仅前后定位词和偏移量表达式测试",
		},
		{
			input: "我們不可能成為戀人！絕對不行。（※似乎可行？） => {[tmdbid=277513;type=tv;s=1]}我们不可能成为恋人！绝对不行。 (※似乎可行？) && 前缀 <> 后缀 >> 2*EP",
			expected: &wordmatch.CustomWord{
				ReplaceFrom: "我們不可能成為戀人！絕對不行。（※似乎可行？）",
				ReplaceTo:   "{[tmdbid=277513;type=tv;s=1]}我们不可能成为恋人！绝对不行。 (※似乎可行？)",
				PrefixWord:  "前缀",
				SuffixWord:  "后缀",
				OffsetExpr:  "2*EP",
			},
			shouldError: false,
			description: "乘数表达式测试",
		},
		{
			input: "地縛少年花子君.2 => {[tmdbid=95269;type=tv;s=2]}地缚少年花子君 S02 && 前缀 <> 后缀 >> EP-1",
			expected: &wordmatch.CustomWord{
				ReplaceFrom: "地縛少年花子君.2",
				ReplaceTo:   "{[tmdbid=95269;type=tv;s=2]}地缚少年花子君 S02",
				PrefixWord:  "前缀",
				SuffixWord:  "后缀",
				OffsetExpr:  "EP-1",
			},
			shouldError: false,
			description: "减法表达式测试",
		},
		{
			input: `Lycoris Recoil 莉可麗絲：友誼是時間的竊賊 - => {[tmdbid=154494;type=tv;s=0]}莉可丽丝 S00 E && E <> \[1080P\] >> EP+4`,
			expected: &wordmatch.CustomWord{
				ReplaceFrom: "Lycoris Recoil 莉可麗絲：友誼是時間的竊賊 -",
				ReplaceTo:   "{[tmdbid=154494;type=tv;s=0]}莉可丽丝 S00 E",
				PrefixWord:  "E",
				SuffixWord:  `\[1080P\]`,
				OffsetExpr:  "EP+4",
			},
			shouldError: false,
			description: "复杂前后缀测试",
		},
		{
			input:       "invalid format without &&",
			expected:    nil,
			shouldError: true,
			description: "无效格式 - 缺少 &&",
		},
		{
			input:       " => replacement && invalid episode format",
			expected:    nil,
			shouldError: true,
			description: "无效格式 - 集数表达式格式错误",
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			result, err := wordmatch.ParseLine(test.input)

			if test.shouldError {
				require.Error(t, err, "Expected error for input: %s", test.input)
				require.Nil(t, result, "Expected nil result for invalid input: %s", test.input)
			} else {
				require.NoError(t, err, "Failed to parse line: %s", test.input)
				require.NotNil(t, result, "Expected non-nil result for input: %s", test.input)

				require.Equal(t, test.expected.ReplaceFrom, result.ReplaceFrom, "ReplaceFrom mismatch")
				require.Equal(t, test.expected.ReplaceTo, result.ReplaceTo, "ReplaceTo mismatch")
				require.Equal(t, test.expected.PrefixWord, result.PrefixWord, "PrefixWord mismatch")
				require.Equal(t, test.expected.SuffixWord, result.SuffixWord, "SuffixWord mismatch")
				require.Equal(t, test.expected.OffsetExpr, result.OffsetExpr, "OffsetExpr mismatch")
			}
		})
	}
}
