const dbConfig = require('../config/dbConfig');

const db = {};
db.url = dbConfig.database;

db.posts = require('../model/User');

module.exports = db;
