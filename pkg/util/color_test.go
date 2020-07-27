package util_test

import (
	"testing"

	"github.com/plumming/chilly/pkg/util"
	"github.com/stretchr/testify/assert"
)

func TestColorInfo(t *testing.T) {
	raw := "this is a message"
	colored := util.ColorInfo(raw)
	t.Logf("raw=%d", len(raw))
	t.Logf("colored=%d", len(colored))
	t.Logf("strip colored=%d", len(util.Strip(colored)))

	assert.Equal(t, 17, len(raw))
	assert.Equal(t, 17, len(util.Strip(colored)))
}
