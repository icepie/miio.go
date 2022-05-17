# miio.go

Go implementation for MiIO protocol.
inpsired by [miio.dart](https://github.com/ctrysbita/miio.dart)

## Install

```sh
go get github.com/icepie/miio.go
```

## Protocol

```
0                   1                   2                   3
 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
| Magic Number = 0x2131         | Packet Length (incl. header)  |
|-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-|
| Unknown                                                       |
|-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-|
| Device ID ("did")                                             |
|-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-|
| Stamp                                                         |
|-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-|
| MD5 Checksum                                                  |
| ... or Device Token in response to the "Hello" packet         |
|-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-|
| Optional variable-sized payload (encrypted)                   |
|...............................................................|


Packet Length: 16 bits unsigned int
    Length in bytes of the whole packet, including header(0x20 bytes).

Unknown: 32 bits
    This value is always 0.
    0xFFFFFFFF in "Hello" packet.

Device ID: 32 bits
    Unique number. Possibly derived from the MAC address.
    0xFFFFFFFF in "Hello" packet.

Stamp: 32 bit unsigned int
    Continously increasing counter.
    Number of seconds since device startup.

MD5 Checksum:
    Calculated for the whole packet including the MD5 field itself,
    which must be initialized with token.

    In "Hello" packet,
    this field contains the 128-bit 0xFF.

    In the response to the first "Hello" packet,
    this field contains the 128-bit device token.

Optional variable-sized payload:
    Payload encrypted with AES-128-CBC (PKCS#7 padding).

        Key = MD5(Token)
        IV  = MD5(Key + Token)
```

## Link

- https://github.com/OpenMiHome/mihome-binary-protocol
- https://github.com/rytilahti/python-miio
- https://home.miot-spec.com/
