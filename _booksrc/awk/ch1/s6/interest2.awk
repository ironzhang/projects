# interest2 - compute compound interest
#	input: amount rate years
#	output: compounded value at the end of each year

{
	print $0
	for (i = 1; i <= $3; i++)
		printf("\t%.2f\n", $1 * (1 + $2) ^ i)
}
