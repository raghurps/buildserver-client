package teamcity

// TCBuildType ...
type TCBuildType struct {
	ID          string `json:"id"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	ProjectName string `json:"projectName,omitempty"`
	ProjectID   string `json:"projectId,omitempty"`
	WebURL      string `json:"webUrl,omitempty"`
}

// TCBuildComment ...
type TCBuildComment struct {
	Text string `json:"text"`
}

// TCBuildProperty ...
type TCBuildProperty struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// TCBuildProperties ...
type TCBuildProperties struct {
	Count    int               `json:"count,omitempty"`
	Property []TCBuildProperty `json:"property"`
}

// TCBuildSnapshotDependencies ...
type TCBuildSnapshotDependencies struct {
	Count  int              `json:"count,omitempty"`
	Builds []TCBuildDetails `json:"build,omitempty"`
}

// TCBuildPayload ...
type TCBuildPayload struct {
	BuildType            TCBuildType                  `json:"buildType"`
	Comment              TCBuildComment               `json:"comment,omitempty"`
	Properties           TCBuildProperties            `json:"properties,omitempty"`
	Personal             string                       `json:"personal"` // "true" or "false"
	BranchName           string                       `json:"branchName,omitempty"`
	SnapshotDependencies *TCBuildSnapshotDependencies `json:"snapshot-dependencies,omitempty"`
	ArtifactDependencies *TCBuildSnapshotDependencies `json:"artifact-dependencies,omitempty"`
}

// TCBuildDetails ...
type TCBuildDetails struct {
	ID                   int                          `json:"id"`
	BuildTypeID          string                       `json:"buildTypeId,omitempty"`
	Number               string                       `json:"number,omitempty"`
	Status               string                       `json:"status,omitempty"`
	State                string                       `json:"state,omitempty"`
	BranchName           string                       `json:"branchName,omitempty"`
	WebURL               string                       `json:"webUrl,omitempty"`
	StatusText           string                       `json:"statusText,omitempty"`
	Comment              TCBuildComment               `json:"comment,omitempty"`
	BuildType            TCBuildType                  `json:"buildType,omitempty"`
	Properties           TCBuildProperties            `json:"properties,omitempty"`
	SnapshotDependencies *TCBuildSnapshotDependencies `json:"snapshot-dependencies,omitempty"`
	ArtifactDependencies *TCBuildSnapshotDependencies `json:"artifact-dependencies,omitempty"`
}

// TCBuildStopPayload ...
type TCBuildStopPayload struct {
	Comment        string `json:"comment"`
	ReaddIntoQueue string `json:"readdIntoQueue"` // "true" or "false"
}
