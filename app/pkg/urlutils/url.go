package urlutils

import (
	"errors"
	"github.com/goware/urlx"
	"net/url"
	"regexp"
	"strings"
)

const URL string = `^(?:http(s)?:\/\/)?[\w.-]+(?:\.[\w\.-]+)+[\w\-\._~/?#[\]@!\$&'\(\)\*\+,;=.]+$`

var RegexpURL, _ = regexp.Compile(URL)

func IsIgnoredURL(link string, ignored []string) bool {
	u, _ := url.Parse(link)
	parts := strings.Split(u.Hostname(), ".")
	domain := parts[len(parts)-2] + "." + parts[len(parts)-1]

	for _, ignoredURL := range ignored {
		if ignoredURL == domain {
			return true
		}
	}
	return false
}

func Normalize(link string) (normalized string, err error) {
	normalized, err = urlx.NormalizeString(link)
	if err != nil {
		return
	}

	return
}

func Validate(link string) error {
	if !RegexpURL.MatchString(link) {
		return errors.New("invalid link")
	}

	_, err := urlx.Parse(link)
	if err != nil {
		return err
	}

	return nil
}
