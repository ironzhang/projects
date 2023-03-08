package restful

import "fmt"

type API struct {
	Method  string
	Path    string
	Handler interface{}
}

func Register(m *ServeMux, apis []API) (err error) {
	for _, a := range apis {
		if err = m.Add(a.Method, a.Path, a.Handler); err != nil {
			return fmt.Errorf("add api(%s, %s, %T): %v", a.Method, a.Path, a.Handler, err)
		}
	}
	return nil
}
