package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

type body struct {
	Host     string `json:"host"`
	Service  string `json:"service"`
	Username string `json:"username"`
	Password string `json:"password"`
	Query    string `json:"query"`
}

func QueryExecute(c *fiber.Ctx) error {
	Body := new(body)

	if err := c.BodyParser(Body); err != nil {
		return err
	}

	conn := OracleDBConf(*Body)

	rows, err := conn.Query(Body.Query)
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}

	defer rows.Close()

	// Fetch columns dynamically
	columns, err := rows.Columns()
	if err != nil {
		log.Fatalf("Failed to get columns: %v", err)
	}

	// Create a slice to hold column values as interface{}
	values := make([]interface{}, len(columns))
	valuePointers := make([]interface{}, len(columns))

	// Iterate through the result set
	type Response []interface{}

	var response Response

	for rows.Next() {
		// Prepare pointers to each value
		for i := range values {
			valuePointers[i] = &values[i]
		}

		// Scan into the pointers
		if err := rows.Scan(valuePointers...); err != nil {
			log.Fatalf("Failed to scan row: %v", err)
		}

		// Print out the row values
		rowData := make(map[string]interface{})
		for i, colName := range columns {
			// Dereference pointers to get actual values
			val := values[i]
			switch v := val.(type) {
			case []byte:
				rowData[colName] = string(v) // Handle []byte as string
			default:
				rowData[colName] = v
			}
		}
		response = append(response, rowData)
	}
	return c.Status(200).JSON(&fiber.Map{
		"response": response,
	})
}
