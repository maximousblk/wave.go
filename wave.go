package main

import (
	// standard
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math/big"
	"regexp"

	// third party
	"github.com/alexflint/go-arg"
	"github.com/mendsley/gojwk"
)

// goreleaser variables
var (
	version = "dev"
	commit  = "untagged"
)

// commandline
type args struct {
	Workers int    `arg:"-w,--workers" default:"4" help:"Number of workers to spawn"`
	Count   int    `arg:"-n,--number" default:"1" help:"Number of wallets to generate"`
	Pattern string `arg:"positional,required" help:"Regex pattern to match the wallet address"`
}

// set version and commit
func (args) Version() string {
	return fmt.Sprintf("wave %v (%v)", version, commit[:8])
}

func main() {
	// parse commandline arguments
	var args args
	arg.MustParse(&args)

	numWorkers := args.Workers          // number of workers to spawn
	numWallets := args.Count            // number of wallets to generate
	vanityPattern := args.Pattern       // regex pattern to match the wallet address
	walletChan := make(chan *Wallet, 1) // channel to get wallets from workers

	fmt.Println("Workers spawned:", numWorkers)
	fmt.Println("Wallets to generate:", numWallets)
	fmt.Println("Wallet pattern:", "/"+vanityPattern+"/")

	for n := 1; n <= numWallets; n++ {
		for w := 1; w <= numWorkers; w++ {
			// spawn a worker
			go worker(w, vanityPattern, walletChan)
		}

		// get wallet from worker
		k := <-walletChan

		fmt.Println("found!:", k.address)
		keyfile, err := gojwk.Marshal(k.key) // get keyfile as json string

		errcheck(err)

		// TODO(@maximousblk): add logic to output to a file
		fmt.Println("keyfile:", string(keyfile))
	}
}

func worker(workerId int, pattern string, walletChan chan<- *Wallet) {
	for {
		// generate wallet
		wallet := GenerateWallet()
		walletAddress := wallet.address

		// check if wallet address matches the provided pattern
		match, err := regexp.MatchString(pattern, walletAddress)

		errcheck(err)

		fmt.Printf("[WORKER%v] address: %v | match: %v]\n", workerId, walletAddress, match)

		// send wallet to main if matched
		if match {
			walletChan <- wallet
			break
		}
	}
}

type Wallet struct {
	address string
	key     *gojwk.Key
}

func GenerateWallet() *Wallet {
	// generate an RSA key
	reader := rand.Reader
	rsaKey, err := rsa.GenerateKey(reader, 4096)

	errcheck(err)

	// create new wallet instance
	wallet := &Wallet{}

	// Generate keyfile from RSA key
	wallet.key = &gojwk.Key{
		Kty: "RSA",
		E:   base64.RawURLEncoding.EncodeToString(big.NewInt(int64(rsaKey.E)).Bytes()),
		N:   base64.RawURLEncoding.EncodeToString(rsaKey.N.Bytes()),
		D:   base64.RawURLEncoding.EncodeToString(rsaKey.D.Bytes()),
	}

	// Generate wallet address
	h := sha256.New()
	h.Write(rsaKey.N.Bytes())                                         // Take the "n", in bytes and hash it using SHA256
	wallet.address = base64.RawURLEncoding.EncodeToString(h.Sum(nil)) // Then base64url encode it to get the wallet address

	return wallet
}

func errcheck(e error) {
	if e != nil {
		panic(e)
	}
}
