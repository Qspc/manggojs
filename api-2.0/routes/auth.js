const { getAll, getOne, getProfile, editUser } = require('../controller/post');
const router = require('express').Router();

router;

//ambil semua data
router.get('/all', getAll);

// ambil satu data
router.get('/:_id', getOne);

// edit data profile
router.put('/:_id', editUser);

// get profile
router.get('/:userName', getProfile);

module.exports = router;
