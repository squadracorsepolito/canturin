package main

import "fmt"

func getNextNewName(baseName string, entities []entity) string {
	takenNames := make(map[string]struct{})
	for _, tmpEnt := range entities {
		takenNames[tmpEnt.Name()] = struct{}{}
	}

	res := "new_" + baseName
	count := 0

	for {
		if _, ok := takenNames[res]; !ok {
			break
		}

		res = fmt.Sprintf("%s_%d", baseName, count)
		count++
	}

	return res
}
