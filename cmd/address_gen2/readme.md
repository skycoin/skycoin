# Address generator

### Start scanning and other options

This utility have several flags:

```
-addrfile string
        command for changing addresses output file (default "addresses")
  -coin string
        address output type: sky/btc (default "sky)
  -infofile string
        create file with date of generation, seed, coin, number of keys generated
  -n int
        Number of addresses to generate (default 1)
  -secfile string
        command for file to write the secret keys
  -seed string
        Seed for deterministic key generation. Will use bip39 as the seed if not provided
  -strict bool
        Checks if input is space separated list of words. (default true)
```
### Example usage

```
go run address_gen.go -n=5 -seed="any lens mango buddy cigar uniform engage owner such nothing express load" 
go run address_gen.go -n=10 -info_out=info -sec_file=secretkeys -coin="btc"
```
