package tag_test

import (
	"testing"

	"github.com/user/portwatch/internal/portscanner"
	"github.com/user/portwatch/internal/tag"
)

func entry(port uint16, proto string) portscanner.PortEntry {
	return portscanner.PortEntry{Port: port, Protocol: proto}
}

func TestNew_InvalidProtocol_ReturnsError(t *testing.T) {
	_, err := tag.New([]tag.Rule{{Port: 80, Protocol: "udx", Label: "web"}})
	if err == nil {
		t.Fatal("expected error for invalid protocol")
	}
}

func TestNew_EmptyLabel_ReturnsError(t *testing.T) {
	_, err := tag.New([]tag.Rule{{Port: 80, Label: ""}})
	if err == nil {
		t.Fatal("expected error for empty label")
	}
}

func TestTag_ExactPortMatch(t *testing.T) {
	tgr, _ := tag.New([]tag.Rule{{Port: 443, Protocol: "tcp", Label: "https"}})
	labels := tgr.Tag(entry(443, "tcp"))
	if len(labels) != 1 || labels[0] != "https" {
		t.Fatalf("expected [https], got %v", labels)
	}
}

func TestTag_NoMatch_ReturnsEmpty(t *testing.T) {
	tgr, _ := tag.New([]tag.Rule{{Port: 22, Protocol: "tcp", Label: "ssh"}})
	labels := tgr.Tag(entry(80, "tcp"))
	if len(labels) != 0 {
		t.Fatalf("expected no labels, got %v", labels)
	}
}

func TestTag_AnyProtocol_MatchesBoth(t *testing.T) {
	tgr, _ := tag.New([]tag.Rule{{Port: 53, Label: "dns"}})
	if l := tgr.Tag(entry(53, "tcp")); len(l) != 1 {
		t.Fatalf("expected match for tcp, got %v", l)
	}
	if l := tgr.Tag(entry(53, "udp")); len(l) != 1 {
		t.Fatalf("expected match for udp, got %v", l)
	}
}

func TestTag_MultipleRules_MultipleLabels(t *testing.T) {
	tgr, _ := tag.New([]tag.Rule{
		{Port: 8080, Label: "http-alt"},
		{Port: 8080, Protocol: "tcp", Label: "proxy"},
	})
	labels := tgr.Tag(entry(8080, "tcp"))
	if len(labels) != 2 {
		t.Fatalf("expected 2 labels, got %v", labels)
	}
}
