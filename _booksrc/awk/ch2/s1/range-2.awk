#!/usr/bin/awk -f

FNR == 1, FNR == 5 { print FILENAME ": " $0}
