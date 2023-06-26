package utils

import (
	"regexp"
	"strconv"
	"strings"
)

var (
	regexSeparator        = regexp.MustCompile(`[\[\]【】]`)
	placeholderSeparator  = "/"
	placeholderSeparators = "/+"
	regexSpecialChar      = regexp.MustCompile(`[\$()\*\+\.?\^{}\|\\]`)
	regexArticle          = regexp.MustCompile(`(\b|\s)((?i)a|an|the)\s`)
	placeholder           = " "
	placeholders          = `\s+`
)

// CleanTitle cleans a given title by replacing certain characters and words with placeholders.
// The `regex` parameter is a regular expression that is used to replace certain patterns
// in the title with a placeholder.
func CleanTitle(title, regex string) string {
	title = regexSeparator.ReplaceAllString(title, placeholderSeparator)
	title = regexSpecialChar.ReplaceAllString(title, placeholder)
	title = regexArticle.ReplaceAllString(title, placeholder)
	cleanTitle := regexp.MustCompile(regex).ReplaceAllString(title, placeholder)
	if strings.TrimSpace(cleanTitle) == "" {
		cleanTitle = title
	}
	cleanTitle = regexp.MustCompile(placeholderSeparators).ReplaceAllString(cleanTitle, placeholderSeparator)
	cleanTitle = regexp.MustCompile(placeholders).ReplaceAllString(cleanTitle, placeholder)
	return strings.ToLower(strings.TrimSpace(cleanTitle))
}

// RemoveYear removes the year (i.e., a four-digit number starting with "19" or "20") from a given title.
func RemoveYear(title string) string {
	if strings.TrimSpace(title) == "" {
		return ""
	}
	return regexp.MustCompile(`\s(19|20)\d{2}$`).ReplaceAllString(title, "")
}

// RemoveEpisode removes the episode number from a given title.
// The episode number is assumed to be a three-digit number (or a two-digit number between 100 and 189)
// preceded by one or more zeros.
func RemoveEpisode(title string) string {
	if strings.TrimSpace(title) == "" {
		return ""
	}
	return regexp.MustCompile(`\s0*(\d{1,3}|1[0-8]\d{2})$`).ReplaceAllString(title, "")
}

// RemoveSeason removes the season number (i.e., a capital letter "S" followed by one or more digits) from a given title.
func RemoveSeason(title string) string {
	if strings.TrimSpace(title) == "" {
		return ""
	}
	return regexp.MustCompile(`\s(S\d+)$`).ReplaceAllString(title, "")
}

// RemoveSeasonEpisode removes both the season number and the episode number from a given title.
func RemoveSeasonEpisode(title string) string {
	if strings.TrimSpace(title) == "" {
		return ""
	}
	return regexp.MustCompile(`\s(S\d+ |)\d+$`).ReplaceAllString(title, "")
}

// ReplaceToken replaces a token (i.e., a string enclosed in curly braces) in a given text with a given value.
// If the `offset` parameter is not nil, it is used to add an integer offset to any digits in the value before replacing the token.
func ReplaceTokenOffset(token, value, text string, offset *int) string {
	if offset != nil && *offset != 0 {
		r := regexp.MustCompile(`(\d+)`)
		value = r.ReplaceAllStringFunc(value, func(s string) string {
			n, _ := strconv.Atoi(s)
			return strconv.Itoa(n + *offset)
		})
	}
	return strings.ReplaceAll(text, "{"+token+"}", value)
}

func ReplaceToken(token, value, text string) string {
	return strings.ReplaceAll(text, token, value)
}

// RemoveToken removes a token from a given text.
func RemoveToken(token, text string) string {
	return strings.ReplaceAll(text, "{"+token+"}", "")
}

// RemoveAllToken removes all tokens from a given text.
func RemoveAllToken(text string) string {
	return regexp.MustCompile(`{[^{}]+}`).ReplaceAllString(text, "")
}
