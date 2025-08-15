package commands

import (
	"os"
	"testing"

	"github.com/onlyati/quadlet-lsp/internal/utils"
)

type mockMessenger struct {
	receivedPriority utils.MessengerLevel
	receivedMessage  string
}

func (m *mockMessenger) SendMessage(level utils.MessengerLevel, text string) {
	m.receivedMessage = text
	m.receivedPriority = level
}

type mockCommander struct{}

func (m mockCommander) Run(name string, args ...string) ([]string, error) {
	return []string{}, nil
}

func TestPullAll_ValidListJobs(t *testing.T) {
	tmpDir := t.TempDir()
	os.Chdir(tmpDir)
	commandExecutor := NewEditorCommandExecutor(tmpDir)
	commandExecutor.syncCall = true
	messenger := mockMessenger{}

	err := commandExecutor.Run(
		"listJobs",
		&messenger,
		mockCommander{},
	)
	if err != nil {
		t.Fatalf("should not get error, but got %s", err.Error())
	}

	if messenger.receivedMessage != "Running tasks: [listJobs]" {
		t.Fatalf("unexpected message: %s", messenger.receivedMessage)
	}

	if messenger.receivedPriority != utils.MessengerInfo {
		t.Fatalf("expected %d priority, but got: %d", utils.MessengerInfo, messenger.receivedPriority)
	}
}

func TestPullAll_InvalidCommand(t *testing.T) {
	tmpDir := t.TempDir()
	os.Chdir(tmpDir)
	commandExecutor := NewEditorCommandExecutor(tmpDir)
	commandExecutor.syncCall = true
	messenger := mockMessenger{}

	err := commandExecutor.Run(
		"nonexist",
		&messenger,
		mockCommander{},
	)
	if err != nil {
		t.Fatalf("should not get error, but got %s", err.Error())
	}

	if messenger.receivedMessage != "Command failed: nonexist, reason: not found" {
		t.Fatalf("unexpected message: %s", messenger.receivedMessage)
	}

	if messenger.receivedPriority != utils.MessengerError {
		t.Fatalf("expected %d priority, but got: %d", utils.MessengerError, messenger.receivedPriority)
	}
}
