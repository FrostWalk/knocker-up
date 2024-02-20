# Wake-on-LAN and Online Checker

This repository contains a Go application that utilizes Wake-on-LAN (WoL) to wake up a computer and then checks if it's online by pinging it.

## Features

- **Wake-on-LAN**: Sends a "magic packet" to the specified MAC address to wake up the target machine.
- **Online Checking**: Pings the specified IP address or hostname to check if the machine is online.
- **Exponential Retry**: Implements exponential backoff between retries to wake up the machine, allowing for more efficient resource usage.

## Installation

To install the application, you can use the provided script `build-and-install.sh`, which automates the compilation, installation, and setup process. Please note that this script requires root privileges.

```bash
sudo ./build-and-install.sh
```

The script performs the following actions:

- **Compilation**: Compiles the Go program, trimming the path and removing debug information.
- **Installation**: Moves the compiled binary to `/usr/local/bin/`.
- **Systemd Setup**: Copies the systemd service file to `/etc/systemd/system/` and enables the service.
- **Unprivileged Ping Setup**: Configures the system to allow unprivileged users to use the ping command.
- **Service Activation**: Reloads systemd services, starts the service, and enables it to start on boot.

### Systemd Unit File

The provided systemd unit file `knocker-up.service` defines the configuration for the application service. Here's a breakdown of its contents:

```ini
[Unit]
Description=Knocker up
After=multi-user.target
Requires=network-online.target

[Service]
Type=oneshot
ExecStart=/usr/local/bin/knocker-up
DynamicUser=yes

[Install]
WantedBy=default.target
```

- **Description**: Describes the service.
- **After/Requires**: Specifies dependencies on other systemd targets.
- **ExecStart**: Defines the command to execute when the service starts.
- **DynamicUser**: Enables dynamic user and group creation for the service.
- **WantedBy**: Specifies the target that this service should be associated with.


### Run the Application
Execute the built binary with the appropriate command-line arguments. Here's the general syntax:

   ```bash
   ./wake-on-lan [flags] mac ip/hostname
   ```

   Replace [flags], mac, and ip/hostname with the desired options and parameters. Available flags are:

    - `-a int`: Number of attempts to wake the host (default 4)
    - `-e int`: Exponential wait between retries in seconds (default 5)
    - `-t int`: Ping Timeout in seconds (default 2)
    - `-w int`: Seconds to wait between sending the wake command and pinging (default 60)

   Example usage:

   ```bash
   ./wake-on-lan -a 4 -e 5 -t 2 -w 60 00:1A:2B:3C:4D:5E 192.168.1.100
   ```

   This command wakes up the host with the MAC address `00:1A:2B:3C:4D:5E` and checks if it's online at `192.168.1.100`, with default and specified options.

## Dependencies and thanks

- [github.com/mkch/wol](https://github.com/mkch/wol): A Go library for sending Wake-on-LAN (WoL) packets.
- [github.com/prometheus-community/pro-bing](https://github.com/prometheus-community/pro-bing): A Go library for pinging hosts.

## License

This software is licensed under the [MIT License](LICENSE).
