# fmtdump

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
