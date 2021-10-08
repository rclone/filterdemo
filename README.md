# Rclone filter demo website

See https://filterdemo.rclone.org/ for the website which allows you to
try out [rclone](https://rclone.org) filters on a set of file names.

## How it works

This compiles the filters part of rclone into a WASM module and
imports it into the browser along with a bit of user interface (also
written in Go).

## Testing ##

Build the code with:

    ./build
    
Run a webserver to serve it:

    ./serve &

Then go to http://localhost:3000/
