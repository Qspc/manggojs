# Delete existing artifacts
rm genesis.block mychannel.tx
rm -rf ../../channel-artifacts/*

#Generate Crypto artifactes for organizations
# cryptogen generate --config=./crypto-config.yaml --output=./crypto-config/



# System channel
SYS_CHANNEL="sys-channel"

# channel name defaults to "mychannel"
CHANNEL_NAME="channel1"`
`CHANNEL_NAME2="channel2"`


echo $CHANNEL_NAME
echo $CHANNEL_NAME2

# Generate System Genesis block
configtxgen -profile OrdererGenesis -configPath . -channelID $SYS_CHANNEL  -outputBlock ./genesis.block


# Generate channel configuration block channel1
configtxgen -profile BasicChannel -configPath . -outputCreateChannelTx ./mychannel.tx -channelID $CHANNEL_NAME

echo "#######    Generating anchor peer update for PenangkarMSP  ##########"
configtxgen -profile BasicChannel -configPath . -outputAnchorPeersUpdate ./PenangkarMSPanchors.tx -channelID $CHANNEL_NAME -asOrg PenangkarMSP

echo "#######    Generating anchor peer update for PetaniMSP  ##########"
configtxgen -profile BasicChannel -configPath . -outputAnchorPeersUpdate ./PetaniMSPanchors.tx -channelID $CHANNEL_NAME -asOrg PetaniMSP

echo "#######    Generating anchor peer update for PengumpulMSP  ##########"
configtxgen -profile BasicChannel -configPath . -outputAnchorPeersUpdate ./PengumpulMSPanchors.tx -channelID $CHANNEL_NAME -asOrg PengumpulMSP

echo "#######    Generating anchor peer update for PedagangMSP  ##########"
configtxgen -profile BasicChannel -configPath . -outputAnchorPeersUpdate ./PedagangMSPanchors.tx -channelID $CHANNEL_NAME -asOrg PedagangMSP

# Generate channel configuration block channel2
configtxgen -profile BasicChannel -configPath . -outputCreateChannelTx ./mychannel.tx -channelID $CHANNEL_NAME2

echo "#######    Generating anchor peer update for PenangkarMSP  ##########"
configtxgen -profile BasicChannel -configPath . -outputAnchorPeersUpdate ./PenangkarMSPanchors.tx -channelID $CHANNEL_NAME2 -asOrg PenangkarMSP

echo "#######    Generating anchor peer update for PetaniMSP  ##########"
configtxgen -profile BasicChannel -configPath . -outputAnchorPeersUpdate ./PetaniMSPanchors.tx -channelID $CHANNEL_NAME2 -asOrg PetaniMSP

echo "#######    Generating anchor peer update for PengumpulMSP  ##########"
configtxgen -profile BasicChannel -configPath . -outputAnchorPeersUpdate ./PengumpulMSPanchors.tx -channelID $CHANNEL_NAME2 -asOrg PengumpulMSP

echo "#######    Generating anchor peer update for PedagangMSP  ##########"
configtxgen -profile BasicChannel -configPath . -outputAnchorPeersUpdate ./PedagangMSPanchors.tx -channelID $CHANNEL_NAME2 -asOrg PedagangMSP


# # Generate channel configuration block channel1
# configtxgen -profile BasicChannel -configPath . -outputCreateChannelTx ./mychannel.tx -channelID channel1

# echo "#######    Generating anchor peer update for PenangkarMSP  ##########"
# configtxgen -profile BasicChannel -configPath . -outputAnchorPeersUpdate ./PenangkarMSPanchors.tx -channelID channel1 -asOrg PenangkarMSP

# echo "#######    Generating anchor peer update for PetaniMSP  ##########"
# configtxgen -profile BasicChannel -configPath . -outputAnchorPeersUpdate ./PetaniMSPanchors.tx -channelID channel1 -asOrg PetaniMSP

# echo "#######    Generating anchor peer update for PengumpulMSP  ##########"
# configtxgen -profile BasicChannel -configPath . -outputAnchorPeersUpdate ./PengumpulMSPanchors.tx -channelID channel1 -asOrg PengumpulMSP

# echo "#######    Generating anchor peer update for PedagangMSP  ##########"
# configtxgen -profile BasicChannel -configPath . -outputAnchorPeersUpdate ./PedagangMSPanchors.tx -channelID channel1 -asOrg PedagangMSP

# # Generate channel configuration block channel2
# configtxgen -profile BasicChannel -configPath . -outputCreateChannelTx ./mychannel.tx -channelID channel2

# echo "#######    Generating anchor peer update for PenangkarMSP  ##########"
# configtxgen -profile BasicChannel -configPath . -outputAnchorPeersUpdate ./PenangkarMSPanchors.tx -channelID channel2 -asOrg PenangkarMSP

# echo "#######    Generating anchor peer update for PetaniMSP  ##########"
# configtxgen -profile BasicChannel -configPath . -outputAnchorPeersUpdate ./PetaniMSPanchors.tx -channelID channel2 -asOrg PetaniMSP

# echo "#######    Generating anchor peer update for PengumpulMSP  ##########"
# configtxgen -profile BasicChannel -configPath . -outputAnchorPeersUpdate ./PengumpulMSPanchors.tx -channelID channel2 -asOrg PengumpulMSP

# echo "#######    Generating anchor peer update for PedagangMSP  ##########"
# configtxgen -profile BasicChannel -configPath . -outputAnchorPeersUpdate ./PedagangMSPanchors.tx -channelID channel2 -asOrg PedagangMSP