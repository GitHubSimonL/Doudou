/*
	This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package network

import (
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
		// log.Erro("[accesscontrol] no whitelist filename")
		return
	}

	fileReader, err := os.Open(filename)
	if err != nil {
		// log.Erro("[accesscontrol] read whitelist failed %v", err)
		return
	}

	defer fileReader.Close()

	content, err := ioutil.ReadAll(fileReader)

	if err != nil {
		// log.Erro("[accesscontrol] read whitelist failed %v", err)
		return
	}

	for _, whitelistIP := range strings.Split(string(content), "\n") {
		this.whiteList = append(this.whiteList, whitelistIP)
	}

	// for _, whitelistIP := range this.whiteList {
	// log.Info("[accesscontrol] loadwhitlist %v", whitelistIP)
	// }

	return true
}

func (this *WhiteList) AccessCheck(ip string) (ok bool) {
	if len(this.whiteList) <= 0 {
		return true
	}

	ipAddr, err := net.ResolveTCPAddr("", ip)

	if err != nil || ipAddr == nil {
		// log.Erro("err %v", err)
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
