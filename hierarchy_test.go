package main

import (
	"testing"
)

func TestMultipleCEO(t *testing.T) {
	_, err := GenerateHierarchyFromCSV("./test_data/TestMultipleCEO.csv")
	if err != MultipleCEOError {
		t.Errorf("Expected error catch MultipleCEOError")
	}
}

func TestMissingCEO(t *testing.T) {
	_, err := GenerateHierarchyFromCSV("./test_data/TestMissingCEO.csv")
	if err != InvalidEmployeesError {
		t.Errorf("Expected error catch InvalidEmployeesError, received %v", err)
	}
}

func TestMissingManager(t *testing.T) {
	_, err := GenerateHierarchyFromCSV("./test_data/TestMissingManager.csv")
	if err != InvalidEmployeesError {
		t.Errorf("Expected error catch InvalidEmployeesError, received %v", err)
	}
}

func TestOwnManager(t *testing.T) {
	_, err := GenerateHierarchyFromCSV("./test_data/TestOwnManager.csv")
	if err != InvalidEmployeesError {
		t.Errorf("Expected error catch InvalidEmployeesError, received %v", err)
	}
}

func TestCircularManager(t *testing.T) {
	_, err := GenerateHierarchyFromCSV("./test_data/TestCircularManager.csv")
	if err != InvalidEmployeesError {
		t.Errorf("Expected error catch InvalidEmployeesError, received %v", err)
	}
}

func TestHierarchyGeneration(t *testing.T) {
	output, err := GenerateHierarchyFromCSV("./test_data/TestHierarchyGeneration.csv")
	if err != nil {
		t.Errorf("Expected no error, received %v", err)
	}

	expectedOutput := []string{
		"jamie",
		"\talan",
		"\t\tmartin",
		"\t\talex",
		"\tsteve",
		"\t\tdavid",
	}

	for i, row := range output {
		if row != expectedOutput[i] {
			t.Errorf("Expected %s, received %s", expectedOutput[i], row)
		}
	}
}
