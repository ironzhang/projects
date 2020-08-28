#!/usr/bin/awk -f

{
	x[NR] = $0
}

END {
	for(i = NR; i > 0; i--)
		print x[i]
}
