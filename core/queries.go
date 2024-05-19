package core

var QueryConsumerPull = `SELECT * FROM kanthorq_consumer_pull(@consumer_name, @size);`
