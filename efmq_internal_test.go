package efmq

import (
	"net"
	"os"
	"reflect"
	"testing"
)

func TestConnect(t *testing.T) {
	const goodInterface = "en1" // Mac
	// set network interface
	var ni *net.Interface
	ni, err := net.InterfaceByName(goodInterface)
	if err != nil {
		t.Fatal(err)
	}
	// connect
	conn, err := connect(ni) // error return unchecked - fails because we don't have permmissions
	if err != nil {
		if os.IsPermission(err) {
			t.Skip(err)
		} else {
			t.Fatal(err)
		}
	}
	check := reflect.TypeOf(conn).String()
	if check != "*net.PacketConn" {
		t.Errorf("Testconnect: Expect *net.PacketConn, got %v", check)
	}
}
