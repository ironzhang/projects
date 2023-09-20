#!/usr/bin/awk -f

BEGIN { FS = "\t" }

{ pop[$4] += $3 }

END {
	print max(pop["Africa"], max(pop["Europe"], pop["Asia"]))
}

func max(m, n) {
	return m > n ? m: n
}
