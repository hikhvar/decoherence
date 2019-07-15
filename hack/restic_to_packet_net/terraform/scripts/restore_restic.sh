#!/usr/bin/env bash

set -xeu

source env.conf

curl -sL https://github.com/restic/restic/releases/download/v0.9.5/restic_0.9.5_linux_amd64.bz2 -o restic_0.9.5_linux_amd64.bz2
bzip2 -d restic_0.9.5_linux_amd64.bz2
mv restic_0.9.5_linux_amd64 /usr/bin/restic
chmod +x /usr/bin/restic

apt-get install -y pip
pip install b2

b2 authorize-account ${B2_ACCOUNT_ID} ${B2_ACCOUNT_KEY}
mkdir -p /data/source
b2 sync ${BACKUP_TARGET} /data/source

mkdir -p /data/restore/${SNAPSHOT}
restic restore ${SNAPSHOT} -r /data/source --target /data/restore/${SNAPSHOT}