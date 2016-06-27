package teamgen

import (
	"reflect"
	"testing"
)

func TestGetCombinations(t *testing.T) {
	result := getCombinations([]string{"A", "B", "C"})
	expectedResult := []string{"A,B,C", "A,C,B", "B,C,A"}
	if !reflect.DeepEqual(expectedResult, result) {
		t.Error(expectedResult, result)
	}
}

func TestBuildPostMessage(t *testing.T) {
	result := buildPostMessage("A,B,C,D,E", 2)
	expectedResult := "Team 1: A, B\nTeam 2: C, D, E"
	if !reflect.DeepEqual(expectedResult, result) {
		t.Error(expectedResult, result)
	}
}
