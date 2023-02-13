package generator

import "time"

func Generator() int64 {
	current_timestamp := time.Now()

	return int64(current_timestamp.Unix())
}
