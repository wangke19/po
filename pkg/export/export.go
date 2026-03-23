package export

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"

	"github.com/wangke19/po/pkg/polarion"
)

func WriteWorkItemsCSV(w io.Writer, items []polarion.WorkItem) error {
	cw := csv.NewWriter(w)
	if err := cw.Write([]string{"id", "type", "status", "author", "title", "url"}); err != nil {
		return fmt.Errorf("write csv header: %w", err)
	}
	for _, item := range items {
		if err := cw.Write([]string{item.ID, item.Type, item.Status, item.Author, item.Title, item.URL}); err != nil {
			return fmt.Errorf("write csv row: %w", err)
		}
	}
	cw.Flush()
	return cw.Error()
}

func WriteWorkItemsJSON(w io.Writer, items []polarion.WorkItem) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(items); err != nil {
		return fmt.Errorf("write json: %w", err)
	}
	return nil
}

func WriteTestResultsCSV(w io.Writer, records []polarion.TestRecord) error {
	cw := csv.NewWriter(w)
	if err := cw.Write([]string{"caseId", "result", "comment"}); err != nil {
		return fmt.Errorf("write csv header: %w", err)
	}
	for _, r := range records {
		if err := cw.Write([]string{r.CaseID, r.Result, r.Comment}); err != nil {
			return fmt.Errorf("write csv row: %w", err)
		}
	}
	cw.Flush()
	return cw.Error()
}

func WriteTestResultsJSON(w io.Writer, records []polarion.TestRecord) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(records); err != nil {
		return fmt.Errorf("write json: %w", err)
	}
	return nil
}

