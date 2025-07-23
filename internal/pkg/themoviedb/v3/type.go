package themoviedb

type Keyword struct {
	ID   int    `json:"id"`   // 关键词ID
	Name string `json:"name"` // 关键词名称
}

type Vote struct {
	VoteAverage float64 `json:"vote_average"` // 平均评分
	VoteCount   int     `json:"vote_count"`   // 评分数量
}
type Genre Keyword // 类型

type Company struct {
	Keyword
	LogoPath      string `json:"logo_path"`      // 公司Logo路径
	OriginCountry string `json:"origin_country"` // 原始国家
}

type Country struct {
	ISO  string `json:"iso_3166_1"` // ISO 3166-1 国家代码
	Name string `json:"name"`       // 国家名称
}

type Language struct {
	EnglishName string `json:"english_name"` // 英文名称
	ISO         string `json:"iso_639_1"`    // ISO 639-1 语言代码
	Name        string `json:"name"`         // 语言名称
}

type Collection struct {
	Keyword
	PosterPath   string `json:"poster_path"`   // 海报图片路径
	BackDropPath string `json:"backdrop_path"` // 背景图片路径
}

type Photo struct {
	Vote
	AspectRatio float64 `json:"aspect_ratio"` // 纵横比
	Height      int     `json:"height"`       // 高度
	ISO         string  `json:"iso_639_1"`    // ISO 639-1 语言代码
	FilePath    string  `json:"file_path"`    // 文件路径
	Width       int     `json:"width"`        // 宽度
}

type Image struct {
	BackDrops []Photo `json:"backdrops"` // 背景图片列表
	ID        int     `json:"id"`        // ID
	Logos     []Photo `json:"logos"`     // Logo图片列表
	Posters   []Photo `json:"posters"`   // 海报图片列表
}
type Network Company // 播出平台

// 创作者
type Creator struct {
	Keyword
	CreditID    string `json:"credit_id"`    // 唯一标识 ID
	Gender      uint   `json:"gender"`       // 性别 0-> 未知 1-> 女性 2-> 男性
	ProfilePath string `json:"profile_path"` // 个人资料图片路径
}

type TVEpisode struct {
	Keyword
	Vote
	Overview       string `json:"overview"`        // 概述
	AirDate        string `json:"air_date"`        // 首播日期
	EpisodeNumber  int    `json:"episode_number"`  // 集数
	ProductionCode string `json:"production_code"` // 制作代码
	Runtime        int    `json:"runtime"`         // 时长
	SeasonNumber   int    `json:"season_number"`   // 季数
	ShowID         int    `json:"show_id"`         // 电视剧ID
	StillPath      string `json:"still_path"`      // 静态图片路径
}

type TVSeason struct {
	Keyword
	AirDate      string  `json:"air_date"`      // 首播日期
	EpisodeCount int     `json:"episode_count"` // 集数
	Overview     string  `json:"overview"`      // 概述
	PosterPath   string  `json:"poster_path"`   // 海报图片路径
	SeasonNumber int     `json:"season_number"` // 季数
	VoteAverage  float64 `json:"vote_average"`  // 平均评分
}

type Crew struct {
	Creator
	Department         string  `json:"department"`           // 部门
	Job                string  `json:"job"`                  // 职位
	Adult              bool    `json:"adult"`                // 是否成人内容
	KnownForDepartment string  `json:"known_for_department"` // 知名领域
	OriginalName       string  `json:"original_name"`        // 原始姓名
	Popularity         float64 `json:"popularity"`           // 人气
}

type GuestStar struct {
	Creator
	Charactor          string  `json:"character"`            // 角色名称
	Order              int     `json:"order"`                // 出场顺序
	Adult              bool    `json:"adult"`                // 是否成人内容
	KnownForDepartment string  `json:"known_for_department"` // 知名领域
	OriginalName       string  `json:"original_name"`        // 原始姓名
	Popularity         float64 `json:"popularity"`           // 人气
}

type TVImage struct {
	Vote
	AspectRatio float64 `json:"aspect_ratio"` // 图片宽高比
	Height      int     `json:"height"`       // 图片高度
	ISO         string  `json:"iso_639_1"`    // ISO 639-1 语言代码
	FilePath    string  `json:"file_path"`    // 图片文件路径
	Width       int     `json:"width"`        // 图片宽度
}

// 标题/别名
type Title struct {
	Title string `json:"title"`      // 标题
	Type  string `json:"type"`       // 类型
	ISO   string `json:"iso_3166_1"` // ISO iso_3166_1 国家代码
}

type Translation struct {
	ISO3166 string `json:"iso_3166_1"`   // ISO 3166-1 国家代码
	ISO639  string `json:"iso_639_1"`    // ISO 639-1 语言代码
	Name    string `json:"name"`         // 名称
	English string `json:"english_name"` // 英文名称
	Data    struct {
		Name     string `json:"name"`     // 名称
		Overview string `json:"overview"` // 概述
		Homepage string `json:"homepage"` // 主页
		Tagline  string `json:"tagline"`  // 标语
	} `json:"data"` // 翻译数据
}
