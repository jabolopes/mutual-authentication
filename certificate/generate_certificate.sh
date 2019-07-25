#!/bin/bash
#
# Generates a certificate and private key signed by root certificate
# authority. The output must be a filename pattern, to which ".crt"
# and ".key" will appended to make the final filenames for the
# certificate and key (respectively). The root ca must be filename
# pattern used to generate the root certificate. The common name must
# be either a fully qualified hostname (for a server) or an email
# address (for a client or person user).

readonly OUTPUT="$1"

if [[ -z "${OUTPUT}" ]]; then
  echo output must be given
  exit 1
fi

readonly ROOT_CA="$2"

if [[ -z "${ROOT_CA}" ]]; then
  echo root CA must be given
  exit 1
fi

readonly COMMON_NAME="$3"

if [[ -z "${COMMON_NAME}" ]]; then
  echo common name must be given
  exit 1
fi

readonly KEY="${OUTPUT}.key"
readonly CSR="${OUTPUT}.csr"
readonly CERTIFICATE="${OUTPUT}.crt"

readonly ROOT_CRT="${ROOT_CA}.crt"
readonly ROOT_KEY="${ROOT_CA}.key"

openssl genrsa -out "${KEY}" 4096
openssl req -new -sha256 -key "${KEY}" -subj "/CN=${COMMON_NAME}" -out "${CSR}"
openssl x509 -req -in "${CSR}" -CA "${ROOT_CRT}" -CAkey "${ROOT_KEY}" -set_serial 100 -out "${CERTIFICATE}" -days 365 -sha256
