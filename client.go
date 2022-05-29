package miio

import (
	"encoding/json"
	"log"
	"runtime"
	"sync"

	// "github.com/goccy/go-json"

	"github.com/icepie/miio.go/proto"
)

// Client is device client that extends protocol connection.
type Client struct {
	sync.Mutex
	proto.Conn
	requestID int
	Did       string
}

// New creates new device client.
//
// Example:
//   New("192.168.0.3")
func New(addr string) *Client {

	conn, err := proto.Dial(addr, nil)
	if err != nil {
		println(err.Error())
	}

	client := &Client{sync.Mutex{}, conn, 1, ""}
	runtime.SetFinalizer(client, func(c *Client) {
		err = c.Close()
		if err != nil {
			println(err.Error())
		}
	})

	return client
}

// Set Token for device.
func (c *Client) SetToken(token string) *Client {
	c.Conn.SetToken(token)
	return c
}

// Set Device ID for device.
func (c *Client) SetDid(did string) *Client {
	c.Did = did
	return c
}

// close connection.
func (c *Client) Close() error {
	return c.Conn.Close()
}

// Send sends request to device.
func (c *Client) Send(method string, params interface{}) ([]byte, error) {
	req := struct {
		RequestID int         `json:"id"`
		Method    string      `json:"method"`
		Params    interface{} `json:"params"`
	}{
		RequestID: c.requestID,
		Method:    method,
		Params:    params,
	}

	payload, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	if _, err := c.Write(payload); err != nil {
		return nil, err
	}

	resp := make([]byte, proto.DefaultReadBufferSize)
	n, err := c.Read(resp)
	if err != nil {
		return nil, err
	}

	if err == nil {
		c.Lock()
		c.requestID++
		c.Unlock()
	}

	return resp[:n], nil

	// // trim non-printable characters
	// return bytes.TrimFunc(resp[:n], func(r rune) bool {
	// 	return !unicode.IsGraphic(r)
	// }), err
}

// ConfigRouter configures wifi network on device.
func (c *Client) ConfigRouter(ssid string, passwd string, uid string) ([]byte, error) {
	v := struct {
		SSID   string `json:"ssid"`
		Passwd string `json:"passwd"`
		UID    string `json:"uid"`
	}{
		SSID:   ssid,
		Passwd: passwd,
		UID:    uid,
	}

	return c.Send("miIO.config_router", v)
}

func (c *Client) Info() ([]byte, error) {
	return c.Send("miIO.info", nil)
}

// OTAProgress requests OTA update progress.
func (c *Client) OTAProgress() ([]byte, error) {
	return c.Send("miIO.get_ota_progress", nil)
}

// OTAState requests available update for device.
func (c *Client) OTAState() ([]byte, error) {
	return c.Send("miIO.get_ota_state", nil)
}

// OTA updates the device.
func (c *Client) OTA(url string, fileMD5 string) ([]byte, error) {
	v := struct {
		Mode    string `json:"mode"`
		Install string `json:"install"`
		AppURL  string `json:"app_url"`
		FileMD5 string `json:"file_md5"`
		Proc    string `json:"proc"`
	}{
		Mode:    "normal",
		Install: "1",
		AppURL:  url,
		FileMD5: fileMD5,
		Proc:    "dnld install",
	}

	return c.Send("miIO.ota", v)
}

// GetProperties gets device propetriest.
func (c *Client) GetProps(params ...PropParam) ([]byte, error) {

	var fParams []PropParam

	for _, p := range params {
		fParams = append(fParams, *p.SetDid(c.Did))
	}

	return c.Send("get_properties", fParams)
}

// SetProperties sets device propetriest.
func (c *Client) SetProps(params ...PropParam) ([]byte, error) {

	var fParams []PropParam

	for _, p := range params {
		fParams = append(fParams, *p.SetDid(c.Did))
	}

	return c.Send("set_properties", fParams)
}

// Action execute device action.
func (c *Client) DoAction(param ActionParam) ([]byte, error) {
	param.SetDid(c.Did)

	log.Println(param)

	return c.Send("action", param)
}
