#!/usr/bin/awk -f

BEGIN { FS = "\t" }

{ pop[$4] += $3 }

END {
	if ("Africa" in pop) {
		print "Africa", pop["Africa"]
	}
	if ("Europe" in pop) {
		print "Europe", pop["Europe"]
	}
	if ("Asia" in pop) {
		print "Asia", pop["Asia"]
	}
}
