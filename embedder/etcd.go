package embedder

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/auth/authpb"
	"go.etcd.io/etcd/embed"
	"go.etcd.io/etcd/etcdserver/etcdserverpb"
	"go.etcd.io/etcd/pkg/types"
	"../log"
	"strconv"
	"time"
)

const (
	Username = "root"
	Password = "613f#8d164df4ACPF49@93a510df49!66f98b*d6"
)

type etcdEmbed struct {
	options *Options
	conf    *embed.Config
	ee      *embed.Etcd
}

func (e *etcdEmbed) Init(ctx context.Context, opts ...Option) (err error) {
	e.options = new(Options)
	for _, opt := range opts {
		opt(e.options)
	}
	e.conf = embed.NewConfig()
	e.conf.Name = e.options.Name
	e.conf.Dir = e.options.Dir
	e.conf.InitialClusterToken = "odin-token"
	e.conf.ClusterState = e.options.ClusterState // "new" or "existing"
	e.conf.EnablePprof = false
	e.conf.TickMs = 200
	e.conf.ElectionMs = 2000
	e.conf.EnableV2 = false

	e.conf.HostWhitelist = e.hostWhitelist(e.options.Cluster)
	e.conf.CORS = e.hostWhitelist(e.options.Cluster)
	e.conf.InitialCluster = e.initialCluster(e.options.PeerPort, e.options.Cluster)

	// Metrics 监控
	//e.conf.Metrics = "basic" //  "extensive"
	//if e.conf.ListenMetricsUrls, err = types.NewURLs([]string{"http://127.0.0.1:1111"}); err != nil {
	//	return
	//}

	//cfg.QuotaBackendBytes : c.cfg.DataQuota
	//cfg.ClusterState : "new"
	e.conf.Logger = "zap"    // Logger is logger options: "zap", "capnslog".
	e.conf.LogLevel = "warn" // "debug" "info" "warn" "error"

	if e.conf.LCUrls, err = types.NewURLs([]string{"http://" + e.options.IP + ":" + e.options.ClientPort}); err != nil {
		return
	}

	if e.conf.ACUrls, err = types.NewURLs([]string{"http://" + e.options.IP + ":" + e.options.ClientPort}); err != nil {
		return
	}

	if e.conf.LPUrls, err = types.NewURLs([]string{"http://" + e.options.IP + ":" + e.options.PeerPort}); err != nil {
		return
	}
	if e.conf.APUrls, err = types.NewURLs([]string{"http://" + e.options.IP + ":" + e.options.PeerPort}); err != nil {
		return
	}
	return
}

func (e *etcdEmbed) Run(ready chan struct{}) (err error) {
	e.ee, err = embed.StartEtcd(e.conf)
	if err != nil {
		log.Sugar.Fatal("embed start failed. error: ", err)
	}

	defer e.ee.Close()

	select {
	case <-e.ee.Server.ReadyNotify():
		ready <- struct{}{}
		log.Sugar.Info("embed server is Ready!")
	case <-time.After(3600 * time.Second):
		e.ee.Server.Stop() // trigger a shutdown
		log.Sugar.Error("embed server took too long to start!")
	}
	log.Sugar.Fatal(<-e.ee.Err())
	return
}

func (e *etcdEmbed) SetAuth(username, password string) (err error) {
	var (
		ul *etcdserverpb.AuthUserListResponse
		rl *etcdserverpb.AuthRoleListResponse
	)
	ee := e.ee
	// 添加用户
	ul, err = ee.Server.AuthStore().UserList(&etcdserverpb.AuthUserListRequest{})
	if ul.Users == nil || len(ul.Users) == 0 || ul.Users[0] != username {
		user := &etcdserverpb.AuthUserAddRequest{
			Name:     username,
			Password: password,
			Options: &authpb.UserAddOptions{
				NoPassword: false,
			},
		}
		_, err = ee.Server.AuthStore().UserAdd(user)
		if err != nil {
			log.Sugar.Error("embed set auth UserAdd failed. error: ", err)
			return
		}
	}

	// 添加角色
	rl, err = ee.Server.AuthStore().RoleList(&etcdserverpb.AuthRoleListRequest{})
	if rl.Roles == nil || len(rl.Roles) == 0 || rl.Roles[0] != username {
		_, err = ee.Server.AuthStore().RoleAdd(&etcdserverpb.AuthRoleAddRequest{Name: username})
		if err != nil {
			log.Sugar.Error("embed set auth RoleAdd failed. error: ", err)
			return
		}
		perm := &etcdserverpb.AuthRoleGrantPermissionRequest{
			Name: username,
			Perm: &authpb.Permission{
				PermType: 2,
				Key:      []byte("/*"),
				RangeEnd: []byte("/*"),
			},
		}
		_, err = ee.Server.AuthStore().RoleGrantPermission(perm)

		if err != nil {
			log.Sugar.Error("embed set auth RoleGrantPermission failed. error: ", err)
			return
		}
	}

	// 关联角色用户
	_, err = ee.Server.AuthStore().UserGrantRole(&etcdserverpb.AuthUserGrantRoleRequest{User: username, Role: username})
	if err != nil {
		log.Sugar.Error("embed set auth UserGrantRole failed. error: ", err)
		return
	}

	// 开启认证
	if !ee.Server.AuthStore().IsAuthEnabled() {
		err = ee.Server.AuthStore().AuthEnable()
		if err != nil {
			log.Sugar.Error("embed set auth AuthEnable failed. error: ", err)
			return
		}
	}
	return
}

func (e *etcdEmbed) IsLeader() bool {
	return e.ee.Server.Leader().String() == e.ee.Server.ID().String()
}

func (e *etcdEmbed) initialCluster(port string, cluster []string) (str string) {
	for i, ip := range cluster {
		str += fmt.Sprintf(",%s=http://%s:%s", "odin"+strconv.Itoa(i), ip, port)
	}
	return str[1:]
}

func (e *etcdEmbed) hostWhitelist(cluster []string) (list map[string]struct{}) {
	list = make(map[string]struct{})
	for _, n := range cluster {
		list[n] = struct{}{}
	}
	return
}

//func NewEmbed(id, dir, ip, cp, pp string, cluster []string) (em *etcdEmbed) {
//	em = new(etcdEmbed)
//
//	em.conf = embed.NewConfig()
//	em.conf.Name = id
//	em.conf.Dir = dir
//
//	em.conf.InitialClusterToken = "odin-token"
//	em.conf.ClusterState = "new"
//	em.conf.EnablePprof = false
//	em.conf.TickMs = 200
//	em.conf.ElectionMs = 2000
//	em.conf.EnableV2 = false
//
//	em.conf.LCUrls, _ = types.NewURLs([]string{"http://" + ip + ":" + cp})
//	em.conf.ACUrls, _ = types.NewURLs([]string{"http://" + ip + ":" + cp})
//
//	em.conf.LPUrls, _ = types.NewURLs([]string{"http://" + ip + ":" + pp})
//	em.conf.APUrls, _ = types.NewURLs([]string{"http://" + ip + ":" + pp})
//
//	em.conf.HostWhitelist = em.hostWhitelist(cluster)
//	em.conf.CORS = em.hostWhitelist(cluster)
//	// []string{"10.0.0.1","10.0.0.2","10.0.0.3"}
//	em.conf.InitialCluster = em.initialCluster(pp, cluster)
//
//	//cfg.ListenMetricsUrls = metricsURLs(c.cfg.PrivateAddress)
//	em.conf.Metrics = "extensive"
//	//cfg.QuotaBackendBytes = c.cfg.DataQuota
//	//cfg.ClusterState = "new"
//	em.Ready = make(chan struct{})
//
//	//em.conf.PeerAutoTLS=true
//	//em.conf.ClientAutoTLS=true
//
//	return
//}

//func main()  {
//	os.Remove("disk")
//	time.Sleep(1*time.Second)
//	e:=NewEmbed("odin0","disk","127.0.0.1","12379","12380",[]string{"127.0.0.1"})
//	e.Run()
//}
//
//func list(e *embed.Etcd)  {
//
//	// 展示user
//	ul,err:=e.Server.AuthStore().UserList(&etcdserverpb.AuthUserListRequest{})
//	fmt.Println("UserList:",err)
//	fmt.Println("UserList:",ul.String())
//	// 展示role
//	rl,err:=e.Server.AuthStore().RoleList(&etcdserverpb.AuthRoleListRequest{})
//	fmt.Println("RoleList:",err)
//	fmt.Println("RoleList:",rl.String())
//}
