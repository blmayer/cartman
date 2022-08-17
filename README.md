# cartman

> A library to help servers identify their users by managing
TLS client certificates


## Running

Import this library, then create a Store:

`cartman, err := NewFileStore("users")`

this will scan the *users* directory for files containing certificate
fingerprints. The accepted structure is as follows:

```
users/
├── me
└── test
```

The fingerprint must be the sha1 of a x509 certificate in DER encoding.

If no error occurs the certificates for *test* and *me* will be identified
if sent by them. You can check if the sent certificate is known using the
*tls* package:

```
connInfo := conn.ConnectionState()
userName, err := cartman.GetClientFromCert(connInfo.PeerCertificates[0])
```


### Testing a client

To successfuly connect to a server using cartman the client must send a
certificate, this example uses openssl directly:

```
openssl s_client --connect localhost:8080 -cert test/cert.pem -key test/cert.key -crlf -CAfile fullchain.pem
```

But you can use curl likewise:

```
curl -vv https://localhost:8080 -E test/cert.pem --key test/cert.key --cacert fullchain.pem
```

Creating certificate files can be done using openssl:

```
openssl req -x509 -newkey rsa:4096 -keyout test/cert.key -out test/cert.pem -days 36500 -nodes
```

To see an example of a server check the [cartman_test.go](cartman_test.go) file.

