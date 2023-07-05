package utils

import (
	"regexp"
	"strings"
)

func XmlCount(xml string) int {
	count := 0
	if strings.TrimSpace(xml) == "" {
		return count
	}
	re := regexp.MustCompile("<item>[^<]+")
	matches := re.FindAllString(xml, -1)
	count = len(matches)
	return count
}

func XmlMerge(xml1 string, xml2 string) string {
	index := strings.Index(xml1, "<item>")
	if index == -1 {
		return xml2
	}
	index = strings.Index(xml2, "<item>")
	if index == -1 {
		return xml1
	}
	builder := strings.Builder{}
	builder.WriteString(xml1[:strings.Index(xml1, "</channel>")])
	builder.WriteString(xml2[strings.Index(xml2, "<item>"):])
	return builder.String()
}
