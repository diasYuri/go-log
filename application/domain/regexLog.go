package domain

var patternMap = map[string]string{
	"PATTERN_01": `(?P<date>\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2},\d{0,3})[ ]{1,2}-[ ]{1,2}(?P<level>[A-Z]* )[ ]{1,2}(?P<header>\[.*?] )[ ]{0,1}-[ ]{0,2}(?P<message>.*)`,
	"PATTERN_02": `(?P<date>\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}),{0,1}\d{0,3}[ ]{1,2}-[ ]{1,2}(?P<level>[A-Z]* )[ ]{1,2}(?P<header>\[.*?] )[ ]{0,1}-[ ]{0,2}(?P<message>.*)`,
}

func GetPattern(key string) string {
	return patternMap[key]
}
