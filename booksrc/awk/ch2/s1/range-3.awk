#!/usr/bin/awk -f

FNR <= 5 { print FILENAME ": " $0}
