package utils_test

import (
	"MediaTools/utils"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

type TestParams struct {
	Query        string  `query:"query"`
	Page         *int    `query:"page"`
	IncludeAdult *bool   `query:"include_adult"`
	Language     *string `query:"language"`
	EmptyStr     string  `query:"empty_str"`
	NilPtr       *int    `query:"nil_ptr"`
}

func TestStructToQuery(t *testing.T) {
	page := 2
	includeAdult := true
	lang := "zh-CN"
	params := TestParams{
		Query:        "test",
		Page:         &page,
		IncludeAdult: &includeAdult,
		Language:     &lang,
		EmptyStr:     "",
		NilPtr:       nil,
	}

	got := utils.StructToQuery(params)
	want := url.Values{
		"query":         []string{"test"},
		"page":          []string{"2"},
		"include_adult": []string{"true"},
		"language":      []string{"zh-CN"},
	}

	require.Equal(t, want, got)
}
