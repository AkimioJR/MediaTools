package meta

import "strings"

// Tokens 结构体用于处理token解析
type Tokens struct {
	tokens []string
	index  int
}

// NewTokens 创建新的Tokens实例
func NewTokens(text string) *Tokens {
	// 简单的分词逻辑，按空格、点、下划线等分割
	tokens := tokenSplitRe.Split(text, -1)

	// 过滤空token
	filteredTokens := make([]string, 0)
	for _, token := range tokens {
		if strings.TrimSpace(token) != "" {
			filteredTokens = append(filteredTokens, strings.TrimSpace(token))
		}
	}

	return &Tokens{
		tokens: filteredTokens,
		index:  -1,
	}
}

// 判断是否结束
func (t *Tokens) IsEnd() bool {
	return t.index >= len(t.tokens)
}

// GetNext 获取下一个token
func (t *Tokens) GetNext() string {
	t.index++
	if t.IsEnd() {
		return ""
	}
	return t.tokens[t.index]
}

// Peek 查看下一个token但不移动索引
func (t *Tokens) Peek() string {
	nextIndex := t.index + 1
	if nextIndex >= len(t.tokens) {
		return ""
	}
	return t.tokens[nextIndex]
}

// Current 获取当前token
func (t *Tokens) Current() string {
	if t.index < 0 || t.index >= len(t.tokens) {
		return ""
	}
	return t.tokens[t.index]
}

// GetPrevious 获取前一个token但不移动索引
func (t *Tokens) GetPrevious() string {
	prevIndex := t.index - 1
	if prevIndex < 0 {
		return ""
	}
	return t.tokens[prevIndex]
}

// GetByIndex 根据索引获取token
func (t *Tokens) GetByIndex(index int) string {
	if index < 0 || index >= len(t.tokens) {
		return ""
	}
	return t.tokens[index]
}

// GetCurrentIndex 获取当前索引
func (t *Tokens) GetCurrentIndex() int {
	return t.index
}

// GetTokensInRange 获取指定范围内的tokens
func (t *Tokens) GetTokensInRange(start, end int) []string {
	if start < 0 {
		start = 0
	}
	if end > len(t.tokens) {
		end = len(t.tokens)
	}
	if start >= end {
		return []string{}
	}
	return t.tokens[start:end]
}

// GetLength 获取tokens的总长度
func (t *Tokens) GetLength() int {
	return len(t.tokens)
}
