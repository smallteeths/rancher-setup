package consts

const (
	ErrorsBasePath  string = "初始化项目根目录失败"
	ErrorInitConfig string = "初始化配置实例失败"
	ErrorInitLogger string = "初始化日志实例失败"
	Scripts                = `#!/bin/bash
# 安装 cert manger
curl -sfL https://charts.rancher.cn/cert-manager/cert-manager-v1.11.3.tgz | tar xvzf - -C ./
sleep 10

helm --kubeconfig /root/.kube/config install cert-manager ./cert-manager \
  --namespace cert-manager \
  --create-namespace \
  --set installCRDs=true \
  --set image.repository=docker.io/cnrancher/mirrored-jetstack-cert-manager-controller \
  --set webhook.image.repository=docker.io/cnrancher/mirrored-jetstack-cert-manager-webhook \
  --set cainjector.image.repository=docker.io/cnrancher/mirrored-jetstack-cert-manager-cainjector \
  --set startupapicheck.image.repository=docker.io/cnrancher/mirrored-jetstack-cert-manager-ctl
sleep 10

# 部署 secret
kubectl --kubeconfig /root/.kube/config create namespace cattle-system

kubectl --kubeconfig /root/.kube/config create -n cattle-system secret docker-registry secret-tiger-docker \
  --docker-username=${HarborUser} \
  --docker-password=${HarborPasswd} \
  --docker-server=docker.io

# 安装 rancher
wget https://charts.rancher.cn/2.7-prime/latest/${Version}.tgz
sleep 10

helm --kubeconfig /root/.kube/config install rancher ./${Version}.tgz \
  --namespace cattle-system \
  --create-namespace \
  --set ingress.ingressClassName=nginx \
  --set hostname=${Host} \
  --set ingress.tls.source=rancher \
  --set global.cattle.psp.enabled=false \
  --set bootstrapPassword=Rancher@123456 \
  --set rancherImage=docker.io/cnrancher/rancher \
  --set imagePullSecrets[0].name=secret-tiger-docker \
  --set systemDefaultRegistry=docker.io
`
)
