package themoviedb

type Vote struct {
	VoteAverage float64 `json:"vote_average"` // 平均评分
	VoteCount   uint64  `json:"vote_count"`   // 评分数量
}
type Genre struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

type Company struct {
	ID            uint64 `json:"id"`             // 公司ID
	LogoPath      string `json:"logo_path"`      // 公司Logo路径
	Name          string `json:"name"`           // 公司名称
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
	ID           uint64 `json:"id"`            // 系列ID
	Name         string `json:"name"`          // 系列名称
	PosterPath   string `json:"poster_path"`   // 海报图片路径
	BackDropPath string `json:"backdrop_path"` // 背景图片路径
}

type Photo struct {
	Vote
	AspectRatio float64 `json:"aspect_ratio"` // 纵横比
	Height      uint64  `json:"height"`       // 高度
	ISO         string  `json:"iso_639_1"`    // ISO 639-1 语言代码
	FilePath    string  `json:"file_path"`    // 文件路径
	Width       uint64  `json:"width"`        // 宽度
}

type Image struct {
	BackDrops []Photo `json:"backdrops"` // 背景图片列表
	ID        uint64  `json:"id"`        // ID
	Logos     []Photo `json:"logos"`     // Logo图片列表
	Posters   []Photo `json:"posters"`   // 海报图片列表
}
