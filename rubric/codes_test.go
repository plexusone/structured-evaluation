package rubric

import "testing"

func TestGetReasonCodeInfo(t *testing.T) {
	// Test existing code
	info := GetReasonCodeInfo(CodeREQAmbiguous)
	if info == nil {
		t.Fatal("Expected info for REQ-AMBIGUOUS")
	}
	if info.Code != CodeREQAmbiguous {
		t.Errorf("Expected code %s, got %s", CodeREQAmbiguous, info.Code)
	}
	if info.Category != CategoryREQ {
		t.Errorf("Expected category '%s', got %s", CategoryREQ, info.Category)
	}
	if info.Description == "" {
		t.Error("Expected non-empty description")
	}
	if info.RepairPrompt == "" {
		t.Error("Expected non-empty repair prompt")
	}

	// Test non-existent code
	unknown := GetReasonCodeInfo("NONEXISTENT_CODE")
	if unknown != nil {
		t.Error("Expected nil for unknown code")
	}
}

func TestGetReasonCodesByCategory(t *testing.T) {
	reqCodes := GetReasonCodesByCategory(CategoryREQ)
	if len(reqCodes) == 0 {
		t.Error("Expected at least one requirement code")
	}

	// Verify all returned codes are in the requirement category
	for _, code := range reqCodes {
		info := GetReasonCodeInfo(code)
		if info == nil || info.Category != CategoryREQ {
			t.Errorf("Code %s should be in %s category", code, CategoryREQ)
		}
	}
}

func TestGetReasonCodesBySpecType(t *testing.T) {
	prdCodes := GetReasonCodesBySpecType("prd")
	if len(prdCodes) == 0 {
		t.Error("Expected at least one code for prd spec type")
	}

	// Verify each code applies to prd
	for _, code := range prdCodes {
		info := GetReasonCodeInfo(code)
		if info == nil {
			continue
		}
		found := false
		for _, st := range info.SpecTypes {
			if st == "prd" {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Code %s should apply to prd spec type", code)
		}
	}
}

func TestAllReasonCodeCategories(t *testing.T) {
	categories := AllReasonCodeCategories()
	if len(categories) == 0 {
		t.Error("Expected at least one category")
	}

	// Verify expected categories exist
	categorySet := make(map[string]bool)
	for _, cat := range categories {
		categorySet[cat] = true
	}

	expected := []string{CategoryREQ, CategoryMETRIC, CategorySEC, CategoryARCH}
	for _, exp := range expected {
		if !categorySet[exp] {
			t.Errorf("Expected category %s to exist", exp)
		}
	}
}

func TestReasonCodeRegistry_Complete(t *testing.T) {
	// Verify all codes in registry have required fields
	for code, info := range ReasonCodeRegistry {
		if info.Code != code {
			t.Errorf("Code mismatch: key=%s, info.Code=%s", code, info.Code)
		}
		if info.Category == "" {
			t.Errorf("Code %s missing category", code)
		}
		if info.Description == "" {
			t.Errorf("Code %s missing description", code)
		}
		if info.RepairPrompt == "" {
			t.Errorf("Code %s missing repair prompt", code)
		}
		if info.DefaultSeverity == "" {
			t.Errorf("Code %s missing default severity", code)
		}
	}
}

func TestSecurityCodesExist(t *testing.T) {
	// Verify critical security codes exist
	securityCodes := []ReasonCode{
		CodeSECGap,
		CodeSECNoAuth,
		CodeSECNoAuthz,
		CodeSECPrivacy,
	}

	for _, code := range securityCodes {
		info := GetReasonCodeInfo(code)
		if info == nil {
			t.Errorf("Expected security code %s to exist", code)
		}
		if info != nil && info.Category != CategorySEC {
			t.Errorf("Code %s should be in %s category", code, CategorySEC)
		}
	}
}

func TestLegacyCodeMapping(t *testing.T) {
	// Test that legacy codes are normalized correctly
	tests := []struct {
		legacy ReasonCode
		newOne ReasonCode
	}{
		{"AMBIGUOUS_REQUIREMENT", CodeREQAmbiguous},
		{"MISSING_ACCEPTANCE_CRITERIA", CodeREQNoCriteria},
		{"SECURITY_GAP", CodeSECGap},
		{"UNMEASURABLE_SUCCESS_METRIC", CodeMETRICUnmeasurable},
		{"INCOMPLETE_ERROR_HANDLING", CodeARCHNoErrorHandling},
	}

	for _, tc := range tests {
		normalized := NormalizeCode(tc.legacy)
		if normalized != tc.newOne {
			t.Errorf("Expected %s to normalize to %s, got %s", tc.legacy, tc.newOne, normalized)
		}

		// Also verify that GetReasonCodeInfo works with legacy codes
		info := GetReasonCodeInfo(tc.legacy)
		if info == nil {
			t.Errorf("Expected to find info for legacy code %s", tc.legacy)
		}
	}
}

func TestGetRepairPrompt(t *testing.T) {
	prompt := GetRepairPrompt(CodeREQAmbiguous)
	if prompt == "" {
		t.Error("Expected non-empty repair prompt")
	}
	if prompt == "" {
		t.Skip()
	}

	// Should contain actionable guidance
	if len(prompt) < 50 {
		t.Error("Repair prompt seems too short to be actionable")
	}
}

func TestRequiresHumanReview(t *testing.T) {
	// Security codes should require human review
	if !RequiresHumanReview(CodeSECNoAuth) {
		t.Error("Security authentication code should require human review")
	}
	if !RequiresHumanReview(CodeSECPrivacy) {
		t.Error("Privacy code should require human review")
	}

	// Some codes should not require human review
	if RequiresHumanReview(CodeREQAmbiguous) {
		t.Error("Ambiguous requirement can be fixed by AI without human review")
	}
	if RequiresHumanReview(CodeDOCNoDiagram) {
		t.Error("Missing diagram can be added by AI without human review")
	}

	// Unknown codes should default to requiring human review
	if !RequiresHumanReview("UNKNOWN_CODE") {
		t.Error("Unknown codes should default to requiring human review")
	}
}

func TestCodePrefixConsistency(t *testing.T) {
	// Verify all codes start with their category prefix
	for code, info := range ReasonCodeRegistry {
		if code == CodeOther {
			continue // OTHER is a special case
		}

		prefix := info.Category + "-"
		if len(code) < len(prefix) || string(code[:len(prefix)]) != prefix {
			t.Errorf("Code %s should start with prefix %s", code, prefix)
		}
	}
}
