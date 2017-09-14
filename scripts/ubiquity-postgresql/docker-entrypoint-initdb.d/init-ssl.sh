#!/usr/bin/env bash
set -e

echo "Configuring Postgres for SSL!";

if [ -z "$POSTGRES_USER" ]; then
  export POSTGRES_USER="postgres";
fi

if [ -z "$POSTGRES_EMAIL" ]; then
  export POSTGRES_EMAIL="user@test.com";
fi

export SSLDIR="/var/lib/postgresql/ssl"

# Update HBA to require SSL and Client Cert auth
head -n -1 /var/lib/postgresql/data/pg_hba.conf > /tmp/pg_hba.conf
echo "hostssl all all all password" >> /tmp/pg_hba.conf
mv /tmp/pg_hba.conf /var/lib/postgresql/data/pg_hba.conf

if [ ! -d $SSLDIR ]
then
    # Create SSL directory
    mkdir $SSLDIR
    chown postgres.postgres $SSLDIR
fi

# Create SSL certificates
cd $SSLDIR

if [ ! -s "$SSLDIR/server.crt" ]
then
    # root CA
    openssl req -new -x509 -nodes -out root.crt -keyout root.key -newkey rsa:4096 -sha512 -subj /CN=TheRootCA
    chown postgres.postgres root.key
    chmod 600 root.key

    # Server certificate
    openssl req -new -out server.req -keyout server.key -nodes -newkey rsa:4096 -subj "/CN=$( hostname )/emailAddress=$POSTGRES_EMAIL"
    openssl x509 -req -in server.req -CAkey root.key -CA root.crt -set_serial $RANDOM -sha512 -out server.crt

    chown postgres.postgres server.key
    chmod 600 server.key
fi

# edit the configuration files

sed -i 's/#ssl/ssl/g' /var/lib/postgresql/data/postgresql.conf
sed -i 's/ssl \= off/ssl \= on/g' /var/lib/postgresql/data/postgresql.conf
sed -i "s/ssl_cert_file = 'server.crt'/ssl_cert_file = '\/var\/lib\/postgresql\/ssl\/server.crt'/g" /var/lib/postgresql/data/postgresql.conf
sed -i "s/ssl_key_file = 'server.key'/ssl_key_file = '\/var\/lib\/postgresql\/ssl\/server.key'/g" /var/lib/postgresql/data/postgresql.conf

echo "SSL Configuration - Done!"

