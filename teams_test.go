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

func TestRandomPairsReturnsPairs(t *testing.T) {
	members := []string{"A", "B", "C", "D", "E"}
	ps := getRandomPairs(members, []Pair{}, []Pair{})
	if !reflect.DeepEqual(3, len(ps)) {
		t.Error(3, len(ps))
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
	expectedResult := []string{"C, A", "B, D", "E"}
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
