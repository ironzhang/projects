package dashboard_module

import (
	"context"
	"net/url"

	"github.com/ironzhang/matrix/errs"
	"github.com/ironzhang/matrix/framework/pkg/model"
	"github.com/ironzhang/matrix/restful"
)

type handlers struct {
	configs *model.Values
}

func (h *handlers) Register(m *restful.ServeMux) error {
	apis := []restful.API{
		{"GET", "/dashboard/configs", h.GetConfigs},
		{"GET", "/dashboard/configs/:module", h.GetModuleConfig},
		{"PUT", "/dashboard/configs/:module", h.PutModuleConfig},
	}
	return restful.Register(m, apis)
}

func (h *handlers) GetConfigs(ctx context.Context, values url.Values, req interface{}, resp *interface{}) error {
	*resp = h.configs.Interfaces()
	return nil
}

func (h *handlers) GetModuleConfig(ctx context.Context, values url.Values, req interface{}, resp *interface{}) error {
	module := values.Get(":module")
	c, ok := h.configs.GetInterface(module)
	if !ok {
		return errs.NotFound("configs", module)
	}
	*resp = c
	return nil
}

func (h *handlers) PutModuleConfig(ctx context.Context, values url.Values, req map[string]interface{}, resp *interface{}) (err error) {
	module := values.Get(":module")
	v, ok := h.configs.GetValue(module)
	if !ok {
		return errs.NotFound("configs", module)
	}
	if err = v.Store(req); err != nil {
		return err
	}
	if err = v.Reload(); err != nil {
		return err
	}
	*resp = v.Load()
	return nil
}
