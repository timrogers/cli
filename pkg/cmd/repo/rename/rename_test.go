package rename

import (
	"bytes"
	"testing"

	"github.com/cli/cli/v2/internal/ghrepo"
	"github.com/cli/cli/v2/pkg/cmdutil"
	"github.com/cli/cli/v2/pkg/iostreams"
	"github.com/stretchr/testify/assert"
)

func TestRepoRename_validation(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		isErr     bool
		wantError string
	}{
		{
			name:      "no arguments",
			args:      []string{},
			isErr:     true,
			wantError: "new name argument required when not running interactively",
		},
		{
			name:      "too many arguments",
			args:      []string{"newname", "invalid"},
			isErr:     true,
			wantError: "accepts at most 1 arg(s), received 2",
		},
		{
			name:      "name contains slash",
			args:      []string{"org/newname"},
			isErr:     true,
			wantError: "new repository name cannot contain '/'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			io, _, _, _ := iostreams.Test()
			io.SetStdinTTY(false)
			io.SetStdoutTTY(false)

			f := &cmdutil.Factory{
				IOStreams: io,
			}

			cmd := NewCmdRename(f, nil)
			cmd.SetArgs(tt.args)
			cmd.SetIn(&bytes.Buffer{})
			cmd.SetOut(io.Out)
			cmd.SetErr(io.ErrOut)

			_, err := cmd.ExecuteC()
			if tt.isErr {
				assert.EqualError(t, err, tt.wantError)
				return
			}
		})
	}
}

func TestRepoRename_run(t *testing.T) {
	tests := []struct {
		name          string
		opts          *RenameOptions
		expectedError string
	}{
		{
			name: "valid new name",
			opts: &RenameOptions{
				newRepoSelector: "newname",
			},
		},
		{
			name: "new name contains slash",
			opts: &RenameOptions{
				newRepoSelector: "org/newname",
			},
			expectedError: "new repository name cannot contain '/'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			io, _, _, _ := iostreams.Test()
			tt.opts.IO = io
			tt.opts.BaseRepo = func() (ghrepo.Interface, error) {
				return ghrepo.New("owner", "repo"), nil
			}

			err := renameRun(tt.opts)
			if tt.expectedError != "" {
				assert.EqualError(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}