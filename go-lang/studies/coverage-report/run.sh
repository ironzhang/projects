#!/bin/bash

./coverage-report MD <coverage.data >coverage.md
./coverage-report XMD <coverage.data >coverage-x.md
./coverage-report HTML <coverage.data >coverage.html
./coverage-report XHTML <coverage.data >coverage-x.html

