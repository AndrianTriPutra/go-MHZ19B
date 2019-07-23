# go-MHZ19B

1. Check your port with dmesg
2. sudo chmod a+rw /dev/tty_your port
3. Change "err := atmhz.Open("/dev/ttyUSB0", 9600, 5*time.Second)" which your port
