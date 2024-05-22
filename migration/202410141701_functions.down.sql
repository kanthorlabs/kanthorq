BEGIN;

DROP FUNCTION IF EXISTS kanthorq_consumer_pull;
DROP FUNCTION IF EXISTS consumer_ensure;
DROP FUNCTION IF EXISTS stream_ensure;

COMMIT;