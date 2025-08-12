package unixcompat

import (
	"log"
	"testing"
)

func TestUnixCompat(t *testing.T) {
	lines := []string{"foo", "bar", "foo", "baz"}

	log.Println("Uniq:", Uniq(lines))
	log.Println("Head 2:", Head(lines, 2))
	log.Println("Basename:", Basename("/path/to/file.txt"))
	log.Println("Reverse:", Reverse("abc"))

	// Test Grep (substring, case-sensitive)
	grep1 := Grep([]string{"foo", "bar", "FooBar"}, "foo", false, false)
	if len(grep1) != 1 || grep1[0] != "foo" {
		t.Errorf("Grep substring failed: got %v", grep1)
	}

	// Test Grep (substring, ignore case)
	grep2 := Grep([]string{"foo", "bar", "FooBar"}, "foo", false, true)
	if len(grep2) != 2 {
		t.Errorf("Grep ignore case failed: got %v", grep2)
	}

	// Test Grep (regex, correct pattern)
	grep3 := Grep([]string{"foo", "bar", "baz123"}, `baz\\d+`, true, false)
	if len(grep3) != 1 || grep3[0] != "baz123" {
		t.Errorf("Grep regex failed: got %v", grep3)
	}

	// Test Grep (regex, Go-style pattern)
	grep3b := Grep([]string{"foo", "bar", "baz123"}, `baz\d+`, true, false)
	if len(grep3b) != 1 || grep3b[0] != "baz123" {
		t.Errorf("Grep regex (Go pattern) failed: got %v", grep3b)
	}

	// Test Grep (regex, no match)
	grep3c := Grep([]string{"foo", "bar", "baz"}, `baz\d+`, true, false)
	if len(grep3c) != 0 {
		t.Errorf("Grep regex (no match) failed: got %v", grep3c)
	}

	// Test Fields
	fields := Fields("a b\tc")
	if len(fields) != 3 || fields[1] != "b" {
		t.Errorf("Fields failed: got %v", fields)
	}

	// Test SplitN
	split := SplitN("a:b:c", ":", 2)
	if len(split) != 2 || split[1] != "b:c" {
		t.Errorf("SplitN failed: got %v", split)
	}

	// Test SortStrings
	toSort := []string{"c", "a", "b"}
	SortStrings(toSort)

	if toSort[0] != "a" || toSort[2] != "c" {
		t.Errorf("SortStrings failed: got %v", toSort)
	}

	// Test SortByValue
	m := map[string]int{"a": 2, "b": 3, "c": 1}

	kvs := SortByValue(m)
	if kvs[0].Key != "b" || kvs[2].Key != "c" {
		t.Errorf("SortByValue failed: got %v", kvs)
	}

	// Test ReplaceAll
	replaced := ReplaceAll("foo bar foo", "foo", "baz")
	if replaced != "baz bar baz" {
		t.Errorf("ReplaceAll failed: got %v", replaced)
	}

	// Test ToLower
	if ToLower("FOO") != "foo" {
		t.Errorf("ToLower failed")
	}

	// Test DateNow (format check)
	date := DateNow("2006-01-02")
	if len(date) != 10 {
		t.Errorf("DateNow failed: got %v", date)
	}
}
