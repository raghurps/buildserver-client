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

// TCBuildPayload ...
type TCBuildPayload struct {
	BuildType  TCBuildType       `json:"buildType"`
	Comment    TCBuildComment    `json:"comment"`
	Properties TCBuildProperties `json:"properties"`
	Personal   string            `json:"personal"` // "true" or "false"
	BranchName string            `json:"branchName"`
}

// TCBuildDetails ...
type TCBuildDetails struct {
	ID          int               `json:"id"`
	BuildTypeID string            `json:"buildTypeId"`
	Number      string            `json:"number"`
	Status      string            `json:"status"`
	State       string            `json:"state"`
	BranchName  string            `json:"branchName"`
	WebURL      string            `json:"webUrl"`
	StatusText  string            `json:"statusText"`
	Comment     TCBuildComment    `json:"comment"`
	BuildType   TCBuildType       `json:"buildType"`
	Properties  TCBuildProperties `json:"properties"`
}

// TCBuildStopPayload ...
type TCBuildStopPayload struct {
	Comment        string `json:"comment"`
	ReaddIntoQueue string `json:"readdIntoQueue"` // "true" or "false"
}
