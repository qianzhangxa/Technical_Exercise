#!/bin/bash

# Exit on any error
set -e

KERNEL_URL="https://kernel.ubuntu.com/mainline/v5.15.100/amd64/linux-image-unsigned-5.15.100-0515100-generic_5.15.100-0515100.202303111431_amd64.deb"
KERNEL_DEB="linux-image-5.15.100.deb"
ROOTFS="rootfs.img"
MOUNT_DIR="mnt"
BUSYBOX_URL="https://busybox.net/downloads/binaries/1.31.0-defconfig-multiarch-musl/busybox-x86_64"
BUSYBOX_BIN="busybox"

# Download and extract the pre-built Linux kernel
fetch_kernel() {
  if [ ! -f "$KERNEL_DEB" ]; then
    echo "Downloading pre-built Linux kernel..."
    wget -O $KERNEL_DEB $KERNEL_URL
  fi

  echo "Extracting the kernel..."
  dpkg-deb -x $KERNEL_DEB ./kernel_extracted
  KERNEL_IMAGE=$(find ./kernel_extracted -name 'vmlinuz*' | head -n 1)  # Get the vmlinuz file
}

# Download the BusyBox binary
fetch_busybox() {
  if [ ! -f "$BUSYBOX_BIN" ]; then
    echo "Downloading BusyBox..."
    wget -O $BUSYBOX_BIN $BUSYBOX_URL
    chmod +x $BUSYBOX_BIN
  fi
}

# Build a minimal root filesystem
create_rootfs() {
  echo "Creating root filesystem image..."
  
  # Create an empty disk image
  dd if=/dev/zero of=$ROOTFS bs=1M count=64
  mkfs.ext4 $ROOTFS

  # Mount the disk image
  mkdir -p $MOUNT_DIR
  sudo mount $ROOTFS $MOUNT_DIR

  # Create essential directories
  sudo mkdir -p $MOUNT_DIR/{dev,bin,sbin,etc,proc,sys,usr/bin}

  # Copy the BusyBox binary to rootfs
  sudo cp $BUSYBOX_BIN $MOUNT_DIR/bin/busybox

  # Create symlinks for BusyBox commands
  sudo ln -s /bin/busybox $MOUNT_DIR/bin/sh
  sudo ln -s /bin/busybox $MOUNT_DIR/bin/init
  sudo ln -s /bin/busybox $MOUNT_DIR/bin/poweroff

  # Create /etc/init.d/rcS to run at boot and print "hello world"
  sudo mkdir -p $MOUNT_DIR/etc/init.d
  echo -e '#!/bin/sh\n\necho "hello world"\npoweroff -f' | sudo tee $MOUNT_DIR/etc/init.d/rcS
  sudo chmod +x $MOUNT_DIR/etc/init.d/rcS

  # Unmount the disk image
  sudo umount $MOUNT_DIR
}

# Run QEMU with the pre-built kernel and the root filesystem
run_qemu() {
  echo "Running QEMU..."
  qemu-system-x86_64 -kernel "$KERNEL_IMAGE" \
    -append "root=/dev/sda console=ttyS0" \
    -drive file="$ROOTFS",format=raw \
    -nographic \
    -m 512M
}

# Start to execute the tasks
fetch_kernel
fetch_busybox
create_rootfs
run_qemu
