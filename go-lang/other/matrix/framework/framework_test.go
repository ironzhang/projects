package framework_test

import (
	"testing"

	"github.com/ironzhang/matrix/framework"
	"github.com/ironzhang/matrix/framework/modules/dashboard-module"
	"github.com/ironzhang/matrix/framework/modules/debug-module"
	"github.com/ironzhang/matrix/framework/modules/etcd-module"
	"github.com/ironzhang/matrix/framework/modules/micro-module"
)

var _ = dashboard_module.Module
var _ = debug_module.Module
var _ = etcd_module.Module
var _ = micro_module.Module

func TestFramework(t *testing.T) {
	framework.Main()
}
