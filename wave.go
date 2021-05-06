package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math/big"
	"regexp"

	"github.com/alexflint/go-arg"
	"github.com/mendsley/gojwk"
)

type args struct {
	Workers int    `arg:"-w,--workers" default:"4" help:"Number of workers to spawn"`
	Count   int    `arg:"-n,--number" default:"1" help:"Number of wallets to generate"`
	Pattern string `arg:"positional,required" help:"Wallet address regex pattern"`
}

func (args) Version() string {
	return "wave 0.1.0"
}

func main() {
	var args args

	arg.MustParse(&args)

	numWorkers := args.Workers
	numWallets := args.Count
	vanityPattern := args.Pattern
	result := make(chan *Wallet, 1)

	fmt.Println("Workers spawned:", numWorkers)
	fmt.Println("Wallets to generate:", numWallets)
	fmt.Println("Wallet pattern:", "/"+vanityPattern+"/")

	for n := 1; n <= numWallets; n++ {
		for w := 1; w <= numWorkers; w++ {
			go worker(w, vanityPattern, result)
		}

		r := <-result

		fmt.Println("found!:", r.address)
		keyfile, _ := gojwk.Marshal(r.key)

		// TODO(@maximousblk): add logic to output to a file
		fmt.Println("keyfile", string(keyfile))
	}
}

type Wallet struct {
	address   string
	key       *gojwk.Key
	publicKey string
	pubKey    *rsa.PublicKey
}

func worker(workerId int, pattern string, result chan<- *Wallet) {
	for {
		wallet := GenerateWallet()

		walletAddress := wallet.address

		match, _ := regexp.MatchString(pattern, walletAddress)

		fmt.Printf("[W%v] address: %v\n", workerId, walletAddress)

		if match {
			result <- wallet
			break
		}
	}
}

func GenerateWallet() *Wallet {
	reader := rand.Reader
	rsaKey, _ := rsa.GenerateKey(reader, 4096)
	w := &Wallet{}

	w.key = &gojwk.Key{
		Kty: "RSA",
		E:   base64.RawURLEncoding.EncodeToString(big.NewInt(int64(rsaKey.E)).Bytes()),
		N:   base64.RawURLEncoding.EncodeToString(rsaKey.N.Bytes()),
		D:   base64.RawURLEncoding.EncodeToString(rsaKey.D.Bytes()),
	}

	w.pubKey = rsaKey.Public().(*rsa.PublicKey)
	// Take the "n", in bytes and hash it using SHA256
	h := sha256.New()
	h.Write(rsaKey.N.Bytes())

	// Finally base64url encode it to have the resulting address
	w.address = base64.RawURLEncoding.EncodeToString(h.Sum(nil))
	w.publicKey = base64.RawURLEncoding.EncodeToString(rsaKey.N.Bytes())
	return w
}
