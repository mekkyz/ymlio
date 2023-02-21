#!/bin/bash

ymlio=~/Work/ymlio/V1.6/ymlio/main
CTEST="unknown"

function fail() {
    echo -e "\033[31;1m ✗ Test [$CTEST] failed:\033[0m $1"
}

function success() {
    echo -e "\033[32;1m ✓ Test [$CTEST] finished successfully.\033[0m"
}

function test1() {
    cleanexit() {
        retval=${1:-0}
        rm -f test.yaml test2.yml
        exit "$retval"
    }

    CTEST=1
    echo "Test 1: basic split"
    rm -f test.yaml test2.yml

    tmpfile=$(mktemp)
    cat <<EOF >"$tmpfile"
test.yaml:
    string_key: value
    int_key: 1
test2.yml:
    string_key2: value2
    int_key2: 2
EOF

    ./main split "$tmpfile"
    ret=$?
    rm -f "$tmpfile"

    if [[ $ret -ne 0 ]]; then
        # if return code not equal 0, we fail immediately
        fail "return code not equal 0"
        return 1
    fi

    if [[ ! -f test.yaml || ! -f test2.yml ]]; then
        fail "output file missing"
        cleanexit 1
    fi

    if ! grep "string_key: value" test.yaml; then
        fail "string_key not found in test.yaml"
        cleanexit 1
    fi
    if ! grep "int_key2: 2" test2.yml; then
        fail "int_key not found in test2.yml"
        cleanexit 1
    fi

    success
    cleanexit 0
}

function test2() {
    CTEST=2
    echo "Test 1: basic combine"

    tmpfile1=$(mktemp).yml
    cat <<EOF >"$tmpfile1"
production:
    string_key: value
    int_key: 1
EOF

    tmpfile2=$(mktemp).yml
    cat <<EOF >"$tmpfile2"
development:
    string_key2: value2
    int_key2: 2
EOF

    ./main combine "$tmpfile1" "$tmpfile2" "./output.yml"
    ret=$?

    rm -f "$tmpfile1" "$tmpfile2"

    if [[ $ret -ne 0 ]]; then
        # if return code not equal 0, we fail immediately
        fail "return code not equal 0"
        return 1
    fi

    if [[ ! -f output.yml ]]; then
        fail "output file missing"
        return 1
    fi

    if ! yq "${tmpfile1}.production.string_key" <output.yml | grep -q 'value'; then
        fail "string_key not found in output.yml"
        rm output.yml
        return 1
    fi

    if ! yq "${tmpfile2}.development.int_key2" <output.yml | grep -q '2'; then
        fail "int_key2 not found in output.yml"
        rm output.yml
        return 1
    fi

    rm output.yml
    success
}

test1
test2
