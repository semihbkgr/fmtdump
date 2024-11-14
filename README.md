# fmtdump

`fmtdump` is a flexible data file dump tool for custom-formatted binary files. Designed to simplify the inspection of database internals and other structured binary files, fmtdump enables the parsing and display of data layouts in files such as WAL, LSM, and other database-specific formats. With fmtdump, binary file contents can be analyzed and troubleshooted by viewing data according to predefined structures.

```txt
Flexible data file dump tool for custom formats

Usage:
  fmtdump --format=<format.json> <data-file>

Flags:
  -f, --format string   path of the format file
  -h, --help            help for fmtdump
```

## Format Specification

- `name`: The name of the field.
- `size`: The size of the field's value in bytes.
- `sizeRef`: Specifies the size of the value based on the value of another field. It references the name of the field that defines the size. This is useful when the size of a value is not predefined, such as when the payload size is determined by the value in a len field.
- `type`: The type of the value and how it will be displayed. Possible values: [`int`, `string`, `bytes`]
- `encoding`: The decoding format for the value. Possible values: [`little-endian`, `big-endian`]

Below is an example of a format specification.

```json
[
    {
        "name": "crc",
        "size": 4,
        "type": "bytes",
        "encoding": "littleEndian"
    },
    {
        "name": "len",
        "size": 2,
        "encoding": "littleEndian",
        "type": "int"
    },
    {
        "name": "type",
        "size": 1,
        "encoding": "littleEndian",
        "type": "int"
    },
    {
        "name": "payload",
        "sizeRef": "len",
        "encoding": "littleEndian",
        "type": "string"
    }
]
```

## Example

```shell
$ fmtdump --format format.json data.bin
```

```txt
75 59 96 0d             |     crc: [117 89 150 13]
07 00                   |     len: 7
00                      |    type: 0
65 6e 74 72 79 2d 31    | payload: entry-1

cf 08 9f 94             |     crc: [207 8 159 148]
07 00                   |     len: 7
00                      |    type: 0
65 6e 74 72 79 2d 32    | payload: entry-2

59 38 98 e3             |     crc: [89 56 152 227]
07 00                   |     len: 7
00                      |    type: 0
65 6e 74 72 79 2d 33    | payload: entry-3
```
