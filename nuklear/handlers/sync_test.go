package handlers

import (
	"testing"

	"github.com/aarzilli/nucular"
	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/app/walletcore"
)

func TestSyncHandler_BeforeRender(t *testing.T) {
	type fields struct {
		err                         error
		isRendering                 bool
		isShowingPercentageProgress bool
		percentageProgress          int
		report                      string
		status                      walletcore.SyncStatus
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := &SyncHandler{
				err:                         test.fields.err,
				isRendering:                 test.fields.isRendering,
				isShowingPercentageProgress: test.fields.isShowingPercentageProgress,
				percentageProgress:          test.fields.percentageProgress,
				report:                      test.fields.report,
				status:                      test.fields.status,
			}
			s.BeforeRender()
		})
	}
}

func TestSyncHandler_Render(t *testing.T) {
	type fields struct {
		err                         error
		isRendering                 bool
		isShowingPercentageProgress bool
		percentageProgress          int
		report                      string
		status                      walletcore.SyncStatus
	}
	tests := []struct {
		name       string
		fields     fields
		window     *nucular.Window
		wallet     app.WalletMiddleware
		changePage func(string)
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := &SyncHandler{
				err:                         test.fields.err,
				isRendering:                 test.fields.isRendering,
				isShowingPercentageProgress: test.fields.isShowingPercentageProgress,
				percentageProgress:          test.fields.percentageProgress,
				report:                      test.fields.report,
				status:                      test.fields.status,
			}
			s.Render(test.window, test.wallet, test.changePage)
		})
	}
}

func TestSyncHandler_syncBlockchain(t *testing.T) {
	type fields struct {
		err                         error
		isRendering                 bool
		isShowingPercentageProgress bool
		percentageProgress          int
		report                      string
		status                      walletcore.SyncStatus
	}
	tests := []struct {
		name   string
		fields fields
		window *nucular.Window
		wallet app.WalletMiddleware
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := &SyncHandler{
				err:                         test.fields.err,
				isRendering:                 test.fields.isRendering,
				isShowingPercentageProgress: test.fields.isShowingPercentageProgress,
				percentageProgress:          test.fields.percentageProgress,
				report:                      test.fields.report,
				status:                      test.fields.status,
			}
			s.syncBlockchain(test.window, test.wallet)
		})
	}
}

func TestSyncHandler_updateStatusWithPercentageProgress(t *testing.T) {
	type fields struct {
		err                         error
		isRendering                 bool
		isShowingPercentageProgress bool
		percentageProgress          int
		report                      string
		status                      walletcore.SyncStatus
	}
	tests := []struct {
		name               string
		fields             fields
		report             string
		status             walletcore.SyncStatus
		percentageProgress int64
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := &SyncHandler{
				err:                         test.fields.err,
				isRendering:                 test.fields.isRendering,
				isShowingPercentageProgress: test.fields.isShowingPercentageProgress,
				percentageProgress:          test.fields.percentageProgress,
				report:                      test.fields.report,
				status:                      test.fields.status,
			}
			s.updateStatusWithPercentageProgress(test.report, test.status, test.percentageProgress)
		})
	}
}

func TestSyncHandler_updateStatus(t *testing.T) {
	type fields struct {
		err                         error
		isRendering                 bool
		isShowingPercentageProgress bool
		percentageProgress          int
		report                      string
		status                      walletcore.SyncStatus
	}
	tests := []struct {
		name   string
		fields fields
		report string
		status walletcore.SyncStatus
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := &SyncHandler{
				err:                         test.fields.err,
				isRendering:                 test.fields.isRendering,
				isShowingPercentageProgress: test.fields.isShowingPercentageProgress,
				percentageProgress:          test.fields.percentageProgress,
				report:                      test.fields.report,
				status:                      test.fields.status,
			}
			s.updateStatus(test.report, test.status)
		})
	}
}
