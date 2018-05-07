for gosrc in `find . -type f -name '*.go'`; do
    $HOME/.gopath/bin/golint $gosrc
done
