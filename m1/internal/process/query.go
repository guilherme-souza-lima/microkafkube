package process

// RegisterQuery returns the SQL for inserting a new record
func RegisterQuery() string {
	return `
        INSERT INTO registros 
            (trace_id, payload, byte_size, total_characters) 
        VALUES 
            ($1, $2, $3, $4)
    `
}

func UpdatePublishedQuery() string {
	return `UPDATE registros SET published_to_queue = true WHERE trace_id = $1`
}

func GetByPublishedFalseQuery() string {
	return `
        SELECT trace_id, payload, byte_size, total_characters 
        FROM registros 
        WHERE published_to_queue = false 
        ORDER BY created_at ASC 
        LIMIT 50
    `
}
