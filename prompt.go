package libsecret

import (
  "fmt"
  "github.com/godbus/dbus"
  "strings"
)


type Prompt struct {
  conn *dbus.Conn
  dbus  dbus.BusObject
}


func NewPrompt(conn *dbus.Conn, path dbus.ObjectPath) *Prompt {
  return &Prompt{
    conn: conn,
    dbus: conn.Object(DBusServiceName, path),
  }
}


func (prompt Prompt) Path() dbus.ObjectPath {
  return prompt.dbus.Path()
}


func isPrompt(path dbus.ObjectPath) bool {
  promptPath := DBusPath + "/prompt/"
  return strings.HasPrefix(string(path), promptPath)
}


// Prompt (IN String window-id);
func (prompt *Prompt) Prompt() (*dbus.Variant, error) {
  // Secret Service API specification doesn't specify whether signals are broadcast or unicast so it depends on the specific server implementation
  // We need to explicitly connect to signal to handle both cases
  args := fmt.Sprintf("type='signal',path='%s',interface='%s',sender='%s'", prompt.Path(), "org.freedesktop.Secret.Prompt", DBusServiceName)
  err := prompt.conn.BusObject().Call("org.freedesktop.DBus.AddMatch", 0, args).Store()
  if err != nil {
    return &dbus.Variant{}, err
  }

  // prompts are asynchronous so we connect to the signal
  // and block with a channel until we get a response
  c := make(chan *dbus.Signal, 10)
  defer close(c)

  prompt.conn.Signal(c)
  defer prompt.conn.RemoveSignal(c)

  err = prompt.dbus.Call("org.freedesktop.Secret.Prompt.Prompt", 0, "").Store()
  if err != nil {
    return &dbus.Variant{}, err
  }

  for {
    if result := <-c; result.Path == prompt.Path() {
      value := result.Body[1].(dbus.Variant)
      return &value, nil
    }
  }
}
