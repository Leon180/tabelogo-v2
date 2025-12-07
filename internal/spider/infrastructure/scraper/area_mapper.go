package scraper

import (
	"fmt"
	"strings"
)

// AreaMapper maps Google Maps addresses to Tabelog area codes
type AreaMapper struct {
	// Map of area names to Tabelog area codes
	areaCodeMap map[string]string
}

// NewAreaMapper creates a new area mapper
func NewAreaMapper() *AreaMapper {
	return &AreaMapper{
		areaCodeMap: buildAreaCodeMap(),
	}
}

// MapToTabelogArea converts a Google Maps address to Tabelog area code
// Example: "Meguro, Tokyo" -> "tokyo/A1316"
func (m *AreaMapper) MapToTabelogArea(address string) string {
	originalAddress := address

	// Normalize address
	address = strings.ToLower(strings.TrimSpace(address))

	// Try to find matching area
	for areaName, areaCode := range m.areaCodeMap {
		if strings.Contains(address, areaName) {
			// Log successful mapping
			fmt.Printf("🗺️  Area Mapper: '%s' -> '%s' (matched '%s')\n", originalAddress, areaCode, areaName)
			return areaCode
		}
	}

	// Default to Tokyo general search if no specific area found
	fmt.Printf("⚠️  Area Mapper: '%s' -> 'tokyo' (no match found, using default)\n", originalAddress)
	return "tokyo"
}

// buildAreaCodeMap builds the mapping of area names to Tabelog area codes
// Based on Tabelog's URL structure: https://tabelog.com/{area_code}/rstLst/
func buildAreaCodeMap() map[string]string {
	return map[string]string{
		// Tokyo 23 Wards (東京23区)
		"chiyoda":    "tokyo/A1301", // 千代田区
		"chuo":       "tokyo/A1302", // 中央区
		"shibuya":    "tokyo/A1303", // 渋谷区
		"shinjuku":   "tokyo/A1304", // 新宿区
		"minato":     "tokyo/A1307", // 港区
		"bunkyo":     "tokyo/A1310", // 文京区
		"taito":      "tokyo/A1311", // 台東区
		"sumida":     "tokyo/A1312", // 墨田区
		"koto":       "tokyo/A1313", // 江東区
		"shinagawa":  "tokyo/A1314", // 品川区
		"meguro":     "tokyo/A1316", // 目黒区
		"ota":        "tokyo/A1317", // 大田区
		"setagaya":   "tokyo/A1318", // 世田谷区
		"nakano":     "tokyo/A1319", // 中野区
		"suginami":   "tokyo/A1320", // 杉並区
		"toshima":    "tokyo/A1321", // 豊島区
		"kita":       "tokyo/A1322", // 北区
		"arakawa":    "tokyo/A1323", // 荒川区
		"itabashi":   "tokyo/A1324", // 板橋区
		"nerima":     "tokyo/A1325", // 練馬区
		"adachi":     "tokyo/A1326", // 足立区
		"katsushika": "tokyo/A1327", // 葛飾区
		"edogawa":    "tokyo/A1328", // 江戸川区

		// Tokyo Cities (東京市部)
		"hachioji":        "tokyo/A1329", // 八王子市
		"tachikawa":       "tokyo/A1330", // 立川市
		"musashino":       "tokyo/A1331", // 武蔵野市
		"mitaka":          "tokyo/A1332", // 三鷹市
		"fuchu":           "tokyo/A1333", // 府中市
		"chofu":           "tokyo/A1334", // 調布市
		"machida":         "tokyo/A1335", // 町田市
		"koganei":         "tokyo/A1336", // 小金井市
		"kodaira":         "tokyo/A1337", // 小平市
		"hino":            "tokyo/A1338", // 日野市
		"higashimurayama": "tokyo/A1339", // 東村山市
		"kokubunji":       "tokyo/A1340", // 国分寺市
		"kunitachi":       "tokyo/A1341", // 国立市

		// Common area names in English
		"roppongi":      "tokyo/A1307", // 六本木 (Minato)
		"ginza":         "tokyo/A1302", // 銀座 (Chuo)
		"asakusa":       "tokyo/A1311", // 浅草 (Taito)
		"ueno":          "tokyo/A1311", // 上野 (Taito)
		"ikebukuro":     "tokyo/A1321", // 池袋 (Toshima)
		"ebisu":         "tokyo/A1303", // 恵比寿 (Shibuya)
		"harajuku":      "tokyo/A1303", // 原宿 (Shibuya)
		"akihabara":     "tokyo/A1301", // 秋葉原 (Chiyoda)
		"nakameguro":    "tokyo/A1316", // 中目黒 (Meguro)
		"daikanyama":    "tokyo/A1303", // 代官山 (Shibuya)
		"jiyugaoka":     "tokyo/A1316", // 自由が丘 (Meguro)
		"shimokitazawa": "tokyo/A1318", // 下北沢 (Setagaya)
		"kichijoji":     "tokyo/A1331", // 吉祥寺 (Musashino)
		"koenji":        "tokyo/A1320", // 高円寺 (Suginami)
		"komaba":        "tokyo/A1316", // 駒場 (Meguro)
	}
}
