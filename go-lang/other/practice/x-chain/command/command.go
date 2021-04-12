package command

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/ironzhang/practice/x-chain/blockchain"
)

func fprintBlock(w io.Writer, b blockchain.Block) {
	b.Write(w)
	fmt.Fprintf(w, "PoW: %t\n", b.Validate())
	fmt.Fprintln(w)
}

func doAddBlockCommand(chain *blockchain.Chain, args []string) (err error) {
	var data string
	var address string

	fs := flag.NewFlagSet("AddBlock", flag.ExitOnError)
	fs.StringVar(&data, "data", "", "block data")
	fs.StringVar(&address, "address", "", "address")
	fs.Parse(args)

	tx := blockchain.NewCoinbaseTX(address, data)
	txs := []*blockchain.Transaction{tx}
	if err = chain.AddBlock(txs); err != nil {
		return err
	}
	if block, err := chain.LastBlock(); err != nil {
		return err
	} else {
		fprintBlock(os.Stdout, block)
	}
	return nil
}

func doPrintBlocksCommand(chain *blockchain.Chain, args []string) error {
	var ok bool
	block, err := chain.LastBlock()
	if err != nil {
		return err
	}
	for {
		fprintBlock(os.Stdout, block)
		block, ok, err = chain.PrevBlock(block)
		if err != nil {
			return err
		}
		if !ok {
			break
		}
	}
	return nil
}

func doGetBalanceCommand(chain *blockchain.Chain, args []string) error {
	var address string

	fs := flag.NewFlagSet("GetBalance", flag.ExitOnError)
	fs.StringVar(&address, "address", "", "account address")
	fs.Parse(args)

	balance := 0
	utxos, err := chain.FindUTXO(address)
	if err != nil {
		return err
	}
	for _, utxo := range utxos {
		balance += utxo.Value
	}

	fmt.Fprintf(os.Stdout, "Balance of %q: %d\n", address, balance)
	return nil
}

func doSendCommand(chain *blockchain.Chain, args []string) error {
	var from, to, mine string
	var amount int

	fs := flag.NewFlagSet("Send", flag.ExitOnError)
	fs.StringVar(&from, "from", "", "from address")
	fs.StringVar(&to, "to", "", "to address")
	fs.StringVar(&mine, "mine", "", "mine address")
	fs.IntVar(&amount, "amount", 0, "amount")
	fs.Parse(args)

	coinbase := blockchain.NewCoinbaseTX(mine, "")
	transaction, err := chain.NewTransaction(from, to, amount)
	if err != nil {
		return err
	}
	txs := []*blockchain.Transaction{coinbase, transaction}
	if err = chain.AddBlock(txs); err != nil {
		return err
	}
	if block, err := chain.LastBlock(); err != nil {
		return err
	} else {
		fprintBlock(os.Stdout, block)
	}
	return nil
}

var commands = map[string]func(chain *blockchain.Chain, args []string) error{
	"AddBlock":    doAddBlockCommand,
	"PrintBlocks": doPrintBlocksCommand,
	"GetBalance":  doGetBalanceCommand,
	"Send":        doSendCommand,
}

func usage() {
	fmt.Fprintln(os.Stderr, "Usage:")
	fmt.Fprintln(os.Stderr, "  Addblock -data BLOCK_DATA - add a block to the blockchain")
	fmt.Fprintln(os.Stderr, "  PrintBlocks - print all the blocks of the blockchain")
	fmt.Fprintln(os.Stderr, "  GetBalance - print balance of address")
	fmt.Fprintln(os.Stderr, "  Send - send balance")
}

type Executer struct {
	chain *blockchain.Chain
}

func NewExecuter(chain *blockchain.Chain) *Executer {
	return &Executer{chain: chain}
}

func (p *Executer) Execute(args []string) {
	if len(args) <= 0 {
		fmt.Fprintf(os.Stderr, "invalid args")
		return
	}

	fs := flag.NewFlagSet(args[0], flag.ExitOnError)
	fs.Usage = usage
	fs.Parse(args[1:])

	fargs := fs.Args()
	if len(fargs) <= 0 {
		fs.Usage()
		return
	}

	name := fargs[0]
	f, ok := commands[name]
	if !ok {
		fs.Usage()
		return
	}

	defer func() {
		if e := recover(); e != nil {
			fmt.Fprintf(os.Stderr, "do %s command: %v\n", name, e)
		}
	}()

	if err := f(p.chain, fargs[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "do %s command: %v\n", name, err)
		return
	}
}
