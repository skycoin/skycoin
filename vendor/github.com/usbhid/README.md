# libusb + hidapi go wrapper

This is a go wrapper around libusb and hidapi.

We have devices that can work either with libusb or hidapi, so we needed to make a go package that can talk with both.

Note that this is necessary only because of macOS; on Linux, hidapi is using libusb; and on Windows, libusb talks to hid devices using the same hid.dll as hidapi. In theory it would be cleaner to add HID API to libusb on macOS, but making this go package was quicker.

The code is mostly copied from https://github.com/karalabe/hid and https://github.com/deadsy/libusb

## License

Code is under GNU LGPL 2.1.

* (C) Karel Bilek 2017
* (c) 2017 Jason T. Harris (also see https://github.com/deadsy/libusb for comprehensive list)
* (C) 2017 Péter Szilágyi (also see https://github.com/karalabe/hid for comprehensive list)
