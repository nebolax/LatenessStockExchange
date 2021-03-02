package general

import (
	"fmt"
	"time"
)

// Delete element from list by index
func Remove(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}

func TimeToString(time time.Time) string {
	var res string = fmt.Sprint(time)

	res = res[:19]

	return res
}

// Check error for null-ness (if there any exception)
func CheckError(err error) bool{
	if err != nil {
		if DebugMode {
			panic(err)
		}
		return false
	}
	return true
}
