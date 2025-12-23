#!/usr/bin/env bash

# Simplified macOS .pkg installer creation for Skycoin CLI tools
# This creates installers for skycoin and skyhw binaries

set -euo pipefail

# Variables
mac_script_dir="./scripts/mac_installer"
installer_build_dir="./mac_build"
git_tag=$(git describe --tags)
go_arch=
sign_binary=false
notarize_binary=false
output=

greent='\033[0;32m'
yellowt='\033[0;33m'
nct='\033[0m' # No Color

# Has to be run from MacOS Host
current_os="$(uname -s)"
if [[ "$current_os" != "Darwin" ]]; then
  echo "Can only be run from MacOS Host"
  exit 1
fi

function print_usage() {
  echo "Usage: sh create_installer.sh [-o|--output output_dir] [-s | --sign] [-n | --notarize]"
  echo "Environment variables for signing/notarizing:"
  echo -e "${greent}MAC_HASH_INSTALLER_ID${nct}  : Hash of Developer ID Installer"
  echo -e "${greent}MAC_DEVELOPER_USERNAME${nct} : Developer Account Email"
  echo -e "${greent}MAC_DEVELOPER_PASSWORD${nct} : Application Password"
}

function build_installer() {
  if [ -z "$output" ]; then
    output="${PWD}/"
    echo "No output flag provided, storing installer to: ${output}"
  else
    if [ "${output:(-1)}" != "/" ]; then
      output="${output}/"
    fi
    echo "Storing installer to ${output}"
  fi

  # Fetch skycoin binaries from the current release
  download_url=$(curl --header "Authorization: Bearer ${GITHUB_TOKEN}" \
    https://api.github.com/repos/skycoin/skycoin/releases | \
    jq -r ".[0].assets[] | select(.name|match(\"skycoin-${git_tag}-darwin-${go_arch}.tar.gz\")) | .browser_download_url")

  if [ -z "$download_url" ] || [ "$download_url" == "null" ]; then
    echo "Error: Could not find darwin-${go_arch} release archive"
    exit 1
  fi

  echo "Downloading ${download_url}"
  wget "${download_url}" -O - | tar -xz

  # Clean previous build
  rm -rf ${installer_build_dir}

  # Create directories for package payload
  mkdir -p ${installer_build_dir}/payload/usr/local/bin

  # Move binaries to payload
  mv ./skycoin ${installer_build_dir}/payload/usr/local/bin/skycoin
  if [ -f ./skyhw ]; then
    mv ./skyhw ${installer_build_dir}/payload/usr/local/bin/skyhw
  fi

  # Make binaries executable
  chmod +x ${installer_build_dir}/payload/usr/local/bin/*

  # Build the package
  package_name=skycoin-installer-${git_tag}-darwin-${go_arch}.pkg

  if [ "$sign_binary" == true ] && [ ! -z ${MAC_HASH_INSTALLER_ID+x} ]; then
    pkgbuild --sign "$MAC_HASH_INSTALLER_ID" \
      --root ${installer_build_dir}/payload \
      --identifier com.skycoin.skycoin \
      --version "${git_tag}" \
      --install-location / \
      "${output}${package_name}"
  else
    pkgbuild --root ${installer_build_dir}/payload \
      --identifier com.skycoin.skycoin \
      --version "${git_tag}" \
      --install-location / \
      "${output}${package_name}"
  fi

  cd "${output}"

  if [ "$notarize_binary" == true ]; then
    if [ -z "$MAC_DEVELOPER_USERNAME" ] || [ -z "$MAC_DEVELOPER_PASSWORD" ]; then
      echo -e "${yellowt}MAC_DEVELOPER_USERNAME and MAC_DEVELOPER_PASSWORD required for notarization${nct}"
      exit 1
    fi
    xcrun altool --notarize-app \
      --primary-bundle-id "com.skycoin.skycoin" \
      --username="$MAC_DEVELOPER_USERNAME" \
      --password="$MAC_DEVELOPER_PASSWORD" \
      --file "${package_name}" && {
      echo -e "${greent}Check your email for notarization status${nct}"
    }
  fi

  rm -rf ${installer_build_dir}
}

# Parse arguments
while :; do
  case "${1:-}" in
  -o | --output)
    if [ -n "${2:-}" ]; then
      output=$2
      shift
    else
      echo 'ERROR: "--output" requires a non-empty option argument.' >&2
      exit 1
    fi
    ;;
  --output=?*)
    output=${1#*=}
    ;;
  --output=)
    echo 'ERROR: "--output" requires a non-empty option argument.' >&2
    exit 1
    ;;
  -h | --help)
    print_usage
    exit 0
    ;;
  -s | --sign)
    sign_binary=true
    ;;
  -n | --notarize)
    notarize_binary=true
    ;;
  -?*)
    printf 'WARN: Unknown option (ignored): %s\n' "$1" >&2
    ;;
  *)
    break
    ;;
  esac
  shift
done

# Build installers for both architectures
echo "Building arm64 installer..."
go_arch=arm64
build_installer

echo "Building amd64 installer..."
go_arch=amd64
build_installer

echo -e "${greent}Installer creation complete!${nct}"
