const mongoose = require('mongoose');
const { stringify } = require('nodemon/lib/utils');

const userSchema = mongoose.Schema(
  {
    userName: {
      type: String,
      required: true,
      unique: true,
    },
    password: {
      type: String,
      required: true,
    },
    email: {
      type: String,
      required: true,
      unique: true,
    },
    namaLengkap: {
      type: String,
      required: true,
    },
    noTelp: {
      type: Number,
      required: true,
    },
    tglLahir: {
      type: String,
      required: true,
    },
    nik: {
      type: Number,
      required: true,
      unique: true,
    },
    role: {
      type: Number,
      required: true,
    },
    alamat: {
      type: String,
      required: true,
    },
  },
  { timestamps: true }
);

const model = mongoose.model('userSchema', userSchema);
module.exports = model;
