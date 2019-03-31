package app

import "testing"

func TestVersion(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "version",
			want: "0.0.1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Version(); got != tt.want {
				t.Errorf("Version() = %v, want %v", got, tt.want)
			}
		})
	}
}
