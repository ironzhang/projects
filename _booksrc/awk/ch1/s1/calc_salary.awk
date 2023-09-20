#!/usr/bin/awk -f

$3 > 0 {
	print $1, $2 * $3
}

$3 == 0 {
	print $1
}
