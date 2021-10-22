pragma solidity >= 0.7.0 < 0.9.0;

contract learnEtherUnits 
{
    
    function test() public 
    {   
    	assert(1000000000000000000 wei == 1 ether); // 10^18 wei = 1 eth
    	assert(1 wei == 1); // 1 wei = 1 
    
    	assert(1 ether == 1e18); // 1 ether == 1e18
    	assert(2 ether == 2000000000000000000 wei); 
    }
    
    
    function exercise() public 
    {
        
        assert(1 minutes == 60 seconds);
        assert(24 hours == 1440 minutes);
        assert(1 days == 24 hours);
        assert(1 weeks == 7 days);
        
        
    }  
}


