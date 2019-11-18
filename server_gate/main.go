package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/zhiguochi/chess/common"
	"github.com/zhiguochi/chess/server_gate/config"
	"github.com/zhiguochi/chess/server_gate/connid"
	"github.com/zhiguochi/chess/server_gate/pkg"
	"github.com/zhiguochi/chess/util/log"
	"github.com/zhiguochi/chess/util/rpc"
	"github.com/zhiguochi/chess/util/services"
)

func refreshBackend() {
	backends := make(map[string]bool)

	for {
		backendList := config.GetBackendConfig()

		hostAndPorts := strings.Split(backendList, ",")
		for i := 0; i < len(hostAndPorts); i++ {
			if _, present := backends[hostAndPorts[i]]; present {
				continue
			}

			if _, _, err := net.SplitHostPort(hostAndPorts[i]); err != nil {
				log.Error("invalid backend addr:%s", hostAndPorts[i])
				continue
			}

			backends[hostAndPorts[i]] = true
			go pkg.DoBackend(hostAndPorts[i])
		}

		time.Sleep(time.Second * 10)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s conf_path\n", os.Args[0])
		return
	}

	log.Info("server start, pid = %d", os.Getpid())

	if !config.Init(os.Args[1]) {
		return
	}

	listenPort := common.GetListenPort()
	gateid := common.GetGateid()

	log.Info("listenPort=%d", listenPort)
	log.Info("gateid=%d", gateid)

	rpc.Add(services.Center, common.GetCenterAddr(), 1)
	if !services.DelConnInfoByGateid(gateid) {
		return
	}

	connid.Init()
	pkg.Init()

	go refreshBackend()

	if err := pkg.Serve(listenPort); err != nil {
		log.Error("%s", err.Error())
		return
	}
}
