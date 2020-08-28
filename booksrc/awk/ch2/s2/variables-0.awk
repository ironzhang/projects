#!/usr/bin/awk -f

$4 == "Asia" {print $1, 100 * $2}
