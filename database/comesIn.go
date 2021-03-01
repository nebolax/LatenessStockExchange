package database

import (
	"fmt"
	"time"
)

const certificationsBoundary = 2

/* Main function of working with certifications of comes in.
 * Adds one certification
 * If come in action is registered, returns true, time of come in and null as the error
 * If certification is registered, but there are not enough certifiers,
 * returns false, null as timestamp and null as error
 * If error occurred, returns false, null as timestamp and error object
 */
func AddCertification(stockId int, certifierId int, timestamp time.Time) (bool, time.Time, error) {
	if !initialized {
		return false, time.Time{}, DatabaseError{"Database is not initialized"}
	}

	_, err := dataBase.Exec(fmt.Sprintf(
		"INSERT INTO comes_in (certifier_id, stock_id, timestamp) values (%d, %d, %s)",
		certifierId, stockId, timestamp.Format(datetimeFormat)))

	if !checkError(err) {
		return false, time.Time{}, err
	}

	var current = time.Now()

	var today = time.Date(current.Year(), current.Month(), current.Day(),
		0, 0, 0, 0, current.Location())

	//Get all today's come in reports of this stock
	result, cringe := dataBase.Query(fmt.Sprintf(
		"SELECT timestamp FROM comes_in WHERE (timestamp > %s) AND (stock_id = %d)",
		today.Format(datetimeFormat), stockId))

	if !checkError(cringe) {
		return false, time.Time{}, cringe
	}

	var count = 0
	var totalTime int64 = 0

	// switch from one to another and check them
	for result.Next() {
		var timeStamp time.Time
		err = result.Scan(&timeStamp)

		if !checkError(err) {
			return false, time.Time{}, err
		}

		totalTime += timeStamp.UnixNano() / int64(time.Millisecond)
		count++
	}

	// If it is exactly (no less - not enough, no more - already registered) needed amount of certifications
	if count == certificationsBoundary {
		totalTime = totalTime / int64(count)
		realTimestamp := time.Unix(0, totalTime * int64(time.Millisecond))

		return true, realTimestamp, nil
	}

	return false, time.Time{}, nil
}
