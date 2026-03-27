package polarion

// WorkItem represents a Polarion work item.
type WorkItem struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Type        string `json:"type"`
	Status      string `json:"status"`
	Author      string `json:"author"`
	Description string `json:"description"`
	URL         string `json:"url"`
}

// WorkItemInput holds fields for creating/updating work items.
type WorkItemInput struct {
	Title       string `json:"title,omitempty"`
	Type        string `json:"type,omitempty"`
	Status      string `json:"status,omitempty"`
	Description string `json:"description,omitempty"`
}

// TestRun represents a Polarion test run.
type TestRun struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Status   string `json:"status"`
	Template string `json:"template"`
	URL      string `json:"url"`
}

// TestRunInput holds fields for creating/updating test runs.
type TestRunInput struct {
	Title    string `json:"title,omitempty"`
	Template string `json:"template,omitempty"`
}

// TestResult represents the outcome of a test execution.
type TestResult struct {
	Result  string `json:"result"` // passed|failed|blocked
	Comment string `json:"comment,omitempty"`
}

// TestRecord links a test case to its result in a test run.
type TestRecord struct {
	CaseID  string `json:"caseId"`
	Result  string `json:"result"` // passed|failed|blocked|""
	Comment string `json:"comment,omitempty"`
}

// TestRunProgress tracks execution statistics for a test run.
type TestRunProgress struct {
	Total   int `json:"total"`
	Passed  int `json:"passed"`
	Failed  int `json:"failed"`
	Blocked int `json:"blocked"`
	NotRun  int `json:"notRun"`
}

// TestStep represents a single step in a test case.
type TestStep struct {
	StepIndex      int    `json:"stepIndex"`
	Action         string `json:"action"`
	ExpectedResult string `json:"expectedResult"`
}

// TestStepInput holds fields for creating/updating test steps.
type TestStepInput struct {
	Action         string `json:"action"`
	ExpectedResult string `json:"expectedResult"`
}

// Attachment represents a file attached to a work item.
type Attachment struct {
	ID          string `json:"id"`
	FileName    string `json:"fileName"`
	Title       string `json:"title"`
	ContentType string `json:"contentType"`
	Size        int64  `json:"size"`
}

// WorkItemLink represents a relationship between two work items.
type WorkItemLink struct {
	TargetID string `json:"targetId"`
	Role     string `json:"role"`
}

// Comment represents a comment on a work item.
type Comment struct {
	ID      string `json:"id"`
	Author  string `json:"author"`
	Created string `json:"created"`
	Body    string `json:"body"`
}

// User represents a Polarion user.
type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Project represents a Polarion project.
type Project struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
