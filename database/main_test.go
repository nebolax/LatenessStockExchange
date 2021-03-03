package database

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	Init(testPath)
	defer clearTest()
	defer dataBase.Close()
	os.Exit(m.Run())
}
