#!/usr/bin/awk -f

{
	$1 = substr($1, 1, 3)
	print
}
