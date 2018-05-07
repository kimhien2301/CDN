!/bin/bash

make clean
make
./bin/mirage -G config/json/ntt.json -S true
