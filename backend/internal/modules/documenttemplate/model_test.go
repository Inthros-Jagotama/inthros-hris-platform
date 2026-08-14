package documenttemplate

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestDocumentTemplateJSONKeysAreSnakeCase(t *testing.T) {
	tpl := DocumentTemplate{
		ID:           "id-1",
		Name:         "Test",
		Code:         "TEST_CODE",
		DocumentType: DocumentTypeContractAgreement,
		Status:       StatusActive,
	}
	b, err := json.Marshal(tpl)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	s := string(b)
	for _, key := range []string{`"id"`, `"name"`, `"code"`, `"document_type"`, `"status"`, `"is_default"`, `"created_at"`} {
		if !strings.Contains(s, key) {
			t.Errorf("expected JSON output to contain %s, got: %s", key, s)
		}
	}
	for _, badKey := range []string{`"ID"`, `"DocumentType"`, `"Status"`} {
		if strings.Contains(s, badKey) {
			t.Errorf("JSON output should not contain PascalCase key %s, got: %s", badKey, s)
		}
	}
}

func TestDocumentTemplateVersionJSONKeysAreSnakeCase(t *testing.T) {
	v := DocumentTemplateVersion{ID: "v-1", TemplateID: "t-1", Version: 1, Content: "<p>x</p>"}
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	s := string(b)
	for _, key := range []string{`"id"`, `"template_id"`, `"version"`, `"content"`, `"paper_size"`, `"orientation"`} {
		if !strings.Contains(s, key) {
			t.Errorf("expected JSON output to contain %s, got: %s", key, s)
		}
	}
}
