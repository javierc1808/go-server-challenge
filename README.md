
This is a simple server that exposes the following endpoints:

## Real-time notifications

    ws://localhost:8080/notifications     

This endpoint emits a notification everytime a new document is created by other user.
Each notification is represented as a JSON document.

Example notifications:

    {"timestamp":"2020-08-12T07:30:08.28093+02:00","userId":"3ffe27e5-fe2c-45ea-8b3c-879b757b0455","userName":"Alicia Wolf","documentId":"f09acc46-3875-4eff-8831-10ccf3356420","documentTitle":"Edmund Fitzgerald Porter","type":"document_created"}
    ...
 

## Documents API 

    http://localhost:8080/documents     

This endpoint exposes the documents API thru a simple JSON over HTTP endpoint.

Example response
    
```json    
    [
        {
            "id": "69517c79-a4b2-4f64-9c83-20e5678e4519",
            "title": "Arrogant Bastard Ale",
            "version": "5.3.15",
            "attachments": [
                "European Amber Lager",
                "Wood-aged Beer"
            ],
            "contributors": [
                {
                    "id": "1b41861e-51e2-4bf4-ba13-b20f01ce81ef",
                    "name": "Jasen Crona",
                    "createdAt": "2020-08-12T07:30:08.28093+02:00",
                    "updatedAt": "2020-08-12T07:30:08.28093+02:00"
                },
                {
                    "id": "2a1d6ed0-7d2d-4dc6-b3ea-436a38fd409e",
                    "name": "Candace Jaskolski",
                    "createdAt": "2020-08-12T07:30:08.28093+02:00",
                    "updatedAt": "2020-08-12T07:30:08.28093+02:00"
                },
                {
                    "id": "9ae28565-4a1c-42e3-9ae8-e39e6f783e14",
                    "name": "Rosemarie Schaden",
                    "createdAt": "2020-08-12T07:30:08.28093+02:00",
                    "updatedAt": "2020-08-12T07:30:08.28093+02:00"
                }
            ],
            "createdAt": "1912-03-08T06:01:39.382278739Z",
            "updatedAt": "1952-02-29T22:21:13.817038244Z"
        },
        ...
    ]
```

# Running the server

To run the server, you need Golang runtime installed in your workspace. Then run the following:

    go run server.go

By default, the server listen to `localhost:8080`, if needed change the listening address using `-addr` flag, for example:

    go run server.go -addr localhost:9090 


## Run with Cloudflare Tunnel

This repository includes a script to run the server and expose it securely via a public HTTPS tunnel using Cloudflare Tunnel. This is useful for Expo/Android, where plain HTTP (localhost/127.0.0.1) may be blocked and SSL is required.

### Requirements
- Go installed (`go` in PATH)
- `curl` installed
- The script will try to install `cloudflared` automatically (Homebrew on macOS; apt/dnf/pacman or direct binary download on Linux). On Windows, WSL is recommended.

### Basic usage
```bash
chmod +x scripts/run_with_tunnel.sh
./scripts/run_with_tunnel.sh
```

By default it starts the server on `http://localhost:8080` and opens a tunnel. The console will display:
- Public HTTPS URL: `https://<something>.trycloudflare.com`
- Secure WebSocket (WSS): `wss://<something>.trycloudflare.com/notifications`

### Choose port
```bash
PORT=8443 ./scripts/run_with_tunnel.sh
```

### Use it from your app (React Native / Expo)
```javascript
const BASE_URL = 'https://<your-subdomain>.trycloudflare.com';

// REST
const res = await fetch(`${BASE_URL}/documents`);

// WebSocket
const ws = new WebSocket(`${BASE_URL.replace('https', 'wss')}/notifications`);
```

### Stop the process
Press Ctrl+C to close both the server and the tunnel.

### Platform notes
- macOS: uses Homebrew to install `cloudflared`. If you don't have Homebrew, install it from `https://brew.sh` or install `cloudflared` manually.
- Linux: tries apt/dnf/pacman, and if not available, downloads the official binary automatically.
- Windows: run the script from WSL or Git Bash. Alternatively, install `cloudflared` manually and run the server and tunnel separately.

### Troubleshooting
- If the tunnel URL doesn't show up, rerun the script and ensure the port is free (try a different `PORT`).
- Ensure you have Internet connectivity so `cloudflared` can create the tunnel.
