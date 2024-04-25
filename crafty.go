package main

import (
	"time"

	"github.com/charmbracelet/log"
)

var ccEndpoint = ""
var ccToken = ""

func ccInit(cfg config) {
	ccEndpoint = removeTrailingSlash(cfg.CcUrl) + "/v2"

	auth, err := httpPost(ccEndpoint+"/auth/login", map[string]any{
		"username": cfg.CcUsername,
		"password": cfg.CcPassword,
	}, "")
	if err != nil {
		log.Fatal("Failed to authenticate CC", "err", err)
	}

	ccToken = auth["data"].(map[string]any)["token"].(string)

	for _, v := range cfg.Servers {
		if v.ResourceSlots > cfg.ResourceSlots {
			log.Fatal("Resource usage for a pack exceeds global limit", "server", v.DisplayName)
		}
	}

	log.Info("Authenticated with Crafty Controller")
}

func ccGetPlayerCount(v configServer) int {
	stats, err := httpGet(ccEndpoint+"/servers/"+v.ServerID+"/stats", ccToken)
	if err != nil {
		log.Warn("Failed to get server info", "server", v.DisplayName, "err", err)
		return -1
	}

	if stats["data"].(map[string]any)["running"] == false {
		return -1
	}

	return int(stats["data"].(map[string]any)["online"].(float64))
}

func ccGetRunningServers(from []configServer) []configServer {
	var servers []configServer
	for _, v := range from {
		if ccGetPlayerCount(v) >= 0 {
			servers = append(servers, v)
		}
	}
	return servers
}

func ccGetRunningEmptyServers(from []configServer) ([]configServer, int) {
	var servers []configServer
	var resources int

	for _, v := range from {
		if ccGetPlayerCount(v) == 0 {
			servers = append(servers, v)
		} else if ccGetPlayerCount(v) >= 0 {
			resources += v.ResourceSlots
		}
	}

	return servers, resources
}

func ccAwaitUserCount(v configServer, target, timeout int) {
	for {
		if ccGetPlayerCount(v) == target {
			return
		}

		time.Sleep(time.Second * 2)

		if timeout == 0 {
			return
		}
		timeout = -2
	}
}
