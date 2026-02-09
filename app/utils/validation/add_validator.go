package validation

import (
	"regexp"

	"github.com/gookit/validate"
)

func AddValidator() {
	validate.AddValidator("alphadotspace", func(val any) bool {
		namePattern := regexp.MustCompile(`^[A-Za-z. ]+$`)
		str, ok := val.(string)
		if !ok {
			return false
		}
		return namePattern.MatchString(str)
	})
}
