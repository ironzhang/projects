package luavm

import "testing"

func TestLVM(t *testing.T) {
	lvm := New()

	var err error
	if err = lvm.Load("./scripts/test1.lua"); err != nil {
		t.Fatalf("load: %v", err)
	}
	if err = lvm.Load("./scripts/test2.lua"); err != nil {
		t.Fatalf("load: %v", err)
	}

	tb, err := lvm.GetGlobalTable("MatchFuncs")
	if err != nil {
		t.Fatalf("get global table: %v", err)
	}

	fn, err := GetFunctionFromTable(tb, "www.test1.com")
	if err != nil {
		t.Fatalf("get function from table: %v", err)
	}

	type Cluster struct {
		Name                 string
		Features             map[string]string
		EndpointNum          int
		AvailableEndpointNum int
	}
	tags := map[string]string{
		"X-Lane-Cluster": "hna-sim100-v",
	}
	clusters := []Cluster{
		{Name: "hna", Features: nil, EndpointNum: 10, AvailableEndpointNum: 9},
		{Name: "hnb", Features: nil, EndpointNum: 8, AvailableEndpointNum: 8},
		{Name: "hnc", Features: nil, EndpointNum: 8, AvailableEndpointNum: 0},
	}
	ltags, err := lvm.NewValue(tags)
	if err != nil {
		t.Fatalf("new value: %v", err)
	}
	lclusters, err := lvm.NewValue(clusters)
	if err != nil {
		t.Fatalf("new value: %v", err)
	}

	_, err = lvm.Call(fn, 1, ltags, lclusters)
	if err != nil {
		t.Fatalf("call: %v", err)
	}
}
