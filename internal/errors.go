package internal

import "errors"

var (
	ErrNoteExists   = errors.New("note already exists")
	ErrNoteNotFound = errors.New("note not found")
	ErrEditorLaunch = errors.New("editor launch failed")
)
