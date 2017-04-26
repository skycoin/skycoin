go build
gotty -w -p 9999 --reconnect ./cli
rm ./cli
