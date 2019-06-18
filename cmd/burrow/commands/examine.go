package commands

import (
	"encoding/json"

	"github.com/hyperledger/burrow/bcm"

	"github.com/hyperledger/burrow/txs"
	cli "github.com/jawher/mow.cli"
	"github.com/tendermint/tendermint/libs/db"
)

func Examine(output Output) func(cmd *cli.Cmd) {
	return func(dump *cli.Cmd) {
		configOpts := addConfigOptions(dump)

		var explorer *bcm.BlockStore

		dump.Before = func() {
			conf, err := configOpts.obtainBurrowConfig()
			if err != nil {
				output.Fatalf("Could not obtain config: %v", err)
			}
			tmConf, err := conf.TendermintConfig()
			if err != nil {
				output.Fatalf("Could not build Tendermint config:", err)
			}

			explorer = bcm.NewBlockExplorer(db.DBBackendType(tmConf.DBBackend), tmConf.DBDir())
		}

		dump.Command("blocks", "dump blocks to stdout", func(cmd *cli.Cmd) {
			rangeArg := cmd.StringArg("RANGE", "", "Range as START_HEIGHT:END_HEIGHT where omitting "+
				"either endpoint implicitly describes the start/end and a negative index counts back from the last block")

			cmd.Spec = "[RANGE]"

			cmd.Action = func() {
				start, end, err := parseRange(*rangeArg)

				err = explorer.Blocks(start, end,
					func(block *bcm.Block) error {
						bs, err := json.Marshal(block)
						if err != nil {
							output.Fatalf("Could not serialise block: %v", err)
						}
						output.Printf(string(bs))
						return nil
					})
				if err != nil {
					output.Fatalf("Error iterating over blocks: %v", err)
				}
			}
		})

		dump.Command("txs", "dump transactions to stdout", func(cmd *cli.Cmd) {
			rangeArg := cmd.StringArg("RANGE", "", "Range as START_HEIGHT:END_HEIGHT where omitting "+
				"either endpoint implicitly describes the start/end and a negative index counts back from the last block")

			cmd.Spec = "[RANGE]"

			cmd.Action = func() {
				start, end, err := parseRange(*rangeArg)

				err = explorer.Blocks(start, end,
					func(block *bcm.Block) error {
						err := block.Transactions(func(txEnv *txs.Envelope) error {
							wrapper := struct {
								Height int64
								Tx     *txs.Envelope
							}{
								Height: block.Height,
								Tx:     txEnv,
							}
							bs, err := json.Marshal(wrapper)
							if err != nil {
								output.Fatalf("Could not deserialise transaction: %v", err)
							}
							output.Printf(string(bs))
							return nil
						})
						if err != nil {
							output.Fatalf("Error iterating over transactions: %v", err)
						}
						// If we stopped transactions stop everything
						return nil
					})
				if err != nil {
					output.Fatalf("Error iterating over blocks: %v", err)
				}
			}
		})
	}
}
