const express = require('express');
const bodyParser = require('body-parser');
const cors = require('cors');
const bcrypt = require('bcrypt');
const jwt = require('jsonwebtoken');
const helper = require('./app/helper');
const user = require('./model/user');

const mongoose = require('mongoose');
require('dotenv').config();
const database = 'mongodb://localhost:27017/manggojs';

const app = express();
let refreshTokens = [];
const { json } = require('express/lib/response');

// pengecekan berjalan di port . . .
const port = 4000;
app.listen(port, function () {
  console.log('server berjalan  di port ' + port);
});

//connection mongoose local
mongoose
  .connect(database, {
    useNewUrlParser: true,
    useUnifiedTopology: true,
  })
  .then(() => console.log('connect mongoDB'))
  .catch((err) => console.log(err));

mongoose.connection.on('connected', () => {
  console.log(`${database} terkoneksi. . .`);
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
  // console.log(pesan)
  // console.log(response)

  try {
    console.log('User success !!', response);
  } catch (err) {
    console.log(err);
    return res.status(400).json({ status: 'error' });
  }

  res.json({ status: 'ok' });
  res.status(201).json({ 'message': pesan, 'data': response });
});

// login
app.post('/api/login', async (req, res) => {
  const { userName, password } = req.body;

  // fungsi buat loginnnya --> helper
  let data = await helper.loginUserMongo(req, res);

  // wallet input. login ke sistem blockchain --> helper
  let isUserRegistered = await helper.isUserRegistered(userName, data.role);

  // const refreshToken = accessToken;

  if (isUserRegistered) {
    const id = { id: data._id };
    if (await bcrypt.compare(password, data.password)) {
      const accessToken = generateAccessToken(id);
      const refreshToken = jwt.sign(id, process.env.REFRESH_TOKEN_SECRET);

      res.status(200).json({ status: 'selamat datang ' + data.userName + ' dengan role ' + data.role, Token: accessToken, refresh: refreshToken });
      // refreshToken.push(refreshTokens)
    } else {
      return res.status(404).json({ status: 'error', error: 'invalid password' });
    }
  } else {
    res.status(404).json({
      message: `User not found, Please register first.`,
    });
  }
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
const ScRegister = require('./routes/sc');
app.use('/api', authenticateToken, allRegister);
app.use('/sc', ScRegister);
