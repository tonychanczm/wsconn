# WebsocketConnection
将websocket.Conn转换为net.Conn!

```go
// 没错，就只有这个
func NewWSConn(conn *websocket.Conn) net.Conn 
// 调用 Close() error 方法会自动发送 close message
// 调用 ForceClose() error 则不会
```
