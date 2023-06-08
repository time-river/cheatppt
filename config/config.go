package config

import (
	"crypto/rand"
	"encoding/base64"
	"flag"
	"io"
	"os"

	"github.com/pelletier/go-toml/v2"
)

func CmdlineParse() {

	cfg := flag.String("config", "./cheatppt.toml", "configuration file")
	help := flag.Bool("help", false, "Display help documentation")

	flag.Parse()

	if *help {
		flag.PrintDefaults()
		os.Exit(0)
	}

	err := cfgParse(cfg)
	if err != nil {
		panic(err)
	}

}

func cfgParse(cfg *string) (err error) {
	var writer *os.File

	_, err = os.Stat(*cfg)
	if err != nil && !os.IsNotExist(err) {
		return err
	} else if err == nil {
		reader, err := os.Open(*cfg)
		if err != nil {
			return err
		}
		err = toml.NewDecoder(reader).Decode(&GlobalCfg)
		if err != nil {
			return err
		}
	}

	writer, err = os.Create(*cfg)
	if err != nil {
		return err
	}
	defer writer.Close()

	return cfgGenerate(writer, &GlobalCfg)
}

func secretCheck(orig *string) (secret *string, err error) {
	var key []byte

	secret = nil
	if len(*orig) > 0 {
		key, err = base64.StdEncoding.DecodeString(*orig)
		if err != nil {
			return
		}
	}

	// HMACSHA256 need 32 bytes
	if len(key) != 32 {
		bytes := make([]byte, 32)
		_, err = rand.Read(bytes)
		if err != nil {
			return
		}
		key = bytes
	}

	str := base64.StdEncoding.EncodeToString(key)
	secret = &str
	return
}

func cfgGenerate(writer io.Writer, cfg *Cfg) error {
	orig := &GlobalCfg.Server.Secret
	secret, err := secretCheck(orig)
	if err != nil {
		return err
	}
	GlobalCfg.Server.Secret = *secret

	if err := toml.NewEncoder(writer).Encode(cfg); err != nil {
		return err
	}

	return nil
}
