## hasil export dari artifact
export CORE_PEER_TLS_ENABLED=true
export ORDERER_CA=${PWD}/artifacts/channel/crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
export PEER0_PENANGKAR_CA=${PWD}/artifacts/channel/crypto-config/peerOrganizations/penangkar.example.com/peers/peer0.penangkar.example.com/tls/ca.crt
export PEER0_PETANI_CA=${PWD}/artifacts/channel/crypto-config/peerOrganizations/petani.example.com/peers/peer0.petani.example.com/tls/ca.crt
export PEER0_PENGUMPUL_CA=${PWD}/artifacts/channel/crypto-config/peerOrganizations/pengumpul.example.com/peers/peer0.pengumpul.example.com/tls/ca.crt
export PEER0_PEDAGANG_CA=${PWD}/artifacts/channel/crypto-config/peerOrganizations/pedagang.example.com/peers/peer0.pedagang.example.com/tls/ca.crt
export FABRIC_CFG_PATH=${PWD}/artifacts/channel/config/

export CHANNEL_NAME=channel1 ## nama channel

setGlobalsForOrderer() {
    export CORE_PEER_LOCALMSPID="OrdererMSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/artifacts/channel/crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
    export CORE_PEER_MSPCONFIGPATH=${PWD}/artifacts/channel/crypto-config/ordererOrganizations/example.com/users/Admin@example.com/msp

}

setGlobalsForPeer0Penangkar() {
    export CORE_PEER_LOCALMSPID="PenangkarMSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_PENANGKAR_CA
    export CORE_PEER_MSPCONFIGPATH=${PWD}/artifacts/channel/crypto-config/peerOrganizations/penangkar.example.com/users/Admin@penangkar.example.com/msp
    export CORE_PEER_ADDRESS=localhost:7051
}

setGlobalsForPeer0Petani() {
    export CORE_PEER_LOCALMSPID="PetaniMSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_PETANI_CA
    export CORE_PEER_MSPCONFIGPATH=${PWD}/artifacts/channel/crypto-config/peerOrganizations/petani.example.com/users/Admin@petani.example.com/msp
    export CORE_PEER_ADDRESS=localhost:8051

}

setGlobalsForPeer0Pengumpul(){
    export CORE_PEER_LOCALMSPID="PengumpulMSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_PENGUMPUL_CA
    export CORE_PEER_MSPCONFIGPATH=${PWD}/artifacts/channel/crypto-config/peerOrganizations/pengumpul.example.com/users/Admin@pengumpul.example.com/msp
    export CORE_PEER_ADDRESS=localhost:10051
    
}

setGlobalsForPeer0Pedagang(){
    export CORE_PEER_LOCALMSPID="PedagangMSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_PEDAGANG_CA
    export CORE_PEER_MSPCONFIGPATH=${PWD}/artifacts/channel/crypto-config/peerOrganizations/pedagang.example.com/users/Admin@pedagang.example.com/msp
    export CORE_PEER_ADDRESS=localhost:11051
    
}

presetup() {
    echo Vendoring Go dependencies ...
    pushd ./artifacts/src/github.com/mangga/go 
    GO111MODULE=on go mod vendor
    popd
    echo Finished vendoring Go dependencies
}
# presetup

CHANNEL_NAME="channel1"
CC_RUNTIME_LANGUAGE="golang"
VERSION="3.22"
SEQUENCE="4"
CC_SRC_PATH="./artifacts/src/github.com/mangga/go" 
CC_NAME="manggach1_cc" ## nama chaincode

packageChaincode() {
    rm -rf ${CC_NAME}.tar.gz
    setGlobalsForPeer0Penangkar
    peer lifecycle chaincode package ${CC_NAME}.tar.gz \
        --path ${CC_SRC_PATH} --lang ${CC_RUNTIME_LANGUAGE} \
        --label ${CC_NAME}_${VERSION}
    echo "===================== Chaincode is packaged ===================== "
}
# packageChaincode

installChaincode() {
    setGlobalsForPeer0Penangkar
    peer lifecycle chaincode install ${CC_NAME}.tar.gz
    echo "===================== Chaincode is installed on peer0.penangkar ===================== "

    setGlobalsForPeer0Petani
    peer lifecycle chaincode install ${CC_NAME}.tar.gz
    echo "===================== Chaincode is installed on peer0.petani ===================== "

    setGlobalsForPeer0Pengumpul
    peer lifecycle chaincode install ${CC_NAME}.tar.gz
    echo "===================== Chaincode is installed on peer0.pengumpul ===================== "

    setGlobalsForPeer0Pedagang
    peer lifecycle chaincode install ${CC_NAME}.tar.gz
    echo "===================== Chaincode is installed on peer0.pedagang ===================== "
}

# installChaincode

queryInstalled() {
    setGlobalsForPeer0Penangkar
    peer lifecycle chaincode queryinstalled >&log.txt
    cat log.txt
    PACKAGE_ID=$(sed -n "/${CC_NAME}_${VERSION}/{s/^Package ID: //; s/, Label:.*$//; p;}" log.txt)
    echo PackageID is ${PACKAGE_ID}
    echo "===================== Query installed successful on peer0.penangkar on channel ===================== "
}

# queryInstalled

# --collections-config ./artifacts/private-data/collections_config.json \
#         --signature-policy "OR('Org1MSP.member','Org2MSP.member')" \

approveForMyPenangkar() {
    setGlobalsForPeer0Penangkar
    # set -x
    peer lifecycle chaincode approveformyorg -o localhost:7050 \
        --ordererTLSHostnameOverride orderer.example.com --tls \
        --cafile $ORDERER_CA --channelID $CHANNEL_NAME --name ${CC_NAME} --version ${VERSION} \
        --init-required --package-id ${PACKAGE_ID} \
        --sequence ${SEQUENCE}
    # set +x

    echo "===================== chaincode approved from penangkar ===================== "

}
# queryInstalled
# approveForMyPenangkar

# --signature-policy "OR ('Org1MSP.member')"
# --peerAddresses localhost:7051 --tlsRootCertFiles $PEER0_ORG1_CA --peerAddresses localhost:9051 --tlsRootCertFiles $PEER0_ORG2_CA
# --peerAddresses peer0.org1.example.com:7051 --tlsRootCertFiles $PEER0_ORG1_CA --peerAddresses peer0.org2.example.com:9051 --tlsRootCertFiles $PEER0_ORG2_CA
#--channel-config-policy Channel/Application/Admins
# --signature-policy "OR ('Org1MSP.peer','Org2MSP.peer')"

checkCommitReadynessPenangkar() {
    setGlobalsForPeer0Penangkar
    peer lifecycle chaincode checkcommitreadiness \
        --channelID $CHANNEL_NAME --name ${CC_NAME} --version ${VERSION} \
        --sequence ${SEQUENCE} --output json --init-required
    echo "===================== checking commit readyness from penangkar ===================== "
}

# checkCommitReadyness

approveForMyPetani() {
    setGlobalsForPeer0Petani

    peer lifecycle chaincode approveformyorg -o localhost:7050 \
        --ordererTLSHostnameOverride orderer.example.com --tls $CORE_PEER_TLS_ENABLED \
        --cafile $ORDERER_CA --channelID $CHANNEL_NAME --name ${CC_NAME} \
        --version ${VERSION} --init-required --package-id ${PACKAGE_ID} \
        --sequence ${SEQUENCE}

    echo "===================== chaincode approved from petani ===================== "
}

# queryInstalled
# approveForMyPetani

checkCommitReadynessPetani() {

    setGlobalsForPeer0Petani
    peer lifecycle chaincode checkcommitreadiness --channelID $CHANNEL_NAME \
        --peerAddresses localhost:8051 --tlsRootCertFiles $PEER0_PETANI_CA \
        --name ${CC_NAME} --version ${VERSION} --sequence ${SEQUENCE} --output json --init-required
    echo "===================== checking commit readyness from petani ===================== "
}

# checkCommitReadyness

approveForMyPengumpul() {
    setGlobalsForPeer0Pengumpul

    peer lifecycle chaincode approveformyorg -o localhost:7050 \
        --ordererTLSHostnameOverride orderer.example.com --tls $CORE_PEER_TLS_ENABLED \
        --cafile $ORDERER_CA --channelID $CHANNEL_NAME --name ${CC_NAME} \
        --version ${VERSION} --init-required --package-id ${PACKAGE_ID} \
        --sequence ${SEQUENCE}

    echo "===================== chaincode approved from pengumpul ===================== "
}

# queryInstalled
# approveForMyPengumpul

checkCommitReadynessPengumpul() {

    setGlobalsForPeer0Pengumpul
    peer lifecycle chaincode checkcommitreadiness --channelID $CHANNEL_NAME \
        --peerAddresses localhost:10051 --tlsRootCertFiles $PEER0_PENGUMPUL_CA \
        --name ${CC_NAME} --version ${VERSION} --sequence ${SEQUENCE} --output json --init-required
    echo "===================== checking commit readyness from pengumpul ===================== "
}

# checkCommitReadyness

approveForMyPedagang() {
    setGlobalsForPeer0Pedagang

    peer lifecycle chaincode approveformyorg -o localhost:7050 \
        --ordererTLSHostnameOverride orderer.example.com --tls $CORE_PEER_TLS_ENABLED \
        --cafile $ORDERER_CA --channelID $CHANNEL_NAME --name ${CC_NAME} \
        --version ${VERSION} --init-required --package-id ${PACKAGE_ID} \
        --sequence ${SEQUENCE}

    echo "===================== chaincode approved from pedagang ===================== "
}

# queryInstalled
# approveForMyPedagang

checkCommitReadynessPedagang() {

    setGlobalsForPeer0Pedagang
    peer lifecycle chaincode checkcommitreadiness --channelID $CHANNEL_NAME \
        --peerAddresses localhost:11051 --tlsRootCertFiles $PEER0_PEDAGANG_CA \
        --name ${CC_NAME} --version ${VERSION} --sequence ${SEQUENCE} --output json --init-required
    echo "===================== checking commit readyness from pedagang ===================== "
}

# checkCommitReadyness

commitChaincodeDefination() {
    setGlobalsForPeer0Penangkar
    peer lifecycle chaincode commit -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com \
        --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA \
        --channelID $CHANNEL_NAME --name ${CC_NAME} \
        --peerAddresses localhost:7051 --tlsRootCertFiles $PEER0_PENANGKAR_CA \
        --peerAddresses localhost:8051 --tlsRootCertFiles $PEER0_PETANI_CA \
        --peerAddresses localhost:10051 --tlsRootCertFiles $PEER0_PENGUMPUL_CA \
        --peerAddresses localhost:11051 --tlsRootCertFiles $PEER0_PEDAGANG_CA \
        --version ${VERSION} --sequence ${SEQUENCE} --init-required

}

# commitChaincodeDefination

queryCommitted() {
    setGlobalsForPeer0Penangkar
    peer lifecycle chaincode querycommitted --channelID $CHANNEL_NAME --name ${CC_NAME}

}

# queryCommitted

chaincodeInvokeInit() {
    setGlobalsForPeer0Penangkar
    peer chaincode invoke -o localhost:7050 \
        --ordererTLSHostnameOverride orderer.example.com \
        --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA \
        -C $CHANNEL_NAME -n ${CC_NAME} \
        --peerAddresses localhost:7051 --tlsRootCertFiles $PEER0_PENANGKAR_CA \
        --peerAddresses localhost:8051 --tlsRootCertFiles $PEER0_PETANI_CA \
        --peerAddresses localhost:10051 --tlsRootCertFiles $PEER0_PENGUMPUL_CA \
        --peerAddresses localhost:11051 --tlsRootCertFiles $PEER0_PEDAGANG_CA \
        --isInit -c '{"Args":[]}'

}

# chaincodeInvokeInit

chaincodeInvoke() {
    setGlobalsForPeer0Penangkar

    # Registrasi Benih
    peer chaincode invoke -o localhost:7050 \
        --ordererTLSHostnameOverride orderer.example.com \
        --tls $CORE_PEER_TLS_ENABLED \
        --cafile $ORDERER_CA \
        -C $CHANNEL_NAME -n ${CC_NAME}  \
        --peerAddresses localhost:7051 --tlsRootCertFiles $PEER0_PENANGKAR_CA \
        --peerAddresses localhost:8051 --tlsRootCertFiles $PEER0_PETANI_CA \
        --peerAddresses localhost:10051 --tlsRootCertFiles $PEER0_PENGUMPUL_CA \
        --peerAddresses localhost:11051 --tlsRootCertFiles $PEER0_PEDAGANG_CA \
        -c '{"function": "PanenMangga","Args":[
        "{\"ukuran\":\"lonjong\",\"pestisida\":\"abcd\",\"kadarAir\":85.5,\"perlakuan\":\"baik\",\"produktivitas\":\"subur\",\"kuantitasManggaKg\": 100}",
        "f05b894dfce8ffcf7e60fa1d78fa7399efdff157fe4e2d6e0e0274c06ec7d061"
    ]}'

}

# chaincodeInvoke

chaincodeQuery() {
    setGlobalsForPeer0Penangkar

    # Registrasi Benih
    peer chaincode invoke -o localhost:7050 \
        --ordererTLSHostnameOverride orderer.example.com \
        --tls $CORE_PEER_TLS_ENABLED \
        --cafile $ORDERER_CA \
        -C $CHANNEL_NAME -n ${CC_NAME}  \
        --peerAddresses localhost:7051 --tlsRootCertFiles $PEER0_PENANGKAR_CA \
        -c '{"function": "GetHistoryForAssetByID","Args":["e2efc7c6bc37a249d57bba7e1366f470a1f05cc482ce6d8009ba0128e9217c59"]}'
}

# chaincodeQuery

# Run this function if you add any new dependency in chaincode
# presetup

# packageChaincode
# installChaincode
# queryInstalled

# approveForMyPenangkar
# checkCommitReadynessPenangkar
# approveForMyPetani
# checkCommitReadynessPetani
# approveForMyPengumpul
# checkCommitReadynessPengumpul
# approveForMyPedagang
# checkCommitReadynessPedagang

# commitChaincodeDefination
# queryCommitted
# chaincodeInvokeInit


# sleep 5
# chaincodeInvoke
# sleep 3
# chaincodeQuery
