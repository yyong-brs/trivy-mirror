# trivy-mirror
将 trivy db 漏洞库做成镜像推送到 docker hub，方便 CICD 使用等场景.

## 原理
1. `main.go` 作为主程序主要实现将 `trivy` 漏洞库缓存文件下载下来，主要是 trivy.db 以及 metadata.json 文件
2. `Dockerfile` 用于基于以上结果，构建镜像
3. 由于网络限制，通过 配置 github action 实现自动下载缓存文件构建镜像并push 到 docker hub
4. 镜像推送地址： https://hub.docker.com/repository/docker/yueyongdada/trivy_ci