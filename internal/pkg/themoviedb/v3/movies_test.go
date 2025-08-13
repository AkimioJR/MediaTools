package themoviedb_test

import (
	"MediaTools/internal/pkg/themoviedb/v3"
	"testing"

	"github.com/stretchr/testify/require"
)

const MovieID = 874745 // 致深愛妳的那個我
var client *themoviedb.Client

func init() {
	var err error
	client, err = themoviedb.NewClient("db55323b8d3e4154498498a75642b381")
	if err != nil {
		panic(err)
	}
}

func TestGetMovieDetail(t *testing.T) {
	movieDetail, err := client.GetMovieDetail(MovieID, nil)
	require.NoError(t, err)
	require.NotNil(t, movieDetail)
	require.Equal(t, MovieID, movieDetail.ID)
}

func TestGetMovieAlternativeTitle(t *testing.T) {
	titles, err := client.GetMovieAlternativeTitle(MovieID, nil)
	require.NoError(t, err)
	require.NotEmpty(t, titles.Titles)
}

func TestGetMovieTranslation(t *testing.T) {
	translation, err := client.GetMovieTranslation(MovieID)
	require.NoError(t, err)
	require.NotEmpty(t, translation.Translations)
}

func TestGetMovieCredit(t *testing.T) {
	credit, err := client.GetMovieCredit(MovieID, nil)
	require.NoError(t, err)
	require.NotNil(t, credit)
	require.NotEmpty(t, credit.Cast)
	require.NotEmpty(t, credit.Crew)
}

func TestGetMovieImage(t *testing.T) {
	_, err := client.GetMovieImage(MovieID, nil, nil)
	require.NoError(t, err)
}

func TestGetMovieKeyword(t *testing.T) {
	keywords, err := client.GetMovieKeyword(MovieID)
	require.NoError(t, err)
	require.NotNil(t, keywords)
	require.NotEmpty(t, keywords.Keywords)
}
