package polarion

type WorkItem struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Type        string `json:"type"`
	Status      string `json:"status"`
	Author      string `json:"author"`
	Description string `json:"description"`
	URL         string `json:"url"`
}

type WorkItemInput struct {
	Title       string `json:"title,omitempty"`
	Type        string `json:"type,omitempty"`
	Status      string `json:"status,omitempty"`
	Description string `json:"description,omitempty"`
}

type TestRun struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Status   string `json:"status"`
	Template string `json:"template"`
	URL      string `json:"url"`
}

type TestRunInput struct {
	Title    string `json:"title,omitempty"`
	Template string `json:"template,omitempty"`
}

type TestResult struct {
	Result  string `json:"result"` // passed|failed|blocked
	Comment string `json:"comment,omitempty"`
}

type TestRecord struct {
	CaseID  string `json:"caseId"`
	Result  string `json:"result"` // passed|failed|blocked|""
	Comment string `json:"comment,omitempty"`
}

type TestRunProgress struct {
	Total   int `json:"total"`
	Passed  int `json:"passed"`
	Failed  int `json:"failed"`
	Blocked int `json:"blocked"`
	NotRun  int `json:"notRun"`
}

type TestStep struct {
	StepIndex      int    `json:"stepIndex"`
	Action         string `json:"action"`
	ExpectedResult string `json:"expectedResult"`
}

type TestStepInput struct {
	Action         string `json:"action"`
	ExpectedResult string `json:"expectedResult"`
}

type Attachment struct {
	ID          string `json:"id"`
	FileName    string `json:"fileName"`
	Title       string `json:"title"`
	ContentType string `json:"contentType"`
	Size        int64  `json:"size"`
}
