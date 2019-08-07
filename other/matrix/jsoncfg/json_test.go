package jsoncfg

import "testing"

func TestJSON(t *testing.T) {
	type T1 struct {
		C64  complex64
		C128 complex128
	}

	//	data, err := json.Marshal(T1{})
	//	if err != nil {
	//		t.Fatal(err)
	//	}
	//	fmt.Println(string(data))
}
