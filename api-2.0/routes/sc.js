const router = require('express').Router({ mergeParams: true });
const { Query, Invoke } = require('../controller/bmController');

router;

router.get('/channels/:channelName/chaincodes/:chaincodeName', Query);
router.post('/channels/:channelName/chaincodes/:chaincodeName', Invoke);

module.exports = router;
