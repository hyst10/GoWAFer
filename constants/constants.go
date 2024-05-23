package constants

import "regexp"

const (
	BlackIPKey   = "blackIPList"
	WhiteIPKey   = "whiteIPList"
	BlackPathKey = "blackPath"
	WhitePathKey = "whitePath"
	SqlInjectKey = "sqlInjectRules"
	XssDetectKey = "xssDetectRules"
)

var SqlInjectRules = map[*regexp.Regexp]bool{
	regexp.MustCompile(`(?i)\bunion\b.*\bselect\b`):     true,
	regexp.MustCompile(`(?i)\bselect\b.*\bfrom\b`):      true,
	regexp.MustCompile(`(?i)\binsert\s+into\b`):         true,
	regexp.MustCompile(`(?i)\bdelete\s+from\b`):         true,
	regexp.MustCompile(`(?i)\bdrop\s+table\b`):          true,
	regexp.MustCompile(`(?i)\bupdate\b.*\bset\b`):       true,
	regexp.MustCompile(`(?i)=\d+\s+(or|and)\s+\d+=\d+`): true,
}
var XssDetectRules = map[*regexp.Regexp]bool{
	regexp.MustCompile("javascript:[^\\\\s]*"):                                   true,
	regexp.MustCompile("<object[^>]*>.*?</object>"):                              true,
	regexp.MustCompile("on\\\\w+=\\\"[^\\\"]*\\\""):                              true,
	regexp.MustCompile("<script[^>]*>.*?</script>"):                              true,
	regexp.MustCompile("on\\\\w+='[^']*'"):                                       true,
	regexp.MustCompile("<iframe[^>]*>.*?</iframe>"):                              true,
	regexp.MustCompile("<embed[^>]*>.*?</embed>"):                                true,
	regexp.MustCompile("srcdoc=\\\"[^\\\"]*\\\""):                                true,
	regexp.MustCompile("<img[^>]*src=\\\"[^\\\"]*javascript:[^\\\"]*\\\"[^>]*>"): true,
}
