#!/usr/bin/env bash

set -xeu
SSH_KEY="../../ssh/id_rsa"
SNAPSHOT=latest

scp -i ${SSH_KEY} init_lvm.sh root@${IP}:/root/init_lvm.sh
scp -i ${SSH_KEY} restore_restic.sh root@${IP}:/root/restore_restic.sh
scp -i ${SSH_KEY} ../../../../decoherence root@${IP}:/root/decoherence
ssh -i ${SSH_KEY} root@${IP} chmod +x /root/restore_restic.sh /root/init_lvm.sh


#ssh -i ${SSH_KEY} root@${IP} /root/init_lvm.sh $@
cat > env.conf <<EOF
export B2_ACCOUNT_ID=${B2_ACCOUNT_ID}
export B2_ACCOUNT_KEY=${B2_ACCOUNT_KEY}
export BACKUP_TARGET=${BACKUP_TARGET}
export RESTIC_PASSWORD='${RESTIC_PASSWORD}'
export SNAPSHOT=${SNAPSHOT}
EOF

scp -i ${SSH_KEY} env.conf root@${IP}:/root/env.conf
ssh -i ${SSH_KEY} root@${IP} /root/restore_restic.sh

ssh -i ${SSH_KEY} /root/decoherence record --store /root/result.json --path /data/restore/${SNAPSHOT}/input
scp -i ${SSH_KEY} root@${IP}:/root/result.json .
cd ..
terraform destroy -auto-approve