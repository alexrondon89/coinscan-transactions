package config

type Config struct {
	AppName           string
	Version           string
	Ethereum          blockChain
	BinanceSmartChain blockChain
	Bitcoin           blockChain
}

type blockChain struct {
	Node  node
	Cache cache
}

type node struct {
	Name string
	Url  string
}

type cache struct {
	TimeToUpdate        uint16
	NumberOfElements    uint16
	MaxNumberOfElements uint16
}
