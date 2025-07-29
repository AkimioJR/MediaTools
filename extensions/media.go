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

// 字幕文件扩展名
var SubtitleExtensions = []string{
	".srt", ".ass", ".ssa", ".sup",
}

// 音轨文件扩展名
var AudioTrackExtensions = []string{
	".mka",
}

// 音频文件扩展名
var AudioExtensions = []string{
	".aac", ".ac3", ".amr", ".caf", ".cda", ".dsf",
	".dff", ".kar", ".m4a", ".mp1", ".mp2", ".mp3",
	".mid", ".mod", ".mka", ".mpc", ".nsf", ".ogg",
	".pcm", ".rmi", ".s3m", ".snd", ".spx", ".tak",
	".tta", ".vqf", ".wav", ".wma",
	".aifc", ".aiff", ".alac", ".adif", ".adts",
	".flac", ".midi", ".opus", ".sfalc",
}
