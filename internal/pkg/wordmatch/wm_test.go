package wordmatch_test

import (
	"MediaTools/internal/pkg/wordmatch"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWordsMatcher_MatchAndProcess(t *testing.T) {
	tests := []struct {
		name     string
		words    []string
		input    string
		expected string
	}{
		{
			name:     "屏蔽词测试 - 替换为空",
			words:    []string{"被驅逐出勇者隊伍的白魔導師，被.S.級冒險者撿到"},
			input:    "被驅逐出勇者隊伍的白魔導師，被.S.級冒險者撿到 第01集",
			expected: " 第01集",
		},
		{
			name: "替换词测试",
			words: []string{
				"被驅逐出勇者隊伍的白魔導師，被.S.級冒險者撿到 => {[tmdbid=284771;type=tv;s=1]}被驅逐出勇者隊伍的白魔導師，被S級冒險者撿到",
			},
			input:    "被驅逐出勇者隊伍的白魔導師，被.S.級冒險者撿到 第01集",
			expected: "{[tmdbid=284771;type=tv;s=1]}被驅逐出勇者隊伍的白魔導師，被S級冒險者撿到 第01集",
		},
		{
			name: "前后定位词和偏移量表达式测试 - 数字",
			words: []string{
				"前缀 <> 后缀 >> 2*EP",
			},
			input:    "前缀5后缀",
			expected: "前缀10后缀",
		},
		{
			name: "前后定位词和偏移量表达式测试 - EP+1",
			words: []string{
				`E <> [1080P] >> EP+1`,
			},
			input:    "E05[1080P]",
			expected: "E6[1080P]",
		},
		{
			name: "前后定位词和偏移量表达式测试 - EP-1",
			words: []string{
				"前缀 <> 后缀 >> EP-1",
			},
			input:    "前缀5后缀",
			expected: "前缀4后缀",
		},
		{
			name: "复合规则测试1",
			words: []string{
				"我們不可能成為戀人！絕對不行。 => 我们不可能成为恋人！绝对不行。 (※似乎可行？) && 前缀 <> 后缀 >> EP+2",
			},
			input:    "我們不可能成為戀人！絕對不行。 前缀3后缀",
			expected: "我们不可能成为恋人！绝对不行。 (※似乎可行？) 前缀5后缀",
		},
		{
			name: "复合规则测试2",
			words: []string{
				"被驅逐出勇者隊伍的白魔導師，被.S.級冒險者撿到 => {[tmdbid=284771;type=tv;s=1]}被驅逐出勇者隊伍的白魔導師，被S級冒險者撿到 && 第 <> 集 >> 3*EP+2",
			},
			input:    "被驅逐出勇者隊伍的白魔導師，被.S.級冒險者撿到 第01集",
			expected: "{[tmdbid=284771;type=tv;s=1]}被驅逐出勇者隊伍的白魔導師，被S級冒險者撿到 第5集",
		},
		{
			name: "中文数字测试",
			words: []string{
				"第 <> 集 >> EP+1",
			},
			input:    "第三集",
			expected: "第4集",
		},
		{
			name: "罗马数字测试",
			words: []string{
				"Season <> Episode >> EP-1",
			},
			input:    "Season III Episode",
			expected: "Season 2 Episode",
		},
		{
			name: "跳过注释和空行",
			words: []string{
				"# 这是注释",
				"",
				"   ",
				"测试 => 替换",
			},
			input:    "测试内容",
			expected: "替换内容",
		},
		{
			name: "无匹配情况",
			words: []string{
				"不存在的词 => 替换",
			},
			input:    "原始内容",
			expected: "原始内容",
		},
		{
			name: "前后定位词不完整",
			words: []string{
				"前缀 <> 后缀 >> EP+1",
			},
			input:    "前缀5", // 缺少后缀
			expected: "前缀5",
		},
		{
			name: "非数字内容在前后定位词之间",
			words: []string{
				"E <> P >> EP+1",
			},
			input:    "E测试P", // 中间不是数字
			expected: "E测试P",
		},
		{
			name: "零集数处理",
			words: []string{
				"第 <> 集 >> EP+1",
			},
			input:    "第0集", // 集数为0应该被跳过
			expected: "第0集",
		},
		{
			name: "多个规则应用",
			words: []string{
				"原词1 => 新词1",
				"原词2 => 新词2",
				"E <> P >> EP+1",
			},
			input:    "原词1 原词2 E5P",
			expected: "新词1 新词2 E6P",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			matcher, err := wordmatch.NewWordsMatcher(test.words)
			require.NoError(t, err, "Failed to create WordsMatcher")

			result, _ := matcher.MatchAndProcess(test.input)
			require.Equal(t, test.expected, result,
				"Input: %s, Expected: %s, Got: %s", test.input, test.expected, result)
		})
	}
}
