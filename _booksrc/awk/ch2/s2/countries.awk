#!/usr/bin/awk -f

{
	s = s substr($1, 1, 3) " "
}

END {
	print substr(s, 1, length(s)-1) # substr 去除末尾空格
}
