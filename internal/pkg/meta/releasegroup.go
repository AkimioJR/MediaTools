package meta

import (
	"regexp"
	"strings"
)

var releaseGroups = map[string][]string{
	"0ff":          {"FF(?:(?:A|WE)B|CD|E(?:DU|B)|TV)"},
	"1pt":          {},
	"52pt":         {},
	"audiences":    {"Audies", "AD(?:Audio|E(?:book|)|Music|Web)"},
	"azusa":        {},
	"beitai":       {"BeiTai"},
	"btschool":     {"Bts(?:CHOOL|HD|PAD|TV)", "Zone"},
	"carpt":        {"CarPT"},
	"chdbits":      {"CHD(?:Bits|PAD|(?:|HK)TV|WEB|)", "StBOX", "OneHD", "Lee", "xiaopie"},
	"discfan":      {},
	"dragonhd":     {},
	"eastgame":     {"(?:(?:iNT|(?:HALFC|Mini(?:S|H|FH)D))-|)TLF"},
	"filelist":     {},
	"gainbound":    {"(?:DG|GBWE)B"},
	"hares":        {"Hares(?:(?:M|T)V|Web|)"},
	"hd4fans":      {},
	"hdarea":       {"HDA(?:pad|rea|TV)", "EPiC"},
	"hdatmos":      {},
	"hdbd":         {},
	"hdchina":      {"HDC(?:hina|TV|)", "k9611", "tudou", "iHD"},
	"hddolby":      {"D(?:ream|BTV)", "(?:HD|QHstudI)o"},
	"hdfans":       {"beAst(?:TV|)"},
	"hdhome":       {"HDH(?:ome|Pad|TV|WEB|)"},
	"hdpt":         {"HDPT(?:Web|)"},
	"hdsky":        {"HDS(?:ky|TV|Pad|WEB|)", "AQLJ"},
	"hdtime":       {},
	"HDU":          {},
	"hdvideo":      {},
	"hdzone":       {"HDZ(?:one|)"},
	"hhanclub":     {"HHWEB"},
	"hitpt":        {},
	"htpt":         {"HTPT"},
	"iptorrents":   {},
	"joyhd":        {},
	"keepfrds":     {"FRDS", "Yumi", "cXcY"},
	"lemonhd":      {"L(?:eague(?:(?:C|H)D|(?:M|T)V|NF|WEB)|HD)", "i18n", "CiNT"},
	"mteam":        {"MTeam(?:TV|)", "MPAD"},
	"nanyangpt":    {},
	"nicept":       {},
	"oshen":        {},
	"ourbits":      {"Our(?:Bits|TV)", "FLTTH", "Ao", "PbK", "MGs", "iLove(?:HD|TV)"},
	"piggo":        {"PiGo(?:NF|(?:H|WE)B)"},
	"ptchina":      {},
	"pterclub":     {"PTer(?:DIY|Game|(?:M|T)V|WEB|)"},
	"pthome":       {"PTH(?:Audio|eBook|music|ome|tv|WEB|)"},
	"ptmsg":        {},
	"ptsbao":       {"PTsbao", "OPS", "F(?:Fans(?:AIeNcE|BD|D(?:VD|IY)|TV|WEB)|HDMv)", "SGXT"},
	"pttime":       {},
	"putao":        {"PuTao"},
	"soulvoice":    {},
	"springsunday": {"CMCT(?:V|)"},
	"sharkpt":      {"Shark(?:WEB|DIY|TV|MV|)"},
	"tccf":         {},
	"tjupt":        {"TJUPT"},
	"totheglory":   {"TTG", "WiKi", "NGB", "DoA", "(?:ARi|ExRE)N"},
	"U2":           {},
	"ultrahd":      {},
	"others":       {"B(?:MDru|eyondHD|TN)", "C(?:fandora|trlhd|MRG)", "DON", "EVO", "FLUX", "HONE(?:yG|)", "N(?:oGroup|T(?:b|G))", "PandaMoon", "SMURF", "T(?:EPES|aengoo|rollHD )"},
	"anime":        {"ANi", "HYSUB", "KTXP", "LoliHouse", "MCE", "Nekomoe kissaten", "SweetSub", "MingY", "(?:Lilith|NC)-Raws", "织梦字幕组", "枫叶字幕组", "猎户手抄部", "喵萌奶茶屋", "漫猫字幕社", "霜庭云花Sub", "北宇治字幕组", "氢气烤肉架", "云歌字幕组", "萌樱字幕组", "极影字幕社", "悠哈璃羽字幕社", "❀拨雪寻春❀", "沸羊羊(?:制作|字幕组)", "(?:桜|樱)都字幕组"},
	"forge":        {"FROG(?:E|Web|)"},
	"ubits":        {"UB(?:its|WEB|TV)"},
}
var releaseGroupsRe *regexp.Regexp = nil

func init() {
	var groups []string
	for _, values := range releaseGroups {
		groups = append(groups, values...)
	}
	// 使用捕获组而不是前后查找
	releaseGroupsRe = regexp.MustCompile(`[-@\[￡【&](` + strings.Join(groups, "|") + `)[@.\s\S\]\[】&]`)
}

func findReleaseGroups(title string) []string {
	title += " "
	matches := releaseGroupsRe.FindAllStringSubmatch(title, -1)

	var groups []string
	for _, match := range matches {
		if len(match) > 1 && match[1] != "" {
			groups = append(groups, match[1])
		}
	}

	return groups
}
