1. Creating a solidity file for transfer of coins from one account to another
2. Only the Creater or minter can generate the coins and send to other accounts/user.
3. Function mint checks that wheather sender is minter or not and then the coins are transfer to the reciever account
4. Error keyword is used to check wheather the balance in the senders account is nsufficient or not to send the coins, 
   if balance is insufficient then the function is revert back to the error message
5. And according to that the coins are send and receive to the accounts.