package scrape_controller

import "encoding/xml"

// 唯一标识（多类型ID，如IMDb、TMDB、TVDB等）
type UniqueID struct {
	Type  string `xml:"type,attr"` // 标识类型（tmdb/imdb/tvdb等）
	Value string `xml:",chardata"` // 标识值
}

// 演员信息
type Actor struct {
	Name string `xml:"name"` // 演员姓名
	Role string `xml:"role"` // 饰演角色
	Type string `xml:"type"` // 演员类型（Actor/GuestStar等）

	Thumb   string `xml:"thumb,omitempty"`   // 演员缩略图URL（可选）
	Profile string `xml:"profile,omitempty"` // 演员个人简介（可选）

	TMDBID    string `xml:"tmdbid,omitempty"`     // TMDB演员ID（可选）
	TVDBID    string `xml:"tvdbid,omitempty"`     // TVDB演员ID（可选）
	IMDbID    string `xml:"imdbid,omitempty"`     // IMDb演员ID（可选）
	BangumiID string `xml:"bangumiid,omitempty"`  // 番组计划ID（可选，仅部分文档有）
	DoubanID  string `xml:"doubanidid,omitempty"` // 豆瓣ID（可选，仅部分文档有）
}

// 评分信息
type Rating struct {
	Default bool    `xml:"default,attr"` // 是否默认评分
	Max     int     `xml:"max,attr"`     // 满分分值（通常为10）
	Name    string  `xml:"name,attr"`    // 评分来源（如themoviedb）
	Value   float64 `xml:"value"`        // 评分值
	Votes   int     `xml:"votes"`        // 投票数
}

// 缩略图/海报
type Thumb struct {
	Aspect string `xml:"aspect,attr"` // 图片类型（poster/banner/landscape等）
	URL    string `xml:",chardata"`   // 图片URL
}

// 制作公司/发行方
type Studio struct {
	Name string `xml:",chardata"` // 公司名称
}

// 创作者（导演/编剧等）
type Creator struct {
	TMDBID  string `xml:"tmdbid,attr,omitempty"` // TMDBID（可选）
	IMDbID  string `xml:"imdbid,attr,omitempty"` // IMDbID（可选）
	Name    string `xml:",chardata"`             // 姓名
	Profile string `xml:"profile,omitempty"`     // 个人简介（可选）
}

// 电视剧元数据
type TVSeriesMetaData struct {
	XMLName       xml.Name   `xml:"tvshow"`
	Title         string     `xml:"title"`                  // 剧集标题
	OriginalTitle string     `xml:"originaltitle"`          // 原始标题
	Plot          string     `xml:"plot"`                   // 剧集总剧情
	Outline       string     `xml:"outline"`                // 剧集总概要
	Season        int        `xml:"season"`                 // 季数（默认-1）
	Episode       int        `xml:"episode"`                // 集数（默认-1）
	LockData      bool       `xml:"lockdata,omitempty"`     // 数据锁定状态
	DateAdded     string     `xml:"dateadded,omitempty"`    // 添加时间
	Actors        []Actor    `xml:"actor,omitempty"`        // 主要演员
	Trailer       string     `xml:"trailer,omitempty"`      // 预告片URL
	Rating        float64    `xml:"rating,omitempty"`       // 剧集总评分
	Year          int        `xml:"year,omitempty"`         // 发行年份
	SortTitle     string     `xml:"sorttitle,omitempty"`    // 排序标题
	MPAA          string     `xml:"mpaa,omitempty"`         // 分级（如TV-14）
	IMDbID        string     `xml:"imdb_id,omitempty"`      // IMDb标识
	TMDBID        int        `xml:"tmdbid,omitempty"`       // TMDB标识
	Premiered     string     `xml:"premiered,omitempty"`    // 首播日期
	ReleaseDate   string     `xml:"releasedate,omitempty"`  // 发布日期
	Runtime       int        `xml:"runtime,omitempty"`      // 单集时长（分钟）
	Countries     []string   `xml:"country,omitempty"`      // 制作国家
	Genres        []string   `xml:"genre,omitempty"`        // 类型（动画/喜剧等）
	Studios       []Studio   `xml:"studio,omitempty"`       // 制作公司
	Tags          []string   `xml:"tag,omitempty"`          // 标签（如anime/romance）
	UniqueIDs     []UniqueID `xml:"uniqueid,omitempty"`     // 多类型唯一标识
	TVDBID        string     `xml:"tvdbid,omitempty"`       // TVDB标识
	WikidataID    string     `xml:"wikidataid,omitempty"`   // 维基数据ID
	Eidrid        string     `xml:"eidrid,omitempty"`       // EIDR标识
	EpisodeGuide  string     `xml:"episodeguide,omitempty"` // 剧集指南JSON
	ID            string     `xml:"id,omitempty"`           // 内部ID
	DisplayOrder  string     `xml:"displayorder,omitempty"` // 排序方式（aired）
	Status        string     `xml:"status,omitempty"`       // 状态（如Continuing）
	ShowTitle     string     `xml:"showtitle,omitempty"`    // 剧集标题
	Top250        int        `xml:"top250,omitempty"`       // 排行榜名次（0为未上榜）
	Ratings       []Rating   `xml:"ratings>rating"`         // 评分详情
	UserRating    int        `xml:"userrating,omitempty"`   // 用户评分
	Thumbs        []Thumb    `xml:"thumb,omitempty"`        // 各类缩略图
	Fanart        []struct {
		Thumb string `xml:"thumb,omitempty"`
	} `xml:"fanart"` // 粉丝创作图
	Certification string `xml:"certification,omitempty"` // 分级认证
	Watched       bool   `xml:"watched,omitempty"`       // 是否观看
	PlayCount     int    `xml:"playcount,omitempty"`     // 播放次数
	UserNote      string `xml:"user_note,omitempty"`     // 用户备注
}

func (t *TVSeriesMetaData) XML() ([]byte, error) {
	return xml.MarshalIndent(t, "", "  ")
}

// 电视剧季信息
type TVSeasonMetaData struct {
	XMLName      xml.Name   `xml:"season"`
	Plot         string     `xml:"plot,omitempty"`         // 剧情描述
	Outline      string     `xml:"outline,omitempty"`      // 剧情概要
	LockData     bool       `xml:"lockdata,omitempty"`     // 数据锁定状态
	DateAdded    string     `xml:"dateadded,omitempty"`    // 添加时间
	Title        string     `xml:"title,omitempty"`        // 季标题
	Year         int        `xml:"year,omitempty"`         // 发行年份
	SortTitle    string     `xml:"sorttitle,omitempty"`    // 排序标题
	IMDbID       string     `xml:"imdbid,omitempty"`       // IMDb标识
	TVDBID       string     `xml:"tvdbid,omitempty"`       // TVDB标识
	TMDBID       int        `xml:"tmdbid,omitempty"`       // TMDB标识
	Premiered    string     `xml:"premiered,omitempty"`    // 首映日期
	ReleaseDate  string     `xml:"releasedate,omitempty"`  // 发布日期
	UniqueIDs    []UniqueID `xml:"uniqueid,omitempty"`     // 多类型唯一标识
	SeasonNumber int        `xml:"seasonnumber,omitempty"` // 季数
	ShowTitle    string     `xml:"showtitle,omitempty"`    // 所属剧集标题
	Thumbs       []Thumb    `xml:"thumb,omitempty"`        // 缩略图
	UserNote     string     `xml:"user_note,omitempty"`    // 用户备注
}

func (t *TVSeasonMetaData) XML() ([]byte, error) {
	return xml.MarshalIndent(t, "", "  ")
}

// 电视剧集信息
type TVEpisodeMetaData struct {
	XMLName          xml.Name   `xml:"episodedetails,omitempty"`
	Plot             string     `xml:"plot,omitempty"`              // 单集剧情
	Outline          string     `xml:"outline,omitempty"`           // 单集概要（可选）
	LockData         bool       `xml:"lockdata,omitempty"`          // 数据锁定状态
	DateAdded        string     `xml:"dateadded,omitempty"`         // 添加时间
	Title            string     `xml:"title,omitempty"`             // 单集标题
	OriginalTitle    string     `xml:"originaltitle,omitempty"`     // 原始标题（日语等）
	Actors           []Actor    `xml:"actor,omitempty"`             // 演员列表
	Directors        []Creator  `xml:"director,omitempty"`          // 导演
	Writers          []Creator  `xml:"writer,omitempty"`            // 编剧列表
	Credits          []Creator  `xml:"credits,omitempty"`           // 制作人员
	Rating           float64    `xml:"rating,omitempty"`            // 单集评分
	Year             int        `xml:"year,omitempty"`              // 发行年份
	SortTitle        string     `xml:"sorttitle"`                   // 排序标题
	IMDbID           string     `xml:"imdbid,omitempty"`            // IMDb标识
	TVDBID           string     `xml:"tvdbid,omitempty"`            // TVDB标识
	TMDBID           int        `xml:"tmdbid,omitempty"`            // TMDB标识
	Studios          []Studio   `xml:"studio,omitempty"`            // 制作公司
	UniqueIDs        []UniqueID `xml:"uniqueid,omitempty"`          // 多类型唯一标识
	Episode          int        `xml:"episode,omitempty"`           // 集数
	Season           int        `xml:"season,omitempty"`            // 所属季数
	Aired            string     `xml:"aired,omitempty"`             // 播出日期
	ShowTitle        string     `xml:"showtitle,omitempty"`         // 所属剧集标题
	Ratings          []Rating   `xml:"ratings>rating"`              // 评分详情
	UserRating       int        `xml:"userrating"`                  // 用户评分
	Thumbs           []Thumb    `xml:"thumb,omitempty"`             // 缩略图
	Watched          bool       `xml:"watched,omitempty"`           // 是否观看
	PlayCount        int        `xml:"playcount,omitempty"`         // 播放次数
	EpBookmark       string     `xml:"epbookmark,omitempty"`        // 播放书签
	Code             string     `xml:"code,omitempty"`              // 编码（可选）
	Source           string     `xml:"source,omitempty"`            // 来源（如UNKNOWN）
	Edition          string     `xml:"edition,omitempty"`           // 版本（如NONE）
	OriginalFilename string     `xml:"original_filename,omitempty"` // 原始文件名
	UserNote         string     `xml:"user_note,omitempty"`         // 用户备注
}

func (t *TVEpisodeMetaData) XML() ([]byte, error) {
	return xml.MarshalIndent(t, "", "  ")
}

// 电影元数据
type MovieMetaData struct {
	XMLName       xml.Name   `xml:"movie"`
	Title         string     `xml:"title"`                    // 电影标题
	OriginalTitle string     `xml:"originaltitle"`            // 原始标题
	Plot          string     `xml:"plot"`                     // 电影剧情
	Outline       string     `xml:"outline"`                  // 电影概要
	LockData      bool       `xml:"lockdata,omitempty"`       // 数据锁定状态
	DateAdded     string     `xml:"dateadded,omitempty"`      // 添加时间
	Actors        []Actor    `xml:"actor,omitempty"`          // 演员列表
	Directors     []Creator  `xml:"director,omitempty"`       // 导演
	Writers       []Creator  `xml:"writer,omitempty"`         // 编剧列表
	Credits       []Creator  `xml:"credits,omitempty"`        // 制作人员
	Trailer       string     `xml:"trailer,omitempty"`        // 预告片URL
	Rating        float64    `xml:"rating,omitempty"`         // 电影评分
	Year          int        `xml:"year,omitempty"`           // 上映年份
	SortTitle     string     `xml:"sorttitle,omitempty"`      // 排序标题
	MPAA          string     `xml:"mpaa,omitempty"`           // 分级（如US:R）
	IMDbID        string     `xml:"imdbid,omitempty"`         // IMDb标识
	TVDBID        string     `xml:"tvdbid,omitempty"`         // TVDB标识
	TMDBID        int        `xml:"tmdbid,omitempty"`         // TMDB标识
	Premiered     string     `xml:"premiered,omitempty"`      // 首映日期
	ReleaseDate   string     `xml:"releasedate,omitempty"`    // 上映日期
	Tagline       string     `xml:"tagline,omitempty"`        // 宣传语
	Countries     []string   `xml:"country,omitempty"`        // 制作国家
	Genres        []string   `xml:"genre,omitempty"`          // 类型（动作/科幻等）
	Studios       []Studio   `xml:"studio,omitempty"`         // 制作公司
	Tags          []string   `xml:"tag,omitempty"`            // 标签（如dystopia/revolution）
	UniqueIDs     []UniqueID `xml:"uniqueid,omitempty"`       // 多类型唯一标识
	ID            string     `xml:"id,omitempty"`             // 内部ID（IMDbID）
	Ratings       []Rating   `xml:"ratings>rating,omitempty"` // 评分详情
	UserRating    float64    `xml:"userrating,omitempty"`     // 用户评分
	Top250        int        `xml:"top250,omitempty"`         // 排行榜名次
	Thumb         Thumb      `xml:"thumb,omitempty"`          // 主海报
	Fanart        []struct {
		Thumb string `xml:"thumb"`
	} `xml:"fanart"` // 粉丝创作图
	Certification    string    `xml:"certification,omitempty"`     // 分级认证
	Status           string    `xml:"status,omitempty"`            // 状态（空为已上映）
	Code             string    `xml:"code,omitempty"`              // 编码
	Watched          bool      `xml:"watched,omitempty"`           // 是否观看
	PlayCount        int       `xml:"playcount,omitempty"`         // 播放次数
	Producers        []Creator `xml:"producer,omitempty"`          // 制片人
	Languages        []string  `xml:"languages,omitempty"`         // 语言
	Source           string    `xml:"source,omitempty"`            // 来源（如BLURAY）
	Edition          string    `xml:"edition,omitempty"`           // 版本
	OriginalFilename string    `xml:"original_filename,omitempty"` // 原始文件名
	UserNote         string    `xml:"user_note,omitempty"`         // 用户备注
}

func (t *MovieMetaData) XML() ([]byte, error) {
	return xml.MarshalIndent(t, "", "  ")
}

var _ InfoData = (*TVSeriesMetaData)(nil)
var _ InfoData = (*TVSeasonMetaData)(nil)
var _ InfoData = (*TVEpisodeMetaData)(nil)
