package main

import (
	"fmt"

	"github.com/wailsapp/wails/v3/pkg/application"
)

func printError(err error) {
	application.Get().Logger.Error(err.Error())
}

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

func newSaveNetworkDialog() *application.SaveFileDialogStruct {
	dialog := application.SaveFileDialog()

	dialog.AddFilter("protobuf binary file", "*.binpb")
	dialog.AddFilter("protobuf JSON file", "*.json")
	dialog.AddFilter("protobuf text file", "*.txtpb")

	return dialog
}

func newOpenNetworkDialog() *application.OpenFileDialogStruct {
	dialog := application.OpenFileDialog()

	dialog.AddFilter("protobuf binary file", "*.binpb")
	dialog.AddFilter("protobuf JSON file", "*.json")
	dialog.AddFilter("protobuf text file", "*.txtpb")

	return dialog
}
