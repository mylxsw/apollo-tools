# apollo-tools

Apollo-tools 是一个命令行的小工具，用于获取 Apollo 配置中心的配置，输出到指定文件。

一次性执行，输出应用配置

```bash
apollo-tools -app-id app -server-addr http://127.0.0.1:8088 -format "%s: %s" -output config.yaml
```

持续运行，当配置变更时，更新应用

```bash
apollo-tools -app-id app -server-addr http://127.0.0.1:8088 -format "%s=%s" -forever -on-change 'cat .env'  -output .env
```
