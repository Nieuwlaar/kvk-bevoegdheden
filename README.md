KVK Extract (Golang lib)
------

GO library that fetches an extract from the KVK Dataservice and extracts all fields from it, that are relevant for a 'Bevoegdheid' of 'Machtiging'.

To get access to the KVK Dataservice, you need a certificate (OV or EV). Assuming you have got that certificate in p12 format:
Extract private key PKCS1 format from p12 file
```
openssl pkcs12 -in yourP12File.pfx -nocerts -out private_key_PKCS1.pem
```

Extract certificate from p12 file
```
openssl pkcs12 -in yourP12File.pfx -clcerts -nokeys -out certificate.pem
```