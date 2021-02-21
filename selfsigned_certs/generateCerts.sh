#! /bin/zsh

echo "Generating Server Pass Key"
openssl genrsa -aes256 -passout pass:mockingServer -out mocking.pass.key 4096

echo "Generating Server Key"
openssl rsa -passin pass:mockingServer -in mocking.pass.key -out mockingServer.key

echo "Clean up Server Pass Key"
rm mocking.pass.key

echo "Generate CSR"
openssl req -new -key mockingServer.key -out mockingServer.csr

echo "Generate Self Signed Server Certificate"
openssl x509 -req -sha256 -days 365 -in mockingServer.csr -signkey mockingServer.key -out mockingServer.crt