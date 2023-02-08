package main

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParsingTimewarriorInput(t *testing.T) {
	var testInput = `temp.db: /home/user/.timewarrior
temp.report.end: 20160401T000000Z
temp.report.start: 20160430T235959Z
temp.report.tags: "This is a multi-word tag",ProjectA,tag123
temp.version: 0.1.0

[{"start":"20160405T162205Z","end":"20160405T162211Z","tags":["This is a multi-word tag","ProjectA","tag123"]}]`

	config, intervals, err := ParseTimewarrior(strings.NewReader(testInput))
	require.NoError(t, err)
	require.Len(t, config, 5)
	require.Len(t, intervals, 1)

	expectedStartTime, err := time.Parse(time.RFC3339, "2016-04-05T16:22:05Z")
	require.NoError(t, err)
	expectedEndTime, err := time.Parse(time.RFC3339, "2016-04-05T16:22:11Z")
	require.NoError(t, err)

	assert.Equal(t, expectedStartTime, intervals[0].Start.Time())
	assert.Equal(t, expectedEndTime, intervals[0].End.Time())
	assert.Len(t, intervals[0].Tags, 3)
}
