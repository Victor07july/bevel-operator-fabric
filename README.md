      ---
id: getting-started
title: Getting started
---

# Hyperledger Fabric Operator

## Recursos

- [x] Criação de certificates authorities (CA)
- [x] Criação de peers
- [x] Criação de ordering services
- [x] Criação de recursos sem modificação manual do material criptográfico
- [x] Roteamento de domínio com SNI usando Istio 
- [x] Execução de chaincode como chaincode externo via Kubernetes
- [x] Suporte a Hyperledger Fabric 2.3+
- [x] Gerenciamento de genesis para Ordering Services
- [x] E2E testing including the execution of chaincodes in KIND
- [x] Renovação de certificados

# Tutorial

Resources:
- [Hyperledger Fabric build ARM](https://www.polarsparc.com/xhtml/Hyperledger-ARM-Build.html)

## Criar Cluster Kubernetes

Para começar o deploy da rede Fabric é necessário criar um cluster Kubernetes. Será utilizado aqui o KinD.

Certifique-se de ter as seguintes portas disponíveis antes de começar:
- 80
- 443

```bash
cat << EOF > resources/kind-config.yaml
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
  image: kindest/node:v1.25.8
  extraPortMappings:
  - containerPort: 30949
    hostPort: 80
  - containerPort: 30950
    hostPort: 443
EOF

kind create cluster --config=./resources/kind-config.yaml

```

## Instalar o operador Kubernetes

Nesta etapa instalaremos o operador Kubernetes para o Fabric. Isso irá instalar:

- CRD (Custom Resource Definitions) to deploy Certification Fabric Peers, Orderers and Authorities
- Deploy the program to deploy the nodes in Kubernetes

Instale o helm: [https://helm.sh/docs/intro/install/](https://helm.sh/docs/intro/install/)

```bash
helm repo add kfs https://kfsoftware.github.io/hlf-helm-charts --force-update

helm install hlf-operator --version=1.9.0 -- kfs/hlf-operator
```


### Instalar o plugin Kubectl

Antes de instalar o plugin Kubectl, instale antes o Krew:
[https://krew.sigs.k8s.io/docs/user-guide/setup/install/](https://krew.sigs.k8s.io/docs/user-guide/setup/install/)

A seguir, instale o Kubectl com o seguinte comando:

```bash
kubectl krew install hlf
```

### Instale o Istio

Install Istio binaries on the machine:
```bash
curl -L https://istio.io/downloadIstio | sh -
```

Install Istio on the Kubernetes cluster:

```bash

kubectl create namespace istio-system

istioctl operator init

kubectl apply -f - <<EOF
apiVersion: install.istio.io/v1alpha1
kind: IstioOperator
metadata:
  name: istio-gateway
  namespace: istio-system
spec:
  addonComponents:
    grafana:
      enabled: false
    kiali:
      enabled: false
    prometheus:
      enabled: false
    tracing:
      enabled: false
  components:
    ingressGateways:
      - enabled: true
        k8s:
          hpaSpec:
            minReplicas: 1
          resources:
            limits:
              cpu: 500m
              memory: 512Mi
            requests:
              cpu: 100m
              memory: 128Mi
          service:
            ports:
              - name: http
                port: 80
                targetPort: 8080
                nodePort: 30949
              - name: https
                port: 443
                targetPort: 8443
                nodePort: 30950
            type: NodePort
        name: istio-ingressgateway
    pilot:
      enabled: true
      k8s:
        hpaSpec:
          minReplicas: 1
        resources:
          limits:
            cpu: 300m
            memory: 512Mi
          requests:
            cpu: 100m
            memory: 128Mi
  meshConfig:
    accessLogFile: /dev/stdout
    enableTracing: false
    outboundTrafficPolicy:
      mode: ALLOW_ANY
  profile: default

EOF

```

### Environment Variables for AMD (Default)

```bash
export PEER_IMAGE=hyperledger/fabric-peer
export PEER_VERSION=2.5.0

export ORDERER_IMAGE=hyperledger/fabric-orderer
export ORDERER_VERSION=2.5.0

export CA_IMAGE=hyperledger/fabric-ca
export CA_VERSION=1.5.6
```


### Environment Variables for ARM (Mac M1)

```bash
export PEER_IMAGE=hyperledger/fabric-peer
export PEER_VERSION=2.5.0

export ORDERER_IMAGE=hyperledger/fabric-orderer
export ORDERER_VERSION=2.5.0

export CA_IMAGE=hyperledger/fabric-ca             
export CA_VERSION=1.5.6

```


### Configurar DNS Interno

```bash
CLUSTER_IP=$(kubectl -n istio-system get svc istio-ingressgateway -o json | jq -r .spec.clusterIP)
kubectl apply -f - <<EOF
kind: ConfigMap
apiVersion: v1
metadata:
  name: coredns
  namespace: kube-system
data:
  Corefile: |
    .:53 {
        errors
        health {
           lameduck 5s
        }
        rewrite name regex (.*)\.localho\.st host.ingress.internal
        hosts {
          ${CLUSTER_IP} host.ingress.internal
          fallthrough
        }
        ready
        kubernetes cluster.local in-addr.arpa ip6.arpa {
           pods insecure
           fallthrough in-addr.arpa ip6.arpa
           ttl 30
        }
        prometheus :9153
        forward . /etc/resolv.conf {
           max_concurrent 1000
        }
        cache 30
        loop
        reload
        loadbalance
    }
EOF
```

### Criação do CA para Org1


```bash

kubectl hlf ca create  --image=$CA_IMAGE --version=$CA_VERSION --storage-class=standard --capacity=1Gi --name=org1-ca \
    --enroll-id=enroll --enroll-pw=enrollpw --hosts=org1-ca.localho.st --istio-port=443

kubectl wait --timeout=180s --for=condition=Running fabriccas.hlf.kungfusoftware.es --all
```

Verifique se o CA foi implementado e funciona:

```bash
curl -k https://org1-ca.localho.st:443/cainfo
```

Registre um usuário peer no CA da Organização 1 (Org1MSP)

```bash
# register user in CA for peers
kubectl hlf ca register --name=org1-ca --user=peer --secret=peerpw --type=peer \
 --enroll-id enroll --enroll-secret=enrollpw --mspid Org1MSP

```

### Deploy de peers para a Org1MSP

```bash
kubectl hlf peer create --statedb=couchdb --image=$PEER_IMAGE --version=$PEER_VERSION --storage-class=standard --enroll-id=peer --mspid=Org1MSP \
        --enroll-pw=peerpw --capacity=5Gi --name=org1-peer0 --ca-name=org1-ca.default \
        --hosts=peer0-org1.localho.st --istio-port=443

kubectl wait --timeout=180s --for=condition=Running fabricpeers.hlf.kungfusoftware.es --all
```

(ALTERNATIVA) O peer acima não é capaz de instalar chaincode local, apenas chaincodes CCAS. 
Para criar peers que instalam chaincode local, será neceessário criar um peer com o atributo kubernetes chaincode builder (k8s builder) com o comando abaixo.
Lembre de escolher apenas um dos peers apresentados para este tutorial.
```bash

export PEER_IMAGE=quay.io/kfsoftware/fabric-peer
export PEER_VERSION=2.4.1-v0.0.3
export MSP_ORG=Org1MSP
export PEER_SECRET=peerpw

kubectl hlf peer create --statedb=couchdb --image=$PEER_IMAGE --version=$PEER_VERSION --storage-class=standard --enroll-id=peer --mspid=$MSP_ORG \
--enroll-pw=$PEER_SECRET --capacity=5Gi --name=org1-peer0 --ca-name=org1-ca.default --k8s-builder=true --hosts=peer0-org1.localho.st --istio-port=443

kubectl wait --timeout=180s --for=condition=Running fabricpeers.hlf.kungfusoftware.es --all

# leva alguns minutos

```

Verifique se os peers foram implementados e funcionam:

```bash
openssl s_client -connect peer0-org1.localho.st:443

```

### Criação do CA para Org2

```bash
kubectl hlf ca create  --image=$CA_IMAGE --version=$CA_VERSION --storage-class=standard --capacity=1Gi --name=org2-ca \
    --enroll-id=enroll --enroll-pw=enrollpw --hosts=org2-ca.localho.st --istio-port=443

kubectl wait --timeout=180s --for=condition=Running fabriccas.hlf.kungfusoftware.es --all
```

Verifique se o CA está funcionando

```bash
curl -k https://org2-ca.localho.st:443/cainfo
```

Registre um usuário no CA da Organização 2 (Org2MSP) para os peers

```bash
# register user in CA for peers
kubectl hlf ca register --name=org2-ca --user=peer --secret=peerpw --type=peer \
 --enroll-id enroll --enroll-secret=enrollpw --mspid Org2MSP
```

### Deploy de peers para Org2 (escolha um apenas)

```bash
kubectl hlf peer create --statedb=couchdb --image=$PEER_IMAGE --version=$PEER_VERSION --storage-class=standard --enroll-id=peer --mspid=Org2MSP \
        --enroll-pw=peerpw --capacity=5Gi --name=org2-peer0 --ca-name=org2-ca.default \
        --hosts=peer0-org2.localho.st --istio-port=443

kubectl wait --timeout=180s --for=condition=Running fabricpeers.hlf.kungfusoftware.es --all
```

Alternativa com k8s-builder

```bash

export PEER_IMAGE=quay.io/kfsoftware/fabric-peer
export PEER_VERSION=2.4.1-v0.0.3
export MSP_ORG=Org2MSP
export PEER_SECRET=peerpw

kubectl hlf peer create --statedb=couchdb --image=$PEER_IMAGE --version=$PEER_VERSION --storage-class=standard --enroll-id=peer --mspid=$MSP_ORG \
--enroll-pw=$PEER_SECRET --capacity=5Gi --name=org2-peer0 --ca-name=org2-ca.default --k8s-builder=true --hosts=peer0-org2.localho.st --istio-port=443

kubectl wait --timeout=180s --for=condition=Running fabricpeers.hlf.kungfusoftware.es --all

# leva alguns minutos

```

Verifique se o peer funciona

```
openssl s_client -connect peer0-org2.localho.st:443
```

### Deploy de uma organização `Orderer`

para fazer o deploy de uma organização orderer, temos que:

1. Criar um certification authority (CA)
2. Registrar usuário `orderer` com senha `ordererpw`
3. Criar orderer

### Criar o CA

```bash

kubectl hlf ca create  --image=$CA_IMAGE --version=$CA_VERSION --storage-class=standard --capacity=1Gi --name=ord-ca \
    --enroll-id=enroll --enroll-pw=enrollpw --hosts=ord-ca.localho.st --istio-port=443

kubectl wait --timeout=180s --for=condition=Running fabriccas.hlf.kungfusoftware.es --all

```

Verifique se a certificação foi implementada e funciona:

```bash
curl -vik https://ord-ca.localho.st:443/cainfo
```

### Registre o usuário `orderer`

```bash
kubectl hlf ca register --name=ord-ca --user=orderer --secret=ordererpw \
    --type=orderer --enroll-id enroll --enroll-secret=enrollpw --mspid=OrdererMSP --ca-url="https://ord-ca.localho.st:443"

```
### Deploy de três orderers

```bash
kubectl hlf ordnode create --image=$ORDERER_IMAGE --version=$ORDERER_VERSION \
    --storage-class=standard --enroll-id=orderer --mspid=OrdererMSP \
    --enroll-pw=ordererpw --capacity=2Gi --name=ord-node0 --ca-name=ord-ca.default \
    --hosts=orderer0-ord.localho.st --istio-port=443

kubectl hlf ordnode create --image=$ORDERER_IMAGE --version=$ORDERER_VERSION \
    --storage-class=standard --enroll-id=orderer --mspid=OrdererMSP \
    --enroll-pw=ordererpw --capacity=2Gi --name=ord-node1 --ca-name=ord-ca.default \
    --hosts=orderer1-ord.localho.st --istio-port=443

kubectl hlf ordnode create --image=$ORDERER_IMAGE --version=$ORDERER_VERSION \
    --storage-class=standard --enroll-id=orderer --mspid=OrdererMSP \
    --enroll-pw=ordererpw --capacity=2Gi --name=ord-node2 --ca-name=ord-ca.default \
    --hosts=orderer2-ord.localho.st --istio-port=443


kubectl wait --timeout=180s --for=condition=Running fabricorderernodes.hlf.kungfusoftware.es --all
```

Verifique se os orderers funcionam:

```bash
kubectl get pods
```

```bash
openssl s_client -connect orderer0-ord.localho.st:443
openssl s_client -connect orderer1-ord.localho.st:443
openssl s_client -connect orderer2-ord.localho.st:443
```


## Create channel

Para criar o canal nós precisamos criar o "wallet secret", que irá conter as identidades usadas pelo bevel operator para gerenciar o canal

### Registrar e matricular identidade OrdererMSP

```bash
# register
kubectl hlf ca register --name=ord-ca --user=admin --secret=adminpw \
    --type=admin --enroll-id enroll --enroll-secret=enrollpw --mspid=OrdererMSP

# enroll
kubectl hlf ca enroll --name=ord-ca --namespace=default \
    --user=admin --secret=adminpw --mspid OrdererMSP \
    --ca-name tlsca  --output resources/orderermsp.yaml
```


### Registrar e matricular identidade Org1MSP

```bash
# register
kubectl hlf ca register --name=org1-ca --namespace=default --user=admin --secret=adminpw \
    --type=admin --enroll-id enroll --enroll-secret=enrollpw --mspid=Org1MSP

# enroll
kubectl hlf ca enroll --name=org1-ca --namespace=default \
    --user=admin --secret=adminpw --mspid Org1MSP \
    --ca-name ca  --output resources/org1msp.yaml
```

### Registrar e matricular identidade Org2MSP

```bash
# register
kubectl hlf ca register --name=org2-ca --namespace=default --user=admin --secret=adminpw \
    --type=admin --enroll-id enroll --enroll-secret=enrollpw --mspid=Org2MSP

# enroll
kubectl hlf ca enroll --name=org2-ca --namespace=default \
    --user=admin --secret=adminpw --mspid Org2MSP \
    --ca-name ca  --output resources/org2msp.yaml

```

### Criar o segredo

```bash


kubectl create secret generic wallet --namespace=default \
        --from-file=org1msp.yaml=$PWD/resources/org1msp.yaml \
        --from-file=org2msp.yaml=$PWD/resources/org2msp.yaml \
        --from-file=orderermsp.yaml=$PWD/resources/orderermsp.yaml
```

### Create main channel

```bash
export PEER_ORG_SIGN_CERT=$(kubectl get fabriccas org1-ca -o=jsonpath='{.status.ca_cert}')
export PEER_ORG_TLS_CERT=$(kubectl get fabriccas org1-ca -o=jsonpath='{.status.tlsca_cert}')
export IDENT_8=$(printf "%8s" "")
export ORDERER_TLS_CERT=$(kubectl get fabriccas ord-ca -o=jsonpath='{.status.tlsca_cert}' | sed -e "s/^/${IDENT_8}/" )
export ORDERER0_TLS_CERT=$(kubectl get fabricorderernodes ord-node0 -o=jsonpath='{.status.tlsCert}' | sed -e "s/^/${IDENT_8}/" )
export ORDERER1_TLS_CERT=$(kubectl get fabricorderernodes ord-node1 -o=jsonpath='{.status.tlsCert}' | sed -e "s/^/${IDENT_8}/" )
export ORDERER2_TLS_CERT=$(kubectl get fabricorderernodes ord-node2 -o=jsonpath='{.status.tlsCert}' | sed -e "s/^/${IDENT_8}/" )

kubectl apply -f - <<EOF
apiVersion: hlf.kungfusoftware.es/v1alpha1
kind: FabricMainChannel
metadata:
  name: demo
spec:
  name: demo
  adminOrdererOrganizations:
    - mspID: OrdererMSP
  adminPeerOrganizations:
    - mspID: Org1MSP
  channelConfig:
    application:
      acls: null
      capabilities:
        - V2_0
      policies: null
    capabilities:
      - V2_0
    orderer:
      batchSize:
        absoluteMaxBytes: 1048576
        maxMessageCount: 120
        preferredMaxBytes: 524288
      batchTimeout: 2s
      capabilities:
        - V2_0
      etcdRaft:
        options:
          electionTick: 10
          heartbeatTick: 1
          maxInflightBlocks: 5
          snapshotIntervalSize: 16777216
          tickInterval: 500ms
      ordererType: etcdraft
      policies: null
      state: STATE_NORMAL
    policies: null
  externalOrdererOrganizations: []
  peerOrganizations:
    - mspID: Org1MSP
      caName: "org1-ca"
      caNamespace: "default"
    - mspID: Org2MSP
      caName: "org2-ca"
      caNamespace: "default" 
  identities:
    OrdererMSP:
      secretKey: orderermsp.yaml
      secretName: wallet
      secretNamespace: default
    Org1MSP:
      secretKey: org1msp.yaml
      secretName: wallet
      secretNamespace: default
  externalPeerOrganizations: []
  ordererOrganizations:
    - caName: "ord-ca"
      caNamespace: "default"
      externalOrderersToJoin:
        - host: ord-node0
          port: 7053
        - host: ord-node1
          port: 7053
        - host: ord-node2
          port: 7053
      mspID: OrdererMSP
      ordererEndpoints:
        - ord-node0:7050
        - ord-node1:7050
        - ord-node2:7050
      orderersToJoin: []
  orderers:
    - host: ord-node0
      port: 7050
      tlsCert: |-
${ORDERER0_TLS_CERT}
    - host: ord-node1
      port: 7050
      tlsCert: |-
${ORDERER1_TLS_CERT}
    - host: ord-node2
      port: 7050
      tlsCert: |-
${ORDERER2_TLS_CERT}

EOF

```

## Inserir peers de Org1 no canal

```bash

export IDENT_8=$(printf "%8s" "")
export ORDERER0_TLS_CERT=$(kubectl get fabricorderernodes ord-node0 -o=jsonpath='{.status.tlsCert}' | sed -e "s/^/${IDENT_8}/" )

kubectl apply -f - <<EOF
apiVersion: hlf.kungfusoftware.es/v1alpha1
kind: FabricFollowerChannel
metadata:
  name: demo-org1msp
spec:
  anchorPeers:
    - host: org1-peer0.default
      port: 7051
  hlfIdentity:
    secretKey: org1msp.yaml
    secretName: wallet
    secretNamespace: default
  mspId: Org1MSP
  name: demo
  externalPeersToJoin: []
  orderers:
    - certificate: |
${ORDERER0_TLS_CERT}
      url: grpcs://ord-node0.default:7050
  peersToJoin:
    - name: org1-peer0
      namespace: default
EOF


```

## Inserir peers de Org2 no canal

```bash

export IDENT_8=$(printf "%8s" "")
export ORDERER0_TLS_CERT=$(kubectl get fabricorderernodes ord-node0 -o=jsonpath='{.status.tlsCert}' | sed -e "s/^/${IDENT_8}/" )

kubectl apply -f - <<EOF
apiVersion: hlf.kungfusoftware.es/v1alpha1
kind: FabricFollowerChannel
metadata:
  name: demo-org2msp
spec:
  anchorPeers:
    - host: org2-peer0.default
      port: 7051
  hlfIdentity:
    secretKey: org2msp.yaml
    secretName: wallet
    secretNamespace: default
  mspId: Org2MSP
  name: demo
  externalPeersToJoin: []
  orderers:
    - certificate: |
${ORDERER0_TLS_CERT}
      url: grpcs://orderer0-ord.localho.st:443
  peersToJoin:
    - name: org2-peer0
      namespace: default
EOF

```

### Instalação de chaincode as a service (ccas): Em breve

-----------
-----------

## Instalação de chaincode Local

### Preparar string / arquivo de conexão para um peer

Para preparar a string de conexão, precisamos:

1. Obter a string de conexão sem usuários para a organização Org1MSP e OrdererMSP
2. Registre um usuário no CA para assinatura (registro)
3. Obter os certificados usando o usuário criado no passo 2 (enroll)
4. Anexar o usuário à string de conexão

(Repetir 2, 3 e 4 para Org2)

--------------

1. Obter a string de conexão sem usuários para a organização Org1MSP e OrdererMSP

```bash
kubectl hlf inspect -c=demo --output resources/network.yaml -o Org1MSP -o Org2MSP -o OrdererMSP
```

### Registrar usuário para Org1

2. Registre um usuário no CA para assinatura (registro)
```bash
kubectl hlf ca register --name=org1-ca --user=admin --secret=adminpw --type=admin \
 --enroll-id enroll --enroll-secret=enrollpw --mspid Org1MSP  
```

3. Obter os certificados usando o usuário criado no passo 2 (enroll)
```bash
kubectl hlf ca enroll --name=org1-ca --user=admin --secret=adminpw --mspid Org1MSP \
        --ca-name ca  --output resources/peer-org1.yaml
```

4. Anexar o usuário à string de conexão
```bash
kubectl hlf utils adduser --userPath=resources/peer-org1.yaml --config=resources/network.yaml --username=admin --mspid=Org1MSP
```

### Registrar usuário para Org2 (repetir passos acima para Org2)

```bash
kubectl hlf ca register --name=org2-ca --user=admin --secret=adminpw --type=admin \
 --enroll-id enroll --enroll-secret=enrollpw --mspid Org2MSP

kubectl hlf ca enroll --name=org2-ca --user=admin --secret=adminpw --mspid Org2MSP \
        --ca-name ca  --output resources/peer-org2.yaml

kubectl hlf utils adduser --userPath=resources/peer-org2.yaml --config=resources/network.yaml --username=admin --mspid=Org2MSP
```

### Instalação do chaincode
Com o arquivo de conexão preparado, vamos instalar o chaincode no peer que possua o atributo k8s-builder, como explicado no passo de deploy de peers

```bash
kubectl hlf chaincode install --path=./fixtures/chaincodes/fabcar/go \
    --config=resources/network.yaml --language=golang --label=fabcar --user=admin --peer=org1-peer0.default

kubectl hlf chaincode install --path=./fixtures/chaincodes/fabcar/go \
    --config=resources/network.yaml --language=golang --label=fabcar --user=admin --peer=org2-peer0.default

# this can take 3-4 minutes
```

Verificação de chaincodes instalados

```bash
kubectl hlf chaincode queryinstalled --config=resources/network.yaml --user=admin --peer=org1-peer0.default

kubectl hlf chaincode queryinstalled --config=resources/network.yaml --user=admin --peer=org2-peer0.default
```

Aprovar chaincode

```bash
#Organização 1

PACKAGE_ID=fabcar:0c616be7eebace4b3c2aa0890944875f695653dbf80bef7d95f3eed6667b5f40 # replace it with the package id of your chaincode
kubectl hlf chaincode approveformyorg --config=resources/network.yaml --user=admin --peer=org1-peer0.default \
    --package-id=$PACKAGE_ID \
    --version "1.0" --sequence 1 --name=fabcar \
    --policy="OR('Org1MSP.member')" --channel=demo

# Organização 2

PACKAGE_ID=fabcar:0c616be7eebace4b3c2aa0890944875f695653dbf80bef7d95f3eed6667b5f40 # replace it with the package id of your chaincode
kubectl hlf chaincode approveformyorg --config=resources/network.yaml --user=admin --peer=org2-peer0.default \
    --package-id=$PACKAGE_ID \
    --version "1.0" --sequence 1 --name=fabcar \
    --policy="AND('Org1MSP.member', 'Org2MSP.member')" --channel=demo
```

Fazer o commit do chaincode

```bash
kubectl hlf chaincode commit --config=resources/network.yaml --mspid=Org1MSP --user=admin \
    --version "1.0" --sequence 1 --name=fabcar \
    --policy="OR('Org1MSP.member')" --channel=demo
```

Testar chaincode

```bash
kubectl hlf chaincode invoke --config=resources/network.yaml \
    --user=admin --peer=org1-peer0.default \
    --chaincode=fabcar --channel=demo \
    --fcn=initLedger -a '[]'
```

Fazer query de todos os carros / assets

```bash
kubectl hlf chaincode query --config=resources/network.yaml \
    --user=admin --peer=org1-peer0.default \
    --chaincode=fabcar --channel=demo \
    --fcn=QueryAllCars -a '[]'
```

## Usando clientes:
EM BREVE SE DEUS QUISER


## Finalizando
A essa altura, você deve ter:


- Um serviço de ordenação com 3 orderers e CA
- Organização com 3 peers e CA
- Um canal chamado "demo"
- Um chaincode instalado no "peer2" da organização aprovado e commitado


## Derrubando o ambiente

```bash
kubectl delete fabricorderernodes.hlf.kungfusoftware.es --all-namespaces --all
kubectl delete fabricpeers.hlf.kungfusoftware.es --all-namespaces --all
kubectl delete fabriccas.hlf.kungfusoftware.es --all-namespaces --all
kubectl delete fabricchaincode.hlf.kungfusoftware.es --all-namespaces --all
kubectl delete fabricmainchannels --all-namespaces --all
kubectl delete fabricfollowerchannels --all-namespaces --all

kind delete cluster
```

## Troubleshooting

### Chaincode installation/build error

Chaincode installation/build can fail due to unsupported local kubertenes version such as [minikube](https://github.com/kubernetes/minikube).

```shell
$ kubectl hlf chaincode install --path=./fixtures/chaincodes/fabcar/go \
        --config=org1.yaml --language=golang --label=fabcar --user=admin --peer=org1-peer0.default

Error: Transaction processing for endorser [192.168.49.2:31278]: Chaincode status Code: (500) UNKNOWN.
Description: failed to invoke backing implementation of 'InstallChaincode': could not build chaincode:
external builder failed: external builder failed to build: external builder 'my-golang-builder' failed:
exit status 1
```

If your purpose is to test the hlf-operator please consider to switch to [kind](https://github.com/kubernetes-sigs/kind) that is tested and supported.
