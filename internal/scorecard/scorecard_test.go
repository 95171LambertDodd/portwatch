package scorecard_test

import (
	"testing"

	"github.com/user/portwatch/internal/portscanner"
	"github.com/user/portwatch/internal/scorecard"
)

func makeEntry(port int, proto string) portscanner.Entry {
	return portscanner.Entry{
		Port:     port,
		Protocol: proto,
		PID:      1234,
	}
}

func TestNew_ValidRules_ReturnsNonNil(t *testing.T) {
	s, err := scorecard.New([]scorecard.Rule{
		{Name: "high-port", Weight: 10, Match: func(e portscanner.Entry) bool { return e.Port > 1024 }},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil {
		t.Fatal("expected non-nil scorer")
	}
}

func TestNew_EmptyName_ReturnsError(t *testing.T) {
	_, err := scorecard.New([]scorecard.Rule{
		{Name: "", Weight: 5, Match: func(e portscanner.Entry) bool { return true }},
	})
	if err == nil {
		t.Fatal("expected error for empty rule name")
	}
}

func TestNew_NilMatch_ReturnsError(t *testing.T) {
	_, err := scorecard.New([]scorecard.Rule{
		{Name: "bad-rule", Weight: 5, Match: nil},
	})
	if err == nil {
		t.Fatal("expected error for nil Match")
	}
}

func TestNew_NegativeWeight_ReturnsError(t *testing.T) {
	_, err := scorecard.New([]scorecard.Rule{
		{Name: "neg", Weight: -1, Match: func(e portscanner.Entry) bool { return true }},
	})
	if err == nil {
		t.Fatal("expected error for negative weight")
	}
}

func TestEvaluate_NoRulesMatch_ScoreZero(t *testing.T) {
	s, _ := scorecard.New([]scorecard.Rule{
		{Name: "never", Weight: 20, Match: func(e portscanner.Entry) bool { return false }},
	})
	res := s.Evaluate(makeEntry(80, "tcp"))
	if res.Score != 0 {
		t.Errorf("expected score 0, got %d", res.Score)
	}
	if len(res.Reasons) != 0 {
		t.Errorf("expected no reasons, got %v", res.Reasons)
	}
}

func TestEvaluate_MultipleRulesMatch_ScoreAccumulates(t *testing.T) {
	s, _ := scorecard.New([]scorecard.Rule{
		{Name: "rule-a", Weight: 10, Match: func(e portscanner.Entry) bool { return e.Port == 8080 }},
		{Name: "rule-b", Weight: 25, Match: func(e portscanner.Entry) bool { return e.Protocol == "tcp" }},
		{Name: "rule-c", Weight: 5, Match: func(e portscanner.Entry) bool { return e.PID == 1234 }},
	})
	res := s.Evaluate(makeEntry(8080, "tcp"))
	if res.Score != 40 {
		t.Errorf("expected score 40, got %d", res.Score)
	}
	if len(res.Reasons) != 3 {
		t.Errorf("expected 3 reasons, got %v", res.Reasons)
	}
}

func TestEvaluate_EmptyRules_ScoreZero(t *testing.T) {
	s, _ := scorecard.New(nil)
	res := s.Evaluate(makeEntry(443, "tcp"))
	if res.Score != 0 {
		t.Errorf("expected score 0, got %d", res.Score)
	}
}
