package types

import "regexp"

var SqlInjectRules = map[*regexp.Regexp]bool{
	regexp.MustCompile("(?i)(union)(.*)(select)"): true,
	regexp.MustCompile("(?i)select(.*)from"):      true,
	regexp.MustCompile("(?i)insert into"):         true,
	regexp.MustCompile("(?i)delete from"):         true,
	regexp.MustCompile("(?i)drop table"):          true,
	regexp.MustCompile("(?i)update(.*)set"):       true,
	regexp.MustCompile("--"):                      true,
	regexp.MustCompile("(\\b|\\')(OR|or|oR|Or)('|\\b)\\s*('\\d+'|'\\d+'--\\s*|'\\d+'(\\s+)(--)?|\\d+)(\\s+)(=|like)(\\s+)(\\b|\\')\\d+('|\\b)"): true,
	regexp.MustCompile("/\\*.*\\*/"): true,
	regexp.MustCompile(";"):          true,
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
