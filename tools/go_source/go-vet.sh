for gosrc in `find . -type f -name '*.go'`; do
    go vet $gosrc
done
