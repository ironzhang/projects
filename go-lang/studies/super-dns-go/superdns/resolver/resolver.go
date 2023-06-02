package resolver

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ironzhang/super-dns-go/superdns/pkg/filewatch"
	"github.com/ironzhang/super-dns-go/superdns/pkg/model"
	"github.com/ironzhang/super-dns-go/superdns/pkg/superutil"
)

// LoadBalancer 负载均衡器接口
type LoadBalancer interface {
	Pickup(endpoints []model.Endpoint) model.Endpoint
}

// Resolver 域名解析程序
type Resolver struct {
	ResourcePath  string        // 资源路径
	WatchInterval time.Duration // 订阅间隔
	LoadBalancer  LoadBalancer  // 负载均衡器

	once      sync.Once
	mu        sync.Mutex           // 并发互斥锁
	watcher   *filewatch.Watcher   // 文件订阅程序
	providers map[string]*provider // 服务提供方映射表，key 为 domain
}

func (r *Resolver) init() {
	r.watcher = filewatch.NewWatcher(r.WatchInterval)
	r.providers = make(map[string]*provider)
}

// LookupEndpoint 查找地址节点
func (r *Resolver) LookupEndpoint(ctx context.Context, domain string, tags map[string]string) (model.Endpoint, error) {
	return model.Endpoint{}, nil
}

// LookupCluster 查找集群节点
func (r *Resolver) LookupCluster(ctx context.Context, domain string, tags map[string]string) (model.Cluster, error) {
	r.once.Do(r.init)

	r.mu.Lock()
	defer r.mu.Unlock()

	p, ok := r.providers[domain]
	if !ok {
		p = r.newProvider(ctx, domain)
		r.providers[domain] = p
	}
	return p.lookupCluster(ctx, tags)
}

func (r *Resolver) newProvider(ctx context.Context, domain string) *provider {
	// TODO 向 agent 发送订阅域名请求

	// 构建新的服务提供方对象
	p := &provider{}

	// 订阅域名文件
	r.watcher.WatchFile(ctx, r.domainFilePath(domain), func(path string) bool {
		var m model.ProviderModel
		err := superutil.ReadJSON(path, &m)
		if err != nil {
			return false
		}
		p.reload(m)
		return false
	})

	return p
}

func (p *Resolver) domainFilePath(domain string) string {
	return fmt.Sprintf("%s/domains/%s.json", p.ResourcePath, domain)
}
