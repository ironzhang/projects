# reverse - print input in reverse order by line

{ line[NR] = $0 }

END {
	for(i = NR; i > 0; i--)
		print line[i]
}
