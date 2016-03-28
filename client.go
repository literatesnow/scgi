package scgi

import (
	"bytes"
	"io"
	"net"
	"strconv"
)

type Client struct {
	network string
	address string
	headers map[string]string
}

//NewClient creates a new SCGI client with default headers which will connect to a specified address.
func NewClient(network string, address string) *Client {
	cl := &Client{network: network, address: address}

	//Default headers.
	cl.headers = map[string]string{
		"SCGI":           "1",
		"REQUEST_METHOD": "POST"}

	return cl
}

//Request sends a request to the server and reads the response.
func (cl *Client) Request(request []byte) (response *bytes.Buffer, err error) {
	conn, err := net.Dial(cl.network, cl.address)
	if err != nil {
		return nil, err
	}

	err = cl.writeRequest(request, conn)
	if err != nil {
		return nil, err
	}

	response, err = cl.readResponse(conn)
	if err != nil {
		return nil, err
	}

	return response, nil
}

//writeRequest creates the underlying request and writes it to the socket.
func (cl *Client) writeRequest(body []byte, conn net.Conn) (err error) {
	header := cl.makeHeaders(len(body))

	conn.Write(cl.netstring(header, body))

	return nil
}

//readResponse reads the response from the server.
func (cl *Client) readResponse(conn net.Conn) (response *bytes.Buffer, err error) {
	var buf bytes.Buffer
	_, err = io.Copy(&buf, conn)
	return &buf, err
}

//makeHeaders creates the headers part of the request.
func (cl *Client) makeHeaders(bodyLen int) []byte {
	var headers = cl.appendHeader([]byte{}, "CONTENT_LENGTH", strconv.Itoa(bodyLen))

	for k, v := range cl.headers {
		headers = cl.appendHeader(headers, k, v)
	}

	return headers
}

//appendHeader appends a headers.
func (cl *Client) appendHeader(buf []byte, key string, value string) []byte {
	buf = append(buf, []byte(key)...)
	buf = append(buf, []byte{0}...)
	buf = append(buf, []byte(value)...)
	buf = append(buf, []byte{0}...)

	return buf
}

//netstring creates a formatted netstring.
func (cl *Client) netstring(headers []byte, body []byte) []byte {
	buf := []byte(strconv.Itoa(len(headers)))
	buf = append(buf, ":"...)
	buf = append(buf, headers...)
	buf = append(buf, ","...)
	buf = append(buf, body...)

	return buf
}

//SetHeader sets a header to a specific value
func (cl *Client) SetHeader(name string, value string) {
	cl.headers[name] = value
}
