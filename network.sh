#!/bin/bash

function chaincode_install() {
  CHAINCODE_LABEL = $1

  echo "Preparando conexão com a organização INMETRO"

  kubectl hlf inspect -c=demo --output resources/network.yaml -o INMETROMSP -o PUCMSP -o OrdererMSP
  kubectl hlf ca register --name=inmetro-ca --user=admin --secret=adminpw --type=admin \
    --enroll-id enroll --enroll-secret=enrollpw --mspid INMETROMSP
  kubectl hlf ca enroll --name=inmetro-ca --user=admin --secret=adminpw --mspid INMETROMSP \
          --ca-name ca  --output resources/peer-inmetro.yaml
  kubectl hlf utils adduser --userPath=resources/peer-inmetro.yaml --config=resources/network.yaml --username=admin --mspid=INMETROMSP

  echo "Preparando conexão com a organização PUC"

  kubectl hlf ca register --name=puc-ca --user=admin --secret=adminpw --type=admin \
    --enroll-id enroll --enroll-secret=enrollpw --mspid PUCMSP
  kubectl hlf ca enroll --name=puc-ca --user=admin --secret=adminpw --mspid PUCMSP \
    --ca-name ca  --output resources/peer-puc.yaml

  kubectl hlf utils adduser --userPath=resources/peer-puc.yaml --config=resources/network.yaml --username=admin --mspid=PUCMSP

  export PACKAGE_ID=$(kubectl hlf chaincode calculatepackageid --path=./examples/chaincodes/$CHAINCODE_LABEL --language=go --label=$CHAINCODE_LABEL)
  echo "PACKAGE_ID=$PACKAGE_ID"

  echo "Esse processo pode levar alguns minutos"
  echo "Instalando chaincode na organização INMETRO"

  kubectl hlf chaincode install --path=./fixtures/chaincodes/$CHAINCODE_LABEL \
    --config=resources/network.yaml --language=golang --label=$CHAINCODE_LABEL --user=admin --peer=inmetro-peer0.default

  echo "Instalando chaincode na organização PUC"

  kubectl hlf chaincode install --path=./fixtures/chaincodes/$CHAINCODE_LABEL \
      --config=resources/network.yaml --language=golang --label=$CHAINCODE_LABEL --user=admin --peer=puc-peer0.default
  # this can take 3-4 minutes

  # esse processo de aguardar evita erros com o approve
  echo "Aguardando 1 minuto para o pod dos chaincodes serem carregado"
  sleep 60

  echo "Aprovando chaincode em ambas as organizações"
  
  #Organização INMETRO
  kubectl hlf chaincode approveformyorg --config=resources/network.yaml --user=admin --peer=inmetro-peer0.default \
      --package-id=$PACKAGE_ID \
      --version "1.0" --sequence 1 --name=$CHAINCODE_LABEL \
      --policy="AND('INMETROMSP.member', 'PUCMSP.member')" --channel=demo

  # Organização PUC
  kubectl hlf chaincode approveformyorg --config=resources/network.yaml --user=admin --peer=puc-peer0.default \
      --package-id=$PACKAGE_ID \
      --version "1.0" --sequence 1 --name=$CHAINCODE_LABEL \
      --policy="AND('INMETROMSP.member', 'PUCMSP.member')" --channel=demo

  echo "Aguarde 10 segundos..."
  sleep 10
  echo "Realizando commit do chaincode"
  
  kubectl hlf chaincode commit --config=resources/network.yaml --mspid=INMETROMSP --user=admin \
    --version "1.0" --sequence 1 --name=$CHAINCODE_LABEL \
    --policy="AND('INMETROMSP.member', 'PUCMSP.member')" --channel=demo

  # evita erro com o invoke
  echo "Aguarde 1 minuto para o pod do chaincode ser carregado"
  echo "Fim"
}

function up() {
  export PEER_IMAGE=hyperledger/fabric-peer
  export PEER_VERSION=2.5.0

  export ORDERER_IMAGE=hyperledger/fabric-orderer
  export ORDERER_VERSION=2.5.0

  export CA_IMAGE=hyperledger/fabric-ca
  export CA_VERSION=1.5.6

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

  echo "Aguardando por 5 segundos..."
  sleep 5
  echo "Prosseguindo."

  helm repo add kfs https://kfsoftware.github.io/hlf-helm-charts --force-update
  helm install hlf-operator --version=1.9.0 -- kfs/hlf-operator

  kubectl krew install hlf


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


  echo "Aguardando 1 minuto para o levantamento do cluster Kubernetes"
  sleep 60
  echo "Prosseguindo."

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

  echo "Criando ca para INMETRO"


  kubectl hlf ca create  --image=$CA_IMAGE --version=$CA_VERSION --storage-class=standard --capacity=1Gi --name=inmetro-ca \
      --enroll-id=enroll --enroll-pw=enrollpw --hosts=inmetro-ca.localho.st --istio-port=443

  kubectl wait --timeout=180s --for=condition=Running fabriccas.hlf.kungfusoftware.es --all

  echo "Registrando usuario no CA para os peers"

  # register user in CA for peers
  kubectl hlf ca register --name=inmetro-ca --user=peer --secret=peerpw --type=peer \
  --enroll-id enroll --enroll-secret=enrollpw --mspid INMETROMSP

  echo "Realizando deploy de 1 peer para a organizacao INMETRO"


  export PEER_IMAGE=quay.io/kfsoftware/fabric-peer
  export PEER_VERSION=2.4.1-v0.0.3
  export MSP_ORG=INMETROMSP
  export PEER_SECRET=peerpw

  kubectl hlf peer create --statedb=couchdb --image=$PEER_IMAGE --version=$PEER_VERSION --storage-class=standard --enroll-id=peer --mspid=$MSP_ORG \
  --enroll-pw=$PEER_SECRET --capacity=5Gi --name=inmetro-peer0 --ca-name=inmetro-ca.default --k8s-builder=true --hosts=peer0-inmetro.localho.st --istio-port=443

  kubectl wait --timeout=180s --for=condition=Running fabricpeers.hlf.kungfusoftware.es --all

  # leva alguns minutos

  echo "Criando CA para organizacao PUC"

  kubectl hlf ca create  --image=$CA_IMAGE --version=$CA_VERSION --storage-class=standard --capacity=1Gi --name=puc-ca \
      --enroll-id=enroll --enroll-pw=enrollpw --hosts=puc-ca.localho.st --istio-port=443

  kubectl wait --timeout=180s --for=condition=Running fabriccas.hlf.kungfusoftware.es --all

  echo "Registrando usuario para os peers"

  # register user in CA for peers
  kubectl hlf ca register --name=puc-ca --user=peer --secret=peerpw --type=peer \
  --enroll-id enroll --enroll-secret=enrollpw --mspid PUCMSP


  echo "Realizando deploy de 1 peer para a organizacao PUC"


  export PEER_IMAGE=quay.io/kfsoftware/fabric-peer
  export PEER_VERSION=2.4.1-v0.0.3
  export MSP_ORG=PUCMSP
  export PEER_SECRET=peerpw

  kubectl hlf peer create --statedb=couchdb --image=$PEER_IMAGE --version=$PEER_VERSION --storage-class=standard --enroll-id=peer --mspid=$MSP_ORG \
  --enroll-pw=$PEER_SECRET --capacity=5Gi --name=puc-peer0 --ca-name=puc-ca.default --k8s-builder=true --hosts=peer0-puc.localho.st --istio-port=443

  kubectl wait --timeout=180s --for=condition=Running fabricpeers.hlf.kungfusoftware.es --all

  # leva alguns minutos


  echo "Criando 3 orderers"

  kubectl hlf ca create  --image=$CA_IMAGE --version=$CA_VERSION --storage-class=standard --capacity=1Gi --name=ord-ca \
      --enroll-id=enroll --enroll-pw=enrollpw --hosts=ord-ca.localho.st --istio-port=443

  kubectl wait --timeout=180s --for=condition=Running fabriccas.hlf.kungfusoftware.es --all

  kubectl hlf ca register --name=ord-ca --user=orderer --secret=ordererpw \
      --type=orderer --enroll-id enroll --enroll-secret=enrollpw --mspid=OrdererMSP --ca-url="https://ord-ca.localho.st:443"


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


  echo "Iniciando processo de levantamento de canal..."

  echo "Registrando identidade das organizacoes..."

  # register
  kubectl hlf ca register --name=ord-ca --user=admin --secret=adminpw \
      --type=admin --enroll-id enroll --enroll-secret=enrollpw --mspid=OrdererMSP

  # enroll
  kubectl hlf ca enroll --name=ord-ca --namespace=default \
      --user=admin --secret=adminpw --mspid OrdererMSP \
      --ca-name tlsca  --output resources/orderermsp.yaml

  # register
  kubectl hlf ca register --name=inmetro-ca --namespace=default --user=admin --secret=adminpw \
      --type=admin --enroll-id enroll --enroll-secret=enrollpw --mspid=INMETROMSP

  # enroll
  kubectl hlf ca enroll --name=inmetro-ca --namespace=default \
      --user=admin --secret=adminpw --mspid INMETROMSP \
      --ca-name ca  --output resources/inmetromsp.yaml

  # register
  kubectl hlf ca register --name=puc-ca --namespace=default --user=admin --secret=adminpw \
      --type=admin --enroll-id enroll --enroll-secret=enrollpw --mspid=PUCMSP

  # enroll
  kubectl hlf ca enroll --name=puc-ca --namespace=default \
      --user=admin --secret=adminpw --mspid PUCMSP \
      --ca-name ca  --output resources/pucmsp.yaml

  echo "Criando segredo / secret"

  kubectl create secret generic wallet --namespace=default \
          --from-file=inmetromsp.yaml=$PWD/resources/inmetromsp.yaml \
          --from-file=pucmsp.yaml=$PWD/resources/pucmsp.yaml \
          --from-file=orderermsp.yaml=$PWD/resources/orderermsp.yaml

  echo "Criando canal principal"

  export PEER_ORG_SIGN_CERT=$(kubectl get fabriccas inmetro-ca -o=jsonpath='{.status.ca_cert}')
  export PEER_ORG_TLS_CERT=$(kubectl get fabriccas inmetro-ca -o=jsonpath='{.status.tlsca_cert}')
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
      - mspID: INMETROMSP
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
      - mspID: INMETROMSP
        caName: "inmetro-ca"
        caNamespace: "default"
      - mspID: PUCMSP
        caName: "puc-ca"
        caNamespace: "default" 
    identities:
      OrdererMSP:
        secretKey: orderermsp.yaml
        secretName: wallet
        secretNamespace: default
      INMETROMSP:
        secretKey: inmetromsp.yaml
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


  echo "Ingressando peers das organizacoes no canal"


  export IDENT_8=$(printf "%8s" "")
  export ORDERER0_TLS_CERT=$(kubectl get fabricorderernodes ord-node0 -o=jsonpath='{.status.tlsCert}' | sed -e "s/^/${IDENT_8}/" )

  kubectl apply -f - <<EOF
  apiVersion: hlf.kungfusoftware.es/v1alpha1
  kind: FabricFollowerChannel
  metadata:
    name: demo-inmetromsp
  spec:
    anchorPeers:
      - host: inmetro-peer0.default
        port: 7051
    hlfIdentity:
      secretKey: inmetromsp.yaml
      secretName: wallet
      secretNamespace: default
    mspId: INMETROMSP
    name: demo
    externalPeersToJoin: []
    orderers:
      - certificate: |
  ${ORDERER0_TLS_CERT}
        url: grpcs://ord-node0.default:7050
    peersToJoin:
      - name: inmetro-peer0
        namespace: default
EOF


  export IDENT_8=$(printf "%8s" "")
  export ORDERER0_TLS_CERT=$(kubectl get fabricorderernodes ord-node0 -o=jsonpath='{.status.tlsCert}' | sed -e "s/^/${IDENT_8}/" )

  kubectl apply -f - <<EOF
  apiVersion: hlf.kungfusoftware.es/v1alpha1
  kind: FabricFollowerChannel
  metadata:
    name: demo-pucmsp
  spec:
    anchorPeers:
      - host: puc-peer0.default
        port: 7051
    hlfIdentity:
      secretKey: pucmsp.yaml
      secretName: wallet
      secretNamespace: default
    mspId: PUCMSP
    name: demo
    externalPeersToJoin: []
    orderers:
      - certificate: |
  ${ORDERER0_TLS_CERT}
        url: grpcs://orderer0-ord.localho.st:443
    peersToJoin:
      - name: puc-peer0
        namespace: default
EOF

  echo "Fim"

}


if [ "$1" == "up" ]; then
    echo "Iniciando processo de levantamento do cluster Kubernetes"
    up
    echo "Verifique o status dos pods com o comando: 'kubectl get pods'"

elif [ "$1" == "chaincode" ]; then
  echo "Iniciando processo de deploy do chaincode"
  echo "Nome do chaincode: $2"
  chaincode_install $2
  echo "Verifique o status dos pods com o comando: 'kubectl get pods'"

elif [ "$1" == "down" ]; then
  echo "Iniciando processo de destruicao do cluster Kubernetes"
    
  kubectl delete fabricorderernodes.hlf.kungfusoftware.es --all-namespaces --all
  kubectl delete fabricpeers.hlf.kungfusoftware.es --all-namespaces --all
  kubectl delete fabriccas.hlf.kungfusoftware.es --all-namespaces --all
  kubectl delete fabricchaincode.hlf.kungfusoftware.es --all-namespaces --all
  kubectl delete fabricmainchannels --all-namespaces --all
  kubectl delete fabricfollowerchannels --all-namespaces --all
  
  kind delete cluster
  
  exit 0

else
    echo "Comando não reconhecido"
    echo "Comandos possíveis: "
    echo "'up' - Inicia o cluster Kubernetes"
    echo "'chaincode <nome do chaincode>' - Realiza o deploy do chaincode"
    echo "'down' - Destrói o cluster Kubernetes e todos os recursos criados"
    exit 1
fi