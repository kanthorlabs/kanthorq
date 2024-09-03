#!/bin/sh
set -e

COVERAGE_EXPECTED=${COVERAGE_EXPECTED:-"80.0"}
COVEROUT_FILE=${COVEROUT_FILE:-"cover.out"}
COVERAGE_FILE=${COVERAGE_FILE:-"coverage.out"}

if test -f $COVEROUT_FILE; then
  go tool cover -func $COVEROUT_FILE | grep total | awk '{print substr($3, 1, length($3)-1)}' > $COVERAGE_FILE

  COVERAGE_ACUTAL=$(cat $COVERAGE_FILE)
  if [ $(echo "${COVERAGE_ACUTAL} < ${COVERAGE_EXPECTED}" | bc) -eq 1 ]; 
  then
    echo "actual:$COVERAGE_ACUTAL < expected:$COVERAGE_EXPECTED"
    exit 1
  fi

  COVERAGE_OLD=$(cat $COVERAGE_FILE)
  # warn if coverage is decreased
  if [ $(echo "${COVERAGE_ACUTAL} < ${COVERAGE_OLD}" | bc) -eq 1 ]; 
  then
    echo "WARN: new:$COVERAGE_ACUTAL < old:$COVERAGE_OLD"
  fi
  
  
else
  echo "$COVEROUT_FILE is not found"
fi