pragma solidity >= 0.7.0 < 0.9.0;

contract LedgerBalance 
{
	mapping(address => uint) balance;

	function updatesBalance(uint newBalance) public 
	{
    		balance[msg.sender] = newBalance;
	}
}

contract Updated 
{
	function updatesBalance() public 
	{
        	LedgerBalance ledgerBalance = new LedgerBalance();
        	ledgerBalance.updatesBalance(30);
    	}
    
}

// Other Global Variable Examples 

contract SimpleStorage {
    
    uint storedData;
    
    function set(uint x) public {
        //storedData = x;
        // storedData = block.difficulty;
        // storedData = block.timestamp;
        storedData = block.number;
    }
    
    function get() public view returns (uint) {
        return storedData;
    }
    
}
