# miio.go

Go implementation for MiIO protocol.
inpsired by [miio.dart](https://github.com/ctrysbita/miio.dart)

## Install

```sh
go get github.com/icepie/miio.go
```

## Usage

> for get device token and device id, you can try [micloud](https://github.com/icepie/micloud)

```go
device := miio.New("192.168.1.12").
	SetToken("ffffffffffffffffffffffffffffffff").
	SetDid("did") // some device must use the true device id

// get device info
info, err := device.Info()
if err != nil {
	panic(err)
}

fmt.Printf("%s\n", info)

// for get siid and piid, see https://home.miot-spec.com/
// get properties
getProps, err := device.GetProps(miio.PropParam{
	Siid: 5,
	Piid: 1,
}, miio.PropParam{
	Siid: 5,
	Piid: 2,
})

if err != nil {
	panic(err)
}

fmt.Printf("%s\n", getProps)

// set properties
setPropsParams := []miio.PropParam{
	{
		Siid:  5,
		Piid:  1,
		Value: true,
	},
	{
		Siid:  5,
		Piid:  2,
		Value: 0,
	},
}

setProps, err := device.SetProps(setPropsParams...)
if err != nil {
	panic(err)
}

// Do action
action, err := device.DoAction(miio.ActionParam{
	Siid: 2,
	Aiid: 1,
	// In:   []any{0, "test"}, if need to send extra data
})

if err != nil {
	panic(err)
}

fmt.Printf("%s\n", action)

// some old devices can not use the standard way to control...
// https://github.com/icepie/miio.go/issues/1
// you can try to send raw data, like:
getProps, err := device.Send("get_prop", []interface{}{"power", "mode", "fan_level", "ver_swing"})
setProps, err := device.Send("set_mode", []interface{}{"cool"})


// or you can try: https://github.com/icepie/micloud

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
