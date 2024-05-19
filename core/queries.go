package core

var QueryConsumerPull = `SELECT cursor_current, cursor_next FROM kanthorq_consumer_pull(@consumer_name, @size);`
