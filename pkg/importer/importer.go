package importer

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"

	"github.com/wangke19/po/pkg/polarion"
)

// ReadWorkItemsCSV reads work items from a CSV file.
// Expected header: title,type (required), description (optional).
func ReadWorkItemsCSV(r io.Reader) ([]polarion.WorkItemInput, error) {
	cr := csv.NewReader(r)
	header, err := cr.Read()
	if err != nil {
		return nil, fmt.Errorf("read csv header: %w", err)
	}
	idx, err := csvIndex(header, "title", "type")
	if err != nil {
		return nil, err
	}
	descIdx, hasDesc := idx["description"]

	var items []polarion.WorkItemInput
	for {
		row, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("read csv row: %w", err)
		}
		item := polarion.WorkItemInput{
			Title: row[idx["title"]],
			Type:  row[idx["type"]],
		}
		if hasDesc && descIdx < len(row) {
			item.Description = row[descIdx]
		}
		items = append(items, item)
	}
	return items, nil
}

// ReadWorkItemsJSON reads work items from a JSON array.
func ReadWorkItemsJSON(r io.Reader) ([]polarion.WorkItemInput, error) {
	var items []polarion.WorkItemInput
	if err := json.NewDecoder(r).Decode(&items); err != nil {
		return nil, fmt.Errorf("decode json: %w", err)
	}
	return items, nil
}

// ReadTestResultsCSV reads test records from a CSV file.
// Expected header: caseId,result (required), comment (optional).
func ReadTestResultsCSV(r io.Reader) ([]polarion.TestRecord, error) {
	cr := csv.NewReader(r)
	header, err := cr.Read()
	if err != nil {
		return nil, fmt.Errorf("read csv header: %w", err)
	}
	idx, err := csvIndex(header, "caseId", "result")
	if err != nil {
		return nil, err
	}
	commentIdx, hasComment := idx["comment"]

	var records []polarion.TestRecord
	for {
		row, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("read csv row: %w", err)
		}
		rec := polarion.TestRecord{
			CaseID: row[idx["caseId"]],
			Result: row[idx["result"]],
		}
		if hasComment && commentIdx < len(row) {
			rec.Comment = row[commentIdx]
		}
		records = append(records, rec)
	}
	return records, nil
}

// ReadTestResultsJSON reads test records from a JSON array.
func ReadTestResultsJSON(r io.Reader) ([]polarion.TestRecord, error) {
	var records []polarion.TestRecord
	if err := json.NewDecoder(r).Decode(&records); err != nil {
		return nil, fmt.Errorf("decode json: %w", err)
	}
	return records, nil
}

// csvIndex returns a map of column name -> index for all columns,
// verifying that all required columns are present.
func csvIndex(header []string, required ...string) (map[string]int, error) {
	idx := make(map[string]int, len(header))
	for i, h := range header {
		idx[h] = i
	}
	for _, col := range required {
		if _, ok := idx[col]; !ok {
			return nil, fmt.Errorf("missing required column %q in CSV header", col)
		}
	}
	return idx, nil
}
