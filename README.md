## Go Verse Interaction

This repository is only for the experiment of the interaction between Go and Oasys Verse Layer.

## How to use

### 1. Install Go

https://go.dev/doc/install

### 2. Setting Environment Variables

```bash
cp .env.example .env
```

Replace the values in `.env` with your own values.

### 3. Run scripts

```bash
make goerli # for goerli testnet
make sand/london # interact with sand-verse with post london upgrade's manner
make sand/legacy # interact with sand-verse with pre london upgrade's manner
```
