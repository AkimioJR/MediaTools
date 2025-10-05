#!/bin/bash

xattr -d com.apple.quarantine ./MediaTools.app
xattr -c ./MediaTools.app