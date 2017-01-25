package teamgen

import (
	"reflect"
	"testing"
)

func TestGetPairReturnsRemainingMembers(t *testing.T) {
	members := []string{"A", "B", "C", "D", "E", "F"}
	_, r := getPair(members)
	if !reflect.DeepEqual(4, len(r)) {
		t.Error(4, len(r))
	}
}

func TestPopRandomMemberReturnsOneElement(t *testing.T) {
	members := []string{"A"}
	v, r := popRandomMember(members)
	if !reflect.DeepEqual("A", v) {
		t.Error("A", v)
	}
	if !reflect.DeepEqual([]string{}, r) {
		t.Error([]string{}, r)
	}

	members = []string{}
	v, r = popRandomMember(members)
	if !reflect.DeepEqual("", v) {
		t.Error("", v)
	}
	if !reflect.DeepEqual([]string{}, r) {
		t.Error([]string{}, r)
	}
}

func TestPairExistsReturnsTrueIfExists(t *testing.T) {
	pair := Pair{First: "A", Second: "B"}
	pairs := []Pair{Pair{First: "B", Second: "A"}}
	result := pairExists(pair, pairs)
	if !reflect.DeepEqual(true, result) {
		t.Error(true, result)
	}
}

func TestPairExistsReturnsFalseIfDoesNotExists(t *testing.T) {
	pair := Pair{First: "A", Second: "B"}
	pairs := []Pair{Pair{First: "C", Second: "A"}}
	result := pairExists(pair, pairs)
	if !reflect.DeepEqual(false, result) {
		t.Error(false, result)
	}
}

func TestGetPairsCSVreturnsListOfCSVPairs(t *testing.T) {
	pairs := []Pair{Pair{First: "C", Second: "A"}, Pair{First: "B", Second: "D"}, Pair{First: "E", Second: ""}}
	result := getPairsCSV(pairs)
	expectedResult := []string{"C, A", "B, D, E"}
	if !reflect.DeepEqual(expectedResult, result) {
		t.Error(expectedResult, result)
	}
}

func TestConvertToMemberExclusionPairsConvertsCSVPairsToProperPairs(t *testing.T) {
	csvPairs := []string{"A,B", "C, D", "E"}
	result := convertToMemberExclusionPairs(csvPairs)
	expectedResult := []Pair{Pair{First: "A", Second: "B"}, Pair{First: "C", Second: "D"}, Pair{First: "E", Second: ""}}
	if !reflect.DeepEqual(expectedResult, result) {
		t.Error(expectedResult, result)
	}
}

func TestGenerateCombinationsReturnsListOfPairsForEvenMembers(t *testing.T) {
	members := []string{"A", "B", "C", "D", "E", "F"}
	result := generateCombinations(members)
	expectedResult := []Pair{
		Pair{First: "A", Second: "B"}, Pair{First: "A", Second: "C"}, Pair{First: "A", Second: "D"}, Pair{First: "A", Second: "E"}, Pair{First: "A", Second: "F"},
		Pair{First: "B", Second: "C"}, Pair{First: "B", Second: "D"}, Pair{First: "B", Second: "E"}, Pair{First: "B", Second: "F"},
		Pair{First: "C", Second: "D"}, Pair{First: "C", Second: "E"}, Pair{First: "C", Second: "F"},
		Pair{First: "D", Second: "E"}, Pair{First: "D", Second: "F"},
		Pair{First: "E", Second: "F"},
	}
	if !reflect.DeepEqual(expectedResult, result) {
		t.Error(expectedResult, result)
	}
}

func TestGenerateCombinationsReturnsListOfPairsForOddMembers(t *testing.T) {
	members := []string{"A", "B", "C", "D", "E"}
	result := generateCombinations(members)
	expectedResult := []Pair{
		Pair{First: "A", Second: "B"}, Pair{First: "A", Second: "C"}, Pair{First: "A", Second: "D"}, Pair{First: "A", Second: "E"},
		Pair{First: "B", Second: "C"}, Pair{First: "B", Second: "D"}, Pair{First: "B", Second: "E"},
		Pair{First: "C", Second: "D"}, Pair{First: "C", Second: "E"},
		Pair{First: "D", Second: "E"},
	}
	if !reflect.DeepEqual(expectedResult, result) {
		t.Error(expectedResult, result)
	}
}

func TestGenerateCombinationsReturnsListOfPairsForThreeMembers(t *testing.T) {
	members := []string{"A", "B", "C"}
	result := generateCombinations(members)
	expectedResult := []Pair{
		Pair{First: "A", Second: "B"}, Pair{First: "A", Second: "C"},
		Pair{First: "B", Second: "C"},
	}
	if !reflect.DeepEqual(expectedResult, result) {
		t.Error(expectedResult, result)
	}
}

func TestPairSubtractionReturnsCorrectResult(t *testing.T) {
	combinations := []Pair{
		Pair{First: "A", Second: "B"}, Pair{First: "A", Second: "C"}, Pair{First: "A", Second: "D"}, Pair{First: "A", Second: "E"},
		Pair{First: "B", Second: "C"}, Pair{First: "B", Second: "D"}, Pair{First: "B", Second: "E"},
		Pair{First: "C", Second: "D"}, Pair{First: "C", Second: "E"},
		Pair{First: "D", Second: "E"},
	}
	excludes := []Pair{
		Pair{First: "A", Second: "B"}, Pair{First: "A", Second: "C"},
		Pair{First: "B", Second: "C"},
	}
	result := pairSubtraction(combinations, excludes)
	expectedResult := []Pair{
		Pair{First: "A", Second: "D"}, Pair{First: "A", Second: "E"},
		Pair{First: "B", Second: "D"}, Pair{First: "B", Second: "E"},
		Pair{First: "C", Second: "D"}, Pair{First: "C", Second: "E"},
		Pair{First: "D", Second: "E"},
	}
	if !reflect.DeepEqual(expectedResult, result) {
		t.Error(expectedResult, result)
	}
}
