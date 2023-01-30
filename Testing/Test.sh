#!/usr/bin/bash

./main combine Test/a.yml Test/b.yml /tmp/output.yml
cat /tmp/output.yml | yq -P

