package urlutils

import (
	"errors"
	"github.com/goware/urlx"
	"regexp"
)

const URL string = `^(?:http(s)?:\/\/)?[a-zA-Zа-яА-Я0-9.-]+(?:\.[a-zA-Zа-яА-Я0-9\.-]+)+[a-zA-Zа-яА-Я0-9\-\._~:/%?#[\]@!\$&'\(\)\*\+,;=.]+$`

var RegexpURL, _ = regexp.Compile(URL)

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
