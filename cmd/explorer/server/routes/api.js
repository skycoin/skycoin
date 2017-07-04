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

// address uxouts!
router.get('/address', (req, res) => {
  axios.get(`${API}/explorer/transactions?address=`+req.query.address)
  .then(blocks => {
  res.status(200).json(blocks.data);
})
.catch(error => {
  res.status(500).send(error)
});
});

// address uxouts!
router.get('/uxout', (req, res) => {
  axios.get(`${API}/uxout?uxid=`+req.query.uxid)
  .then(blocks => {
  res.status(200).json(blocks.data);
})
.catch(error => {
  res.status(500).send(error)
});
});


// address uxouts!
router.get('/transaction', (req, res) => {
  axios.get(`${API}/transaction?txid=`+req.query.txid)
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

// Get the block details
router.get('/currentBalance', (req, res) => {
  if(req.query.address){
  axios.get(`${API}/outputs?addrs=`+req.query.address)
  .then(blocks => {
    res.status(200).json(blocks.data);
})
.catch(error => {
    res.status(500).send(error)
});
}

});

// Get the block details
router.get('/coinSupply', (req, res) => {
  axios.get("http://127.0.0.1:6420/explorer/getEffectiveOutputs")
  .then(blocks => {
    res.status(200).json(blocks.data);
})
.catch(error => {
    res.status(500).send(error)
});
});



/*
 http.HandleFunc("/api/blocks", getBlocks)
 http.HandleFunc("/api/blockchain/metadata", getBlockChainMetaData)
 http.HandleFunc("/api/address", getAddress)
 http.HandleFunc("/api/currentBalance",getCurrentBalance)
 http.HandleFunc("/api/uxout", getUxID)
 http.HandleFunc("/api/transaction", getTransaction)
 http.HandleFunc("/api/block", getBlock)
 http.HandleFunc("/api/coinSupply", getSupply)
 */

module.exports = router;
