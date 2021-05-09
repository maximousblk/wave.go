package main

import (
	// Standard
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"regexp"

	// Third Party
	"github.com/alexflint/go-arg"
	"github.com/mendsley/gojwk"
)

// GoReleaser
var (
	version = "dev"
	commit  = "untagged"
)

// go-arg
type args struct {
	Workers int    `arg:"-w,--workers" default:"4" help:"Number of workers to spawn"`
	Count   int    `arg:"-n,--number" default:"1" help:"Number of wallets to generate"`
	Output  string `arg:"-o,--output" default:"./keyfiles" help:"Output directory"`
	Pattern string `arg:"positional,required" help:"Regex pattern to match the wallet address"`
}

// set go-arg version and commit from GoReleaser
func (args) Version() string {
	return fmt.Sprintf("wave %v (%v)", version, commit[:8])
}

func main() {
	// parse commandline arguments
	var args args
	arg.MustParse(&args)

	numWorkers := args.Workers           // Number of workers to spawn
	numWallets := args.Count             // Number of wallets to generate
	vanityPattern := args.Pattern        // Regex pattern to match the wallet address
	outDir := filepath.Join(args.Output) // Output directory
	walletChan := make(chan Wallet, 1)   // Channel to get wallets from workers

	fmt.Println("Pattern:", "/"+vanityPattern+"/")
	fmt.Println("Outputs:", outDir)
	fmt.Println("Workers:", numWorkers)
	fmt.Println("Wallets:", numWallets)

	for n := 1; n <= numWallets; n++ {
		// spawn workers
		for w := 1; w <= numWorkers; w++ {
			go worker(w, vanityPattern, walletChan)
		}

		// get wallet from worker
		k := <-walletChan

		fmt.Println("[MATCH] address:", k.address)

		// get keyfile as json byte slice
		keyfile, err := gojwk.Marshal(k.key)
		errcheck(err)

		keyfilePath := filepath.Join(outDir, "arweave-keyfile-"+k.address+".json")

		// Check if output directory exists
		if _, err := os.Stat(outDir); os.IsNotExist(err) {
			// If not, create it
			derr := os.Mkdir(outDir, 0777)
			errcheck(derr)
		}

		// Write keyfile to file
		ioutil.WriteFile(keyfilePath, keyfile, 0666)
		fmt.Println("[EMIT] keyfile:", keyfilePath)
	}
}

func worker(workerId int, pattern string, walletChan chan<- Wallet) {
	for {
		// Generate wallet
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

func GenerateWallet() Wallet {
	// generate an RSA key
	reader := rand.Reader
	rsaKey, err := rsa.GenerateKey(reader, 4096)
	errcheck(err)

	// create new wallet instance
	wallet := Wallet{}

	// Generate keyfile from RSA key
	wallet.key = &gojwk.Key{
		Kty: "RSA",                                                                     // Key type
		N:   base64.RawURLEncoding.EncodeToString(rsaKey.N.Bytes()),                    // Modulus for the public key and the private keys
		D:   base64.RawURLEncoding.EncodeToString(rsaKey.D.Bytes()),                    // Private key exponent
		E:   base64.RawURLEncoding.EncodeToString(big.NewInt(int64(rsaKey.E)).Bytes()), // public key exponent
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
