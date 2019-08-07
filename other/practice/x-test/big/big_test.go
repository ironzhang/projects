package big

import (
	"encoding/json"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

func TestBigRat(t *testing.T) {
	r := big.NewRat(100, 1000)
	fmt.Printf("string: %s\n", r.String())

	data, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("json: %s\n", data)

	if err = json.Unmarshal(data, r); err != nil {
		t.Fatal(err)
	}
	fmt.Printf("rat: %s\n", r)

	data = []byte(`"1/1000000000000000000000000000000000000000"`)
	if err = json.Unmarshal(data, r); err != nil {
		t.Fatal(err)
	}
	fmt.Printf("rat: %s\n", r)

	r.SetString("1/1000000000000000000000000000000000000")
	fmt.Printf("rat: %s\n", r)

	r.SetInt(big.NewInt(1000000))
	fmt.Printf("rat: %s\n", r)

	r = nil
	data, err = json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("json: %s\n", data)
}

func TestBigInt(t *testing.T) {
	i := big.NewInt(0)
	fmt.Printf("%#d\n", i)
}

func TestHexBig(t *testing.T) {
	var a = struct {
		A *hexutil.Big
		B *big.Int
	}{
		A: (*hexutil.Big)(big.NewInt(0)),
		B: big.NewInt(0),
	}

	data, err := json.Marshal(a)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%s\n", data)
}

type T struct {
	A *big.Int
	B *big.Int
}

func TestIntJSON(t *testing.T) {
	var a = T{
		A: big.NewInt(0),
		B: big.NewInt(0),
	}
	a.B.SetString("10000000000000000000000001234", 10)
	data, err := json.Marshal(a)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%s\n", data)
}
