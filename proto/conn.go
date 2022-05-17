package proto

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"net"
	"time"
	"unicode"
)

const (
	// DefaultReadBufferSize is maximum size of read buffer.
	// Use SetReadBufferSize to override.
	DefaultReadBufferSize = 4096
	// DefaultDeadline is default connection deadline
	// that is used when none is set.
	// Use SetReadDeadline, SetWriteDeadline or SetDeadline to override.
	DefaultPort     = 54321
	DefaultDeadline = 5 * time.Second
	TokenLength     = 16
	proto           = "udp"
)

// Conn is protocol connection.
// Implements net.Conn interface https://pkg.go.dev/net#Conn.
type Conn struct {
	token          []byte
	readBufferSize int
	// requestID      int
	conn          net.Conn
	keys          deviceKeys
	readDeadline  bool
	writeDeadline bool
}

// Dial connects to the device with given token.
// Token should be 16 bytes in lenght.
//
// Example:
//   Dial("192.168.0.3:54321", []byte{...})
//   Dial("192.168.0.3:54321", nil)
func Dial(addr string, token []byte) (Conn, error) {
	conn, err := net.Dial(proto, fmt.Sprintf("%s:%d", addr, DefaultPort))
	if err != nil {
		return Conn{}, err
	}

	return Conn{
		token:          token,
		readBufferSize: DefaultReadBufferSize,
		conn:           conn,
		keys:           newDeviceKeys(token),
	}, nil
}

// Read reads data from the connection.
// Read can be made to time out and return an error after a fixed
// time limit; see SetDeadline and SetReadDeadline.
func (c *Conn) Read(b []byte) (int, error) {
	buff := make([]byte, c.readBufferSize)

	if !c.readDeadline {
		c.conn.SetDeadline(time.Now().Add(DefaultDeadline))
	}

	n, err := c.conn.Read(buff)
	if err != nil {
		return n, err
	}

	decrypted := c.keys.decrypt(buff[32:n])

	// trim non-printable characters
	decrypted = bytes.TrimFunc(decrypted, func(r rune) bool {
		return !unicode.IsGraphic(r)
	})

	copy(b, decrypted)

	return len(decrypted), nil
}

// Write writes data to the connection.
// Write can be made to time out and return an error after a fixed
// time limit; see SetDeadline and SetWriteDeadline.
func (c *Conn) Write(b []byte) (int, error) {
	h, err := c.handshake()
	if err != nil {
		return 0, err
	}

	if c.Token() == "" {
		c.SetToken(h.Token())
	}

	encrypted := c.keys.encrypt(b)
	req := prepareRequest(c.token, h, encrypted)

	if !c.writeDeadline {
		c.conn.SetDeadline(time.Now().Add(DefaultDeadline))
	}

	return c.conn.Write(req)
}

// Close closes the connection.
// Any blocked Read or Write operations will be unblocked and return errors.
func (c *Conn) Close() error {
	return c.conn.Close()
}

// LocalAddr returns the local network address.
func (c *Conn) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

// RemoteAddr returns the remote network address.
func (c *Conn) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

// SetDeadline sets the read and write deadlines associated
// with the connection. It is equivalent to calling both
// SetReadDeadline and SetWriteDeadline.
//
// A deadline is an absolute time after which I/O operations
// fail instead of blocking. The deadline applies to all future
// and pending I/O, not just the immediately following call to
// Read or Write. After a deadline has been exceeded, the
// connection can be refreshed by setting a deadline in the future.
//
// If the deadline is exceeded a call to Read or Write or to other
// I/O methods will return an error that wraps os.ErrDeadlineExceeded.
// This can be tested using errors.Is(err, os.ErrDeadlineExceeded).
// The error's Timeout method will return true, but note that there
// are other possible errors for which the Timeout method will
// return true even if the deadline has not been exceeded.
//
// An idle timeout can be implemented by repeatedly extending
// the deadline after successful Read or Write calls.
//
// A zero value for t means I/O operations will not time out.
func (c *Conn) SetDeadline(t time.Time) error {
	c.readDeadline = true
	c.writeDeadline = true
	return c.conn.SetDeadline(t)
}

// SetReadDeadline sets the deadline for future Read calls
// and any currently-blocked Read call.
// A zero value for t means Read will not time out.
func (c *Conn) SetReadDeadline(t time.Time) error {
	c.readDeadline = true
	return c.conn.SetReadDeadline(t)
}

// SetWriteDeadline sets the deadline for future Write calls
// and any currently-blocked Write call.
// Even if write times out, it may return n > 0, indicating that
// some of the data was successfully written.
// A zero value for t means Write will not time out.
func (c *Conn) SetWriteDeadline(t time.Time) error {
	c.writeDeadline = true
	return c.conn.SetWriteDeadline(t)
}

// SetReadBufferSize sets read buffer
func (c *Conn) SetReadBufferSize(size int) {
	c.readBufferSize = size
}

// Token shows string representation of token
// used by connection
func (c *Conn) Token() string {
	return fmt.Sprintf("%x", c.token)
}

// SetToken sets token
func (c *Conn) SetToken(token string) {
	tokenBytes, err := hex.DecodeString(token)
	if err != nil {
		println(err)
	}

	tokenLength := len(tokenBytes)
	if tokenLength != TokenLength {
		println(fmt.Errorf("expected 16 bytes token, got %d", tokenLength))
	}

	c.token = tokenBytes
	c.keys = newDeviceKeys(tokenBytes)
}

type handshakeResponse struct {
	hello       []byte
	deviceID    []byte
	serverStamp []byte
	token       []byte
}

func (h *handshakeResponse) ServerStamp() time.Duration {
	return time.Duration(binary.BigEndian.Uint32(h.serverStamp)) * time.Second
}

func (h *handshakeResponse) Token() string {
	return fmt.Sprintf("%x", h.token)
}

func (h handshakeResponse) String() string {
	return fmt.Sprintf("{helo:%x deviceID:%x serverStamp:%v token:%s}", h.hello, h.deviceID, h.ServerStamp(), h.Token())
}

func parseHandshakeResponse(data []byte) (*handshakeResponse, error) {
	length := len(data)
	if length != 32 {
		return nil, fmt.Errorf("unable to parse handshake, expected 32 bytes, got %d", length)
	}

	return &handshakeResponse{
		hello:       data[:8],
		deviceID:    data[8:12],
		serverStamp: data[12:16],
		token:       data[16:],
	}, nil
}

func (c *Conn) handshake() (*handshakeResponse, error) {
	hello := []byte{
		// Magic number
		0x21, 0x31,
		// Length
		0x00, 0x20,
		// All the Fs
		0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff,
	}

	_, err := c.conn.Write(hello)
	if err != nil {
		return nil, err
	}

	resp := make([]byte, c.readBufferSize)
	n, err := c.conn.Read(resp)
	if err != nil {
		return nil, err
	}

	return parseHandshakeResponse(resp[:n])
}

func prepareRequest(token []byte, handshake *handshakeResponse, requestBody []byte) []byte {
	header := [32]byte{0x21, 0x31}
	binary.BigEndian.PutUint16(header[2:], uint16(32+len(requestBody)))
	binary.BigEndian.PutUint32(header[8:], binary.BigEndian.Uint32(handshake.deviceID))
	binary.BigEndian.PutUint32(header[12:], binary.BigEndian.Uint32(handshake.serverStamp))
	checksum := md5(header[:16], token, requestBody)
	copy(header[16:], checksum)
	return append(header[:], requestBody...)
}
