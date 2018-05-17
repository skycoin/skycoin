const PROXY_CONFIG = {
  "/api/*": {
    "target": "http://127.0.0.1:46420",
    "secure": false,
    "logLevel": "debug",
    "bypass": function (req) {
      req.headers["host"] = '127.0.0.1:46420';
      req.headers["referer"] = 'http://127.0.0.1:46420';
      req.headers["origin"] = 'http://127.0.0.1:46420';
    }
  }
};

module.exports = PROXY_CONFIG;
