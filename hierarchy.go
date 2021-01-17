package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

var MissingManagerError = errors.New("Manager not found in company")
var InvalidEmployeesError = errors.New("Invalid employees were attempted to be imported")
var MultipleCEOError = errors.New("Multiple CEOs found")
var MissingCEOError = errors.New("CEO not found in company")

type RawEmployee struct {
	Name    string
	ID      string
	Manager string
}

type Employee struct {
	Name     string
	ID       string
	Manager  *Employee
	Managees *[]Employee
}

type Company struct {
	Employees map[string]Employee
	CEO       *Employee
}

// AddEmployee adds a single employee to the company from raw employee data
func (company *Company) AddEmployee(rawEmployee RawEmployee) (*Employee, error) {
	managees := make([]Employee, 0)
	employee := Employee{
		Name:     rawEmployee.Name,
		ID:       rawEmployee.ID,
		Managees: &managees,
	}

	// Check if the employee has a manager and add as CEO if they do not
	if rawEmployee.Manager == "" {
		if company.CEO != nil {
			return nil, MultipleCEOError
		}

		company.CEO = &employee
	} else {

		// Otherwise check if this new employee's manager exists in our current employee map
		if manager, ok := company.Employees[rawEmployee.Manager]; ok {

			employee.Manager = &manager
			*manager.Managees = append(*manager.Managees, employee)
		} else {

			// If not, fail with MissingManagerError so we can try and re-add if the manager is found later
			return nil, MissingManagerError
		}
	}

	company.Employees[employee.ID] = employee

	return &employee, nil
}

// ImportRawEmployeeData imports and organises bulk employee data retrieved from csv
func (company *Company) ImportRawEmployeeData(rawEmployees []RawEmployee) error {
	manageeMap := make(map[string][]RawEmployee)

	for _, rawEmployee := range rawEmployees {
		employee, err := company.AddEmployee(rawEmployee)
		if err == MissingManagerError {

			// If an employee was attempted to be added without their manager already
			// existing in our 'database', we add them to this map and wait if/until
			// their manager is found
			manager := rawEmployee.Manager
			manageeMap[manager] = append(manageeMap[manager], rawEmployee)

			continue
		} else if err != nil {
			return err
		}

		// Check our managee map to see if this new employee manages an employee we
		// were not able to previously add
		err = CheckManageeMap(manageeMap, employee, company)
		if err != nil {
			return err
		}
	}

	// If there are still items in the managee map then there are employees without a valid manager
	if len(manageeMap) > 0 {
		return InvalidEmployeesError
	}

	return nil
}

// GenerateHierarchyTable generates and prints the hierarchy structure in a basic tabular form
func (company *Company) GenerateHierarchyTable() ([]string, error) {
	ceo := company.CEO
	if ceo == nil {
		// Data import will catch this but check anyway
		return nil, MissingCEOError
	}

	output := EmployeeWalk(*ceo, 0)

	return output, nil
}

// CheckManageeMap checks our managee map to see if an employee manages anyone that could not previously be added
// If they do, add the managee/s as employees and perform this check recursively
func CheckManageeMap(manageeMap map[string][]RawEmployee, manager *Employee, company *Company) error {
	if managees, ok := manageeMap[manager.ID]; ok {

		for _, managee := range managees {
			// We can now add these managees as the manager has been added to the database
			employee, err := company.AddEmployee(managee)
			if err != nil {
				return err
			}

			// Recursively check our map, in case these new employee/s also incomplete employees
			CheckManageeMap(manageeMap, employee, company)
		}

		// Delete key so we can check the length of the map after our import has finished
		delete(manageeMap, manager.ID)
	}

	return nil
}

// GenerateTabularRow generates a single table row for our hierarchy terminal output
func GenerateTabularRow(name string, depth int) string {
	return strings.Repeat("\t", depth) + name
}

// EmployeeWalk recursively walk through our employee list to generate our string table output
func EmployeeWalk(employee Employee, depth int) []string {
	output := []string{GenerateTabularRow(employee.Name, depth)}

	depth++

	for _, managee := range *employee.Managees {
		output = append(output, EmployeeWalk(managee, depth)...)
	}

	return output
}

func main() {
	if len(os.Args) < 2 {
		return
	}

	path := os.Args[1]

	output, err := GenerateHierarchyFromCSV(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, row := range output {
		fmt.Println(row)
	}
}

// GenerateHierarchyFromCSV generates our hierarchy table when given a path to a valid CSV
func GenerateHierarchyFromCSV(path string) ([]string, error) {
	data, err := ImportCSVFile(path)
	if err != nil {
		return nil, err
	}

	rawEmployees, err := ReadRawEmployeeData(data)
	if err != nil {
		return nil, err
	}

	company := &Company{
		Employees: make(map[string]Employee),
	}

	err = company.ImportRawEmployeeData(rawEmployees)
	if err != nil {
		return nil, err
	}

	return company.GenerateHierarchyTable()
}

// ImportCSVFile Imports CSV file given path
func ImportCSVFile(path string) ([][]string, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, err
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	r := csv.NewReader(file)

	return r.ReadAll()
}

// Read CSV [][]string data into our RawEmployee struct
func ReadRawEmployeeData(data [][]string) ([]RawEmployee, error) {

	rawEmployees := new([]RawEmployee)
	for _, line := range data {
		rawEmployee := RawEmployee{
			Name:    line[0],
			ID:      line[1],
			Manager: line[2],
		}

		*rawEmployees = append(*rawEmployees, rawEmployee)
	}

	return *rawEmployees, nil
}
