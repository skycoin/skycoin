# SkyWallet Hardware Wallet Utilities

Utilities for managing SkyWallet hardware wallet operations.

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

1. Copy udev rules (from skycoin main repo):
   ```bash
   sudo cp udev/51-skywallet.rules /etc/udev/rules.d/
   ```

2. Reload udev:
   ```bash
   sudo udevadm control --reload-rules
   sudo udevadm trigger
   ```

3. Unplug and replug your SkyWallet device

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

1. Check device detected:
   ```bash
   lsusb | grep 313a:0001
   ```
   Should show: `ID 313a:0001 SkycoinFoundation SKYWALLET`

2. Check permissions (replace XXX/YYY with bus/device numbers from lsusb):
   ```bash
   ls -l /dev/bus/usb/XXX/YYY
   ```
   Should show: `crw-rw-rw-` or `crw-rw-r--+`

**Troubleshooting "libusb: bad access [code -3]":**

If you still get access errors, manually unbind the kernel driver:

```bash
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

Verify libraries installed:
```bash
brew list libusb hidapi
```

Check library version:
```bash
pkg-config --modversion libusb-1.0
```

### Windows: DLL Setup

1. Download latest Windows binaries:
   - [libusb-1.0.dll](https://github.com/libusb/libusb/releases) (from MinGW64/dll/)
   - [hidapi.dll](https://github.com/libusb/hidapi/releases)

2. Place both DLLs in:
   - Same directory as `skyhw.exe`, OR
   - `C:\Windows\System32\`, OR
   - Any directory in your PATH

3. For WinUSB devices, install driver using [Zadig](https://zadig.akeo.ie/) if needed

## Usage

### Start the daemon

Run with default settings (port 9510):
```bash
skyhw daemon
```

Run with debug logging:
```bash
skyhw daemon -l debug
```

Specify custom port:
```bash
skyhw daemon -p 9510
```

### Available Commands

Show help:
```bash
skyhw help
```

Show daemon help:
```bash
skyhw daemon --help
```

## Integration with Skycoin Wallet

**Step 1: Start the hardware wallet daemon**

```bash
skyhw daemon -l debug
```

**Step 2: Start the Skycoin wallet daemon (in another terminal)**

```bash
go run . daemon --enable-gui=true --enable-all-api-sets=true
```

**Step 3: Open the wallet GUI in your browser**

Navigate to: http://127.0.0.1:6420

**Step 4: Access hardware wallet features**

Click "SkyWallet" in the GUI

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

Check if usbhid is still bound:
```bash
find /sys/bus/usb/devices -name "*313a*" -type d -exec ls -l {}/driver \; 2>/dev/null
```
Should show "No such file" (driver unbound) or nothing.

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

Linux/macOS:
```bash
go build -o skyhw cmd/hardware-wallet/skycoin.go
```

Windows (with MinGW):
```bash
set CGO_ENABLED=1
go build -o skyhw.exe cmd/hardware-wallet/skycoin.go
```

### Testing

Start daemon:
```bash
skyhw daemon
```

In another terminal, test API:
```bash
curl http://localhost:9510/api/v1/available
```

## Security Notes

- **Linux MODE="0666"**: Makes device accessible to all users. Safe for hardware wallets requiring physical confirmation.
- **TAG+="uaccess"**: On systemd systems, restricts access to currently logged-in users (more secure).
- **Windows DLLs**: Download from official sources only to avoid malware.
- **Hardware Confirmation**: Sensitive operations require physical button press on device.

## Support

- Hardware Wallet Repository: https://github.com/skycoin/hardware-wallet-go
- Daemon Repository: https://github.com/skycoin/hardware-wallet-daemon
