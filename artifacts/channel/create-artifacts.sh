# Delete existing artifacts
rm genesis.block mychannel.tx
rm -rf ../../channel-artifacts/*

#Generate Crypto artifactes for organizations
# cryptogen generate --config=./crypto-config.yaml --output=./crypto-config/



# System channel
SYS_CHANNEL="sys-channel"
SYS_CHANNEL2="sys-channel2"


# channel name defaults to "mychannel"
CHANNEL_NAME="channel1"
CHANNEL_NAME2="channel2"


echo $CHANNEL_NAME
echo $CHANNEL_NAME2

# Generate System Genesis block
configtxgen -profile OrdererGenesis -configPath . -channelID $SYS_CHANNEL  -outputBlock ./genesis.block
configtxgen -profile OrdererGenesis -configPath . -channelID $SYS_CHANNEL2  -outputBlock ./genesis2.block


# Generate channel configuration block channel1
configtxgen -profile BasicChannel -configPath . -outputCreateChannelTx ./channel1.tx -channelID $CHANNEL_NAME

echo "#######    Generating anchor peer update for PenangkarMSP  ##########"
configtxgen -profile BasicChannel -configPath . -outputAnchorPeersUpdate ./PenangkarMSPanchors_${CHANNEL_NAME}.tx -channelID $CHANNEL_NAME -asOrg PenangkarMSP

echo "#######    Generating anchor peer update for PetaniMSP  ##########"
configtxgen -profile BasicChannel -configPath . -outputAnchorPeersUpdate ./PetaniMSPanchors_${CHANNEL_NAME}.tx -channelID $CHANNEL_NAME -asOrg PetaniMSP

echo "#######    Generating anchor peer update for PengumpulMSP  ##########"
configtxgen -profile BasicChannel -configPath . -outputAnchorPeersUpdate ./PengumpulMSPanchors_${CHANNEL_NAME}.tx -channelID $CHANNEL_NAME -asOrg PengumpulMSP

echo "#######    Generating anchor peer update for PedagangMSP  ##########"
configtxgen -profile BasicChannel -configPath . -outputAnchorPeersUpdate ./PedagangMSPanchors_${CHANNEL_NAME}.tx -channelID $CHANNEL_NAME -asOrg PedagangMSP

# Generate channel configuration block channel2
configtxgen -profile ChannelDua -configPath . -outputCreateChannelTx ./channel2.tx -channelID $CHANNEL_NAME2

echo "#######    Generating anchor peer update for PenangkarMSP  ##########"
configtxgen -profile ChannelDua -configPath . -outputAnchorPeersUpdate ./PenangkarMSPanchors_${CHANNEL_NAME2}.tx -channelID $CHANNEL_NAME2 -asOrg PenangkarMSP

echo "#######    Generating anchor peer update for PetaniMSP  ##########"
configtxgen -profile ChannelDua -configPath . -outputAnchorPeersUpdate ./PetaniMSPanchors_${CHANNEL_NAME2}.tx -channelID $CHANNEL_NAME2 -asOrg PetaniMSP

# echo "#######    Generating anchor peer update for PengumpulMSP  ##########"
# configtxgen -profile ChannelDua -configPath . -outputAnchorPeersUpdate ./PengumpulMSPanchors_${CHANNEL_NAME2}.tx -channelID $CHANNEL_NAME2 -asOrg PengumpulMSP

# echo "#######    Generating anchor peer update for PedagangMSP  ##########"
# configtxgen -profile ChannelDua -configPath . -outputAnchorPeersUpdate ./PedagangMSPanchors_${CHANNEL_NAME2}.tx -channelID $CHANNEL_NAME2 -asOrg PedagangMSP

