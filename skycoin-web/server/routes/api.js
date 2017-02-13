const express = require('express');
const router = express.Router();

const axios = require('axios');
const API = 'http://127.0.0.1:6420';

router.get('/', (req, res) => {
  res.send('api works');
});

// Get all blocks
router.get('/blocks', (req, res) => {

  axios.get(`${API}/blocks?start=`+req.query.start+`&end=`+req.query.end)
    .then(blocks => {
      res.status(200).json(blocks.data);
    })
    .catch(error => {
      res.status(500).send(error)
    });
});

// Get the block metadata
router.get('/blockchain/metadata', (req, res) => {
  axios.get(`${API}/blockchain/metadata`)
  .then(blocks => {
  res.status(200).json(blocks.data);
})
.catch(error => {
  res.status(500).send(error)
});
});

// Get the block details
router.get('/block', (req, res) => {
  if(req.query.hash){
  axios.get(`${API}/block?hash=`+req.query.hash)
    .then(blocks => {
    res.status(200).json(blocks.data);
})
.catch(error => {
    res.status(500).send(error)
});
}

});

module.exports = router;
