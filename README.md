Simple TUN VPN

Usage:

    ./src -h
    Usage of ./src:
      -i string
            tunnel interface (default "tun0")
      -l string
            local ip and netmask (default "10.0.0.1/24")
      -p string
            application port (local and remote) (default "1234")
      -r string
            remote ip (default "8.8.8.8")



In PC1:

    $ ./src -i tun0 -l 10.0.0.1/24 -p 4142 -s pc1.name


In PC2:

    $ ./src -i tun0 -l 10.0.0.2/24 -p 4142 -s pc2.name


After that, you can communicate PC1 and PC2 through using the IPs 10.0.0.1 and
10.0.0.2. For example, you can send pings:

    pc1> ping 10.0.0.2






