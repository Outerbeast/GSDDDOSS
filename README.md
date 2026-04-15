# GSDDDOSS (GoldSource Denied DDoS Script)

Go port of the original Python script. Monitors your GoldSrc/SvEngine game server's UDP log output for attacks and blocks offending IPs in the firewall.

## Requirements

- Go 1.21 or greater
- Administrator/root privileges

## Building

**Windows:**
```
build.cmd
```

**Linux:**
```
chmod +x build.sh
./build.sh
```

## Usage

```
GSDDOSS                           Defaults: 127.0.0.1:8008
GSDDOSS --host 0.0.0.0           Custom host
GSDDOSS --port 9000               Custom port
GSDDOSS --grouped-rules           Windows: group IPs into single firewall rule
GSDDOSS --firewall-cmd "cmd {ip}" Custom firewall command
GSDDOSS --help
```

## Configuration

On first run, `config.json` and `blocked.json` are created automatically.

**config.json:**
```json
{
    "listener_addr": "127.0.0.1",
    "listener_port": 8008,
    "windows_rule_ip_grouped": true,
    "command_add_block": ""
}
```

**blocked.json:**
```json
{
    "ips": []
}
```

### Migration

If you have an existing `ips.json` from the Python version, it will be automatically migrated to `blocked.json` and renamed to `ips.json.bak`.

## Server Configuration

Add to your `server.cfg`:
```
logaddress_add 127.0.0.1 8008
```

Restart the server or run the same command in the server console.

## Detected Attacks

- Bad RCON attempts
- Split packet abuse
- A2S/Info responses (reflected DDoS)
- Rate limit violations
- Port blasting

## Features

- Graceful shutdown (Ctrl+C)
- Localhost protection (won't block 127.0.0.1 or ::1)
- Windows grouped firewall rules (avoids rule limits)
- Custom firewall commands via config or CLI
- Automatic migration from Python version

## Windows Service (Optional)

To run as a background service, use [NSSM](https://nssm.cc/):
```
nssm install GSDDDOSS C:\path\to\GSDDOSS.exe
nssm start GSDDDOSS
```

## Support

- Windows (tested)
- Linux (should work)
