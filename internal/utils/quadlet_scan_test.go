package utils_test

import (
	"testing"

	"github.com/onlyati/quadlet-lsp/internal/utils"
)

func TestFindItems(t *testing.T) {
	text := `[Unit]
Description=description

[Container]
Image=docker.io/library/debian:bookworm-slim
AutoUpdate=registry
Environment=ENV1=VALUE1
Volume=dev-db-volume:/app:rw
Environment=ENV2=VALUE2
# Environment=ENV3=VALUE3

[Service]
Restart=on-failure
RestartSec=5
StartLimitBurst=5
`
	findings := utils.FindItems(text, "[Container]", "Environment")

	if len(findings) != 2 {
		t.Fatalf("Expected 2 founds, got %d", len(findings))
	}

	if findings[0].Property != "Environment" {
		t.Fatalf("Expected to get 'Environment' property, but got %s", findings[0].Property)
	}

	if findings[0].Value != "ENV1=VALUE1" {
		t.Fatalf("Expected to get 'ENV1=VALUE1' value, but got %s", findings[0].Value)
	}

	if findings[0].LineNumber != 6 {
		t.Fatalf("Expected to find in 6th line, but got %d", findings[0].LineNumber)
	}
}
