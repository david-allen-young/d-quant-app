# d-quant-app

**d-quant-app** is a prototype frontend and server backend for generating expressive MIDI phrases and corresponding visualizations based on user-selected musical attributes. It serves as a user interface wrapper for the [`d-quant`](https://github.com/david-allen-young/d-quant) performance modeling system, which uses real-world musical recordings to inform expressive playback.

> ⚠️ **Note:** This version is in early development and currently generates **dummy PNG and MIDI files** to demonstrate the web-to-backend interaction.

---

## 🎯 Project Goals

This app is designed to:

- Let users visually configure note parameters such as articulation, accent, pitch, and dynamic changes.
- Generate a MIDI phrase and Breath Controller CC visualization based on user input.
- Support future integration with the `d-quant` C++ CLI to synthesize expressive musical performance.

---

## 🧰 Current Functionality

The app includes:

- A **Go-based backend** (`main.go`) that listens on port `8080` and executes a command-line tool to generate output.
- A **frontend HTML interface** (`preview.html`) that allows users to configure parameters and preview symbolic notation.
- A **dummy CLI tool** (`dquant_cli.go`) that writes placeholder `.mid` and `.png` files to simulate output.

---

## 🚀 Quickstart Guide

### 1. Clone the repository

```bash
git clone https://github.com/david-allen-young/d-quant-app.git
cd d-quant-app
```

### 2. Build the dummy CLI

```bash
cd cli
go build -o dquant_cli.exe dquant_cli.go
cd ..
```

### 3. Build and run the Go web server

```bash
go run server/main.go
```

> Output files will be written to the `output/` folder at the root level.

### 4. Open the frontend

Launch a browser and visit:

```
http://localhost:8080/html/preview.html
```

### 5. Try it out!

- Select musical attributes using the dropdown menus.
- Click **"Generate"** to simulate a performance.
- A new window will appear with:
  - A dummy **PNG image** preview.
  - A link to **download a dummy MIDI file**.

---

## 📁 Project Structure

```
.
├── cli/
│   └── dquant_cli.go         # Dummy CLI generator (MIDI + PNG)
├── server/
│   └── main.go               # Go web server
├── html/
│   └── preview.html          # Frontend UI for user interaction
├── images/                   # Icons for articulations, dynamics, noteheads
├── output/                   # Auto-generated MIDI/PNG files
```

---

## 🔮 Coming Soon

- Real MIDI generation via integrated `d-quant` C++ CLI.
- Expression envelope visualizations based on performance data.
- Full phrase generation with legato/slur and dynamic phrasing.

---

## 🧠 Background

This app is part of a larger research and development effort to humanize quantized MIDI by statistically modeling real musical performances. See [d-quant](https://github.com/david-allen-young/d-quant) for details.
