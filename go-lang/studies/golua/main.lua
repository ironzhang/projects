function fib(n)
	if n < 2 then
		return n
	else
		return fib(n-1) + fib(n-2)
	end
end

function MatchClusters(tags, clusters)
	base = "hna-sim00-v"
	target = tags["X-Lane-Cluster"]
	for i, cluster in ipairs(clusters) do
		if cluster.Name == target then
			if cluster.AvailableEndpointNum > 0 then
				return { {cluster.Name, 1} }
			end
			break
		end
	end
	return { {base, 0.5}, {"hna-sim100-v", 0.5} }
end

if MatchFuncs == nil then
	MatchFuncs = {}
end

MatchFuncs["www.superdns.com"] = function (tags, clusters)
	base = "hna-sim00-v"
	target = tags["X-Lane-Cluster"]
	for i, cluster in ipairs(clusters) do
		if cluster.Name == target then
			if cluster.AvailableEndpointNum > 0 then
				return { {cluster.Name, 1} }
			end
			break
		end
	end
	return { {base, 1.0} }
end

function main()
	tags = {
		["X-Lane-Cluster"]="hne"
	}
	clusters = {
		{Name="hna", EndpointNum=10, AvailableEndpointNum=9},
		{Name="hnb", EndpointNum=8, AvailableEndpointNum=8},
		{Name="hnc", EndpointNum=8, AvailableEndpointNum=0},
	}

	f = MatchFuncs["www.superdns.com"]

	dests = f(tags, clusters)

	print("name="..dests[1][1]..", percent="..dests[1][2])
end

