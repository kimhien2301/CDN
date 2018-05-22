!/bin/bash

make clean
make

#./bin/mirage -G config/json/3node-iris.json -S true

#./bin/mirage -G config/json/ntt.json -S -GA true

#./bin/mirage -G config/json/ntt.json -S true

# Insert new contents
./bin/mirage -G config/json/ntt-color-only.json -S -I true
./bin/mirage -G config/json/ntt.json -S -I true

