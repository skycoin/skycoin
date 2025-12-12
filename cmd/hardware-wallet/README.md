# SkyWallet Hardware Wallet Daemon

Command-line daemon for managing SkyWallet hardware wallet operations.

## Prerequisites

- SkyWallet hardware device
- libusb library installed on your system
  - **Debian/Ubuntu**: `sudo apt install libusb-1.0-0-dev`
  - **Fedora/RHEL**: `sudo dnf install libusb-devel`
  - **Arch Linux**: `sudo pacman -S libusb`
  - **macOS**: `brew install libusb`

## USB Permissions Setup (Linux)

On Linux, you need to configure udev rules to allow non-root access to the SkyWallet device.

### Quick Setup

1. **Copy the udev rules file:**
   ```bash
   sudo cp ../../51-skywallet.rules /etc/udev/rules.d/
   ```

2. **Reload udev rules:**
   ```bash
   sudo udevadm control --reload-rules
   sudo udevadm trigger
   ```

3. **Unplug and replug your SkyWallet device**

4. **Verify the permissions:**
   ```bash
   # Check that the device is detected
   lsusb | grep -i skywallet
   # Output should show: Bus XXX Device YYY: ID 313a:0001 SkycoinFoundation SKYWALLET
   
   # Check device file permissions (replace XXX and YYY with values from lsusb)
   ls -l /dev/bus/usb/XXX/YYY
   # Should show: crw-rw-rw- ... (mode 0666, world-writable)
   ```

### Troubleshooting Permissions

If the daemon still shows "libusb: bad access [code -3]" errors:

1. **Verify the udev rule is installed:**
   ```bash
   cat /etc/udev/rules.d/51-skywallet.rules
   ```

2. **Check the actual device permissions:**
   ```bash
   lsusb | grep 313a:0001
   # Note the Bus and Device numbers
   ls -l /dev/bus/usb/BUS/DEVICE
   ```
   The file should have `rw-rw-rw-` permissions (0666).

3. **If permissions are still wrong, try manual trigger:**
   ```bash
   # Find the device path
   udevadm info --name=/dev/bus/usb/BUS/DEVICE --attribute-walk | grep KERNEL
   
   # Trigger udev for this device
   sudo udevadm trigger --action=change --attr-match=idVendor=313a
   ```

4. **As a last resort, you can temporarily change permissions manually:**
   ```bash
   sudo chmod 666 /dev/bus/usb/BUS/DEVICE
   ```
   (Note: This is temporary and will be reset when the device is unplugged)

## Usage

### Start the daemon

```bash
# Run with default settings
go run cmd/hardware-wallet/skycoin.go daemon

# Run with debug logging
go run cmd/hardware-wallet/skycoin.go daemon -l debug

# Specify custom port
go run cmd/hardware-wallet/skycoin.go daemon -p 9510
```

### Available Commands

```bash
# Show help
go run cmd/hardware-wallet/skycoin.go help

# Show daemon help
go run cmd/hardware-wallet/skycoin.go daemon --help
```

## Integration with Skycoin Wallet

1. **Start the hardware wallet daemon:**
   ```bash
   go run cmd/hardware-wallet/skycoin.go daemon -l debug
   ```

2. **Start the Skycoin wallet daemon** (in another terminal):
   ```bash
   go run . daemon --enable-gui=true --enable-all-api-sets=true
   ```

3. **Open the wallet GUI** in your browser:
   ```
   http://127.0.0.1:6420
   ```

4. **Click "SkyWallet"** in the GUI to access hardware wallet features

## Troubleshooting

### "libusb: bad access [code -3]" Error

This means your user doesn't have permission to access the USB device. Follow the "USB Permissions Setup" section above.

### Device Not Detected

1. Ensure the SkyWallet is properly connected
2. Check if it appears in lsusb: `lsusb | grep -i skywallet`
3. Try a different USB port or cable
4. Restart the hardware wallet daemon

### libusb Library Not Found

Install libusb development package for your distribution (see Prerequisites section).

## Development

### Building

```bash
go build -o skywallet-daemon cmd/hardware-wallet/skycoin.go
```

### Testing

```bash
# Start daemon and check it responds
curl http://localhost:9510/api/v1/available
```
