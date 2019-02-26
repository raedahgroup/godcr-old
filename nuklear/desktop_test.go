package nuklear

import (
	"context"
	"testing"

	"github.com/raedahgroup/godcr/app"
)

func TestLaunchApp(t *testing.T) {
	tests := []struct {
		name             string
		ctx              context.Context
		walletMiddleware app.WalletMiddleware
		wantErr          bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := LaunchApp(tt.ctx, tt.walletMiddleware); (err != nil) != tt.wantErr {
				t.Errorf("LaunchApp() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
