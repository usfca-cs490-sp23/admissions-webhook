CANAME=ca
DNS_Name=the-captains-hook.default.svc

# generate aes encrypted private key
openssl genrsa -aes256 -out $CANAME.key 4096

# create certificate, for 1 year
openssl req -x509 -new -nodes -key $CANAME.key -sha256 -days 365 -out $CANAME.crt -subj '/CN=TestCa CA/C=US/ST=California/L=San Francisco/O=USF'

# assuming linux behaviour
#sudo cp $CANAME.crt /etc/pki/ca-trust/source/anchors/$CANAME.crt
#sudo update-ca-trust

MYCERT=server
openssl req -new -nodes -out $MYCERT.csr -newkey rsa:4096 -keyout $MYCERT.key -subj '/CN=TestServ CA/C=US/ST=California/L=San Francisco/O=USF'

# create a v3 ext file for SAN properties
cat > $MYCERT.v3.ext << EOF
authorityKeyIdentifier=keyid,issuer
basicConstraints=CA:FALSE
keyUsage = digitalSignature, nonRepudiation, keyEncipherment, dataEncipherment
subjectAltName = @alt_names
[alt_names]
DNS.1 = $DNS_Name
DNS.2 = the-captains-hook1.default.svc
IP.1 = 192.168.1.1
IP.2 = 192.168.2.1
EOF
openssl x509 -req -in $MYCERT.csr -CA $CANAME.crt -CAkey $CANAME.key -CAcreateserial -out $MYCERT.crt -days 730 -sha256 -extfile $MYCERT.v3.ext

# joined the two scripts together here

# takes the server cert and server key and makes the secret yaml file
echo
echo ">> Generating kube secrets..."
kubectl create secret tls the-captains-hook-tls \
  --cert=server.crt \
  --key=server.key \
  --dry-run=client -o yaml \
  > ./pkg/webhook/deploy-rules/webhook.tls.secret.yaml
#   ^ Go up to pkg, then into webhook, then into deploy-rules and write the file

echo

rm ca.key ca.srl server.crt server.csr server.key server.v3.ext

echo ">> MutatingWebhookConfiguration caBundle:"
cat ca.crt | base64 | fold > ./pkg/tls/cab64.crt
rm ca.crt
