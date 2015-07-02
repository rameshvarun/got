#!/bin/bash

RED=$(tput setaf 1)
GREEN=$(tput setaf 2)
NORMAL=$(tput sgr0)

PASSED=0
FAILED=0
TOTAL=0

for file in $(find ./tests -type f -name "*.rb");
do
    printf "Running test ${file}...\t"
    ruby ${file}

    if [[ $? == 0 ]]
    then
	printf "${GREEN}passed...${NORMAL}\n"
	let PASSED++
    else
	printf "${RED}failed...${NORMAL}\n"
	let FAILED++
    fi
    let TOTAL++
done

echo $PASSED "/" $TOTAL "passed..." $FAILED "/" $TOTAL "failed..."

if (( FAILED > 0 ))
then
  exit 1
else
  exit 0
fi
