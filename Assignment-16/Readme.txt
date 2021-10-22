Exercise : 

1. Create a contract Oracle with an address datatype called admin and a public integer called rand. 
2. Create a constructor and set the admin to the current caller. 
3. Write a function which takes the input of an integer and sets the state variable rand to that integer. 
4. Require that the current caller must equal the admin.
5. Set the oracle contract to a new variable called oracle in the GenerateRandomNumber contract (hint calling contracts)
6. Write a constructor in the GenerateRandomNumber contract which sets the oracle to a deployment address of the Oracle 
7. Modify the hash return so that the miners greatly lesson control manipulation to the random generation. 

