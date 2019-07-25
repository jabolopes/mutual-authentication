#!/bin/bash
#
# Generates the root certificate and private key for a custom
# certificate authority. The output must be a filename pattern, to
# which ".crt" and ".key" will appended to make the final filenames
# for the certificate and key (respectively).

readonly OUTPUT="$1"

if [[ -z "${OUTPUT}" ]]; then
  echo output must be given
  exit 1
fi

readonly HOST="${2:-$(hostname -f):27388}"

if [[ -z "${HOST}" ]]; then
  echo host must be given
  exit 1
fi

readonly KEY="${OUTPUT}.key"
readonly CERTIFICATE="${OUTPUT}.crt"

openssl genrsa -out "${KEY}" 4096
openssl req -x509 -new -nodes -key "${KEY}" -sha256 -days 365 -out "${CERTIFICATE}" -subj "/CN=${HOST}"
