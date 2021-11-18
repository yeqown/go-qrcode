## `qrcode` CLI Application

The `qrcode` is a command line application that generates QR codes.

### Usage

```sh
$ qrcode -h
NAME:
   qrcode - qrcode [options] [source text]

USAGE:
   qrcode [global options] command [command options] [arguments...]

VERSION:
   v2.0.0-beta

DESCRIPTION:
   QR code generator

AUTHOR:
   yeqown <yeqown@gmail.com>

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --terminal                 --terminal (default: false)
   --output value, -o value   --output=<output file> (default: qrcode.jpg)
   --block value, -s value    --block=<block size> (default: 5)
   --borders value, -b value  --borders=<borders> (default: 0,0,0,0)
   --circle                   --circle (default: false)
   --help, -h                 show help (default: false)
   --version, -v              print the version (default: false)

COPYRIGHT:
   Copyright (c) 2018 yeqown

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.
```

### Examples

Provides some examples of how to use the `qrcode` command.

```sh
# Generate a QR code into file as default
qrcode "Hello, World!"

# Generate a QR code into file with block size and borders (unit: pixel)
qrcode -o qrcode.png -s  20 -b 20,20,20,20 -m "Hello, World!"

# Generate a QR code into terminal
qrcode --terminal "Hello, World!"
```