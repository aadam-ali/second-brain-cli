package cmd

import (
	"testing"
)

func TestVersionCmd(t *testing.T) {
	var wantError error
	wantOutput := "sb development\n"

	gotOutput, gotError := captureStdout(versionCmdFunction, versionCmd, []string{})

	if gotError != wantError {
		t.Errorf("got %q, want %q", gotError, wantError)
	}

	if gotOutput != wantOutput {
		t.Errorf("got %q, want %q", gotOutput, wantOutput)
	}
}
