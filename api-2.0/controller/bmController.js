require('dotenv').config();
const express = require('express');
const app = express();
const log4js = require('log4js');
const logger = log4js.getLogger('BasicNetwork');

const query = require('../app/query');
const invoke = require('../app/invoke');

logger.level = 'debug';

app.set('secret', process.env.ACCESS_TOKEN_SECRET);

function getErrorMessage(field) {
  var response = {
    success: false,
    message: field + ' field is missing or Invalid in the request',
  };
  return response;
}

exports.Query = async (req, res) => {
  // params = isinya link yg dipakai
  console.log(req.params);
  // query = isinya key & value
  console.log(req.query);

  // data yang di request/dikirim : channelName, chaincodeName, peer (satu aktor aja), fcn, args
  try {
    logger.debug('==================== QUERY BY CHAINCODE ==================');

    var channelName = req.params.channelName;
    var chaincodeName = req.params.chaincodeName;
    console.log(`chaincode name is :${chaincodeName}`);
    let args = req.query.args;
    let fcn = req.query.fcn;
    let peer = req.query.peer;

    logger.debug('channelName : ' + channelName);
    logger.debug('chaincodeName : ' + chaincodeName);
    logger.debug('fcn : ' + fcn);
    logger.debug('args : ' + args);

    if (!chaincodeName) {
      res.status(400).json(getErrorMessage("'chaincodeName'"));
      return;
    }
    if (!channelName) {
      res.status(400).json(getErrorMessage("'channelName'"));
      return;
    }
    if (!fcn) {
      res.status(400).json(getErrorMessage("'fcn'"));
      return;
    }
    if (!args) {
      res.status(400).json(getErrorMessage("'args'"));
      return;
    }

    console.log('args ==========', args);
    console.log(`type of argument ${typeof args}`);
    args = args.replace(/'/g, '"');
    args = JSON.parse(args);
    console.log('args parse ==========', args);

    // proses GET/perintah di proses --> query
    let message = await query.queryTranasction(channelName, chaincodeName, args, fcn, req.role, req.userName);

    logger.debug(message);
    // TODO do error handling here, meesage.contain error -->
    if (typeof message == 'string') {
      const response_payload = {
        result: message,
        error: 'Error in chaincode',
      };

      res.status(500).send(response_payload);
    } else if (typeof message == 'object') {
      if (message.length == 0) {
        const response_payload = {
          result: message,
          error: '404',
          errorData: 'Data not found',
        };

        res.status(404).send(response_payload);
      } else {
        const response_payload = {
          result: message,
          error: null,
          errorData: null,
        };
        res.status(200).send(response_payload);
      }
    }
  } catch (error) {
    const response_payload = {
      result: null,
      error: error.name,
      errorData: error.message,
    };
    res.status(500).send(response_payload);
  }
};

exports.Invoke = async (req, res) => {
  // fe mengirim data : fcn, peers, chaincodeName, channelName, args
  try {
    logger.debug('==================== INVOKE ON CHAINCODE ==================');
    var peers = req.body.peers;
    var chaincodeName = req.params.chaincodeName;
    var channelName = req.params.channelName;
    var fcn = req.body.fcn;
    var args = req.body.args;
    var transient = req.body.transient;
    console.log(`Transient data is ;${transient}`);
    logger.debug('channelName  : ' + channelName);
    logger.debug('chaincodeName : ' + chaincodeName);
    logger.debug('fcn  : ' + fcn);
    logger.debug('args  : ' + args);

    if (!chaincodeName) {
      res.json(getErrorMessage("'chaincodeName'"));
      return;
    }
    if (!channelName) {
      res.json(getErrorMessage("'channelName'"));
      return;
    }
    if (!fcn) {
      res.json(getErrorMessage("'fcn'"));
      return;
    }
    if (!args) {
      res.json(getErrorMessage("'args'"));
      return;
    }

    // masuk ek fungsi post transaksi ---> invoke
    let message = await invoke.invokeTransaction(channelName, chaincodeName, fcn, args, req.userName, req.role, transient);
    console.log(`message result is :`);

    console.log(message);
    //console.log(`${typeof message}`)

    if (typeof message == 'string') {
      const response_payload = {
        result: message,
        error: 'Error in chaincode',
      };

      res.status(500).send(response_payload);
    } else if (typeof message == 'object') {
      const response_payload = {
        result: message,
        error: null,
        errorData: null,
      };

      // TODO: change to 201
      res.status(201).send(response_payload);
    }

    // const response_payload = {
    //   result: message,
    //   error: null,
    //   errorData: null,
    // };
    // res.status(200).send(response_payload);
  } catch (error) {
    const response_payload = {
      result: null,
      error: error.name,
      errorData: error.message,
    };
    res.status(500).send(response_payload);
  }
};
