package generator

import "time"

// generates current timestamp for unique id generation
func Generator() int64 {
	current_timestamp := time.Now()

	return int64(current_timestamp.Unix())
}
