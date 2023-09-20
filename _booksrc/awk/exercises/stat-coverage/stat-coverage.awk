#!/usr/bin/awk -f

BEGIN {
	FS = "\t"
}

NR == 1 {
	for(i=2; i <= NF; i++)
		first_business[i] = $i
}
NR == 2 {
	for(i=2; i <= NF; i++) {
		first_total[first_business[i]] = $i
	}
}
NR == 3 {
	for(i=2; i <= NF; i++)
		first_used[first_business[i]] = $i
}
NR == 4 {
	for(i=2; i <= NF; i++)
		first_cover[first_business[i]] = $i
}

NR == 5 {
	for(i=2; i <= NF; i++)
		core_business[i] = $i
}
NR == 6 {
	for(i=2; i <= NF; i++) {
		core_total[core_business[i]] = $i
	}
}
NR == 7 {
	for(i=2; i <= NF; i++) {
		core_used[core_business[i]] = $i
	}
}
NR == 8 {
	for(i=2; i <= NF; i++) {
		core_cover[core_business[i]] = $i
	}
}

END {
	order[1] = "地图"
	order[2] = "引擎"
	order[3] = "国际化"
	order[4] = "2轮车"
	order[5] = "车服"
	order[6] = "安全"
	order[7] = "业务中台"
	order[8] = "业务平台"
	order[9] = "顺风车"
	order[10] = "金融"
	order[11] = "代驾"
	order[12] = "企业"
	order[13] = "汇总"

	owner["地图"] = "黄楚宏"
	owner["引擎"] = "张慧"
	owner["车服"] = "陈洁"
	owner["安全"] = "吴迪"
	owner["业务中台"] = "曹智轶"
	owner["顺风车"] = "张慧"
	owner["金融"] = "曹智轶"
	owner["企业"] = "吴迪"
	owner["汇总"] = "张慧"

	printf("%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
		   "业务线",
		   "一级服务总数", "已接入一级服务数", "一级服务覆盖率",
		   "核心服务总数", "已接入核心服务数", "核心服务覆盖率",
		   "owner")

	for(i=1; i <= length(order); i++) {
		name = order[i]
		printf("%s\t%d\t%d\t%g\t%d\t%d\t%g\t%s\n",
			   name,
			   first_total[name], first_used[name], first_cover[name],
			   core_total[name], core_used[name], core_cover[name],
			   owner[name])
	}
}
