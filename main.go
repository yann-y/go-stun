// Copyright 2016 Cong Ding
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/ccding/go-stun/stun"
	"log"
	"net"
	"os/exec"
	"strings"
	"time"
)

func main() {
	//interface
	var i = flag.String("i", "", "either interface name or address")
	var serverAddr = flag.String("s", stun.DefaultServerAddr, "STUN server address")
	var p = flag.Int("p", stun.DefaultPort, "port to listen on for client")
	var v = flag.Bool("info", false, "verbose mode")
	var vv = flag.Bool("debug", false, "double verbose mode (includes -info)")
	var lp = flag.Int("loop", 0, "loop interval (company: s)")
	var ttl = flag.Bool("ttl", false, "loop interval (company: s)")
	flag.Parse()
	// 参数检验
	serverAddr = parse(*serverAddr)
	fmt.Println("stun server addr :", *serverAddr)
	// Creates a STUN client. NewClientWithConnection can also be used if you want to handle the UDP listener by yourself.
	client := stun.NewClient()
	// The default addr (stun.DefaultServerAddr) will be used unless we call SetServerAddr.
	client.SetServerAddr(*serverAddr)
	// Non verbose mode will be used by default unless we call SetVerbose(true) or SetVVerbose(true).
	client.SetVerbose(*v || *vv)
	client.SetVVerbose(*vv)
	client.SetClientPort(*p)
	if *i != "" {
		ipaddr := net.ParseIP(*i)
		if ipaddr == nil {
			ipaddr = getLocalIpV4(*i)
			fmt.Printf("interface name: %s, addr: %s \n", *i, ipaddr.To4().String())
		}
		client.SetClientAddr(ipaddr)
	}
	if *ttl == true {
		err := client.NatTTL()
		if err != nil {
			fmt.Println(err)
		}
	} else {
		test(client, *lp)
	}
	//externalIP()
	//sig := make(chan os.Signal)
	//signal.Notify(sig)
	//fmt.Println(<-sig)
}
func test(client *stun.Client, loop int) {
	// Discover the NAT and return the result.
	for {
		nat, host, err := client.Discover()
		if err != nil {
			fmt.Printf("Discover err: %v\n", err)
			return
		}
		fmt.Println("NAT Type:", nat)
		if host != nil {
			fmt.Println("External IP Family:", host.Family())
			fmt.Println("External IP:", host.IP())
			fmt.Println("External Port:", host.Port())
		}
		if loop == 0 {
			break
		}
		time.Sleep(time.Duration(loop) * time.Second)
	}
}
func parse(serverAddr string) *string {
	l := len(strings.Split(serverAddr, ":"))
	switch l {
	case 2:
		return &serverAddr
	default:
		s := serverAddr + ":3478"
		return &s
	}
}
func externalIP() {
	cmd := exec.Command("curl", "cip.cc")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout // 标准输出
	cmd.Stderr = &stderr // 标准错误
	err := cmd.Run()
	outStr, _ := string(stdout.Bytes()), string(stderr.Bytes())
	fmt.Printf("ip info:\n%s\nerr:\n", outStr)
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
}
func behaviorTest(c *stun.Client) {
	natBehavior, err := c.BehaviorTest()
	if err != nil {
		fmt.Println(err)
	}

	if natBehavior != nil {
		fmt.Println("  Mapping Behavior:", natBehavior.MappingType)
		fmt.Println("Filtering Behavior:", natBehavior.FilteringType)
		fmt.Println("   Normal NAT Type:", natBehavior.NormalType())
	}
}
func getLocalIpV4(name string) net.IP {
	inters, err := net.Interfaces()
	if err != nil {
		panic(err)
	}
	for _, inter := range inters {
		if inter.Name == name {
			addrs, err := inter.Addrs()
			if err != nil {
				break
			}
			for _, addr := range addrs {
				if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					//判断是否存在IPV4 IP 如果没有过滤
					if ipnet.IP.To4() != nil {
						return ipnet.IP
					}
				}
			}
			break
		}
	}
	return nil
}
