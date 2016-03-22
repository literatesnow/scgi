package scgi

import (
  "os"
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

  /*
  cl.SetHeader("SCGI", "1")
  cl.SetHeader("SERVER_PROTOCOL", "HTTP/1.1")
  cl.SetHeader("REQUEST_METHOD", "POST")
  */

  cl.SetHeader("SCGI", "1")
  cl.SetHeader("REQUEST_METHOD", "POST")
  cl.SetHeader("REQUEST_URI", "/deepthought")

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
  //var buf bytes.Buffer

  var headerPart = new(bytes.Buffer)
  var bodyPart = []byte(body)
  var sz = 0

  sz += cl.appendHeader(headerPart, "CONTENT_LENGTH", strconv.Itoa(len(bodyPart)))

  for k, v := range cl.headers {
    fmt.Printf("-- %s -> %s\n", k, v)
    sz += cl.appendHeader(headerPart, k, v)
  }

  headerPart.WriteByte(44) //,

  fmt.Printf("!!!!!\n")

  fmt.Printf("%d:", sz)
  headerPart.WriteTo(os.Stdout)
  fmt.Printf("%s", bodyPart);
  fmt.Printf("\n")

  return "nothing"
}

func (cl *Client) appendHeader(buf *bytes.Buffer, key string, value string) int {
  k := []byte(key)
  v := []byte(value)
  buf.Write(k)
  buf.WriteByte(95)
  buf.Write(v)
  buf.WriteByte(95)

  return len(k) + len(v) + 2
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
