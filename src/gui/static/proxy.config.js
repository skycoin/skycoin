const PROXY_CONFIG = {
  "/api/*": {
    "target": "http://127.0.0.1:6420",
    "secure": false,
    "logLevel": "debug",
    "bypass": function (req) {
      req.headers["host"] = '127.0.0.1:6420';
      req.headers["referer"] = 'http://127.0.0.1:6420';
      req.headers["origin"] = 'http://127.0.0.1:6420';
    }
},
  "/teller/*": {
    "target": "http://127.0.0.1:7071",
    "pathRewrite": {
      "^/teller" : "api/"
    },
    "secure": true,
    "logLevel": "debug"
  },
  "/hw-daemon/*": {
    "target": "http://127.0.0.1:9510",
    "pathRewrite": {'^/hw-daemon' : ''},
    "secure": true,
    "logLevel": "debug",
    "bypass": function (req) {
      req.headers["host"] = '127.0.0.1:9510';
      req.headers["referer"] = 'http://127.0.0.1:9510';
      req.headers["origin"] = 'http://127.0.0.1:9510';
    }
  }
};

module.exports = PROXY_CONFIG;
