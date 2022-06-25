'use strict';

var { Gateway, Wallets } = require('fabric-network');
const path = require('path');
const FabricCAServices = require('fabric-ca-client');
const fs = require('fs');
const User = require('../models/user.model');
const jwt = require('jsonwebtoken');
const bcryptjs = require('bcryptjs');
const bcrypt = require('bcrypt');

const util = require('util');

// identifikasi role dari user ---> invoke
const getCCP = async (org) => {
  let ccpPath;
  if (org == 'Penangkar') {
    ccpPath = path.resolve(__dirname, '..', 'config', 'connection-penangkar.json');
  } else if (org == 'Petani') {
    ccpPath = path.resolve(__dirname, '..', 'config', 'connection-petani.json');
  } else if (org == 'Pengumpul') {
    ccpPath = path.resolve(__dirname, '..', 'config', 'connection-pengumpul.json');
  } else if (org == 'Pedagang') {
    ccpPath = path.resolve(__dirname, '..', 'config', 'connection-pedagang.json');
  } else return null;
  const ccpJSON = fs.readFileSync(ccpPath, 'utf8');
  const ccp = JSON.parse(ccpJSON);
  return ccp;
};

// certificate authority
const getCaUrl = async (org, ccp) => {
  let caURL;
  if (org == 'Penangkar') {
    caURL = ccp.certificateAuthorities['ca.penangkar.example.com'].url;
  } else if (org == 'Petani') {
    caURL = ccp.certificateAuthorities['ca.petani.example.com'].url;
  } else if (org == 'Pengumpul') {
    caURL = ccp.certificateAuthorities['ca.pengumpul.example.com'].url;
  } else if (org == 'Pedagang') {
    caURL = ccp.certificateAuthorities['ca.pedagang.example.com'].url;
  } else return null;
  return caURL;
};

// mengambil path/jalur url saat ini ditambah dengan role nya. misal sc/channel/penangkar-wallet --> invoke
const getWalletPath = async (org) => {
  let walletPath;
  if (org == 'Penangkar' || org == 'penangkar') {
    walletPath = path.join(process.cwd(), 'penangkar-wallet');
  } else if (org == 'Petani' || org == 'petani') {
    walletPath = path.join(process.cwd(), 'petani-wallet');
  } else if (org == 'Pengumpul' || org == 'pengumpul') {
    walletPath = path.join(process.cwd(), 'pengumpul-wallet');
  } else if (org == 'Pedagang' || org == 'pengumpul') {
    walletPath = path.join(process.cwd(), 'pedagang-wallet');
  } else return null;
  return walletPath;
};

// TODO: make affiliation
const getAffiliation = async (org) => {
  if (org == 'Penangkar') {
    return 'org1.department1';
  } else if (org == 'Petani') {
    return 'org1.department1';
  } else if (org == 'Pengumpul') {
    return 'org1.department1';
  } else if (org == 'Pedagang') {
    return 'org1.department1';
  } else return null;
};

// register dengan blockchain --> register
const getRegisteredUser = async (userName, role, isJson) => {
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

  let ccp = await getCCP(roleAktor);

  const caURL = await getCaUrl(roleAktor, ccp);

  const ca = new FabricCAServices(caURL);

  const walletPath = await getWalletPath(roleAktor);

  const wallet = await Wallets.newFileSystemWallet(walletPath);
  console.log(`Wallet path: ${walletPath}`);

  const userIdentity = await wallet.get(userName);
  if (userIdentity) {
    console.log(`An identity for the user ${userName} already exists in the wallet`);
    var response = {
      success: true,
      message: userName + ' enrolled Successfully',
    };
    return response;
  }

  // Check to see if we've already enrolled the admin user.
  let adminIdentity = await wallet.get('admin');
  if (!adminIdentity) {
    console.log('An identity for the admin user "admin" does not exist in the wallet');
    await enrollAdmin(roleAktor, ccp);
    adminIdentity = await wallet.get('admin');
    console.log('Admin Enrolled Successfully');
  }

  // build a user object for authenticating with the CA
  const provider = wallet.getProviderRegistry().getProvider(adminIdentity.type);
  const adminUser = await provider.getUserContext(adminIdentity, 'admin');
  let secret;
  try {
    // Register the user, enroll the user, and import the new identity into the wallet.
    secret = await ca.register(
      {
        affiliation: await getAffiliation(roleAktor),
        enrollmentID: userName,
        role: 'client',
      },
      adminUser
    );
    // const secret = await ca.register({ affiliation: 'org1.department1', enrollmentID: userName, role: 'client', attrs: [{ name: 'role', value: 'approver', ecert: true }] }, adminUser);
  } catch (error) {
    return error.message;
  }

  const enrollment = await ca.enroll({
    enrollmentID: userName,
    enrollmentSecret: secret,
  });
  // const enrollment = await ca.enroll({ enrollmentID: userName, enrollmentSecret: secret, attr_reqs: [{ name: 'role', optional: false }] });

  let x509Identity;
  if (roleAktor == 'Penangkar') {
    x509Identity = {
      credentials: {
        certificate: enrollment.certificate,
        privateKey: enrollment.key.toBytes(),
      },
      mspId: 'PenangkarMSP',
      type: 'X.509',
    };
  } else if (roleAktor == 'Petani') {
    x509Identity = {
      credentials: {
        certificate: enrollment.certificate,
        privateKey: enrollment.key.toBytes(),
      },
      mspId: 'PetaniMSP',
      type: 'X.509',
    };
  } else if (roleAktor == 'Pengumpul') {
    x509Identity = {
      credentials: {
        certificate: enrollment.certificate,
        privateKey: enrollment.key.toBytes(),
      },
      mspId: 'PengumpulMSP',
      type: 'X.509',
    };
  } else if (roleAktor == 'Pedagang') {
    x509Identity = {
      credentials: {
        certificate: enrollment.certificate,
        privateKey: enrollment.key.toBytes(),
      },
      mspId: 'PedagangMSP',
      type: 'X.509',
    };
  }

  await wallet.put(userName, x509Identity);
  console.log(`Successfully registered and enrolled admin user ${userName} and imported it into the wallet`);

  console.log(`${userName} has been successfully enrolled`);

  var response = {
    success: true,
    message: userName + ' enrolled Successfully',
  };
  return response;
};

// login dengan blockchain --> login
const isUserRegistered = async (userName, role) => {
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

  const walletPath = await getWalletPath(roleAktor);
  const wallet = await Wallets.newFileSystemWallet(walletPath);
  console.log(`Wallet path: ${walletPath}`);

  const userIdentity = await wallet.get(userName);
  if (userIdentity) {
    console.log(`An identity for the user ${userName} exists in the wallet`);
    return true;
  }
  return false;
};

const getCaInfo = async (org, ccp) => {
  let caInfo;
  if (org == 'Penangkar') {
    caInfo = ccp.certificateAuthorities['ca.penangkar.example.com'];
  } else if (org == 'Petani') {
    caInfo = ccp.certificateAuthorities['ca.petani.example.com'];
  } else if (org == 'Pengumpul') {
    caInfo = ccp.certificateAuthorities['ca.pengumpul.example.com'];
  } else if (org == 'Pedagang') {
    caInfo = ccp.certificateAuthorities['ca.pedagang.example.com'];
  } else return null;
  return caInfo;
};

const enrollAdmin = async (org, ccp) => {
  console.log('calling enroll Admin method');

  try {
    const caInfo = await getCaInfo(org, ccp); //ccp.certificateAuthorities['ca.org1.example.com'];
    const caTLSCACerts = caInfo.tlsCACerts.pem;
    const ca = new FabricCAServices(caInfo.url, { trustedRoots: caTLSCACerts, verify: false }, caInfo.caName);

    // Create a new file system based wallet for managing identities.
    const walletPath = await getWalletPath(org); //path.join(process.cwd(), 'wallet');
    const wallet = await Wallets.newFileSystemWallet(walletPath);
    console.log(`Wallet path: ${walletPath}`);

    // Check to see if we've already enrolled the admin user.
    const identity = await wallet.get('admin');
    if (identity) {
      console.log('An identity for the admin user "admin" already exists in the wallet');
      return;
    }

    // Enroll the admin user, and import the new identity into the wallet.
    const enrollment = await ca.enroll({
      enrollmentID: 'admin',
      enrollmentSecret: 'adminpw',
    });
    let x509Identity;
    if (org == 'Penangkar') {
      x509Identity = {
        credentials: {
          certificate: enrollment.certificate,
          privateKey: enrollment.key.toBytes(),
        },
        mspId: 'PenangkarMSP',
        type: 'X.509',
      };
    } else if (org == 'Petani') {
      x509Identity = {
        credentials: {
          certificate: enrollment.certificate,
          privateKey: enrollment.key.toBytes(),
        },
        mspId: 'PetaniMSP',
        type: 'X.509',
      };
    } else if (org == 'Pengumpul') {
      x509Identity = {
        credentials: {
          certificate: enrollment.certificate,
          privateKey: enrollment.key.toBytes(),
        },
        mspId: 'PengumpulMSP',
        type: 'X.509',
      };
    } else if (org == 'Pedagang') {
      x509Identity = {
        credentials: {
          certificate: enrollment.certificate,
          privateKey: enrollment.key.toBytes(),
        },
        mspId: 'PedagangMSP',
        type: 'X.509',
      };
    }

    await wallet.put('admin', x509Identity);
    console.log('Successfully enrolled admin user "admin" and imported it into the wallet');
    return;
  } catch (error) {
    console.error(`Failed to enroll admin user "admin": ${error}`);
  }
};

const registerAndGetSecret = async (username, userOrg) => {
  let ccp = await getCCP(userOrg);

  const caURL = await getCaUrl(userOrg, ccp);
  const ca = new FabricCAServices(caURL);

  const walletPath = await getWalletPath(userOrg);
  const wallet = await Wallets.newFileSystemWallet(walletPath);
  console.log(`Wallet path: ${walletPath}`);

  const userIdentity = await wallet.get(username);
  if (userIdentity) {
    console.log(`An identity for the user ${username} already exists in the wallet`);
    var response = {
      success: true,
      message: username + ' enrolled Successfully',
    };
    return response;
  }

  // Check to see if we've already enrolled the admin user.
  let adminIdentity = await wallet.get('admin');
  if (!adminIdentity) {
    console.log('An identity for the admin user "admin" does not exist in the wallet');
    await enrollAdmin(userOrg, ccp);
    adminIdentity = await wallet.get('admin');
    console.log('Admin Enrolled Successfully');
  }

  // build a user object for authenticating with the CA
  const provider = wallet.getProviderRegistry().getProvider(adminIdentity.type);
  const adminUser = await provider.getUserContext(adminIdentity, 'admin');
  let secret;
  try {
    // Register the user, enroll the user, and import the new identity into the wallet.
    secret = await ca.register(
      {
        affiliation: await getAffiliation(userOrg),
        enrollmentID: username,
        role: 'client',
      },
      adminUser
    );
    // const secret = await ca.register({ affiliation: 'org1.department1', enrollmentID: username, role: 'client', attrs: [{ name: 'role', value: 'approver', ecert: true }] }, adminUser);
  } catch (error) {
    return error.message;
  }

  var response = {
    success: true,
    message: username + ' enrolled Successfully',
    secret: secret,
  };
  return response;
};

// fungsi semua register --> register
const registerUserMongo = async (req, res) => {
  const { userName, password: plainTextPassword, email, namaLengkap, noTelp, tglLahir, nik, role, alamat } = req.body;

  if (!userName || typeof userName !== 'string') {
    return res.status(400).json({ status: 'error', error: 'invalid username' });
  }
  if (!plainTextPassword || typeof plainTextPassword !== 'string') {
    return res.status(400).json({ status: 'error', error: 'invalid password' });
  }
  if (plainTextPassword.length < 5) {
    return res.status(400).json({ status: 'error', error: 'password must more than five character' });
  }
  if (!email || typeof email !== 'string') {
    return res.status(400).json({ status: 'error', error: 'invalid email' });
  }
  if (!namaLengkap || typeof namaLengkap !== 'string') {
    return res.status(400).json({ status: 'error', error: 'invalid namaLengkap' });
  }
  if (!noTelp || typeof noTelp !== 'string') {
    return res.status(400).json({ status: 'error', error: 'invalid noTelp' });
  }
  if (!tglLahir || typeof tglLahir !== 'string') {
    return res.status(400).json({ status: 'error', error: 'invalid tglLahir' });
  }
  if (!nik || typeof nik !== 'string') {
    return res.status(400).json({ status: 'error', error: 'invalid nik' });
  }
  if (nik.length !== 16) {
    return res.status(400).json({ status: 'error', error: 'nik is not found' });
  }
  if (!role || typeof role !== 'string') {
    return res.status(400).json({ status: 'error', error: 'invalid role' });
  }
  if (role > 4 || role < 1) {
    return res.status(400).json({ status: 'error', error: 'role is not found' });
  }
  if (!alamat || typeof alamat !== 'string') {
    return res.status(400).json({ status: 'error', error: 'invalid alamat' });
  }

  // enkripsi password
  const password = await bcrypt.hash(plainTextPassword, 10);
  const biodata = { userName, password, email, namaLengkap, noTelp, tglLahir, nik, role, alamat };
  const response = await User.create(biodata);

  return response;
};

// fungsi semua login --> login
const loginUserMongo = async (req, res) => {
  const { userName, password } = req.body;

  const data = await user.findOne({ userName }).lean();
  if (!data) {
    return res.status(404).json({ status: 'error', error: 'invalid username/password' });
  }
  return data;

  // if (userDB) {
  //   const passwordDB = await bcryptjs.compare(password, userDB.password);
  //   if (passwordDB) {
  //     // const data = {
  //     //   id: userDB._id,
  //     // };
  //     // const maxAge = 2 * 24 * 60 * 60; //2 hari
  //     // const options = {
  //     //   expiresIn: maxAge,
  //     // };
  //     //const token = await jwt.sign(data, process.env.JWT_SECRET, options);
  //     //res.cookie("jwt", token, { httpOnly: true, maxAge: maxAge * 1000 });

  //     console.log(`token in loginUserMongo ${token}`);
  //     console.log(`User in database:`);
  //     console.log(`${userDB}`);

  //     return userDB;
  //   } else {
  //     return false;
  //   }
  // } else {
  //   return false;
  // }
};

exports.getRegisteredUser = getRegisteredUser;

module.exports = {
  getCCP: getCCP,
  getWalletPath: getWalletPath,
  getRegisteredUser: getRegisteredUser,
  isUserRegistered: isUserRegistered,
  registerAndGetSecret: registerAndGetSecret,
  registerUserMongo: registerUserMongo,
  loginUserMongo: loginUserMongo,
};
