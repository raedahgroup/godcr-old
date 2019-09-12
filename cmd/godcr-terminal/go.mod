module github.com/raedahgroup/godcr/cmd/godcr-terminal

go 1.12

require (
	github.com/decred/slog v1.0.0
	github.com/jrick/logrotate v1.0.0
	github.com/raedahgroup/godcr/app v0.0.0-00010101000000-000000000000
	github.com/raedahgroup/godcr/cli v0.0.0-20190912053213-48fdd185f0dd
	github.com/raedahgroup/godcr/terminal v0.0.0-20190912053213-48fdd185f0dd
)

replace github.com/raedahgroup/godcr/app => ../../app
