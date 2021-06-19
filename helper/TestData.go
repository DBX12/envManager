package helper

import (
	"os"
	"path"
	"testing"
)

func GetTestDataFile(t *testing.T, name string) string {
	t.Helper()
	appPath, err := os.Getwd()
	if err != nil {
		t.Fatalf("FAiled to get cwd with error %s", err.Error())
	}
	testDataDir := path.Join(appPath, "..", "testData")
	return path.Join(testDataDir, name)
}

//
//func copyFixtureFile(t *testing.T, name string) {
//	testDataDir := getTestDataDir(t)
//	sourcePath := path.Join(testDataDir, name)
//	destPath := path.Join(os.TempDir(), name)
//
//	originalFile, err := os.Open(sourcePath)
//	if err != nil {
//		t.Fatalf("Failed to open source file while copying fixture file %s. Error message: %s", sourcePath, err.Error())
//	}
//	defer originalFile.Close()
//
//	newFile, err := os.Create(destPath)
//	if err != nil {
//		t.Fatalf("Failed to open destination file while copying fixture file %s. Error message: %s", destPath, err.Error())
//	}
//	defer newFile.Close()
//
//	_, err = io.Copy(newFile, originalFile)
//	if err != nil {
//		t.Fatalf("Failed to copy fixture file %s. Error message: %s", name, err.Error())
//	}
//	t.Cleanup(func() {
//		err = os.Remove(destPath)
//		t.Fatalf("Failed to delete fixture file %s", destPath)
//	})
//}
