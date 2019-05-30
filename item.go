package libsecret

import (
  "time"

  "github.com/godbus/dbus"
)


type Item struct {
  conn *dbus.Conn
  dbus  dbus.BusObject
}


func NewItem(conn *dbus.Conn, path dbus.ObjectPath) *Item {
  return &Item{
    conn: conn,
    dbus: conn.Object(DBusServiceName, path),
  }
}


func (item Item) Path() dbus.ObjectPath {
  return item.dbus.Path()
}


// READWRITE String Label;
func (item *Item) Label() (string, error) {
  val, err := item.dbus.GetProperty("org.freedesktop.Secret.Item.Label")
  if err != nil {
    return "", err
  }

  return val.Value().(string), nil
}


// READ Boolean Locked;
func (item *Item) Locked() (bool, error) {
  val, err := item.dbus.GetProperty("org.freedesktop.Secret.Item.Locked")
  if err != nil {
    return true, err
  }

  return val.Value().(bool), nil
}


// READ Uint64 Created;
func (item *Item) Created() (time.Time, error) {
	val, err := item.dbus.GetProperty("org.freedesktop.Secret.Item.Created")
	if err != nil {
		return time.Time{}, err
	}
	v := val.Value().(uint64)
	tm := time.Unix(int64(v), 0)

	return tm, nil
}


// READ Uint64 Modified;
func (item *Item) Modified() (time.Time, error) {
	val, err := item.dbus.GetProperty("org.freedesktop.Secret.Item.Modified")
	if err != nil {
		return time.Time{}, err
	}
	v := val.Value().(uint64)
	tm := time.Unix(int64(v), 0)

	return tm, nil
}


// READWRITE Dict<String,String> Attributes;
func (item *Item) Attributes() (map[string]string, error) {
	attributes := make(map[string]string)
	val, err := item.dbus.GetProperty("org.freedesktop.Secret.Item.Attributes")
	if err != nil {
		return attributes, err
	}
	attributes = val.Value().(map[string]string)

	return attributes, nil
}


// GetSecret (IN ObjectPath session, OUT Secret secret);
func (item *Item) GetSecret(session *Session) (*Secret, error) {
  secret := Secret{}

  err := item.dbus.Call("org.freedesktop.Secret.Item.GetSecret", 0, session.Path()).Store(&secret)
  if err != nil {
    return &Secret{}, err
  }

  return &secret, nil
}


// Delete (OUT ObjectPath Prompt);
func (item *Item) Delete() error {
  var prompt dbus.ObjectPath

  err := item.dbus.Call("org.freedesktop.Secret.Item.Delete", 0).Store(&prompt)
  if err != nil {
    return err
  }

  if isPrompt(prompt) {
    prompt := NewPrompt(item.conn, prompt)
    if _, err := prompt.Prompt(); err != nil {
      return err
    }
  }

  return nil
}
