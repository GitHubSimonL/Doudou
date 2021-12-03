package itr

import (
	"Doudou/lib/logger"
	"io/ioutil"
	"net"
	"os"
	"strings"
)

type WhiteList struct {
	whiteList []string
}

func (this *WhiteList) LoadWhiteList(filename string) (ok bool) {
	if len(filename) <= 0 {
		logger.LogErrf("[accesscontrol] no whitelist filename")
		return
	}

	fileReader, err := os.Open(filename)
	if err != nil {
		logger.LogErrf("[accesscontrol] read whitelist failed %v", err)
		return
	}

	defer fileReader.Close()

	content, err := ioutil.ReadAll(fileReader)

	if err != nil {
		logger.LogErrf("[accesscontrol] read whitelist failed %v", err)
		return
	}

	for _, whitelistIP := range strings.Split(string(content), "\n") {
		this.whiteList = append(this.whiteList, whitelistIP)
	}

	for _, whitelistIP := range this.whiteList {
		logger.LogInfof("[accesscontrol] loadwhitlist %v", whitelistIP)
	}

	return true
}

func (this *WhiteList) AccessCheck(ip string) (ok bool) {
	if len(this.whiteList) <= 0 {
		return true
	}

	ipAddr, err := net.ResolveTCPAddr("", ip)

	if err != nil || ipAddr == nil {
		logger.LogErrf("err %v", err)
		return
	}

	ipv4, ipv6 := ipAddr.IP.To4(), ipAddr.IP.To16()

	switch {
	case ipv4 != nil:
		hostFields := strings.Split(ipv4.String(), ".")

		if len(hostFields) != 4 {
			return
		}

		for _, whitelistIP := range this.whiteList {
			marchFields := strings.Split(whitelistIP, ".")

			if len(marchFields) != 4 {
				continue
			}

			if IPFieldMarch(marchFields, hostFields) {
				return true
			}
		}

		return
	case ipv6 != nil:
		hostFields := IPv6Full(strings.Split(ipv6.String(), ":"))

		if len(hostFields) != 8 {
			return
		}

		for _, whitelistIP := range this.whiteList {
			marchFields := IPv6Full(strings.Split(whitelistIP, ":"))

			if len(marchFields) != 8 {
				continue
			}

			if IPFieldMarch(marchFields, hostFields) {
				return true
			}
		}

		return
	}

	return
}
