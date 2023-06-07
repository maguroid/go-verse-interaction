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
make goerli/london # interact with goerli with post london upgrade's manner
make goerli/legacy # interact with goerli with pre london upgrade's manner
make sand/london # interact with sand-verse with post london upgrade's manner
make sand/legacy # interact with sand-verse with pre london upgrade's manner
```

### 4. Caveats

- You cannot use the post london upgrade's manner for the verse layer which does not support the EIP-1559.
- When compiling contracts, you must specify the `--evm-version` as before than `paris`.
