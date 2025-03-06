# hex

A simple tool to encode and decode hexadecimal data.

## Install

```sh
go install github.com/maraino/hex
```

## Usage

No brainer functionality, just three flags and four features.

```sh
$ hex --help
Usage: hex [<filename>]
  -c      encodes the input as hexadecimal followed by characters
  -d      decodes input
  -p      encoded using a prettier format aa:bb, pass -p=n to print using n columns
  -go     encoded using as Go's []byte, pass -go=n to print using n columns
  -n      interprets input as a number
  -b int  base used to when -n is used (default: 10)
```

### Encode

Encodes a file or the standard input to hexadecimal.

```sh
$ echo Hello World! >  hello.txt
$ hex hello.txt
48656c6c6f20576f726c64210a
```

```sh
$ echo Hello World! | hex
48656c6c6f20576f726c64210a
```

Encode to a prettier format:

```sh
$ echo Hello World! | hex -p
48:65:6c:6c:6f:20:57:6f:72:6c:64:21:0a
$ echo Hello World! | hex -p=8
48:65:6c:6c:6f:20:57:6f:
72:6c:64:21:0a
```

Encode to a Go format:

```sh
$ echo Hello World! | hex -go
[]byte{
	0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x57, 0x6f,
	0x72, 0x6c, 0x64, 0x21, 0x0a,
}
$ echo Hello World! | hex -go=16
[]byte{
	0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x57, 0x6f, 0x72, 0x6c, 0x64, 0x21, 0x0a,
}
```

### Decode

Decodes an hexadecimal string:

```sh
$ echo 48656c6c6f20576f726c64210a | hex -d
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

### Decode and dump

Combines `-d` and `-c` to first decode and then dump:

```sh
$ echo 48:65:6C:6C:6F:20:57:6F:72:6C:64:21:0A | hex -d -c
00000000  48 65 6c 6c 6f 20 57 6f  72 6c 64 21 0a           |Hello World!.|
```
