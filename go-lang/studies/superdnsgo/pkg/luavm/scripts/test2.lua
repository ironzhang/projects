if MatchFuncs == nil then
	MatchFuncs = {}
end

MatchFuncs["www.test2.com"] = function (tags, clusters)
	base = "hna-sim100-v"
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
