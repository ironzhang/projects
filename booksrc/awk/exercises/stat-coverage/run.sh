#!/bin/bash

./stat-coverage.awk coverage.data >coverage.out
./markdown.awk coverage.out >coverage.md
