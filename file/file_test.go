package file

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewFileSync(t *testing.T) {
	file := NewFileSync(nil)
	require.NotNil(t, file)
	require.IsType(t, &FileSync{}, file)
}

func TestNewReaderSync(t *testing.T) {
	reader := NewReaderSync(nil)
	require.NotNil(t, reader)
	require.IsType(t, &ReaderSync{}, reader)
}

func TestCountLines(t *testing.T) {
	var tests = []struct {
		fileName string
		want     int
	}{
		{"../test_files/empty_file.txt", 0},
		{"../test_files/one_line_file.txt", 1},
		{"../test_files/ten_line_file.txt", 10},
		{"err", 0},
	}

	for _, test := range tests {
		got, err := CountLines(test.fileName)
		if test.fileName == "err" {
			require.Error(t, err)
			continue
		}
		require.NoError(t, err)
		require.Equal(t, test.want, got)
	}
}
