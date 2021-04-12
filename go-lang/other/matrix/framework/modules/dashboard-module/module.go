package dashboard_module

import (
	"github.com/ironzhang/matrix/framework"
	"github.com/ironzhang/matrix/framework/modules/backend-module"
	"github.com/ironzhang/matrix/restful"
	"github.com/ironzhang/matrix/tlog"
)

var Module = &M{}

func init() {
	framework.Register(Module, nil, nil)
}

type M struct {
}

func (m *M) Name() string {
	return "dashboard-module"
}

func (m *M) Init() (err error) {
	h := handlers{configs: framework.Configs()}
	mux := restful.NewServeMux(nil)
	if err = h.Register(mux); err != nil {
		return err
	}
	backend_module.Module.Handle("/dashboard/", mux)
	backend_module.Module.Handle("/dashboard/log/level", tlog.Level())
	return nil
}

func (m *M) Fini() error {
	return nil
}
