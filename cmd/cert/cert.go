// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Modifications Copyright 2014 Skycoin authors.

// +build ignore

// Generate a self-signed X.509 certificate for a TLS server. Outputs to
// certFile and keyFile and will overwrite existing files.

package main

import (
    "flag"
    "fmt"
    "github.com/skycoin/skycoin/src/util"
    "log"
    "os"
    "time"
)

var (
    host = flag.String("host", "",
        "Comma-separated hostnames and IPs to generate a certificate for")
    organization = flag.String("organization", "Acme Co",
        "Certification organization")
    validFrom = flag.String("start-date", "",
        "Creation date formatted as Jan 1 15:04:05 2011")
    validFor = flag.Duration("duration", 365*24*time.Hour,
        "Duration that certificate is valid for")
    isCA = flag.Bool("ca", false,
        "whether this cert should be its own Certificate Authority")
    rsaBits  = flag.Int("rsa-bits", 2048, "Size of RSA key to generate")
    certFile = flag.String("cert", "cert.pem",
        "Name of file to write certificate to")
    keyFile = flag.String("key", "key.pem",
        "Name of file to write key to")
)

func main() {
    flag.Parse()

    if len(*host) == 0 {
        log.Fatalf("Missing required -host parameter")
    }

    var err error
    var notBefore time.Time
    if len(*validFrom) == 0 {
        notBefore = util.Now()
    } else {
        notBefore, err = time.Parse("Jan 2 15:04:05 2006", *validFrom)
        if err != nil {
            fmt.Fprintf(os.Stderr, "Failed to parse creation date: %v\n", err)
            os.Exit(1)
        }
    }

    err = util.GenerateCert(*certFile, *keyFile, *host, *organization, *rsaBits,
        *isCA, notBefore, *validFor)
    if err == nil {
        fmt.Printf("Created %s and %s\n", *certFile, *keyFile)
    } else {
        fmt.Fprintln(os.Stderr, "Failed to create cert and key")
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}
