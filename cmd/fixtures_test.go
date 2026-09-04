package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestFixturesAddRequiresRoundFlag(t *testing.T) {
	flag := add.Flag("round")
	if flag == nil {
		t.Fatal("expected round flag on fixtures add command")
	}
	if flag.Shorthand != "r" {
		t.Fatalf("expected round shorthand r, got %q", flag.Shorthand)
	}
	if _, ok := flag.Annotations[cobra.BashCompOneRequiredFlag]; !ok {
		t.Fatal("expected round flag to be marked as required")
	}
}

func TestFixturesAddResultRequiresRoundFlag(t *testing.T) {
	flag := addResult.Flag("round")
	if flag == nil {
		t.Fatal("expected round flag on fixtures add-result command")
	}
	if flag.Shorthand != "r" {
		t.Fatalf("expected round shorthand r, got %q", flag.Shorthand)
	}
	if _, ok := flag.Annotations[cobra.BashCompOneRequiredFlag]; !ok {
		t.Fatal("expected round flag to be marked as required")
	}
}
