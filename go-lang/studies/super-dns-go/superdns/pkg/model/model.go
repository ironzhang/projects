package model

// ConsumerModel 消费方模型
type ConsumerModel struct {
}

// ProviderModel 提供方模型
type ProviderModel struct {
	Domain              string             // 域名
	Clusters            map[string]Cluster // 集群映射表
	DefaultDestinations []Destination      // 默认目标节点
}
