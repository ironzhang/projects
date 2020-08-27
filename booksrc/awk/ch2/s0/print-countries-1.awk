#!/usr/bin/awk -f

{
	print \
		  $1, # country name
		  $2, # area in thousands of square miles
		  $3  # population in millions
}
