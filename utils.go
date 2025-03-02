package main

import "fmt"

func getNewName(baseName string, takenNames map[string]struct{}) string {
	res := ""
	count := 0

	for {
		res = fmt.Sprintf("new_%s_%d", baseName, count)

		if _, ok := takenNames[res]; !ok {
			break
		}

		count++
	}

	return res
}
