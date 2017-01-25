package teamgen

import (
	"reflect"
	"testing"
)

func TestParseCommandCorrectly(t *testing.T) {
	cmdType, cmdArgs := parseCommand("member-add Apple Banana Orange")
	expectedType := "member-add"
	if cmdType != expectedType {
		t.Error("Expected member-add, got ", cmdType)
	}
	expectedArgs := []string{"Apple", "Banana", "Orange"}
	if !reflect.DeepEqual(expectedArgs, cmdArgs) {
		t.Error(expectedArgs, cmdArgs)
	}
}

func TestParseCommandHandlesNoArgs(t *testing.T) {
	cmdType, cmdArgs := parseCommand("member-list")
	expectedType := "member-list"
	if !reflect.DeepEqual(expectedType, cmdType) {
		t.Error("Expected member-add, got ", cmdType)
	}
	expectedArgs := []string{}
	if !reflect.DeepEqual(expectedArgs, cmdArgs) {
		t.Error(expectedArgs, cmdArgs)
	}
}
