package jsstring

import "regexp"

func Regexp_replace(data string, pattern string, replace string) string {
	re := regexp.MustCompile(pattern)
	return re.ReplaceAllString(data, replace)
}

func Regexp_match(data string, pattern string) bool {
	re := regexp.MustCompile(pattern)
	return re.MatchString(data)
}

func Regexp_match_all(data string, pattern string) []string {
	re := regexp.MustCompile(pattern)
	return re.FindAllString(data, -1)
}
