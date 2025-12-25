package pattern

import (
	"reflect"
	"testing"
)

func TestExpandPattern(t *testing.T) {
	tests := []struct {
		name     string
		pattern  string
		want     []string
		wantErr  bool
	}{
		{
			name:    "simple range",
			pattern: "EFTA{10724-10726}.pdf",
			want:    []string{"EFTA10724.pdf", "EFTA10725.pdf", "EFTA10726.pdf"},
			wantErr: false,
		},
		{
			name:    "padded range",
			pattern: "EFTA{00010724-00010726}.pdf",
			want:    []string{"EFTA00010724.pdf", "EFTA00010725.pdf", "EFTA00010726.pdf"},
			wantErr: false,
		},
		{
			name:    "colon separator",
			pattern: "file{1:3}.pdf",
			want:    []string{"file1.pdf", "file2.pdf", "file3.pdf"},
			wantErr: false,
		},
		{
			name:    "no pattern",
			pattern: "simple.pdf",
			want:    []string{"simple.pdf"},
			wantErr: false,
		},
		{
			name:    "invalid range",
			pattern: "file{10-5}.pdf",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "URL with pattern",
			pattern: "https://example.com/file{1-2}.pdf",
			want:    []string{"https://example.com/file1.pdf", "https://example.com/file2.pdf"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExpandPattern(tt.pattern)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExpandPattern() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExpandPattern() = %v, want %v", got, tt.want)
			}
		})
	}
}

