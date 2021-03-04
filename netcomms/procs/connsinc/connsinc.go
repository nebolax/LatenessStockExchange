package connsinc

import "sync"

var curID = 0
var mu sync.Mutex

//NewID returns new unique conn id
func NewID() int {
	curID++
	return curID
}
