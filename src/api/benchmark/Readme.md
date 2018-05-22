
[List of benchmark tools](https://gist.github.com/denji/8333630)



**Compare**
<table class="table table-bordered table-striped table-condensed">
   <tr>
      <td>tool</td>
      <td>language</td>
      <td>embed in Go code</td>
      <td>HTTP methods </td>
      <td>keep-alive</td>
      <td>https</td>
      <td>http/2</td>
      <td>multi-target</td>
      <td>result-show</td>
      <td>RPS</td>
   </tr>
   <tr>
     <td>[wrk](https://github.com/wg/wrk)</td>
     <td> C, Lua</td>
     <td>GET</td>
     <td>NO</td>
     <td>NO</td>
     <td>NO</td>
     <td>NO</td>
     <td>NO</td>
     <td>standard output</td>
     <td>~15000</td>
   </tr>
   <tr>
      <td>[vegeta](https://github.com/tsenart/vegeta)</td>
      <td>Go</td>
      <td>ALL</td>
      <td>YES</td>
      <td>YES</td>
      <td>YES</td>
      <td>YES</td>
      <td>NO</td>
      <td>Go, js+html5，standard output</td>
      <td> ~7000 - reached files count limitation</td>
   </tr>
   <tr>
     <td>[hey](https://github.com/rakyll/hey)</td>
     <td>Go</td>
     <td>GET, POST, PUT, DELETE, HEAD, OPTIONS</td>
     <td>NO</td>
     <td>YES</td>
     <td>NO</td>
     <td>YES</td>
     <td>NO</td>
     <td>standard output</td>
     <td>~11000</td>
   </tr>
   <tr>
     <td>[bombarider](https://github.com/codesenberg/bombardier)</td>
     <td>Go</td>
     <td>ALL</td>
     <td>NO</td>
     <td>NO</td>
     <td>YES</td>
     <td>YES</td>
     <td>NO</td>
     <td>standard output</td>
     <td>~15000</td>
   </tr>
   <tr>
      <td>[sniper](https://github.com/btfak/sniper)</td>
      <td>Go</td>
      <td/>GET, POST</td>
      <td>YES</td>
      <td>YES</td>
      <td>NO</td>
      <td>YES</td>
      <td>NO</td>
      <td>js+html5，standard output</td>
      <td>~1600</td>
   </tr>
   <tr>
     <td>[gobench](https://github.com/cmpxchg16/gobench)</td>
     <td>Go</td>
     <td>GET, POST</td>
     <td>NO</td>
     <td>YES</td>
     <td>NO</td>
     <td>NO</td>
     <td>NO</td>
     <td>standard output</td>
     <td>~12500</td>
  </tr>
</table>


I researched usage of several benchmark tools from this [list](https://gist.github.com/denji/8333630).
The most part of them is just CLI tools with basic functionality. [ab](https://en.wikipedia.org/wiki/ApacheBench) is a typical example.
Some of CLI tools are able to have a deal with GET HTTP requests. So, it does not allow us to use them for API tests.
Usually, CLI tools use plain text report format. It is not very convenient to parse this kind of reports.
There are too complicated tools, such as [tsung](https://github.com/processone/tsung)
which requires the availability of Erlang VM and (yandex-tank)[https://github.com/yandex/yandex-tank] which is a scalable system with a lot of config params.
The last one can be used for periodic stress testing of large distributed systems.

Another kind of tools is SAAS for benchmark testing. A good example is [Nginx Amplify](https://www.nginx.com/blog/setting-up-nginx-amplify-in-10-minutes/):
"NGINX Amplify is a SaaS monitoring tool for NGINX and underlying system components. NGINX Amplify is free to use for up to 5 monitored servers."
The more optimal tool for API benchmark testing is [vegeta](https://github.com/tsenart/vegeta).
The tool is written on Go and can be embedded in skycoin code.
I reached a limit on ~7000 RPS and the reason was the limitation in opened files on my Mac.

**Usage (Library)**
```go
package main

import (
  "fmt"
  "time"

  vegeta "github.com/tsenart/vegeta/lib"
)

func main() {
  rate := uint64(8000) // per second
  	duration := 10 * time.Second
  	targeter := vegeta.NewStaticTargeter(vegeta.Target{
  		Method: "GET",
  		URL:    "http://localhost:6420/api/v1/wallets/folderName",
  	})
  	attacker := vegeta.NewAttacker()

  	var metrics vegeta.Metrics
  	for res := range attacker.Attack(targeter, rate, duration) {
  		if res.Error != "" {
  			log.Printf("getWalletDir. err: %v", res.Error)
  		}
  		metrics.Add(res)
  	}
  	metrics.Close()
  	reporter := vegeta.NewTextReporter(&metrics)
  	reporter.Report(os.Stdout)
}
```

##### `text`
```console
Requests      [total, rate]             1200, 120.00
Duration      [total, attack, wait]     10.094965987s, 9.949883921s, 145.082066ms
Latencies     [mean, 50, 95, 99, max]   113.172398ms, 108.272568ms, 140.18235ms, 247.771566ms, 264.815246ms
Bytes In      [total, mean]             3714690, 3095.57
Bytes Out     [total, mean]             0, 0.00
Success       [ratio]                   55.42%
Status Codes  [code:count]              0:535  200:665
Error Set:
Get http://localhost:6060: dial tcp 127.0.0.1:6060: connection refused
Get http://localhost:6060: read tcp 127.0.0.1:6060: connection reset by peer
Get http://localhost:6060: dial tcp 127.0.0.1:6060: connection reset by peer
Get http://localhost:6060: write tcp 127.0.0.1:6060: broken pipe
Get http://localhost:6060: net/http: transport closed before response was received
Get http://localhost:6060: http: can't write HTTP request on broken connection
```

##### `json`
```json
{
  "latencies": {
    "total": 237119463,
    "mean": 2371194,
    "50th": 2854306,
    "95th": 3478629,
    "99th": 3530000,
    "max": 3660505
  },
  "bytes_in": {
    "total": 606700,
    "mean": 6067
  },
  "bytes_out": {
    "total": 0,
    "mean": 0
  },
  "earliest": "2015-09-19T14:45:50.645818631+02:00",
  "latest": "2015-09-19T14:45:51.635818575+02:00",
  "end": "2015-09-19T14:45:51.639325797+02:00",
  "duration": 989999944,
  "wait": 3507222,
  "requests": 100,
  "rate": 101.01010672380401,
  "success": 1,
  "status_codes": {
    "200": 100
  },
  "errors": []
}
```





