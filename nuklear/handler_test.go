package nuklear

import "testing"

func Test_getHandlers(t *testing.T) {
	tests := []struct {
		name string
		want []handlersData
	}{
		// TODO: add test cases
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getHandlers(); !reflect(got, tt.want) {
				t.Errorf("getHandlers() = %v, want %v", got, tt.want)
			}
		})
	}
}
