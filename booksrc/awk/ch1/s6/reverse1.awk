# reverse - print input in reverse order by line

{ line[NR] = $0 } # remember each input line

# print lines in reverse order
END {
	i = NR
	while (i > 0) {
		print line[i]
		i = i - 1
	}
}
