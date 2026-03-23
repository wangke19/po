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
