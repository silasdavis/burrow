package config

import (
	"fmt"

	acm "github.com/hyperledger/burrow/account"
	"github.com/hyperledger/burrow/config/source"
	"github.com/hyperledger/burrow/consensus/tendermint"
	"github.com/hyperledger/burrow/consensus/tendermint/validator"
	"github.com/hyperledger/burrow/core"
	"github.com/hyperledger/burrow/genesis"
	"github.com/hyperledger/burrow/genesis/spec"
	"github.com/hyperledger/burrow/keys"
	logging_config "github.com/hyperledger/burrow/logging/config"
	"github.com/hyperledger/burrow/logging/lifecycle"
	"github.com/hyperledger/burrow/logging/loggers"
)

const DefaultBurrowConfigTOMLFileName = "burrow.toml"
const DefaultBurrowConfigJSONEnvironmentVariable = "BURROW_CONFIG_JSON"
const DefaultGenesisDocJSONFileName = "genesis.json"

type BurrowConfig struct {
	ValidatorAddress *acm.Address                       `json:",omitempty" toml:",omitempty"`
	GenesisDoc       *genesis.GenesisDoc                `json:",omitempty" toml:",omitempty"`
	Tendermint       *tendermint.BurrowTendermintConfig `json:",omitempty" toml:",omitempty"`
	Keys             *keys.KeysConfig                   `json:",omitempty" toml:",omitempty"`
	Logging          *logging_config.LoggingConfig      `json:",omitempty" toml:",omitempty"`
}

func DefaultBurrowConfig() *BurrowConfig {
	return &BurrowConfig{
		Tendermint: tendermint.DefaultBurrowTendermintConfig(),
		Keys:       keys.DefaultKeysConfig(),
		Logging:    logging_config.DefaultNodeLoggingConfig(),
	}
}

func (conf *BurrowConfig) Kernel() (*core.Kernel, error) {
	if conf.GenesisDoc == nil {
		return nil, fmt.Errorf("no GenesisDoc defined in config, cannot make Kernel")
	}
	if conf.ValidatorAddress == nil {
		return nil, fmt.Errorf("no validator address in config, cannot make Kernel")
	}
	logger, err := lifecycle.NewLoggerFromLoggingConfig(conf.Logging)
	if err != nil {
		return nil, fmt.Errorf("could not generate logger from logging config: %v", err)
	}
	keyClient := keys.NewBurrowKeyClient(conf.Keys.URL, logger)
	val, err := keys.Addressable(keyClient, *conf.ValidatorAddress)
	if err != nil {
		return nil, fmt.Errorf("could not get validator addressable from keys client: %v", err)
	}
	privValidator := validator.NewPrivValidatorMemory(val, keys.Signer(keyClient, val.Address()))

	return core.NewKernel(privValidator, conf.GenesisDoc, conf.Tendermint.TendermintConfig(), logger)
}

func (conf *BurrowConfig) JSONString() string {
	return source.JSONString(conf)
}

func (conf *BurrowConfig) TOMLString() string {
	return source.TOMLString(conf)
}

func (conf *BurrowConfig) RealiseGenesisSpec(genesisSpec *spec.GenesisSpec) (*genesis.GenesisDoc, error) {
	if conf.Keys == nil {
		return nil, fmt.Errorf("keys config required for GenesisSpecConfigProvider")
	}
	keyClient := keys.NewBurrowKeyClient(conf.Keys.URL, loggers.NewNoopInfoTraceLogger())
	genesisDoc, err := genesisSpec.GenesisDoc(keyClient)
	if err != nil {
		return nil, fmt.Errorf("could not create GenesisDoc from GenesisSpec in config provider: %v", err)
	}
	return genesisDoc, nil

}
