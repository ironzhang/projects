#!/usr/bin/awk -f

BEGIN {
	srand(2)
}

{
	print rand()
}
