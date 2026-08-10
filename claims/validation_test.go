package claims

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

// Regression: AccessedAt/ValidatedAt were plain time.Time, so
// `omitempty` did not apply (the zero value is not the empty value for a
// non-pointer struct) and unset fields serialized as the zero-value
// timestamp "0001-01-01T00:00:00Z" instead of being omitted.
func TestExternalValidation_AccessedAtOmittedWhenUnset(t *testing.T) {
	e := ExternalValidation{URL: "https://example.com", SourceType: ExternalReputableVendor}
	b, err := json.Marshal(e)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(b), "accessedAt") {
		t.Errorf("expected accessedAt to be omitted, got %s", b)
	}
}

func TestExternalValidation_AccessedAtSerializedWhenSet(t *testing.T) {
	e := ExternalValidation{URL: "https://example.com", SourceType: ExternalReputableVendor}
	e.WithAccessedAt(time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC))
	b, err := json.Marshal(e)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(b), `"accessedAt":"2026-01-15T00:00:00Z"`) {
		t.Errorf("expected accessedAt to be serialized, got %s", b)
	}
}

func TestInternalValidation_ValidatedAtOmittedWhenUnset(t *testing.T) {
	i := InternalValidation{Method: MethodCodeExecution}
	b, err := json.Marshal(i)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(b), "validatedAt") {
		t.Errorf("expected validatedAt to be omitted, got %s", b)
	}
}

func TestInternalValidation_ValidatedAtSerializedWhenSet(t *testing.T) {
	i := InternalValidation{Method: MethodCodeExecution}
	i.WithValidatedAt(time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC))
	b, err := json.Marshal(i)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(b), `"validatedAt":"2026-01-15T00:00:00Z"`) {
		t.Errorf("expected validatedAt to be serialized, got %s", b)
	}
}

func TestNewExternalValidation_SetsAccessedAt(t *testing.T) {
	v := NewExternalValidation("https://example.com", ExternalReputableVendor)
	if v.External.AccessedAt == nil {
		t.Fatal("expected AccessedAt to be set")
	}
}

func TestNewInternalValidation_SetsValidatedAt(t *testing.T) {
	v := NewInternalValidation(MethodCodeExecution, "path/to/evidence", true)
	if v.Internal.ValidatedAt == nil {
		t.Fatal("expected ValidatedAt to be set")
	}
}
