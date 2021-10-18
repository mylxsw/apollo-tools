package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/mylxsw/asteria/log"
	"github.com/philchia/agollo/v4"
)

func main() {

	var appID, serverAddr, accessSecret, cluster, namespace, output, format, onChange string
	var forever bool

	flag.StringVar(&appID, "app-id", "", "app id")
	flag.StringVar(&serverAddr, "server-addr", "http://apollo-config:8080", "apollo server addr")
	flag.StringVar(&accessSecret, "secret", "", "apollo access secret")
	flag.StringVar(&cluster, "cluster", "default", "apollo cluster")
	flag.StringVar(&namespace, "namespace", "application", "apollo namespace")
	flag.StringVar(&output, "output", "stdout", "输出文件路径")
	flag.StringVar(&format, "format", "%s=%s", "输出格式: %s: %s")
	flag.BoolVar(&forever, "forever", false, "持续运行，监听 key 的变化，实时更新文件")
	flag.StringVar(&onChange, "on-change", "", "配置变更时执行的命令")

	flag.Parse()

	client := agollo.NewClient(&agollo.Conf{
		AppID:           appID,
		MetaAddr:        serverAddr,
		AccesskeySecret: accessSecret,
		Cluster:         cluster,
		NameSpaceNames:  []string{namespace},
		CacheDir:        os.TempDir(),
	})

	if err := client.Start(); err != nil {
		panic(err)
	}

	var updateFunc = func() {
		lines := make([]string, 0)
		for _, k := range client.GetAllKeys() {
			lines = append(lines, fmt.Sprintf(format, k, client.GetString(k)))
		}

		if output == "stdout" {
			fmt.Println(strings.Join(lines, "\n"))
		} else {
			if err := ioutil.WriteFile(output, []byte(strings.Join(lines, "\n")), os.ModePerm); err != nil {
				panic(err)
			}
		}
	}

	updateFunc()

	if forever {
		client.OnUpdate(func(ce *agollo.ChangeEvent) {
			log.With(ce).Debug("change event received")
			updateFunc()
			if onChange != "" {
				cmd := exec.Command("/bin/sh", "-c", onChange)
				cmd.Env = os.Environ()
				if data, err := cmd.CombinedOutput(); err != nil {
					log.Errorf("on-chagne command failed: %v", err)
				} else {
					log.Infof("on-change data: %s", string(data))
				}

			}
		})

		stop := make(chan interface{})
		<-stop
	}
}
