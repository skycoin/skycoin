# SkyWallet Hardware Wallet Daemon

Command-line daemon for managing SkyWallet hardware wallet operations.

## Prerequisites

Install the required USB libraries for your platform:

### Linux
```bash
# Debian/Ubuntu
sudo apt-get install libusb-1.0-0-dev libudev-dev

# Fedora/RHEL
sudo dnf install libusb-devel systemd-devel

# Arch Linux
sudo pacman -S libusb systemd
```

### macOS
```bash
brew install libusb hidapi
```

### Windows
1. Install MinGW-w64: https://www.mingw-w64.org/
2. Download pre-built DLLs:
   - libusb-1.0: https://github.com/libusb/libusb/releases
   - hidapi: https://github.com/libusb/hidapi/releases
3. Place DLLs in same directory as executable or in system PATH

## Platform Setup

### Linux: USB Permissions & Driver Setup

On Linux, you need udev rules to allow non-root access and unbind the kernel driver.

**Quick Setup:**
```bash
# 1. Copy udev rules (from skycoin main repo)
sudo cp udev/51-skywallet.rules /etc/udev/rules.d/

# 2. Reload udev
sudo udevadm control --reload-rules
sudo udevadm trigger

# 3. Unplug and replug your SkyWallet device
```

**Manual udev rule creation:**
```bash
sudo tee /etc/udev/rules.d/51-skywallet.rules > /dev/null <<'EOF'
# SkyWallet Hardware Wallet - Permissions and driver unbind
SUBSYSTEM=="usb", ATTR{idVendor}=="313a", ATTR{idProduct}=="0001", MODE="0666", TAG+="uaccess", \
  RUN+="/bin/sh -c 'for iface in /sys/bus/usb/devices/$kernel:*; do \
    if [ -e $iface/driver ]; then \
      echo $kernel:$(basename $iface | cut -d: -f2) > /sys/bus/usb/drivers/usbhid/unbind; \
    fi; \
  done'"
EOF

sudo udevadm control --reload-rules
sudo udevadm trigger
```

**Verify setup:**
```bash
# Check device detected
lsusb | grep 313a:0001
# Should show: ID 313a:0001 SkycoinFoundation SKYWALLET

# Check permissions (replace XXX/YYY with bus/device numbers from lsusb)
ls -l /dev/bus/usb/XXX/YYY
# Should show: crw-rw-rw- or crw-rw-r--+
```

**Troubleshooting "libusb: bad access [code -3]":**

If you still get access errors, manually unbind the kernel driver:
```bash
# Find and unbind usbhid driver
INTERFACE=$(find /sys/bus/usb/devices -type l -name "driver" 2>/dev/null | \
  while read link; do \
    iface=$(dirname "$link"); \
    if grep -q "313a" "$iface/../idVendor" 2>/dev/null && \
       grep -q "0001" "$iface/../idProduct" 2>/dev/null; then \
      basename "$iface"; break; \
    fi; \
  done)

echo "$INTERFACE" | sudo tee /sys/bus/usb/drivers/usbhid/unbind
```

### macOS: No Special Setup Required

macOS allows direct USB HID access without special permissions. Just ensure libusb and hidapi are installed via Homebrew (see Prerequisites).

**If running into issues:**
```bash
# Verify libraries installed
brew list libusb hidapi

# Check library paths
pkg-config --modversion libusb-1.0
```

### Windows: DLL Setup

1. Download latest Windows binaries:
   - [libusb-1.0.dll](https://github.com/libusb/libusb/releases) (from MinGW64/dll/)
   - [hidapi.dll](https://github.com/libusb/hidapi/releases)

2. Place both DLLs in:
   - Same directory as `skyhw-daemon.exe`, OR
   - `C:\Windows\System32\`, OR
   - Any directory in your PATH

3. For WinUSB devices, install driver using [Zadig](https://zadig.akeo.ie/) if needed

## Usage

### Start the daemon

```bash
# Run with default settings (port 9510)
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

### Linux: "libusb: bad access [code -3]"

**Cause:** User doesn't have permission to access USB device, or kernel driver still bound.

**Solution:**
1. Verify udev rules installed: `cat /etc/udev/rules.d/51-skywallet.rules`
2. Reload udev and reconnect device
3. Manually unbind kernel driver (see Linux setup section above)
4. Add user to plugdev group: `sudo usermod -a -G plugdev $USER` (then log out/in)

### All Platforms: Device Not Detected

1. Ensure SkyWallet is properly connected
2. Try different USB port or cable
3. Check if device appears in system:
   - Linux: `lsusb | grep 313a`
   - macOS: `system_profiler SPUSBDataType | grep -A5 313a`
   - Windows: Device Manager → Universal Serial Bus devices

### Linux: Check Kernel Driver Status

```bash
# See if usbhid is still bound
find /sys/bus/usb/devices -name "*313a*" -type d -exec ls -l {}/driver \; 2>/dev/null
# Should show "No such file" (driver unbound) or nothing
```

### Windows: Missing DLL Errors

**Error:** "The code execution cannot proceed because libusb-1.0.dll was not found"

**Solution:** Copy `libusb-1.0.dll` and `hidapi.dll` to executable directory

### macOS: Build Errors

**Error:** `ld: library not found for -lhidapi`

**Solution:**
```bash
brew link libusb hidapi
export CGO_LDFLAGS="-L/usr/local/lib"
export CGO_CFLAGS="-I/usr/local/include"
```

## Development

### Building

```bash
# Linux/macOS
go build -o skyhw-daemon cmd/hardware-wallet/skycoin.go

# Windows (with MinGW)
set CGO_ENABLED=1
go build -o skyhw-daemon.exe cmd/hardware-wallet/skycoin.go
```

### Testing

```bash
# Start daemon
./skyhw-daemon daemon

# In another terminal, test API
curl http://localhost:9510/api/v1/available
```

## Security Notes

- **Linux MODE="0666"**: Makes device accessible to all users. Safe for hardware wallets requiring physical confirmation.
- **TAG+="uaccess"**: On systemd systems, restricts access to currently logged-in users (more secure).
- **Windows DLLs**: Download from official sources only to avoid malware.
- **Hardware Confirmation**: Sensitive operations require physical button press on device.

## Support

- Hardware Wallet Repository: https://github.com/0pcom/hardware-wallet-go
- Daemon Repository: https://github.com/0pcom/hardware-wallet-daemon
- Skycoin Main Repository: https://github.com/0pcom/skycoin
