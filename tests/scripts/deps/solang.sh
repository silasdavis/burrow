#!/usr/bin/env bash
set -e
# Solang using temporary one from ewasm PR https://github.com/hyperledger-labs/solang/pull/378
SOLANG_URL="https://solang.io/solang-ewasm"
SOLANG_BIN="$1"

wget -O "$SOLANG_BIN" "$SOLANG_URL"

chmod +x "$SOLANG_BIN"
