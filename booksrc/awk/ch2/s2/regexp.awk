#!/usr/bin/awk -f

$1 ~ $2 {
	print "true"
}
$1 !~ $2 {
	print "false"
}
