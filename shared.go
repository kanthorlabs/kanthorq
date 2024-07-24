package kanthorq

import (
	_ "embed"
)

//go:embed shared_consumer_lock.sql
var ConsumerLockSql string

//go:embed shared_consumer_update_cursor.sql
var ConsumerUpdateCursorSql string
