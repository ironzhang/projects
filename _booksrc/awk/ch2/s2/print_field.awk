#!/usr/bin/awk -f

#{
#	i = 1
#	while(i <= NF) {
#		print $i
#		i++
#	}
#}

#{
#	i = 1
#	do {
#		print $i
#		i++
#	} while(i <= NF)
#}

{
	for(i = 1; i <= NF; i++)
		print $i
}
