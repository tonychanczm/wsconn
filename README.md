# wsconn
将websocket.Conn转换为net.Conn!

```go
// 没错，就只有这个
func NewWSConn(conn *websocket.Conn) net.Conn 
```
