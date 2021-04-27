package zoutil

import (
	"reflect"
	"testing"

	"github.com/coreos/etcd/client"
	"github.com/ironzhang/x-pearls/govern/etcdv2"
	"github.com/ironzhang/zerone"
)

func TestOptions(t *testing.T) {
	tests := []struct {
		o Options
		z zerone.Options
	}{
		{
			o: Options{
				Zerone:   "SZerone",
				Filename: "./router.json",
			},
			z: zerone.SOptions{
				Filename: "./router.json",
			},
		},
		{
			o: Options{
				Zerone:    "DZerone",
				Namespace: "zerone",
				Driver:    etcdv2.DriverName,
				Etcd:      EtcdOptions{Endpoints: []string{"http://localhost:2379"}},
			},
			z: zerone.DOptions{
				Namespace: "zerone",
				Driver:    etcdv2.DriverName,
				Config:    client.Config{Endpoints: []string{"http://localhost:2379"}},
			},
		},
	}
	for i, tt := range tests {
		z, err := tt.o.zerone()
		if err != nil {
			t.Fatalf("%d: zerone: %v", i, err)
		}
		if got, want := z, tt.z; !reflect.DeepEqual(got, want) {
			t.Fatalf("%d: options: got %v, want %v", i, got, want)
		}
	}
}
