package buildserver

// BuildServer ...
type BuildServer interface {
	GetBuild(int, interface{}) error
	StartBuild(string, string, map[string]string) (int, error)
	CancelQueuedBuild(int) error
	StopBuild(int) error
	GetArtifactTextFile(string, int) ([]byte, string, error)
}
