const express = require('express');
const bodyParser = require('body-parser');
const mongoose = require('mongoose');
const cors = require('cors');
const user = require('./model/User');
const bcrypt = require('bcrypt');
const jwt = require('jsonwebtoken');
const helper = require('../app/helper');

require('dotenv').config();

const app = express();
let refreshTokens = [];

// sambungan ke database
const db = require('./model/index');
const { json } = require('express/lib/response');
// const database = process.env.MONGO_URI || "mongodb://localhost:27017/mangojs"
const port = process.env.PORT || 5050;

// pengecekan berjalan di port . . .
app.listen(port, function () {
  console.log('server berjalan  di port ' + port);
});

app.use(cors());
app.use(
  bodyParser.json({
    extended: true,
    limit: '50mb',
  })
);
app.use(
  bodyParser.urlencoded({
    extended: true,
    limit: '50mb',
  })
);

// register
app.post('/api/registrasi/', async (req, res) => {
  // get data nama lengkap, username, password, role
  var namaLengkap = req.body.namaLengkap;
  var userName = req.body.userName;
  var password = req.body.password;
  var role = req.body.role;

  if (!namaLengkap) {
    res.status(400).json(getErrorMessage("'namaLengkap'"));
    return;
  }
  if (!userName) {
    res.status(400).json(getErrorMessage("'userName'"));
    return;
  }
  if (!password) {
    res.status(400).json(getErrorMessage("'password'"));
    return;
  }
  if (!role) {
    res.status(400).json(getErrorMessage("'role'"));
    return;
  }

  // wallet input. register ke sistem blockchain --> helper
  let pesan = await helper.getRegisteredUser(userName, role, true);

  // fungsi buat register --> helper
  let response = await helper.registerUserMongo(req, res);

  try {
    console.log('User success !!', response);
  } catch (err) {
    console.log(err);
    return res.status(400).json({ status: 'error' });
  }

  res.json({ status: 'ok' });
  res.status(201).json({ message: pesan, data: response });
});

// login
app.post('/api/login', async (req, res) => {
  const { userName, password, role } = req.body;

  // wallet input. login ke sistem blockchain --> helper
  let isUserRegistered = await helper.isUserRegistered(userName, role);

  //   const { userName, password } = req.body;
  if (isUserRegistered) {
    // fungsi buat loginnnya --> helper
    let data = await helper.loginUserMongo(req, res);

    const id = { id: data._id };
    if (await bcrypt.compare(password, data.password)) {
      const accessToken = generateAccessToken(id);
      const refreshToken = jwt.sign(id, process.env.REFRESH_TOKEN_SECRET);
      // refreshToken.push(refreshTokens)

      return res.status(201).json({ status: 'selamat datang ' + data.userName, accessToken: accessToken, refreshToken: refreshToken });
    }

    const refreshToken = accessToken;
  }

  //   const data = await user.findOne({ userName }).lean();
  //   // res.json(data);

  //   if (!data) {
  //     return res.status(404).json({ status: 'error', error: 'invalid username/password' });
  //   }
});

//logout
app.delete('/api/logout/', async (req, res) => {
  refreshTokens = await refreshTokens.filter((token) => token !== req.body.token);
  res.json({ status: 'error', error: req.body.taken });
});

//get profile user login
app.get('/api/profile/:userName', authenticateToken, async (req, res) => {
  const { userName } = req.params;

  user
    .find({
      userName,
    })

    .then((result) => {
      res.send(result);
    })
    .catch((err) => {
      res.send(err);
    });
});

//request token
app.post('/api/token', async (req, res) => {
  const refreshToken = req.body.token;
  if (refreshToken === null) return res.json({ status: 'user belum login' });
  if (!refreshTokens.includes(refreshToken)) return res.json({ status: 'akses ditolak' });
  jwt.verify(refreshToken, process.env.REFRESH_TOKEN_SECRET, (err, id) => {
    if (err) return res.sendStatus(403);
    const accessToken = generateAccessToken(id);
    res.json({ accessToken: accessToken });
  });
});

//autentikasi token jwt
function authenticateToken(req, res, next) {
  const authHeader = req.headers['authorization'];
  const token = authHeader && authHeader.split(' ')[1];
  if (token === null) {
    return res.status(401).json({ message: 'user not login' });
  }

  jwt.verify(token, process.env.ACCESS_TOKEN_SECRET, (err, id) => {
    if (err) return res.status(401).json({ message: 'token not valid' });
    req.id = id;
    next();
  });
}

function generateAccessToken(id) {
  return jwt.sign(id, process.env.ACCESS_TOKEN_SECRET, { expiresIn: '3600s' });
}

//penghubung ke routes
const allRegister = require('./routes/auth');
app.use('/api', authenticateToken, allRegister);
