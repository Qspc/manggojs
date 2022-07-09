const mongoose = require('mongoose');
require('dotenv').config();
// const database = process.env.MONGO_URI || 'mongodb://localhost:27017/';
const database = 'mongodb://localhost:27017/';

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

module.exports = { mongoose };
