

Endpoint Service
====

An endpoint is a "connection", like TCP/IP or UDP
- an endpoint accepts packets (length prefixed messages) to be sent to the destination
- an endpoint receives packets (legnth prefixed messages) from the destination

Route Setup
===

- feed in 33 byte cipher.PubKey to connect to
- attempts to find a route to the Public key, through the node
- attempts to maintain routes and connections to the end point
- will reconnect to end-point if route fails

Behavior
=== 

- accepts length prefixed messages and may block
- has a queue of incoming messages