package database

import "fmt"

// Register new type of stock
func AddStock(name string, traderId int) error {
	if !initialized {
		return DatabaseError{"Database is not initialized!"}
	}

	_, err := dataBase.Exec(fmt.Sprintf(
		"INSERT INTO stocks (name, user_id, come_in_time) values ('%s', '%d', 'NULL')",
		name, traderId))

	return err
}
