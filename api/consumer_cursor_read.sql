-- consumer_cursor_read
SELECT cursor FROM %s WHERE name = @consumer_name FOR UPDATE SKIP LOCKED