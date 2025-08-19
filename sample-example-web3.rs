use web3::{
    futures::Future,
    types::{Address, H256, U256, BlockNumber, TransactionReceipt},
    Web3, Transport, contract::{Contract, Options},
    ethabi::{Contract as ContractABI, Function, Token, ParamType},
};
use std::str::FromStr;
use tokio;

#[derive(Debug)]
pub struct YieldFarmingClient {
    web3: Web3<web3::transports::Http>,
    contract: Contract<web3::transports::Http>,
}

impl YieldFarmingClient {
    pub fn new(rpc_url: &str, contract_address: Address, contract_abi: &[u8]) -> Result<Self, Box<dyn std::error::Error>> {
        let transport = web3::transports::Http::new(rpc_url)?;
        let web3 = Web3::new(transport);
        
        // Parse ABI and create contract instance
        let abi = ContractABI::load(contract_abi)?;
        let contract = Contract::new(web3.eth(), contract_address, abi);
        
        Ok(Self { web3, contract })
    }

    /// Deposit tokens into the yield farming pool
    pub async fn deposit(&self, amount: U256, account: Address) -> Result<H256, Box<dyn std::error::Error>> {
        let options = Options::default();
        
        let result = self.contract
            .call("deposit", (amount,), options)
            .from(account)
            .await?;
            
        Ok(result)
    }

    /// Withdraw tokens from the yield farming pool
    pub async fn withdraw(&self, amount: U256, account: Address) -> Result<H256, Box<dyn std::error::Error>> {
        let options = Options::default();
        
        let result = self.contract
            .call("withdraw", (amount,), options)
            .from(account)
            .await?;
            
        Ok(result)
    }

    /// Claim rewards from the yield farming pool
    pub async fn claim_rewards(&self, account: Address) -> Result<H256, Box<dyn std::error::Error>> {
        let options = Options::default();
        
        let result = self.contract
            .call("claimRewards", (), options)
            .from(account)
            .await?;
            
        Ok(result)
    }

    /// Get user's staked balance
    pub async fn get_staked_balance(&self, account: Address) -> Result<U256, Box<dyn std::error::Error>> {
        let result: U256 = self.contract
            .query("balanceOf", (account,), None, Options::default(), None)
            .await?;
            
        Ok(result)
    }

    /// Get pending rewards for a user
    pub async fn get_pending_rewards(&self, account: Address) -> Result<U256, Box<dyn std::error::Error>> {
        let result: U256 = self.contract
            .query("pendingRewards", (account,), None, Options::default(), None)
            .await?;
            
        Ok(result)
    }

    /// Get total value locked in the pool
    pub async fn get_total_value_locked(&self) -> Result<U256, Box<dyn std::error::Error>> {
        let result: U256 = self.contract
            .query("totalValueLocked", (), None, Options::default(), None)
            .await?;
            
        Ok(result)
    }

    /// Get current APY (Annual Percentage Yield)
    pub async fn get_current_apy(&self) -> Result<U256, Box<dyn std::error::Error>> {
        let result: U256 = self.contract
            .query("getCurrentAPY", (), None, Options::default(), None)
            .await?;
            
        Ok(result)
    }

    /// Wait for transaction confirmation
    pub async fn wait_for_transaction(&self, tx_hash: H256) -> Result<TransactionReceipt, Box<dyn std::error::Error>> {
        let receipt = self.web3.eth()
            .transaction_receipt(tx_hash)
            .await?;
            
        match receipt {
            Some(receipt) => Ok(receipt),
            None => Err("Transaction receipt not found".into()),
        }
    }

    /// Get latest block number
    pub async fn get_latest_block(&self) -> Result<u64, Box<dyn std::error::Error>> {
        let block = self.web3.eth()
            .block(BlockNumber::Latest)
            .await?;
            
        match block {
            Some(block) => Ok(block.number.unwrap().as_u64()),
            None => Err("Latest block not found".into()),
        }
    }
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    // Configuration
    let rpc_url = "https://mainnet.infura.io/v3/YOUR_PROJECT_ID"; // Replace with your Infura project ID
    let contract_address = Address::from_str("0x1234567890123456789012345678901234567890")?; // Replace with actual contract address
    
    // Initialize yield farming client
    let client = YieldFarmingClient::new(rpc_url, contract_address, b"[]")?; // Replace with actual ABI
    
    println!("🚀 Yield Farming Client Initialized");
    println!("Connected to: {}", rpc_url);
    println!("Contract: {:?}", contract_address);
    
    // Example operations (commented out for safety)
    /*
    let account = Address::from_str("YOUR_ACCOUNT_ADDRESS")?;
    let amount = U256::from(1000000000000000000u64); // 1 ETH in wei
    
    // Get current pool information
    let tvl = client.get_total_value_locked().await?;
    let apy = client.get_current_apy().await?;
    let latest_block = client.get_latest_block().await?;
    
    println!("📊 Pool Statistics:");
    println!("Total Value Locked: {} wei", tvl);
    println!("Current APY: {}%", apy);
    println!("Latest Block: {}", latest_block);
    
    // Check user's current position
    let staked_balance = client.get_staked_balance(account).await?;
    let pending_rewards = client.get_pending_rewards(account).await?;
    
    println!("👤 User Position:");
    println!("Staked Balance: {} wei", staked_balance);
    println!("Pending Rewards: {} wei", pending_rewards);
    
    // Example deposit (uncomment to execute)
    // let tx_hash = client.deposit(amount, account).await?;
    // println!("Deposit transaction: {:?}", tx_hash);
    // 
    // let receipt = client.wait_for_transaction(tx_hash).await?;
    // if receipt.status.unwrap().as_u64() == 1 {
    //     println!("✅ Deposit successful!");
    // } else {
    //     println!("❌ Deposit failed!");
    // }
    */
    
    println!("✅ Yield farming client ready for operations!");
    Ok(())
}

#[cfg(test)]
mod tests {
    use super::*;
    use web3::types::Address;
    
    #[tokio::test]
    async fn test_client_initialization() {
        // This is a mock test - in real scenarios you'd use a testnet
        let rpc_url = "https://goerli.infura.io/v3/YOUR_PROJECT_ID";
        let contract_address = Address::from_str("0x1234567890123456789012345678901234567890").unwrap();
        
        let result = YieldFarmingClient::new(rpc_url, contract_address, b"[]");
        assert!(result.is_ok());
    }
}
