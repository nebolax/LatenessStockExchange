package database

import (
	"testing"

	"github.com/nebolax/LatenessStockExcahnge/general"
)

func TestErrorConstructor(t *testing.T) {
	var dataString = "This is test string!"

	var testError = DatabaseError{dataString}

	if testError.Error() != dataString {
		t.Error("Strings do not match: " + dataString + " / " + testError.Error())
	}
}

func TestCheckErrorOk(t *testing.T) {
	if !general.CheckError(nil) {
		t.Error("Exception on non-exception object")
	}
}

func TestCheckErrorBad(t *testing.T) {
	if general.CheckError(DatabaseError{"Test"}) {
		t.Error("No exception on exception object")
	}
}
