package main

import (
	"fmt"
	"log"

	lua "github.com/yuin/gopher-lua"
)

func main() {
	RunExample1()
	// RunExample2()
}

func RunExample1() {
	ls := lua.NewState()
	defer ls.Close()

	err := ls.DoFile("main.lua")
	if err != nil {
		log.Printf("do file: %v", err)
	}
	err = ls.DoFile("main.lua")
	if err != nil {
		log.Printf("do file: %v", err)
	}

	err = ls.CallByParam(lua.P{
		Fn:      ls.GetGlobal("main"),
		NRet:    0,
		Protect: true,
	})
	if err != nil {
		log.Printf("call by param: %v", err)
	}
}

func RunExample2() {
	tags := map[string]string{
		"X-Lane-Cluster": "hnc",
	}
	clusters := []Cluster{
		{Name: "hna", EndpointNum: 10, AvailableEndpointNum: 9},
		{Name: "hnb", EndpointNum: 8, AvailableEndpointNum: 8},
		{Name: "hnc", EndpointNum: 8, AvailableEndpointNum: 0},
	}
	dests, err := MatchClusters(tags, clusters)
	if err != nil {
		log.Fatalf("match cluster: %v", err)
	}

	for _, dest := range dests {
		log.Printf("Cluster=%s, Percent=%f\n", dest.Cluster, dest.Percent)
	}
}

type Cluster struct {
	Name                 string
	EndpointNum          int
	AvailableEndpointNum int
}

type Destination struct {
	Cluster string
	Percent float64
}

func NewTagsTable(ls *lua.LState, tags map[string]string) *lua.LTable {
	t := ls.NewTable()
	for key, value := range tags {
		t.RawSetString(key, lua.LString(value))
	}
	return t
}

func NewClusterTable(ls *lua.LState, cluster Cluster) *lua.LTable {
	t := ls.NewTable()
	t.RawSetString("Name", lua.LString(cluster.Name))
	t.RawSetString("EndpointNum", lua.LNumber(cluster.EndpointNum))
	t.RawSetString("AvailableEndpointNum", lua.LNumber(cluster.AvailableEndpointNum))
	return t
}

func NewClustersTable(ls *lua.LState, clusters []Cluster) *lua.LTable {
	t := ls.NewTable()
	for _, c := range clusters {
		t.Append(NewClusterTable(ls, c))
	}
	return t
}

func MatchClusters(tags map[string]string, clusters []Cluster) ([]Destination, error) {
	ls := lua.NewState()
	defer ls.Close()

	err := ls.DoFile("main.lua")
	if err != nil {
		log.Printf("do file: %v", err)
	}

	err = ls.CallByParam(lua.P{
		Fn:      ls.GetGlobal("MatchClusters"),
		NRet:    1,
		Protect: true,
	}, NewTagsTable(ls, tags), NewClustersTable(ls, clusters))
	if err != nil {
		log.Printf("call by param: %v", err)
	}

	ret := ls.Get(-1)
	ls.Pop(1)

	tb, ok := ret.(*lua.LTable)
	if !ok {
		return nil, fmt.Errorf("invalid return value")
	}

	var dests []Destination
	for i := 1; i <= tb.Len(); i++ {
		value := tb.RawGetInt(i)
		vtb, ok := value.(*lua.LTable)
		if !ok {
			return nil, fmt.Errorf("invalid return value")
		}
		if vtb.Len() != 2 {
			return nil, fmt.Errorf("invalid return value")
		}

		var dest Destination
		dest.Cluster = lua.LVAsString(vtb.RawGetInt(1))
		dest.Percent = float64(lua.LVAsNumber(vtb.RawGetInt(2)))
		dests = append(dests, dest)
	}
	return dests, nil
}
