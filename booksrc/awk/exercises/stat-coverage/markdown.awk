#!/usr/bin/awk -f

BEGIN {
	FS = "\t"
	OFS = "|"
	ORS = "|\n"
}
NR == 2 {
	printf("|")
	for(i = 0; i < NF; i++)
		printf("---|")
	printf("\n")
}
{
	$1 = $1
	printf("|")
	print
}
