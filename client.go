package scgi

import (
  "io"
  "net"
  "fmt"
  "bytes"
  "strconv"
)

type Client struct {
  headers map[string]string
  conn net.Conn
}

func NewClient() *Client {
  cl := &Client{ headers: make(map[string]string) }

  cl.SetHeader("SCGI", "1")
  cl.SetHeader("SERVER_PROTOCOL", "HTTP/1.1")
  cl.SetHeader("REQUEST_METHOD", "POST")

  return cl
}

func (cl *Client) Connect(network string, address string) (err error) {
  cl.conn, err = net.Dial(network, address)
  if err != nil {
    return err
  }

  fmt.Printf("conn: %s\n", cl.conn);

  return nil
}

func (cl *Client) Send(body string) string {
  var bodyPart = []byte(body)
  var headerPart = cl.requestHeader(len(bodyPart))
  var buf = cl.netstring(headerPart, bodyPart)

  _, err := cl.conn.Write(buf)
  if err != nil {
    fmt.Printf("%s\n", err);
    return ""
  }

  var fub bytes.Buffer
  io.Copy(&fub, cl.conn)

  fmt.Printf("%d: %s\n", fub.Len(), fub.Bytes());

  return "nothing"
}

func (cl *Client) requestHeader(bodyLen int) []byte {
  var headerPart = cl.appendHeader([]byte{}, "CONTENT_LENGTH", strconv.Itoa(bodyLen))

  for k, v := range cl.headers {
    headerPart = cl.appendHeader(headerPart, k, v)
  }

  return headerPart
}

func (cl *Client) appendHeader(buf []byte, key string, value string) []byte {
  buf = append(buf, []byte(key) ...)
  buf = append(buf, []byte{0} ...)
  buf = append(buf, []byte(value) ...)
  buf = append(buf, []byte{0} ...)

  return buf
}

func (cl *Client) netstring(header []byte, body []byte) []byte {
  buf := []byte(strconv.Itoa(len(header)))
  buf = append(buf, ":" ...)
  buf = append(buf, header ...)
  buf = append(buf, "," ...)
  buf = append(buf, body ...)

  return buf
}

func (cl *Client) SetHeader(name string, value string) {
  cl.headers[name] = value
}

func (cl *Client) Close() {
  if cl.conn != nil {
    cl.conn.Close()
    cl.conn = nil
  }
}

/*
"70:"
    "CONTENT_LENGTH" <00> "27" <00>
    "SCGI" <00> "1" <00>
    "REQUEST_METHOD" <00> "POST" <00>
    "REQUEST_URI" <00> "/deepthought" <00>
","
"What is the answer to life?"
*/
