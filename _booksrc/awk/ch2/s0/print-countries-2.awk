#!/usr/bin/awk -f

BEGIN {
	print "country", "area", "population"
}
{
	print $1, $2, $3
}
