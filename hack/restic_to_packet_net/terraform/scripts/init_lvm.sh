#!/usr/bin/env bash

apt-get update -y
apt-get install -y parted lvm2

systemctl start lvm2-monitor

partitions=""
for dev in ${@}; do
  parted -s  ${dev} mklabel gpt
  parted -s  ${dev} unit mib mkpart primary 1 100% set 1 lvm on
  pvcreate ${dev}1
  partitions=$(echo ${partitions} ${dev}1)
done

vgcreate vgdata ${partitions}
lvcreate -L 21.8T -i 12 -n lvdata vgdata

mkfs.ext4 /dev/vgdata/lvdata

mkdir /data
mount /dev/vgdata/lvdata /data