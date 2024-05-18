package core

var QueryConsumerPull = `SELECT * FROM kanthorq_consumer_pull(@consumer_name, cast(@size AS SMALLINT));`
