package runner

// CliOptions defines options available only when running in cli interface mode
type CliOptions struct {
	SyncBlockchain bool `long:"sync" description:"Syncs blockchain when running in cli mode. If used with a command, command is executed after blockchain syncs"`
}
