package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/p2eengineering/kalp-sdk-public/kalpsdk"
)

// Key names for token properties
const (
	nameKey       = "name"
	symbolKey     = "symbol"
	decimalsKey   = "decimals"
	totalSupplyKey = "totalSupply"
)

// SmartContract manages the transfer and minting of tokens
type SmartContract struct {
	kalpsdk.Contract
	transactions []event // Slice to hold successful transactions
}

// event provides an organized struct for emitting events
type event struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Value int    `json:"value"`
}

// NewSmartContract initializes a new SmartContract with an empty transactions list
func NewSmartContract(contract kalpsdk.Contract) *SmartContract {
	return &SmartContract{
		Contract:     contract,
		transactions: make([]event, 0),
	}
}

// Claim mints new tokens to the specified address and triggers a Transfer event
func (s *SmartContract) Claim(sdk kalpsdk.TransactionContextInterface, amount int, address string) error {
	initialized, err := checkInitialized(sdk)
	if err != nil {
		return fmt.Errorf("failed to check if contract is initialized: %v", err)
	}
	if !initialized {
		return fmt.Errorf("contract options need to be set before calling any function, call Initialize() to initialize contract")
	}

	clientMSPID, err := sdk.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("failed to get MSPID: %v", err)
	}
	if clientMSPID != "mailabs" {
		return fmt.Errorf("client is not authorized to mint tokens")
	}

	if amount <= 0 {
		return fmt.Errorf("mint amount must be a positive integer")
	}

	currentBalance, err := getBalance(sdk, address)
	if err != nil {
		return fmt.Errorf("failed to get balance: %v", err)
	}

	updatedBalance, err := add(currentBalance, amount)
	if err != nil {
		return fmt.Errorf("failed to update balance: %v", err)
	}

	err = sdk.PutStateWithoutKYC(address, []byte(strconv.Itoa(updatedBalance)))
	if err != nil {
		return fmt.Errorf("failed to update state for address: %v", err)
	}

	totalSupply, err := getTotalSupply(sdk)
	if err != nil {
		return fmt.Errorf("failed to get total supply: %v", err)
	}

	totalSupply, err = add(totalSupply, amount)
	if err != nil {
		return fmt.Errorf("failed to update total supply: %v", err)
	}

	err = sdk.PutStateWithoutKYC(totalSupplyKey, []byte(strconv.Itoa(totalSupply)))
	if err != nil {
		return fmt.Errorf("failed to update total supply in state: %v", err)
	}

	transferEvent := event{"0x0", address, amount}
	transferEventJSON, err := json.Marshal(transferEvent)
	if err != nil {
		return fmt.Errorf("failed to serialize transfer event: %v", err)
	}
	err = sdk.SetEvent("Transfer", transferEventJSON)
	if err != nil {
		return fmt.Errorf("failed to set event: %v", err)
	}

	log.Printf("Minted %d tokens to address %s", amount, address)
	s.recordTransaction("0x0", address, amount)

	return nil
}

// BalanceOf returns the balance of a specific account
func (s *SmartContract) BalanceOf(sdk kalpsdk.TransactionContextInterface, account string) (int, error) {
	initialized, err := checkInitialized(sdk)
	if err != nil {
		return 0, fmt.Errorf("failed to check if contract is initialized: %v", err)
	}
	if !initialized {
		return 0, fmt.Errorf("contract options need to be set before calling any function, call Initialize() to initialize contract")
	}

	return getBalance(sdk, account)
}

// TotalSupply returns the total supply of tokens in the contract
func (s *SmartContract) TotalSupply(sdk kalpsdk.TransactionContextInterface) (int, error) {
	initialized, err := checkInitialized(sdk)
	if err != nil {
		return 0, fmt.Errorf("failed to check if contract is initialized: %v", err)
	}
	if !initialized {
		return 0, fmt.Errorf("contract options need to be set before calling any function, call Initialize() to initialize contract")
	}

	return getTotalSupply(sdk)
}

// TransferFrom transfers a specific amount of tokens from one account to another
func (s *SmartContract) TransferFrom(sdk kalpsdk.TransactionContextInterface, from string, to string, value int) error {
	initialized, err := checkInitialized(sdk)
	if err != nil {
		return fmt.Errorf("failed to check if contract is initialized: %v", err)
	}
	if !initialized {
		return fmt.Errorf("contract options need to be set before calling any function, call Initialize() to initialize contract")
	}

	if from == to {
		return fmt.Errorf("cannot transfer to the same account")
	}

	err = transferHelper(sdk, from, to, value)
	if err != nil {
		return fmt.Errorf("transfer failed: %v", err)
	}

	transferEvent := event{from, to, value}
	transferEventJSON, err := json.Marshal(transferEvent)
	if err != nil {
		return fmt.Errorf("failed to serialize transfer event: %v", err)
	}
	err = sdk.SetEvent("Transfer", transferEventJSON)
	if err != nil {
		return fmt.Errorf("failed to set event: %v", err)
	}

	log.Printf("Transferred %d tokens from %s to %s", value, from, to)
	s.recordTransaction(from, to, value)

	return nil
}

// GetTransactions returns all successful transactions stored in the smart contract
func (s *SmartContract) GetTransactions() []event {
    return s.transactions
}

// Initialize sets the name, symbol, and decimals for the contract
func (s *SmartContract) Initialize(sdk kalpsdk.TransactionContextInterface, name, symbol, decimals string) (bool, error) {
	clientMSPID, err := sdk.GetClientIdentity().GetMSPID()
	if err != nil {
		return false, fmt.Errorf("failed to get MSPID: %v", err)
	}
	if clientMSPID != "mailabs" {
		return false, fmt.Errorf("client is not authorized to initialize contract")
	}

	bytes, err := sdk.GetState(nameKey)
	if err != nil {
		return false, fmt.Errorf("failed to get token name: %v", err)
	}
	if bytes != nil {
		return false, fmt.Errorf("contract is already initialized")
	}

	err = sdk.PutStateWithoutKYC(nameKey, []byte(name))
	if err != nil {
		return false, fmt.Errorf("failed to set token name: %v", err)
	}

	err = sdk.PutStateWithoutKYC(symbolKey, []byte(symbol))
	if err != nil {
		return false, fmt.Errorf("failed to set token symbol: %v", err)
	}

	err = sdk.PutStateWithoutKYC(decimalsKey, []byte(decimals))
	if err != nil {
		return false, fmt.Errorf("failed to set token decimals: %v", err)
	}

	log.Println("Contract initialized")
	return true, nil
}

// Helper functions

// getBalance retrieves the balance of a given account from the ledger
func getBalance(sdk kalpsdk.TransactionContextInterface, account string) (int, error) {
	balanceBytes, err := sdk.GetState(account)
	if err != nil {
		return 0, fmt.Errorf("failed to get state for account %s: %v", account, err)
	}
	if balanceBytes == nil {
		return 0, nil
	}

	balance, _ := strconv.Atoi(string(balanceBytes))
	return balance, nil
}

// getTotalSupply retrieves the total supply from the ledger
func getTotalSupply(sdk kalpsdk.TransactionContextInterface) (int, error) {
	totalSupplyBytes, err := sdk.GetState(totalSupplyKey)
	if err != nil {
		return 0, fmt.Errorf("failed to get total supply: %v", err)
	}
	if totalSupplyBytes == nil {
		return 0, nil
	}

	totalSupply, _ := strconv.Atoi(string(totalSupplyBytes))
	return totalSupply, nil
}

// transferHelper handles the transfer of tokens between two accounts
func transferHelper(sdk kalpsdk.TransactionContextInterface, from, to string, value int) error {
	if value <= 0 {
		return fmt.Errorf("transfer amount must be positive")
	}

	fromBalance, err := getBalance(sdk, from)
	if err != nil {
		return err
	}
	if fromBalance < value {
		return fmt.Errorf("insufficient funds")
	}

	toBalance, err := getBalance(sdk, to)
	if err != nil {
		return err
	}

	fromUpdatedBalance, err := sub(fromBalance, value)
	if err != nil {
		return err
	}

	toUpdatedBalance, err := add(toBalance, value)
	if err != nil {
		return err
	}

	err = sdk.PutStateWithoutKYC(from, []byte(strconv.Itoa(fromUpdatedBalance)))
	if err != nil {
		return err
	}

	err = sdk.PutStateWithoutKYC(to, []byte(strconv.Itoa(toUpdatedBalance)))
	if err != nil {
		return err
	}

	return nil
}

// recordTransaction stores a successful transaction in the transactions log
func (s *SmartContract) recordTransaction(from string, to string, value int) {
	s.transactions = append(s.transactions, event{From: from, To: to, Value: value})
}

// checkInitialized verifies whether the contract is initialized
func checkInitialized(sdk kalpsdk.TransactionContextInterface) (bool, error) {
	tokenName, err := sdk.GetState(nameKey)
	if err != nil {
		return false, fmt.Errorf("failed to get token name: %v", err)
	}
	return tokenName != nil, nil
}

// add adds two numbers with overflow checking
func add(a, b int) (int, error) {
	sum := a + b
	if (sum < a || sum < b) == (a >= 0 && b >= 0) {
		return 0, fmt.Errorf("addition overflow")
	}
	return sum, nil
}

// sub subtracts two numbers with underflow checking
func sub(a, b int) (int, error) {
	if a < b {
		return 0, fmt.Errorf("underflow occurred: %d is less than %d", a, b)
	}
	return a - b, nil
}
