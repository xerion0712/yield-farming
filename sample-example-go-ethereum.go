package main

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/ecdsa"
)

// YieldFarmingClient represents a client for interacting with yield farming contracts
type YieldFarmingClient struct {
	client          *ethclient.Client
	contractAddress common.Address
	contractABI     abi.ABI
	privateKey      *ecdsa.PrivateKey
	auth            *bind.TransactOpts
}

// PoolInfo represents information about a yield farming pool
type PoolInfo struct {
	TotalValueLocked *big.Int
	CurrentAPY       *big.Int
	RewardRate       *big.Int
	LastUpdateTime   *big.Int
}

// UserPosition represents a user's position in the yield farming pool
type UserPosition struct {
	StakedBalance   *big.Int
	PendingRewards  *big.Int
	LastClaimTime   *big.Int
	RewardDebt      *big.Int
}

// NewYieldFarmingClient creates a new yield farming client
func NewYieldFarmingClient(rpcURL string, contractAddress common.Address, privateKeyHex string) (*YieldFarmingClient, error) {
	// Connect to Ethereum client
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum client: %w", err)
	}

	// Parse private key
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	// Create auth for transactions
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1)) // Mainnet
	if err != nil {
		return nil, fmt.Errorf("failed to create transactor: %w", err)
	}

	// Load contract ABI (you would typically load this from a file)
	contractABI, err := abi.JSON(strings.NewReader(`[]`)) // Replace with actual ABI
	if err != nil {
		return nil, fmt.Errorf("failed to load contract ABI: %w", err)
	}

	return &YieldFarmingClient{
		client:          client,
		contractAddress: contractAddress,
		contractABI:     contractABI,
		privateKey:      privateKey,
		auth:            auth,
	}, nil
}

// Deposit tokens into the yield farming pool
func (c *YieldFarmingClient) Deposit(ctx context.Context, amount *big.Int) (*types.Transaction, error) {
	// Prepare transaction data
	data, err := c.contractABI.Pack("deposit", amount)
	if err != nil {
		return nil, fmt.Errorf("failed to pack deposit data: %w", err)
	}

	// Get gas price
	gasPrice, err := c.client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get gas price: %w", err)
	}

	// Get nonce
	nonce, err := c.client.PendingNonceAt(ctx, c.auth.From)
	if err != nil {
		return nil, fmt.Errorf("failed to get nonce: %w", err)
	}

	// Estimate gas
	msg := ethereum.CallMsg{
		From:  c.auth.From,
		To:    &c.contractAddress,
		Value: big.NewInt(0),
		Data:  data,
	}
	gasLimit, err := c.client.EstimateGas(ctx, msg)
	if err != nil {
		return nil, fmt.Errorf("failed to estimate gas: %w", err)
	}

	// Create transaction
	tx := types.NewTransaction(nonce, c.contractAddress, big.NewInt(0), gasLimit, gasPrice, data)
	
	// Sign transaction
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(1)), c.privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %w", err)
	}

	// Send transaction
	err = c.client.SendTransaction(ctx, signedTx)
	if err != nil {
		return nil, fmt.Errorf("failed to send transaction: %w", err)
	}

	return signedTx, nil
}

// Withdraw tokens from the yield farming pool
func (c *YieldFarmingClient) Withdraw(ctx context.Context, amount *big.Int) (*types.Transaction, error) {
	data, err := c.contractABI.Pack("withdraw", amount)
	if err != nil {
		return nil, fmt.Errorf("failed to pack withdraw data: %w", err)
	}

	gasPrice, err := c.client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get gas price: %w", err)
	}

	nonce, err := c.client.PendingNonceAt(ctx, c.auth.From)
	if err != nil {
		return nil, fmt.Errorf("failed to get nonce: %w", err)
	}

	msg := ethereum.CallMsg{
		From:  c.auth.From,
		To:    &c.contractAddress,
		Value: big.NewInt(0),
		Data:  data,
	}
	gasLimit, err := c.client.EstimateGas(ctx, msg)
	if err != nil {
		return nil, fmt.Errorf("failed to estimate gas: %w", err)
	}

	tx := types.NewTransaction(nonce, c.contractAddress, big.NewInt(0), gasLimit, gasPrice, data)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(1)), c.privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %w", err)
	}

	err = c.client.SendTransaction(ctx, signedTx)
	if err != nil {
		return nil, fmt.Errorf("failed to send transaction: %w", err)
	}

	return signedTx, nil
}

// Claim rewards from the yield farming pool
func (c *YieldFarmingClient) ClaimRewards(ctx context.Context) (*types.Transaction, error) {
	data, err := c.contractABI.Pack("claimRewards")
	if err != nil {
		return nil, fmt.Errorf("failed to pack claim rewards data: %w", err)
	}

	gasPrice, err := c.client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get gas price: %w", err)
	}

	nonce, err := c.client.PendingNonceAt(ctx, c.auth.From)
	if err != nil {
		return nil, fmt.Errorf("failed to get nonce: %w", err)
	}

	msg := ethereum.CallMsg{
		From:  c.auth.From,
		To:    &c.contractAddress,
		Value: big.NewInt(0),
		Data:  data,
	}
	gasLimit, err := c.client.EstimateGas(ctx, msg)
	if err != nil {
		return nil, fmt.Errorf("failed to estimate gas: %w", err)
	}

	tx := types.NewTransaction(nonce, c.contractAddress, big.NewInt(0), gasLimit, gasPrice, data)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(1)), c.privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %w", err)
	}

	err = c.client.SendTransaction(ctx, signedTx)
	if err != nil {
		return nil, fmt.Errorf("failed to send transaction: %w", err)
	}

	return signedTx, nil
}

// GetPoolInfo retrieves information about the yield farming pool
func (c *YieldFarmingClient) GetPoolInfo(ctx context.Context) (*PoolInfo, error) {
	// This would typically call contract view functions
	// For now, returning mock data
	return &PoolInfo{
		TotalValueLocked: big.NewInt(1000000000000000000000), // 1000 ETH
		CurrentAPY:       big.NewInt(1500),                   // 15%
		RewardRate:       big.NewInt(1000000000000000000),    // 1 token per second
		LastUpdateTime:   big.NewInt(time.Now().Unix()),
	}, nil
}

// GetUserPosition retrieves the user's position in the yield farming pool
func (c *YieldFarmingClient) GetUserPosition(ctx context.Context, userAddress common.Address) (*UserPosition, error) {
	// This would typically call contract view functions
	// For now, returning mock data
	return &UserPosition{
		StakedBalance:  big.NewInt(10000000000000000000), // 10 ETH
		PendingRewards: big.NewInt(500000000000000000),   // 0.5 tokens
		LastClaimTime:  big.NewInt(time.Now().Unix() - 3600),
		RewardDebt:     big.NewInt(0),
	}, nil
}

// WaitForTransaction waits for a transaction to be mined
func (c *YieldFarmingClient) WaitForTransaction(ctx context.Context, tx *types.Transaction) (*types.Receipt, error) {
	fmt.Printf("Waiting for transaction %s to be mined...\n", tx.Hash().Hex())
	
	receipt, err := bind.WaitMined(ctx, c.client, tx)
	if err != nil {
		return nil, fmt.Errorf("failed to wait for transaction: %w", err)
	}
	
	if receipt.Status == 0 {
		return nil, fmt.Errorf("transaction failed")
	}
	
	fmt.Printf("Transaction mined in block %d\n", receipt.BlockNumber)
	return receipt, nil
}

// GetLatestBlock retrieves the latest block number
func (c *YieldFarmingClient) GetLatestBlock(ctx context.Context) (uint64, error) {
	block, err := c.client.BlockByNumber(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to get latest block: %w", err)
	}
	return block.NumberU64(), nil
}

func main() {
	ctx := context.Background()
	
	// Configuration - replace with your actual values
	rpcURL := "https://mainnet.infura.io/v3/YOUR_PROJECT_ID"
	contractAddress := common.HexToAddress("0x1234567890123456789012345678901234567890")
	privateKeyHex := "YOUR_PRIVATE_KEY_HERE" // Be careful with private keys!
	
	// Initialize yield farming client
	client, err := NewYieldFarmingClient(rpcURL, contractAddress, privateKeyHex)
	if err != nil {
		fmt.Printf("Failed to create client: %v\n", err)
		return
	}
	
	fmt.Println("🚀 Yield Farming Client Initialized")
	fmt.Printf("Connected to: %s\n", rpcURL)
	fmt.Printf("Contract: %s\n", contractAddress.Hex())
	
	// Get pool information
	poolInfo, err := client.GetPoolInfo(ctx)
	if err != nil {
		fmt.Printf("Failed to get pool info: %v\n", err)
		return
	}
	
	fmt.Println("\n📊 Pool Statistics:")
	fmt.Printf("Total Value Locked: %s wei\n", poolInfo.TotalValueLocked.String())
	fmt.Printf("Current APY: %s%%\n", poolInfo.CurrentAPY.String())
	fmt.Printf("Reward Rate: %s wei/sec\n", poolInfo.RewardRate.String())
	
	// Get user position
	userAddress := client.auth.From
	userPosition, err := client.GetUserPosition(ctx, userAddress)
	if err != nil {
		fmt.Printf("Failed to get user position: %v\n", err)
		return
	}
	
	fmt.Println("\n👤 User Position:")
	fmt.Printf("Staked Balance: %s wei\n", userPosition.StakedBalance.String())
	fmt.Printf("Pending Rewards: %s wei\n", userPosition.PendingRewards.String())
	
	// Get latest block
	latestBlock, err := client.GetLatestBlock(ctx)
	if err != nil {
		fmt.Printf("Failed to get latest block: %v\n", err)
		return
	}
	fmt.Printf("Latest Block: %d\n", latestBlock)
	
	// Example operations (commented out for safety)
	/*
	amount := big.NewInt(1000000000000000000) // 1 ETH
	
	// Deposit example
	fmt.Println("\n💰 Depositing 1 ETH...")
	tx, err := client.Deposit(ctx, amount)
	if err != nil {
		fmt.Printf("Deposit failed: %v\n", err)
		return
	}
	
	receipt, err := client.WaitForTransaction(ctx, tx)
	if err != nil {
		fmt.Printf("Transaction failed: %v\n", err)
		return
	}
	
	fmt.Printf("✅ Deposit successful! Gas used: %d\n", receipt.GasUsed)
	*/
	
	fmt.Println("\n✅ Yield farming client ready for operations!")
}
