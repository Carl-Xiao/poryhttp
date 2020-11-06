package handler

import (
	"net/http"
)

//RR 基于权重的随机路由算法
type RR interface {
	Next() interface{}
	Add(node interface{}, weight int)
	RemoveAll()
	Reset()
}

// WNGINX nginx策略
type WNGINX struct {
	nodes []*WeightNginx
	n     int
}

// WeightNginx nginx策略
type WeightNginx struct {
	Node            interface{}
	Weight          int
	CurrentWeight   int
	EffectiveWeight int
}

//Next 获取下一个节点
func (nginx *WNGINX) Next() interface{} {
	if nginx.n == 0 {
		return nil
	}
	if nginx.n == 1 {
		return nginx.nodes[0].Node
	}
	return nextWeightedNode(nginx.nodes).Node
}

//nextWeightedNode 每次选择数组里面最大的一个节点, 因为节点不会超过50个。所以不考虑用其他排序算法
func nextWeightedNode(nodes []*WeightNginx) (best *WeightNginx) {
	total := 0

	for i := 0; i < len(nodes); i++ {
		w := nodes[i]

		if w == nil {
			continue
		}

		w.CurrentWeight += w.EffectiveWeight
		total += w.EffectiveWeight
		if w.EffectiveWeight < w.Weight {
			w.EffectiveWeight++
		}

		if best == nil || w.CurrentWeight > best.CurrentWeight {
			best = w
		}
	}

	if best == nil {
		return nil
	}
	best.CurrentWeight -= total
	return best
}

//Add 添加下一个节点
func (nginx *WNGINX) Add(node interface{}, weight int) {
	weighted := &WeightNginx{
		Node:            node,
		Weight:          weight,
		EffectiveWeight: weight}
	nginx.nodes = append(nginx.nodes, weighted)
	nginx.n++
}

//RemoveAll 移除全部节点
func (nginx *WNGINX) RemoveAll() {
	nginx.nodes = nginx.nodes[:0]
	nginx.n = 0
}

//Reset 重置初始化状态
func (nginx *WNGINX) Reset() {
	for _, s := range nginx.nodes {
		s.EffectiveWeight = s.Weight
		s.CurrentWeight = 0
	}
}

// WeightedRR 初始化负载策略
func WeightedRR() RR {
	return &WNGINX{}
}

var wnginx RR

func init() {
	wnginx = WeightedRR()
	//TODO 负载策略
	wnginx.Add("127.0.0.1:8088", 1)
	wnginx.Add("127.0.0.1:8088", 1)
}

//LoadBalancing 负载均衡策略
func (ps *Server) LoadBalancing(req *http.Request) {
	var proxyHost string
	proxyHost = wnginx.Next().(string)
	req.Host = proxyHost
	req.URL.Host = proxyHost
	req.URL.Scheme = "http"
}
