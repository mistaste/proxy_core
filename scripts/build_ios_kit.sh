#!/bin/bash
clear

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[1;34m'
CYAN='\033[0;36m'
NC='\033[0m'

echo -e "${BLUE}=== Building ProxyCoreKit.xcframework for iOS ===${NC}"

OUTPUT_DIR="build"
ZIP_NAME="ProxyCoreKit.xcframework.zip"
XCFRAMEWORK_NAME="ProxyCoreKit.xcframework"
CMD_DIR="../src/"
SCRIPT_DIR="$(pwd)"

handle_error() {
    echo -e "${RED}Error: $1${NC}"
    exit 1
}

ensure_build_dir() {
    echo -e "${CYAN}Ensuring build directory exists...${NC}"
    mkdir -p "$OUTPUT_DIR" || handle_error "Failed to create build directory."
}

change_directory() {
    echo -e "${CYAN}Changing directory to $CMD_DIR...${NC}"
    cd "$CMD_DIR" || handle_error "Failed to change directory to $CMD_DIR."
}

build_ios_framework() {
    echo -e "${YELLOW}Building $XCFRAMEWORK_NAME for iOS...${NC}"
    go get golang.org/x/mobile/bind
    export GOROOT_FINAL="/go"

    gomobile bind -v -trimpath -target=ios -o "$SCRIPT_DIR/$OUTPUT_DIR/$XCFRAMEWORK_NAME" ./ios

    if [ $? -ne 0 ]; then
        handle_error "Failed to build iOS framework."
    else
        echo -e "${GREEN}Successfully built $XCFRAMEWORK_NAME in $OUTPUT_DIR directory!${NC}"
    fi
}

zip_xcframework() {
    echo -e "${CYAN}Zipping $XCFRAMEWORK_NAME into $ZIP_NAME...${NC}"

    cd "$SCRIPT_DIR/$OUTPUT_DIR/" || handle_error "Failed to change directory to $SCRIPT_DIR."
    zip -r "$ZIP_NAME" "$XCFRAMEWORK_NAME"
}

main() {
    ensure_build_dir
    change_directory

    build_ios_framework
    zip_xcframework
}

main
