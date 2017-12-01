# Address stats

This utility generate address and when analize statics. Histogram of distrubions are placed in file.


### Command line parametrs


```
-n <int> number of generated addresses
-addrs <bool> analize addresses
-pkeys <bool> analize public keys
-hashes <bool> analize public key hashes
```
### Example usage

```
go run address_stats.go -n=1000 -addrs=true
go run address_stats.go -n=100 -hashes=true -addrs=true -pkeys=true
```
