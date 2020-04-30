# hex

A simple tool to encode and decode to hexadecimal.

## Install

```sh
go install github.com/maraino/hex
```

## Usage

No brainer functionality, just two flags and four features.

```sh
$ hex --help
Usage: hex [<filename>]
 -c encodes the input as hexadecimal followed by characters
 -d decodes input
```

### Encode

Encodes a file or the standard input to hexadecimal.

```sh
$ echo Hello World! >  hello.txt
$ hex hello.txt
48656C6C6F20576F726C64210A
```

```sh
$ echo Hello World! | hex
48656C6C6F20576F726C64210A
```

### Decode

Decodes an hexadecimal string:

```sh
$ echo 48656C6C6F20576F726C64210A | hex -d
Hello World!
```

You can also use a file or a different formatting, it will ignore any
non-hexadecimal character:

```sh
$ echo 48:65:6C:6C:6F:20:57:6F:72:6C:64:21:0A > hello.hex
$ hex -d hello.hex
Hello World!
```

### Hex dump

Returns an hex dump of the given data, like `hexdump -C`:

```sh
$ hex -c hello.txt
00000000  48 65 6c 6c 6f 20 57 6f  72 6c 64 21 0a           |Hello World!.|
```

## Decode and dump

Combines `-d` and `-c` to first decode and then dump:

```sh
$ echo 48:65:6C:6C:6F:20:57:6F:72:6C:64:21:0A | hex -d -c
00000000  48 65 6c 6c 6f 20 57 6f  72 6c 64 21 0a           |Hello World!.|
```
