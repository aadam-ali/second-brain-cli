package cmd

import (
	"testing"
	"time"

	"github.com/aadam-ali/second-brain-cli/config"
	"github.com/stretchr/testify/assert"
)

func TestVersionCmd(t *testing.T) {
	wantStdout := "sb development\n"

	gotStdout, _, gotError := captureOutput(versionCmdFunction, versionCmd, []string{})

	assert.NoError(t, gotError)
	assert.Equal(t, wantStdout, gotStdout)
}

func TestDailyCmd(t *testing.T) {
	config.Now = func() time.Time {
		return time.Date(2025, 7, 13, 20, 0, 0, 0, time.UTC)
	}

	wantStdout := `### 2025-07-13 Sunday
@ <location>

* Todo:

* Notes:
---
`

	gotStdout, _, gotError := captureOutput(dailyCmdFunction, dailyCmd, []string{})

	assert.Equal(t, wantStdout, gotStdout)
	assert.NoError(t, gotError)
}
