#!/usr/bin/awk -f

BEGIN {
	FS = OFS = "\t"
}

{
	$1 = substr($1, 1, 3)
	print
}
