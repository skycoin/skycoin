# SkyDEX - Client UI

SkyDEX - Client application UI - Built with React + Vite + Bootstrap

## Features
- 🎨 Modern design with Skywire color scheme (Navy, White, Blue)
- 📱 Single Page Application (SPA)
- 🔌 Market connection via dmsg
- 💰 Product listing and SKY buy/sell functionality
- 📋 Order management
- ⚙️ Wallet settings

## Color Scheme
- Navy: `#101F34`
- White: `#FFFFFF`
- Blue: `#0273FF`

## Setup

### Prerequisites
- Node.js 18+
- npm or yarn

### Install dependencies
```bash
npm install
```

### Run in development mode
```bash
npm run dev
```
Server runs on `http://localhost:8787`.

### Build for production
```bash
npm run build
```
Build files are output to the `dist` folder.

## Embedding in Go Binary

UI files are served embedded in the Go binary:

```go
//go:embed dist/*
var uiFS embed.FS

http.Handle("/", http.FileServer(http.FS(uiFS)))
```

## Project Structure
```
skydex-client/
├── src/
│   ├── components/
│   │   ├── Header.jsx       # Header with connection status
│   │   ├── Market.jsx       # Market page
│   │   ├── Orders.jsx       # My orders
│   │   └── Settings.jsx     # Settings
│   ├── App.jsx              # Main component
│   ├── main.jsx             # Entry point
│   └── index.css            # Styles
├── public/
├── dist/                    # Build output
├── package.json
├── vite.config.js
└── index.html
```

## API Integration

This UI communicates with the market app via dmsg. Message structure is JSON-based:

```javascript
{
  "type": "client.get_products",
  "id": "msg_123",
  "timestamp": 1234567890,
  "data": {}
}
```

## TODO
- [ ] Real dmsg connection
- [ ] Polling for order updates
- [ ] Payment form with non-round amount display
- [ ] Countdown timer for 15-minute windows
- [ ] Error handling and loading states
- [ ] Go backend integration
