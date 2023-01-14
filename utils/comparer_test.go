package utils

import (
	"github.com/google/go-cmp/cmp"
	"strings"
)

var (
	WithoutSpaces = cmp.Transformer("SpacesIgnored", func(in string) string {
		return strings.TrimSpace(in)
	})
)
