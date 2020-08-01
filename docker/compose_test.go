package docker

import (
	"fmt"
	"path"
	"reflect"
	"testing"

	"github.com/lyimmi/laradock-manager/vuex"
)

func TestReturnResponse(t *testing.T) {
	t.Run("returnResponse(true)", func(t *testing.T) {
		message := "testing to true"
		got := returnResponse(true, message)
		expected := `{"success":true,"message":"` + message + `"}`
		if got != expected {
			t.Errorf("returnResponse(true, "+message+") = %s; want %s", got, expected)
		}
	})
	t.Run("returnResponse(false)", func(t *testing.T) {
		message := "testing to false"
		got := returnResponse(false, message)
		expected := `{"success":false,"message":"` + message + `"}`
		if got != expected {
			t.Errorf("returnResponse(false, "+message+") = %s; want %s", got, expected)
		}
	})
}

func newTestVuex() *vuex.State {
	vuexState := &vuex.State{}
	vuexState.StorePath = path.Join("testdata", "laradock-manager-vuex.json")
	return vuexState
}

func newTestCompose(vxs *vuex.State) *Compose {
	compose := &Compose{}
	compose.vuexState = vxs
	return compose
}

func TestNewDockerCompose(t *testing.T) {
	vxs := newTestVuex()
	expected := newTestCompose(vxs)
	got := NewDockerCompose(vxs)
	fmt.Println(reflect.DeepEqual(got, expected))
	if !reflect.DeepEqual(got, expected) {
		t.Error("NewDockerCompose() = {}; want {}", *got, expected)
	}
}

func TestSetLaradockPath(t *testing.T) {
	vxs := newTestVuex()
	c := newTestCompose(vxs)
	got := c.SetLaradockPath("path")
	if !got {
		t.Errorf("SetLaradockPath() = %t; want %t", got, true)
	}
}

func TestSetTerminalPath(t *testing.T) {
	vxs := newTestVuex()
	c := newTestCompose(vxs)
	got := c.SetTerminalPath("path")
	if !got {
		t.Errorf("SetTerminalPath() = %t; want %t", got, true)
	}
}
