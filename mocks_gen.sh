#!/bin/bash

SRC_DIR=$1
MOCK_DIR="mocks"

mkdir -p $MOCK_DIR

# Find all Go files that don't end with _test.go, and generate mocks for them
find $SRC_DIR -type f -name '*.go' ! -name '*_test.go' | while read -r file; do
    src_file=$(basename "$file")
    mock_file="mock_${src_file}"
    package=$(basename $MOCK_DIR)

    echo "Generating mock for $src_file..."
    mockgen -source=$file -destination=$MOCK_DIR/$mock_file -package=$package
done

echo "Mock generation complete."
