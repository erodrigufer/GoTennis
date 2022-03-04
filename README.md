## TLS Certificate
Populate the `./tls/` directory with a self-signed certificate (or a certificate from a CA) and a private key for that certificate.

### Self-signed certificate
Use Go's `generate_cert.go` tool.
If using Mac OS and having installed go through brew, run the following command inside the tls folder:

```
go run /usr/local/Cellar/go/<version>/libexec/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost

with <version> being the version number of Go in your computer 
```



