#!/bin/bash

# http://redsymbol.net/articles/unofficial-bash-strict-mode/
set -euo pipefail
set -x
IFS='\t\n'

function error() {
	echo "$@" >&2
}

function log() {
	echo "[$(date)] $@" #>/dev/null
}

# true if a URL (passed in via $1)
# returns a 2XX status code when cURL'ed
#
# false otherwise
#
# should time out after 3 seconds
function urlValid() {
	if [[ -z "$1" ]]; then
		error "No argument provided to urlValid"
		return 1
	fi

        local status=$(curl -o /dev/null --silent --write-out '%{http_code}\n' --max-time 3 "$1")

	log "Got status code of $status for URL $1"
        if [[ $status -ge 200 && $status -lt 300 ]]; then
            return 0
        else
            return 1
        fi
}

function testMultipleFiles() {
	echo "hello world 1" > /tmp/urlify-test-1.txt
	echo "hello world 2" > /tmp/urlify-test-2.txt

	set +u
	urls=$(go run main.go /tmp/urlify-test-1.txt /tmp/urlify-test-2.txt)
	if [[ $! -ne 0 ]]; then
		error "testMultipleFiles() failed - could not run script"
		return 1
	fi
	set -u

	for url in "$urls"; do
		echo "URL is $url"
            #urlValid "$urls"
	done
        #if [[ $? -ne 0 ]]; then
        #        error "testMultipleFiles() failed - returned URL not valid"
        #        return 1
        #fi
}

function testBasicUsage() {
	echo "hello world" > /tmp/urlify-test.txt
	url=$(go run main.go /tmp/urlify-test.txt)
	if [[ $! -ne 0 ]]; then
		error "testBasicUsage() failed - could not run script"
		return 1
	fi

        urlValid "$url"
        if [[ $? -ne 0 ]]; then
                error "testBasicUsage() failed - returned URL not valid"
                return 1
        fi
}

function testProcessSubstitution() {
	contents="hello world"
	url=$(go run main.go <(echo "$contents"))
	if [[ $! -ne 0 ]]; then
		error "testProcessSubstitution() failed - could not run script"
		return 1
	fi

        urlValid "$url"
        if [[ $? -ne 0 ]]; then
                error "testProcessSubstitution() failed - returned URL not valid"
                return 1
        fi
}

function runTests () {
    failures=0

    #if ! testBasicUsage; then
    #        echo "basic usage test failed"
    #        ((failures++))
    #fi

    #if ! testProcessSubstitution; then
    #        echo "process substitution test failed"
    #        ((failures++))
    #fi

    if ! testMultipleFiles; then
	    echo "multiple file test failed"
	    ((failures++))
    fi

    echo "$failures failures"
}

runTests
#urlValid "www.google.com"
