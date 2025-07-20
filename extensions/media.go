package extensions

// 媒体文件扩展名
var MediaExtensions = []string{
	".mp4",
	".mkv",
	".avi",
	".mov",
	".wmv",
	".flv",
	".webm",
	".m4v",
	".mpg",
	".mpeg",
	".3gp",
	".asf",
	".rm",
	".rmvb",
	".vob",
	".ts",
	".strm",
}

// MediaExtensionsExtended 包含所有媒体扩展名以及额外的 .strm 扩展名
var MediaExtensionsExtended = append(MediaExtensions, ".strm")
