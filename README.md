# cartman

> A library to help servers identify their users by managing
TLS client certificates


## Running

Import this library, then create a Store:

`cartman, err := NewStore("users")`

this will scan the *users* directory for folders containing certificates.
The accepted structure is as follows:

```
users/
├── me
│   └── cert.pem
└── test
    └── cert.pem
```

The certificate file must be a x509 certificate encoded em *pem* format.

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
openssl s_client --connect localhost:8080 -cert users/test/cert.pem -key users/test/cert.key -crlf -CAfile fullchain.pem
```

But you can use curl likewise:

```
curl -vv https://localhost:8080 -E users/test/cert.pem --key users/test/cert.key --cacert fullchain.pem
```

To see an example of a server check the [cartman_test.go](cartman_test.go) file.

