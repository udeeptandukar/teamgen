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

func TestBuildPostMessageForOddMembers(t *testing.T) {
	result := buildPostMessage("A,B,C,D,E", 2)
	expectedResult := []string{"A, B", "C, D, E"}
	if !reflect.DeepEqual(expectedResult, result) {
		t.Error(expectedResult, result)
	}
}

func TestBuildPostMessageForEvenMembers(t *testing.T) {
	result := buildPostMessage("A,B,C,D,E,F", 3)
	expectedResult := []string{"A, B", "C, D", "E, F"}
	if !reflect.DeepEqual(expectedResult, result) {
		t.Error(expectedResult, result)
	}

	result = buildPostMessage("A,B,C,D", 2)
	expectedResult = []string{"A, B", "C, D"}
	if !reflect.DeepEqual(expectedResult, result) {
		t.Error(expectedResult, result)
	}
}
