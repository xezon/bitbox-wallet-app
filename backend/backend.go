// Copyright 2018 Shift Devices AG
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package backend

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/cloudfoundry-attic/jibber_jabber"
	"github.com/digitalbitbox/bitbox-wallet-app/backend/accounts"
	"github.com/digitalbitbox/bitbox-wallet-app/backend/arguments"
	"github.com/digitalbitbox/bitbox-wallet-app/backend/coins/btc"
	"github.com/digitalbitbox/bitbox-wallet-app/backend/coins/btc/electrum"
	"github.com/digitalbitbox/bitbox-wallet-app/backend/coins/btc/electrum/client"
	"github.com/digitalbitbox/bitbox-wallet-app/backend/coins/coin"
	"github.com/digitalbitbox/bitbox-wallet-app/backend/coins/eth"
	"github.com/digitalbitbox/bitbox-wallet-app/backend/coins/ltc"
	"github.com/digitalbitbox/bitbox-wallet-app/backend/config"
	"github.com/digitalbitbox/bitbox-wallet-app/backend/devices/device"
	"github.com/digitalbitbox/bitbox-wallet-app/backend/devices/usb"
	"github.com/digitalbitbox/bitbox-wallet-app/backend/keystore"
	"github.com/digitalbitbox/bitbox-wallet-app/backend/keystore/software"
	"github.com/digitalbitbox/bitbox-wallet-app/backend/signing"
	"github.com/digitalbitbox/bitbox-wallet-app/util/errp"
	"github.com/digitalbitbox/bitbox-wallet-app/util/jsonrpc"
	"github.com/digitalbitbox/bitbox-wallet-app/util/locker"
	"github.com/digitalbitbox/bitbox-wallet-app/util/logging"
	"github.com/digitalbitbox/bitbox-wallet-app/util/observable"
	"github.com/digitalbitbox/bitbox-wallet-app/util/rpc"
	"github.com/ethereum/go-ethereum/params"
	"github.com/sirupsen/logrus"
	"golang.org/x/text/language"
)

const (
	coinBTC  = "btc"
	coinTBTC = "tbtc"
	coinLTC  = "ltc"
	coinTLTC = "tltc"
	coinETH  = "eth"
	coinTETH = "teth"
	coinRETH = "reth"
)

type backendEvent struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

type deviceEvent struct {
	DeviceID string `json:"deviceID"`
	Type     string `json:"type"`
	Data     string `json:"data"`
	// TODO: rename Data to Event, Meta to Data.
	Meta interface{} `json:"meta"`
}

// AccountEvent models an event triggered by an account.
type AccountEvent struct {
	Type string `json:"type"`
	Code string `json:"code"`
	Data string `json:"data"`
}

// Backend ties everything together and is the main starting point to use the BitBox wallet library.
type Backend struct {
	arguments *arguments.Arguments

	config *config.Config

	events chan interface{}

	devices         map[string]device.Interface
	keystores       *keystore.Keystores
	onAccountInit   func(accounts.Interface)
	onAccountUninit func(accounts.Interface)
	onDeviceInit    func(device.Interface)
	onDeviceUninit  func(string)

	coins     map[string]coin.Coin
	coinsLock locker.Locker

	accounts     []accounts.Interface
	accountsLock locker.Locker

	log *logrus.Entry
}

// NewBackend creates a new backend with the given arguments.
func NewBackend(arguments *arguments.Arguments) *Backend {
	log := logging.Get().WithGroup("backend")
	backend := &Backend{
		arguments: arguments,
		config:    config.NewConfig(arguments.AppConfigFilename(), arguments.AccountsConfigFilename()),
		events:    make(chan interface{}, 1000),

		devices:   map[string]device.Interface{},
		keystores: keystore.NewKeystores(),
		coins:     map[string]coin.Coin{},
		accounts:  []accounts.Interface{},
		log:       log,
	}
	GetRatesUpdaterInstance().Observe(func(event observable.Event) { backend.events <- event })
	return backend
}

// addAccount adds the given account to the backend.
func (backend *Backend) addAccount(account accounts.Interface) {
	defer backend.accountsLock.Lock()()
	backend.accounts = append(backend.accounts, account)
	backend.onAccountInit(account)
	backend.events <- backendEvent{Type: "backend", Data: "accountsStatusChanged"}
}

// CreateAndAddAccount creates an account with the given parameters and adds it to the backend.
func (backend *Backend) CreateAndAddAccount(
	coin coin.Coin,
	code string,
	name string,
	scriptType signing.ScriptType,
	getSigningConfiguration func() (*signing.Configuration, error),
) {
	switch specificCoin := coin.(type) {
	case *btc.Coin:
		onEvent := func(code string) func(accounts.Event) {
			return func(event accounts.Event) {
				backend.events <- AccountEvent{Type: "account", Code: code, Data: string(event)}
			}
		}
		account := btc.NewAccount(specificCoin, backend.arguments.CacheDirectoryPath(), code, name,
			getSigningConfiguration, backend.keystores, onEvent(code), backend.log)
		backend.addAccount(account)
	case *eth.Coin:
		onEvent := func(event eth.Event) {
			backend.events <- AccountEvent{Type: "account", Code: code, Data: string(event)}
		}
		account := eth.NewAccount(specificCoin, backend.arguments.CacheDirectoryPath(), code, name,
			getSigningConfiguration, backend.keystores, onEvent, backend.log)
		backend.addAccount(account)
	default:
		panic("unknown coin type")
	}
}

func (backend *Backend) createAndAddAccount(
	coin coin.Coin,
	code string,
	name string,
	keypath string,
	scriptType signing.ScriptType,
) {
	if !backend.arguments.Multisig() && !backend.config.AppConfig().Backend.AccountActive(code) {
		backend.log.WithField("code", code).WithField("name", name).Info("skipping inactive account")
		return
	}
	backend.log.WithField("code", code).WithField("name", name).Info("init account")
	absoluteKeypath, err := signing.NewAbsoluteKeypath(keypath)
	if err != nil {
		panic(err)
	}
	getSigningConfiguration := func() (*signing.Configuration, error) {
		return backend.keystores.Configuration(scriptType, absoluteKeypath, backend.keystores.Count())
	}
	if backend.arguments.Multisig() {
		name += " Multisig"
	}
	backend.CreateAndAddAccount(coin, code, name, scriptType, getSigningConfiguration)
}

// Config returns the app config.
func (backend *Backend) Config() *config.Config {
	return backend.config
}

// DefaultAppConfig returns the default app config.y
func (backend *Backend) DefaultAppConfig() config.AppConfig {
	return config.NewDefaultAppConfig()
}

func (backend *Backend) defaultProdServers(code string) []*rpc.ServerInfo {
	switch code {
	case coinBTC:
		return backend.config.AppConfig().Backend.BTC.ElectrumServers
	case coinTBTC:
		return backend.config.AppConfig().Backend.TBTC.ElectrumServers
	case coinLTC:
		return backend.config.AppConfig().Backend.LTC.ElectrumServers
	case coinTLTC:
		return backend.config.AppConfig().Backend.TLTC.ElectrumServers
	default:
		panic(errp.Newf("The given code %s is unknown.", code))
	}
}

func defaultDevServers(code string) []*rpc.ServerInfo {
	const devShiftCA = `-----BEGIN CERTIFICATE-----
MIIGGjCCBAKgAwIBAgIJAO1AEqR+xvjRMA0GCSqGSIb3DQEBDQUAMIGZMQswCQYD
VQQGEwJDSDEPMA0GA1UECAwGWnVyaWNoMR0wGwYDVQQKDBRTaGlmdCBDcnlwdG9z
ZWN1cml0eTEzMDEGA1UECwwqU2hpZnQgQ3J5cHRvc2VjdXJpdHkgQ2VydGlmaWNh
dGUgQXV0aG9yaXR5MSUwIwYDVQQDDBxTaGlmdCBDcnlwdG9zZWN1cml0eSBSb290
IENBMB4XDTE4MDMwNzE3MzUxMloXDTM4MDMwMjE3MzUxMlowgZkxCzAJBgNVBAYT
AkNIMQ8wDQYDVQQIDAZadXJpY2gxHTAbBgNVBAoMFFNoaWZ0IENyeXB0b3NlY3Vy
aXR5MTMwMQYDVQQLDCpTaGlmdCBDcnlwdG9zZWN1cml0eSBDZXJ0aWZpY2F0ZSBB
dXRob3JpdHkxJTAjBgNVBAMMHFNoaWZ0IENyeXB0b3NlY3VyaXR5IFJvb3QgQ0Ew
ggIiMA0GCSqGSIb3DQEBAQUAA4ICDwAwggIKAoICAQDlz32VZk/D3rfm7Qwx6WkE
Fp9cdQV2FNYTeTjWVErVeTev02ctHHXV1fR3Svk8iIJWaALSJy7phdEDwC/3gDIQ
Ylm15kpntCibOWiQPZZxGq7Udts20fooccdZqtG/PKFRCPWZ2MOgHAOWDKGk6Kb+
siqkr55hkxwtiHuwkCcTh/Q2orEIuteSRbbYwgURZwd6dDIQq4ty7reC3j32xphh
edbnVBoDE6DSdebSS5SJL/gb6LxUdio98XdJPwkaD8292uEODxx0DKw/Ou2e1f5Q
Iv1WBl+LBaSrZ3sJSFUqoSvCQwBQmMAPoPJ1O13jCnFz1xoNygxUfz2eiKRL5E2l
VTmTh7zIez4oniOh5MOmDnKMVgTUGP1II2UU5r6PAq2tDpw4lVwyezhyLaBegwMc
pg/LinbABxUJrP8c8G2tve0yuTAhsir7r+Koo+nAE7FwcuIkD0UTyQcoag2IMS8O
dKZdYMGXjfUPJRBWg60LfXJeqMyU1oHpDrsRoa5iaYPt7ZApxc41kyynqfuuuIRD
du8327gd1nJ6ExMxGHY7dYelE4GNkOg3R0+5czykm/RxnGyDuDcO/RcYBJTChN1L
HYq+dTt0dYPAzBtiXnfuvjDyOsDK5f65pbrDgoOr6AQ4lvDJabcXFsWPrulM9Dyu
p0Y4+fuwXOCd8cr1Zm34MQIDAQABo2MwYTAdBgNVHQ4EFgQU486X86LMbNNSDw7J
NcT2U30NrikwHwYDVR0jBBgwFoAU486X86LMbNNSDw7JNcT2U30NrikwDwYDVR0T
AQH/BAUwAwEB/zAOBgNVHQ8BAf8EBAMCAYYwDQYJKoZIhvcNAQENBQADggIBAN0N
IPVBv8aaKDHDK9Nsu5fwiGp8GgkAN0B1+D34CbxTuzCDurToVMHCPEdo9tk/AzE4
Aa1p/kMW9X3XP8IyCFFj+BpEVkBRr9fXTVuh3XRHbyN6tXFbkKWQ/6QeUcnefq2k
DCpqEGjJQWsujZ4tJKkJl2HLIBZL6FAa/kaDLFHd3LeV1immC66CiN3ieHejCJL1
zZXiWi8pNxvEanTLPBaBjCw/AAl/owg/ySu2hGZzL0wsFboPrUbo4J+KvL1pvwql
PCT8AylJKCu+cn/N9zZDtUsgZJQBIq7btoakC3mCSnfVTlcbxfHVef0DbfohFqoV
ZpdmIuy0/njw7o+2uL/ArPJscPOhNl60ocDbdFIyYvc85oxyts8yMvKDdWV9Bm//
kl7lv4QUAvjqjb7ZgUhYibVk3Eu6n1MGZOP40l1/mm922/Wcd2n/HZVk/LsJs4tt
B6DLMDpf5nzeI1Yz/QtDGvNyb4aiJoRV5tQb9KkFfIeSzBS/ORZto4tVHKS37lxV
d1r8kFyCgpL9KASdahfyLBWCC7awlcOQP1QJA5QoO9u5Feq3lU0VnJF0YCZh8GOy
py3n1TR6S59eT495BiKDjWnhdVchEa8zMGIW/wFW7EX/LyW2zX3hQsdfnmMWUPVr
O3nOxjgSfRAfKWQ2Ny1APKcn6I83P5PFLhtO5I12
-----END CERTIFICATE-----`

	switch code {
	case coinBTC:
		return []*rpc.ServerInfo{{Server: "dev.shiftcrypto.ch:50002", TLS: true, PEMCert: devShiftCA}}
	case coinTBTC:
		return []*rpc.ServerInfo{
			{Server: "s1.dev.shiftcrypto.ch:51003", TLS: true, PEMCert: devShiftCA},
			{Server: "s2.dev.shiftcrypto.ch:51003", TLS: true, PEMCert: devShiftCA},
		}
	case coinLTC:
		return []*rpc.ServerInfo{{Server: "dev.shiftcrypto.ch:50004", TLS: true, PEMCert: devShiftCA}}
	case coinTLTC:
		return []*rpc.ServerInfo{{Server: "dev.shiftcrypto.ch:51004", TLS: true, PEMCert: devShiftCA}}
	default:
		panic(errp.Newf("The given code %s is unknown.", code))
	}
}

func (backend *Backend) defaultElectrumXServers(code string) []*rpc.ServerInfo {
	if backend.arguments.DevMode() {
		return defaultDevServers(code)
	}

	return backend.defaultProdServers(code)
}

// Coin returns the coin with the given code or an error if no such coin exists.
func (backend *Backend) Coin(code string) (coin.Coin, error) {
	defer backend.coinsLock.Lock()()
	coin, ok := backend.coins[code]
	if ok {
		return coin, nil
	}
	dbFolder := backend.arguments.CacheDirectoryPath()
	switch code {
	case "rbtc":
		servers := []*rpc.ServerInfo{{Server: "127.0.0.1:52001", TLS: false, PEMCert: ""}}
		coin = btc.NewCoin("rbtc", "RBTC", &chaincfg.RegressionNetParams, dbFolder, servers, "")
	case coinTBTC:
		servers := backend.defaultElectrumXServers(code)
		coin = btc.NewCoin(coinTBTC, "TBTC", &chaincfg.TestNet3Params, dbFolder, servers,
			"https://blockstream.info/testnet/tx/")
	case coinBTC:
		servers := backend.defaultElectrumXServers(code)
		coin = btc.NewCoin(coinBTC, "BTC", &chaincfg.MainNetParams, dbFolder, servers,
			"https://blockstream.info/tx/")
	case coinTLTC:
		servers := backend.defaultElectrumXServers(code)
		coin = btc.NewCoin(coinTLTC, "TLTC", &ltc.TestNet4Params, dbFolder, servers,
			"http://explorer.litecointools.com/tx/")
	case coinLTC:
		servers := backend.defaultElectrumXServers(code)
		coin = btc.NewCoin(coinLTC, "LTC", &ltc.MainNetParams, dbFolder, servers,
			"https://insight.litecore.io/tx/")
	case coinETH:
		coin = eth.NewCoin(code, params.MainnetChainConfig,
			"https://etherscan.io/tx/", backend.config.AppConfig().Backend.ETH.NodeURL)
	case coinRETH:
		coin = eth.NewCoin(code, params.RinkebyChainConfig,
			"https://rinkeby.etherscan.io/tx/", backend.config.AppConfig().Backend.RETH.NodeURL)
	case coinTETH:
		coin = eth.NewCoin(code, params.TestnetChainConfig,
			"https://ropsten.etherscan.io/tx/", backend.config.AppConfig().Backend.TETH.NodeURL)
	default:
		return nil, errp.Newf("unknown coin code %s", code)
	}
	backend.coins[code] = coin
	coin.Observe(func(event observable.Event) { backend.events <- event })
	return coin, nil
}

func (backend *Backend) initAccounts() {
	// Since initAccounts replaces all previous accounts, we need to properly close them first.
	backend.uninitAccounts()

	if backend.arguments.Testing() {
		switch {
		case backend.arguments.Multisig():
			TBTC, _ := backend.Coin(coinTBTC)
			backend.createAndAddAccount(TBTC, "tbtc-multisig", "Bitcoin Testnet", "m/48'/1'/0'",
				signing.ScriptTypeP2PKH)
			TLTC, _ := backend.Coin(coinTLTC)
			backend.createAndAddAccount(TLTC, "tltc-multisig", "Litecoin Testnet", "m/48'/1'/0'",
				signing.ScriptTypeP2PKH)
		case backend.arguments.Regtest():
			RBTC, _ := backend.Coin("rbtc")
			backend.createAndAddAccount(RBTC, "rbtc-p2pkh", "Bitcoin Regtest Legacy", "m/44'/1'/0'",
				signing.ScriptTypeP2PKH)
			backend.createAndAddAccount(RBTC, "rbtc-p2wpkh-p2sh", "Bitcoin Regtest Segwit", "m/49'/1'/0'",
				signing.ScriptTypeP2WPKHP2SH)
		default:
			TBTC, _ := backend.Coin(coinTBTC)
			backend.createAndAddAccount(TBTC, "tbtc-p2wpkh-p2sh", "Bitcoin Testnet", "m/49'/1'/0'",
				signing.ScriptTypeP2WPKHP2SH)
			backend.createAndAddAccount(TBTC, "tbtc-p2wpkh", "Bitcoin Testnet: bech32", "m/84'/1'/0'",
				signing.ScriptTypeP2WPKH)
			backend.createAndAddAccount(TBTC, "tbtc-p2pkh", "Bitcoin Testnet Legacy", "m/44'/1'/0'",
				signing.ScriptTypeP2PKH)

			TLTC, _ := backend.Coin(coinTLTC)
			backend.createAndAddAccount(TLTC, "tltc-p2wpkh-p2sh", "Litecoin Testnet", "m/49'/1'/0'",
				signing.ScriptTypeP2WPKHP2SH)
			backend.createAndAddAccount(TLTC, "tltc-p2wpkh", "Litecoin Testnet: bech32", "m/84'/1'/0'",
				signing.ScriptTypeP2WPKH)

			if backend.arguments.DevMode() {
				TETH, _ := backend.Coin(coinTETH)
				backend.createAndAddAccount(TETH, "teth", "Ethereum Ropsten", "m/44'/1'/0'/0/0", signing.ScriptTypeP2WPKH)
				RETH, _ := backend.Coin(coinRETH)
				backend.createAndAddAccount(RETH, "reth", "Ethereum Rinkeby", "m/44'/1'/0'/0/0", signing.ScriptTypeP2WPKH)
			}
		}
	} else {
		if backend.arguments.Multisig() {
			BTC, _ := backend.Coin(coinBTC)
			backend.createAndAddAccount(BTC, "btc-multisig", "Bitcoin", "m/48'/0'/0'",
				signing.ScriptTypeP2PKH)
			LTC, _ := backend.Coin(coinLTC)
			backend.createAndAddAccount(LTC, "ltc-multisig", "Litecoin", "m/48'/2'/0'",
				signing.ScriptTypeP2PKH)
		} else {
			BTC, _ := backend.Coin(coinBTC)
			backend.createAndAddAccount(BTC, "btc-p2wpkh-p2sh", "Bitcoin", "m/49'/0'/0'",
				signing.ScriptTypeP2WPKHP2SH)
			backend.createAndAddAccount(BTC, "btc-p2wpkh", "Bitcoin: bech32", "m/84'/0'/0'",
				signing.ScriptTypeP2WPKH)
			backend.createAndAddAccount(BTC, "btc-p2pkh", "Bitcoin Legacy", "m/44'/0'/0'",
				signing.ScriptTypeP2PKH)

			LTC, _ := backend.Coin(coinLTC)
			backend.createAndAddAccount(LTC, "ltc-p2wpkh-p2sh", "Litecoin", "m/49'/2'/0'",
				signing.ScriptTypeP2WPKHP2SH)
			backend.createAndAddAccount(LTC, "ltc-p2wpkh", "Litecoin: bech32", "m/84'/2'/0'",
				signing.ScriptTypeP2WPKH)

			if backend.arguments.DevMode() {
				ETH, _ := backend.Coin(coinETH)
				backend.createAndAddAccount(ETH, "eth", "Ethereum", "m/44'/60'/0'/0/0", signing.ScriptTypeP2WPKH)
			}
		}
	}
}

// AccountsStatus returns whether the accounts have been initialized.
func (backend *Backend) AccountsStatus() string {
	if len(backend.accounts) > 0 {
		return "initialized"
	}
	return "uninitialized"
}

// Testing returns whether this backend is for testing only.
func (backend *Backend) Testing() bool {
	return backend.arguments.Testing()
}

// Accounts returns the current accounts of the backend.
func (backend *Backend) Accounts() []accounts.Interface {
	return backend.accounts
}

// UserLanguage returns the language the UI should be presented in to the user.
func (backend *Backend) UserLanguage() language.Tag {
	userLocale, err := jibber_jabber.DetectIETF()
	if err != nil {
		return language.English
	}
	languages := []language.Tag{
		language.English,
		language.German,
	}
	tag, _, _ := language.NewMatcher(languages).Match(language.Make(userLocale))
	backend.log.WithField("user-language", tag).Debug("Detected user language")
	return tag
}

// OnAccountInit installs a callback to be called when an account is initialized.
func (backend *Backend) OnAccountInit(f func(accounts.Interface)) {
	backend.onAccountInit = f
}

// OnAccountUninit installs a callback to be called when an account is stopped.
func (backend *Backend) OnAccountUninit(f func(accounts.Interface)) {
	backend.onAccountUninit = f
}

// OnDeviceInit installs a callback to be called when a device is initialized.
func (backend *Backend) OnDeviceInit(f func(device.Interface)) {
	backend.onDeviceInit = f
}

// OnDeviceUninit installs a callback to be called when a device is uninitialized.
func (backend *Backend) OnDeviceUninit(f func(string)) {
	backend.onDeviceUninit = f
}

// Start starts the background services. It returns a channel of events to handle by the library
// client.
func (backend *Backend) Start() <-chan interface{} {
	usb.NewManager(backend.arguments.MainDirectoryPath(), backend.Register, backend.Deregister).Start()
	return backend.events
}

// Events returns the push notifications channel.
func (backend *Backend) Events() <-chan interface{} {
	return backend.events
}

// DevicesRegistered returns a map of device IDs to device of registered devices.
func (backend *Backend) DevicesRegistered() map[string]device.Interface {
	return backend.devices
}

func (backend *Backend) uninitAccounts() {
	defer backend.accountsLock.Lock()()
	for _, account := range backend.accounts {
		account := account
		backend.onAccountUninit(account)
		account.Close()
	}
	backend.accounts = []accounts.Interface{}
	backend.events <- backendEvent{Type: "backend", Data: "accountsStatusChanged"}
}

// Keystores returns the keystores registered at this backend.
func (backend *Backend) Keystores() *keystore.Keystores {
	return backend.keystores
}

// RegisterKeystore registers the given keystore at this backend.
func (backend *Backend) RegisterKeystore(keystore keystore.Keystore) {
	backend.log.Info("registering keystore")
	if err := backend.keystores.Add(keystore); err != nil {
		backend.log.Panic("Failed to add a keystore.", err)
	}
	if backend.arguments.Multisig() && backend.keystores.Count() != 2 {
		return
	}
	backend.initAccounts()
}

// DeregisterKeystore removes the registered keystore.
func (backend *Backend) DeregisterKeystore() {
	backend.log.Info("deregistering keystore")
	backend.keystores = keystore.NewKeystores()
	backend.uninitAccounts()
}

// Register registers the given device at this backend.
func (backend *Backend) Register(theDevice device.Interface) error {
	backend.devices[theDevice.Identifier()] = theDevice
	backend.onDeviceInit(theDevice)
	theDevice.Init(backend.Testing())

	mainKeystore := len(backend.devices) == 1
	theDevice.SetOnEvent(func(event device.Event, data interface{}) {
		switch event {
		case device.EventKeystoreGone:
			backend.DeregisterKeystore()
		case device.EventKeystoreAvailable:
			// absoluteKeypath := signing.NewEmptyAbsoluteKeypath().Child(44, signing.Hardened)
			// extendedPublicKey, err := backend.device.ExtendedPublicKey(absoluteKeypath)
			// if err != nil {
			// 	panic(err)
			// }
			// configuration := signing.NewConfiguration(absoluteKeypath,
			// 	[]*hdkeychain.ExtendedKey{extendedPublicKey}, 1)
			if backend.arguments.Multisig() {
				backend.RegisterKeystore(
					theDevice.KeystoreForConfiguration(nil, backend.keystores.Count()))
			} else if mainKeystore {
				// HACK: for device based, only one is supported at the moment.
				backend.keystores = keystore.NewKeystores()

				backend.RegisterKeystore(
					theDevice.KeystoreForConfiguration(nil, backend.keystores.Count()))
			}
		}
		backend.events <- deviceEvent{
			DeviceID: theDevice.Identifier(),
			Type:     "device",
			Data:     string(event),
			Meta:     data,
		}
	})
	select {
	case backend.events <- backendEvent{
		Type: "devices",
		Data: "registeredChanged",
	}:
	default:
	}
	return nil
}

// Deregister deregisters the device with the given ID from this backend.
func (backend *Backend) Deregister(deviceID string) {
	if _, ok := backend.devices[deviceID]; ok {
		backend.onDeviceUninit(deviceID)
		delete(backend.devices, deviceID)
		backend.DeregisterKeystore()
		backend.events <- backendEvent{Type: "devices", Data: "registeredChanged"}
	}
}

// Rates return the latest rates.
func (backend *Backend) Rates() map[string]map[string]float64 {
	return GetRatesUpdaterInstance().Last()
}

// DownloadCert downloads the first element of the remote certificate chain.
func (backend *Backend) DownloadCert(server string) (string, error) {
	var pemCert []byte
	conn, err := tls.Dial("tcp", server, &tls.Config{
		VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
			if len(rawCerts) == 0 {
				return errp.New("no remote certs")
			}

			certificatePEM := &pem.Block{Type: "CERTIFICATE", Bytes: rawCerts[0]}
			certificatePEMBytes := &bytes.Buffer{}
			if err := pem.Encode(certificatePEMBytes, certificatePEM); err != nil {
				panic(err)
			}
			pemCert = certificatePEMBytes.Bytes()
			return nil
		},
		InsecureSkipVerify: true,
	})
	if err != nil {
		return "", err
	}
	_ = conn.Close()
	return string(pemCert), nil
}

// CheckElectrumServer checks if a tls connection can be established with the electrum server, and
// whether the server is an electrum server.
func (backend *Backend) CheckElectrumServer(server string, pemCert string) error {
	backends := []rpc.Backend{
		electrum.NewElectrum(backend.log, &rpc.ServerInfo{Server: server, TLS: true, PEMCert: pemCert}),
	}
	conn, err := backends[0].EstablishConnection()
	if err != nil {
		return err
	}
	_ = conn.Close()
	// Simple check if the server is an electrum server.
	jsonrpcClient := jsonrpc.NewRPCClient(backends, backend.log)
	electrumClient := client.NewElectrumClient(jsonrpcClient, backend.log)
	defer electrumClient.Close()
	_, err = electrumClient.ServerVersion()
	return err
}

// RegisterTestKeystore adds a keystore derived deterministically from a PIN, for convenience in
// devmode.
func (backend *Backend) RegisterTestKeystore(pin string) {
	softwareBasedKeystore := software.NewKeystoreFromPIN(
		backend.keystores.Count(), pin)
	backend.RegisterKeystore(softwareBasedKeystore)
}
