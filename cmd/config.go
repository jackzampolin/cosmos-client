package cmd

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/strangelove-ventures/lens/client"
	"gopkg.in/yaml.v2"
)

//createConfig idempotently creates the config.
func createConfig(home string, debug bool) error {
	cfgPath := path.Join(home, "config.yaml")

	// If the config doesn't exist...
	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		// And the config folder doesn't exist...
		// And the home folder doesn't exist
		if _, err := os.Stat(home); os.IsNotExist(err) {
			// Create the home folder
			if err = os.Mkdir(home, os.ModePerm); err != nil {
				return err
			}
		}
	}

	// Then create the file...
	content := defaultConfig(path.Join(home, "keys"), debug)
	if err := os.WriteFile(cfgPath, content, 0600); err != nil {
		return err
	}

	return nil
}

func overwriteConfig(v *viper.Viper, cfg *Config) error {
	home := v.GetString("home")
	cfgPath := path.Join(home, "config.yaml")
	if err := os.WriteFile(cfgPath, cfg.MustYAML(), 0600); err != nil {
		return err
	}

	log.Printf("updated lens configuration at %s", cfgPath)
	return nil
}

// Config represents the config file for the relayer
type Config struct {
	DefaultChain string                               `yaml:"default_chain" json:"default_chain"`
	Chains       map[string]*client.ChainClientConfig `yaml:"chains" json:"chains"`

	cl map[string]*client.ChainClient
}

func (c *Config) GetDefaultClient() *client.ChainClient {
	return c.GetClient(c.DefaultChain)
}

func (c *Config) GetClient(chainID string) *client.ChainClient {
	if v, ok := c.cl[chainID]; ok {
		return v
	}
	return nil
}

// Called to initialize the relayer.Chain types on Config
func validateConfig(c *Config) error {
	for _, chain := range c.Chains {
		if err := chain.Validate(); err != nil {
			return err
		}
	}
	if c.GetDefaultClient() == nil {
		return fmt.Errorf("default chain (%s) configuration not found", c.DefaultChain)
	}
	return nil
}

// MustYAML returns the yaml string representation of the Paths
func (c Config) MustYAML() []byte {
	out, err := yaml.Marshal(c)
	if err != nil {
		panic(err)
	}
	return out
}

func defaultConfig(keyHome string, debug bool) []byte {
	return Config{
		DefaultChain: "cosmoshub",
		Chains: map[string]*client.ChainClientConfig{
			"cosmoshub": client.GetCosmosHubConfig(keyHome, debug),
			"osmosis":   client.GetOsmosisConfig(keyHome, debug),
		},
	}.MustYAML()
}

// initConfig reads in config file and ENV variables if set.
// This is called as a persistent pre-run command of the root command.
func initConfig(cmd *cobra.Command, v *viper.Viper, lc *lensConfig) error {
	home, err := cmd.PersistentFlags().GetString(flags.FlagHome)
	if err != nil {
		return err
	}

	debug, err := cmd.Flags().GetBool("debug")
	if err != nil {
		return err
	}

	lc.config = Config{}
	cfgPath := path.Join(home, "config.yaml")
	_, err = os.Stat(cfgPath)
	if err != nil {
		err = createConfig(home, debug)
		if err != nil {
			return err
		}
	}
	v.SetConfigFile(cfgPath)
	err = v.ReadInConfig()
	if err != nil {
		return fmt.Errorf("failed to read in config: %w", err)
	}

	// read the config file bytes
	file, err := os.ReadFile(v.ConfigFileUsed())
	if err != nil {
		return fmt.Errorf("error reading config file: %w", err)
	}

	// unmarshall them into the struct
	if err = yaml.Unmarshal(file, &lc.config); err != nil {
		return fmt.Errorf("error unmarshalling config: %w", err)
	}

	// instantiate chain client
	// TODO: this is a bit of a hack, we should probably have a
	// better way to inject modules into the client
	lc.config.cl = make(map[string]*client.ChainClient)
	for name, chain := range lc.config.Chains {
		chain.Modules = append([]module.AppModuleBasic{}, ModuleBasics...)
		cl, err := client.NewChainClient(chain, home, cmd.InOrStdin(), cmd.OutOrStdout())
		if err != nil {
			return fmt.Errorf("error creating chain client: %w", err)
		}
		lc.config.cl[name] = cl
	}

	// override chain if needed
	if cmd.PersistentFlags().Changed("chain") {
		defaultChain, err := cmd.PersistentFlags().GetString("chain")
		if err != nil {
			return err
		}

		lc.config.DefaultChain = defaultChain
	}

	if cmd.PersistentFlags().Changed("output") {
		output, err := cmd.PersistentFlags().GetString("output")
		if err != nil {
			return err
		}

		// Should output be a global configuration item?
		for chain := range lc.config.Chains {
			lc.config.Chains[chain].OutputFormat = output
		}
	}

	// validate configuration
	if err := validateConfig(&lc.config); err != nil {
		return fmt.Errorf("error validating config: %w", err)
	}
	return nil
}
