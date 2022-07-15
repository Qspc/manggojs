const { Gateway, Wallets, TxEventHandler, GatewayOptions, DefaultEventHandlerStrategies, TxEventHandlerFactory } = require('fabric-network');
const fs = require('fs');
const EventStrategies = require('fabric-network/lib/impl/event/defaulteventhandlerstrategies');
const path = require('path');
const log4js = require('log4js');
const logger = log4js.getLogger('BasicNetwork');
const util = require('util');

const helper = require('./helper');
const { blockListener, contractListener } = require('./Listeners');

// fungsi menerima data seputar blockchain berdasarkan parameter yang dikirim
const invokeTransaction = async (channelName, chaincodeName, fcn, args, userName, role, transientData) => {
  try {
    var roleAktor;
    if (role == 1) {
      roleAktor = 'Penangkar';
    } else if (role == 2) {
      roleAktor = 'Petani';
    } else if (role == 3) {
      roleAktor = 'Pengumpul';
    } else if (role == 4) {
      roleAktor = 'Pedagang';
    }
    logger.debug('role is  : ' + roleAktor);
    logger.debug('username is  : ' + userName);



    //identifikasi role
    const ccp = await helper.getCCP(roleAktor);
    // console.log(ccp.peers);
    // mengambil jalur path
    const walletPath = await helper.getWalletPath(roleAktor);
    // menegecek role di fabric
    const wallet = await Wallets.newFileSystemWallet(walletPath);
    console.log(`Wallet path: ${walletPath}`);

    // identifikasi nama username di blockchain
    let identity = await wallet.get(userName);
    if (!identity) {
      console.log(`An identity for the user ${userName} does not exist in the wallet, so registering user`);
      await helper.getRegisteredUser(userName, role, true);
      identity = await wallet.get(userName);
      console.log('Register first before retrying ');
      return;
    }

    // console.log(ccp)
    // console.log(wallet)


    //if (orgName != "")

    const connectOptions = {
      wallet,
      identity: userName,
      discovery: { enabled: false, asLocalhost: true },
      // eventHandlerOptions: EventStrategies.NONE,
    };

    // gateaway = koneksi untuk mengambil jalur jaringan fabric
    const gateway = new Gateway();
    await gateway.connect(ccp, connectOptions);
    // console.log('cek eror 1')


    // console.log(channelName)
    // console.log(chaincodeName)
    //ambil nama channel
    const network = await gateway.getNetwork(channelName);
    //ambil nama chaincode
    const contract = network.getContract(chaincodeName);
    // console.log(contract)

    // await contract.addContractListener(contractListener);
    // await network.addBlockListener(blockListener);
    // console.log('cek eror 2')


    // Multiple smartcontract in one chaincode
    let result;
    let message;
    // console.log('fcn')
    // console.log(args[0])

    // identifikasi fungsi/perintah apa yang digunakan
    switch (fcn) {
      case 'RegistrasiBenih':
        // result = 'berhasil'
        result = await contract.submitTransaction('ManggaContract:'+fcn, args[0]);
        console.log(result.toString());
        result = {txid: result.toString()}
        break;
      case 'TanamBenih':
        result = await contract.submitTransaction('ManggaContract:' + fcn, args[0], args[1]); //info mangga, data dari block sebelumnya
        console.log(result.toString());
        result = { txid: result.toString() };
        break;
      case 'PanenMangga':
        result = await contract.submitTransaction('ManggaContract:' + fcn, args[0], args[1]);
        console.log(result.toString());
        result = { txid: result.toString() };
        break;
      case 'CreateUser':
        result = await contract.submitTransaction('UserContract:' + fcn, args[0]);
        console.log(result.toString());
        result = { txid: result.toString() };
        break;
      case 'CreateTrxManggaByPenangkar':
        result = await contract.submitTransaction('ManggaContract:' + fcn, args[0], args[1]);
        console.log(result.toString());
        result = { txid: result.toString() };
        break;
      case 'CreateTrxManggaByPetani':
        result = await contract.submitTransaction('ManggaContract:' + fcn, args[0], args[1]);
        console.log(result.toString());
        result = { txid: result.toString() };
        break;
      case 'CreateTrxManggaByPengumpul':
        result = await contract.submitTransaction('ManggaContract:' + fcn, args[0], args[1]);
        console.log(result.toString());
        result = { txid: result.toString() };
        break;
      case 'CreateTrxManggaByPedagang':
        result = await contract.submitTransaction('ManggaContract:' + fcn, args[0], args[1]);
        console.log(result.toString());
        result = { txid: result.toString() };
        break;
      case 'AddKuantitasBenihByID':
        result = await contract.submitTransaction('ManggaContract:' + fcn, args[0], args[1]);
        var data = JSON.parse(result.toString());
        console.log(result.toString());
        result = { bawang: data };
        break;
      case 'AddManggaKuantitasByID':
        result = await contract.submitTransaction('ManggaContract:' + fcn, args[0], args[1]);
        var data = JSON.parse(result.toString());
        console.log(result.toString());
        result = { bawang: data };
        break;
      case 'ConfirmTrxByID':
        result = await contract.submitTransaction('ManggaContract:' + fcn, args[0]);
        console.log(result.toString());
        result = { txid: result.toString() };
        break;
      case 'RejectTrxByID':
        result = await contract.submitTransaction('ManggaContract:' + fcn, args[0], args[1], args[2], args[3]);
        console.log(result.toString());
        result = { txid: result.toString() };
        break;
      case 'CreateUser':
        result = await contract.submitTransaction('UserContract:' + fcn, args[0]);
        console.log(result.toString());
        result = { txid: result.toString() };
        break;

      default:
        break;
    }

    console.log('cek eror 3')


    // let result
    // let message;
    // if (fcn === "createCar" || fcn === "createPrivateCarImplicitForOrg1"
    //     || fcn == "createPrivateCarImplicitForOrg2") {
    //     result = await contract.submitTransaction(fcn, args[0], args[1], args[2], args[3], args[4]);
    //     message = `Successfully added the car asset with key ${args[0]}`

    // } else if (fcn === "changeCarOwner") {
    //     result = await contract.submitTransaction(fcn, args[0], args[1]);
    //     message = `Successfully changed car owner with key ${args[0]}`
    // } else if (fcn == "createPrivateCar" || fcn =="updatePrivateData") {
    //     console.log(`Transient data is : ${transientData}`)
    //     let carData = JSON.parse(transientData)
    //     console.log(`car data is : ${JSON.stringify(carData)}`)
    //     let key = Object.keys(carData)[0]
    //     const transientDataBuffer = {}
    //     transientDataBuffer[key] = Buffer.from(JSON.stringify(carData.car))
    //     result = await contract.createTransaction(fcn)
    //         .setTransient(transientDataBuffer)
    //         .submit()
    //     message = `Successfully submitted transient data`
    // }
    // else {
    //     return `Invocation require either createCar or changeCarOwner as function but got ${fcn}`
    // }

    await gateway.disconnect();

    // result = JSON.parse(result.toString());

    let response = {
      // message: message,
      message: `transaction successfully`,
      result,
    };

    return response;
  } catch (error) {
    console.log(`Getting error: ${error}`);
    return error.message;
  }
};

exports.invokeTransaction = invokeTransaction;
