package actor

import (
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	"github.com/funkygao/fae/servant/proxy"
	"github.com/funkygao/golib/ip"
	"strconv"
	"sync/atomic"
)

type FaeExecutor struct {
	client *proxy.Proxy

	myIp string
	txn  int64
}

func NewFaeExecutor() *FaeExecutor {
	this := new(FaeExecutor)
	this.client = proxy.NewWithDefaultConfig()
	this.myIp = ip.LocalIpv4Addrs()[0]
	return this
}

func (this *FaeExecutor) Start() {
	go this.client.StartMonitorCluster()
	this.client.AwaitClusterTopologyReady()
}

func (this *FaeExecutor) Context(reason string) *rpc.Context {
	ctx := rpc.NewContext()
	ctx.Reason = reason
	rid := atomic.AddInt64(&this.txn, 1)
	ctx.Rid = strconv.FormatInt(rid, 10)
	ctx.Host = this.myIp
	return ctx
}
