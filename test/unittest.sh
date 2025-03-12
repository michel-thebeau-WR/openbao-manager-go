#!/bin/bash

declare -a EXISTING_TESTS
declare -a REQUESTED_TESTS

readarray -t EXISTING_TESTS < <( ls -1 unittests/ )
readarray -d " " -t REQUESTED_TESTS < <( echo $@ )

result=0
for tname in "${REQUESTED_TESTS[@]}"; do
    if [ -z "$tname" ]; then
        continue;
    fi
    echo "$tname"
    if [ -f "unittests/${tname}.sh" ]; then
      result=1
    fi
done

exit $result
