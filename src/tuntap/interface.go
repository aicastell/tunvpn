package tuntap

// a TUN network interface
type Interface interface {
    Name() string
    Read(p []byte) (n int, err error)
    Write(p []byte) (n int, err error)
    Close() error
    String() string
}

func Tun(name string) (Interface, error) {
    return newTUN(name)
}

