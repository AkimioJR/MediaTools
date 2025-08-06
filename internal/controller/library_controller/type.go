package library_controller

type Category struct {
	Name              string   `json:"name"`               // 分类名称
	GenreIDs          []int    `json:"genre_ids"`          // 分类对应的流派ID列表
	OriginalLanguages []string `json:"original_languages"` // 分类对应的原始语言列表
	OriginalCountries []string `json:"original_countries"` // 分类对应的原始国家列表 ISO 3166-1
}

type CategoryConfig struct {
	MovieCategories []Category `yaml:"movie_categories"` // 电影分类列表
	TVCategories    []Category `yaml:"tv_categories"`    // 电视剧分类列表
}

var categoryConfig CategoryConfig = CategoryConfig{
	MovieCategories: []Category{
		{
			Name:     "动画电影",
			GenreIDs: []int{16}, // 动画
		},
		{
			Name:              "华语电影",
			OriginalLanguages: []string{"zh", "cn", "bo", "za"}, // 包括中文、藏语、壮语等
		},
		{ // 未匹配以上条件时，分类为外语电影
			Name: "外语电影",
		},
	},
	TVCategories: []Category{
		{
			Name:              "国漫",
			GenreIDs:          []int{16},                  // 动画
			OriginalCountries: []string{"CN", "TW", "HK"}, // 包括中国大陆、台湾、香港
		},
		{
			Name:              "日漫",
			GenreIDs:          []int{16},      // 动画
			OriginalCountries: []string{"JP"}, // 日本
		},
		{
			Name:     "纪录片",
			GenreIDs: []int{99}, // 纪录片
		},
		{
			Name:     "儿童",
			GenreIDs: []int{10762}, // 儿童
		},
		{
			Name:     "综艺",
			GenreIDs: []int{10764, 10767}, // 综艺
		},
		{
			Name:              "国产剧",
			OriginalCountries: []string{"CN", "TW", "HK"}, // 包括中国大陆、台湾、香港
		},
		{
			Name: "欧美剧",
			OriginalCountries: []string{
				"US", // 美国
				"FR", // 法国
				"GB", // 英国
				"DE", // 德国
				"ES", // 西班牙
				"IT", // 意大利
				"NL", // 荷兰
				"PT", // 葡萄牙
				"RU", // 俄罗斯
				"UA", // 乌克兰
			},
		},
		{
			Name: "日韩剧",
			OriginalCountries: []string{
				"JP",       // 日本
				"KP", "KR", // 韩国
				// "TH", // 泰国
				// "IN", // 印度
				// "SG", // 新加坡
			},
		},
		{
			Name: "未分类",
		},
	},
}
