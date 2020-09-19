package buildserver

// BuildServer ...
type BuildServer interface {
	GetBuild(int, interface{}) error
	StartBuild(string, string, string, map[string]string, map[string]int, map[string]int) (int, error)
	CancelQueuedBuild(int, string) error
	StopBuild(int, string) error
	GetArtifactTextFile(string, int) ([]byte, string, error)
}
