# Hierarchy Generation
The challenge can be viewed in `Coding Challenge.pdf`

## Overview
I approached this problem with the thought that if the database of the company has been initialised well then the generation of our hierarchy table should be straightforward. The bulk of the program's work is done during the import of employee data, where we create a tree of employees in our company by adding them sequentially from the top down. This is done by pending employees until their manager has been successfully added to the company. Then, when we want to generate our hierarchy table, we just walk through the employee tree and append to our output string as we go.

This sort of structure is how I imagine it would look in a web app type scenario, where this code could be seen as the 'C' implementation of a CRUD system. 

## Run

`go run hierarchy.go employees.csv`

## Test

`go test`
