package exitcode

const (
	Success = 0

	InvalidInput = 2
	ConfigError  = 3
	AuthError    = 4
	ToolingError = 5
	CompatError  = 6
	BuildError   = 7
	DeployError  = 8
	RuntimeError = 9
)
