#!/usr/bin/env bash
EXIT_CODE=0

ANALYZE_CMD=`gofmt -d -e -s .`
if [[ "$ANALYZE_CMD" != "" ]]; then
	echo "Format error:"
	echo "${ANALYZE_CMD}"
	EXIT_CODE=1
fi

ANALYZE_CMD=`go vet .`
if [[ "$ANALYZE_CMD" != "" ]]; then
	echo "Vet error:"
	echo "${ANALYZE_CMD}"
	EXIT_CODE=1
fi

#sudo apt install golint
ANALYZE_CMD=`golint .`
if [[ "$ANALYZE_CMD" != "" ]]; then
	echo "GoLint:"
	echo "${ANALYZE_CMD}"
	EXIT_CODE=1
fi

#go get -u github.com/kisielk/errcheck
ANALYZE_CMD=`errcheck .`
if [[ "$ANALYZE_CMD" != "" ]]; then
	echo "Error check:"
	echo "${ANALYZE_CMD}"
	EXIT_CODE=1
fi

#gocyclo -top 10  .
exit ${EXIT_CODE}