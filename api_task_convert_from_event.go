package kanthorq

import (
	_ "embed"
)

//go:embed api_task_convert_from_event.sql
var TaskConvertFromEventSql string
