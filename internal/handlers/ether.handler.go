package handlers

import (
	"context"
	"encoding/json"
	"ether/internal/config"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type IBlock interface {
	GetBlockHandler(w http.ResponseWriter, r *http.Request) (Block, error)
}

type Block struct {
	Number       uint64   `json:"number"`
	Transactions []string `json:"transactions"`
}

type Address struct {
	Address string `json:"address"`
	Balance string `json:"balance"`
}

func GetAndGenerateKey(w http.ResponseWriter, r *http.Request) {
	// Generar una nueva clave privada (clave ECDSA)
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal("Error generando la clave privada:", err)
	}

	privateKeyHex := crypto.PubkeyToAddress(privateKey.PublicKey)
	privateKeyBytes := crypto.FromECDSA(privateKey)

	resp := struct {
		Address    string `json:"address"`     // Dirección de Ethereum
		PrivateKey string `json:"private_key"` // Guardar la clave privada como string hexadecimal (¡esto debe ser almacenado de forma segura!)
	}{
		Address:    privateKeyHex.Hex(),
		PrivateKey: fmt.Sprintf("%x", privateKeyBytes),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func GetBlockHandler(w http.ResponseWriter, r *http.Request) {
	defer getClient().Close()

	block, err := getClient().BlockByNumber(context.Background(), nil)
	if err != nil {
		log.Fatalf("Error obteniendo el bloque", err)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(block.NumberU64())
}

func GetTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	defer getClient().Close()

	block, err := getClient().BlockByNumber(context.Background(), nil)
	if err != nil {
		log.Fatalf("Error obteniendo el bloque", err)
	}

	transactions := block.Transactions()
	transactionHashes := make([]string, len(transactions))
	for i, tx := range transactions {
		transactionHashes[i] = tx.Hash().String()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactionHashes)
}

func GetBalanceHandler(w http.ResponseWriter, r *http.Request) {
	defer getClient().Close()

	// vars := mux.Vars(r)
	// address := vars["address"]
	address := "0x00000000219ab540356cBB839Cbe05303d7705Fa" // ETH FUNDATION

	// Convertir la dirección en un tipo de datos de Ethereum
	accountAddress := common.HexToAddress(address)

	// Consultar el balance
	balance, err := getClient().BalanceAt(context.Background(), accountAddress, nil)
	if err != nil {
		log.Fatalf("Error obteniendo el balance: %v", err)
	}

	// Mostrar el balance en ETH
	// `balance` está en Wei, así que lo convertimos a ETH
	ethBalance := new(big.Float).SetInt(balance)
	ethBalance.Quo(ethBalance, big.NewFloat(1e18)) // Convertir Wei a ETH

	fmt.Printf("El balance de la cuenta %s es %s ETH\n", address, ethBalance.String())

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ethBalance)
}

func getClient() *ethclient.Client {
	infuraURL := config.GetConfig().BaseURL + os.Getenv("INFURA_API_KEY")
	client, err := ethclient.Dial(infuraURL)
	if err != nil {
		log.Fatalf("Error creando cliente: %v", err)
	}
	return client
}

func GetClient() *ethclient.Client {
	return getClient()
}
