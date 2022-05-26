const express = require('express');
const { Wallets, Gateway } = require('fabric-network');
const enrollAdmin = require('./app/enrollAdmin');
const registerUser = require('./app/registerUser');
const { buildCCP, buildWallet } = require('./utils/AppUtil.js');
const path = require('path');
const walletPath = path.join(__dirname, './wallet');
const app = express();
app.use(express.json());
const cors = require('cors');
app.use(cors());

app.post('/api/enrolladmin/:org', async (req, res) => {
    let response;
    if (req.params.org == 'org1') {
        response = await enrollAdmin();
    }

    if (response && response.success) {

        console.log(`Enroll was Success: ${response.message}`);
        res.status(200).json(response);

    } else {
        console.log(`Enroll was Failure: ${response.message}`);
        res.status(401).json(response);
        return;
    }
});

app.post('/api/register/:org', async (req, res) => {
    let response;
    if (req.params.org == 'org1') {
        response = await registerUser(req.body.userId, req.body.userAffiliation);
    }

    if (response && response.success) {

        console.log(`Enroll was Success: ${response.message}`);
        res.status(200).json(response);

    } else {
        console.log(`Enroll was Failure: ${response.message}`);
        res.status(401).json(response);
        return;
    }
});

app.get('/api/queryAllTickets', async function (req, res) {
    try {

        const ccp = buildCCP();
        const wallet = await buildWallet(Wallets, walletPath);
        const identity = await wallet.get('user');
        if (!identity) {
            console.log('An identity for the user "user" does not exist in the wallet');
            console.log('Run the registerUser.js application before retrying');
            return;
        }
        const gateway = new Gateway();
        await gateway.connect(ccp, { wallet, identity: 'user', discovery: { enabled: true, asLocalhost: true } });
        const network = await gateway.getNetwork('ticketchannel');
        const contract = network.getContract('ticketing_chaincode');
        const result = await contract.evaluateTransaction('QueryAllTickets');
        //console.log("Result :",result.toString());
        res.status(200).json({ response: result.toString() });

    } catch (error) {
        console.error(`Failed to evaluate transaction: ${error}`);
        res.status(500).json({ error: error });

    }
});

app.post('/api/createTicket/', async (req, res) => {
    try {
        const ccp = buildCCP();
        const wallet = await buildWallet(Wallets, walletPath);
        //console.log(`Wallet path: ${walletPath}`);
        const identity = await wallet.get('user');
        if (!identity) {
            console.log('An identity for the user "user" does not exist in the wallet');
            console.log('Run the registerUser.js application before retrying');
            return;
        }
        const gateway = new Gateway();
        await gateway.connect(ccp, { wallet, identity: 'user', discovery: { enabled: true, asLocalhost: true } });
        const network = await gateway.getNetwork('ticketchannel');
        const contract = network.getContract('ticketing_chaincode');
        console.log("req.body :",req.body);
        
        for (var i = 1; i <= req.body.TotalTickets; i++) {
            await contract.submitTransaction('CreateTicket', req.body.TicketType, req.body.Price, req.body.TotalTickets, req.body.Owner, req.body.EventId);
        }
        res.send('Transaction has been submitted');
    
      
        await gateway.disconnect();

    } catch (error) {
        console.error(`******** FAILED to run the application: ${error}`);

    }

});

app.get('/api/readTicket/:ticketId', async function (req, res) {
    try {
        const ccp = buildCCP();
        const wallet = await buildWallet(Wallets, walletPath);
        const identity = await wallet.get('user');
        if (!identity) {
            console.log('An identity for the user "user" does not exist in the wallet');
            console.log('Run the registerUser.js application before retrying');
            return;
        }
        const gateway = new Gateway();
        await gateway.connect(ccp, { wallet, identity: 'user', discovery: { enabled: true, asLocalhost: true } });
        const network = await gateway.getNetwork('ticketchannel');
        const contract = network.getContract('ticketing_chaincode');
        const result = await contract.evaluateTransaction('ReadTicket', req.params.ticketId);
        console.log("Result:",result.toString())
        res.status(200).json({ response: result.toString() });

    } catch (error) {
        console.error(`Failed to evaluate transaction: ${error}`);
        res.status(500).json({ error: error });

    }
});

app.get('/api/showTicketByEventId/:EventId', async function (req, res) {
    try {
        const ccp = buildCCP();
        const wallet = await buildWallet(Wallets, walletPath);
        const identity = await wallet.get('user');
        if (!identity) {
            console.log('An identity for the user "user" does not exist in the wallet');
            console.log('Run the registerUser.js application before retrying');
            return;
        }
        const gateway = new Gateway();
        await gateway.connect(ccp, { wallet, identity: 'user', discovery: { enabled: true, asLocalhost: true } });
        const network = await gateway.getNetwork('ticketchannel');
        const contract = network.getContract('ticketing_chaincode');
        const result = await contract.evaluateTransaction('ShowTicketByEventId', req.params.EventId);
       // console.log("Result:",result.toString())
        res.status(200).json({ response: result.toString() });

       

    } catch (error) {
        console.error(`Failed to evaluate transaction: ${error}`);
        res.status(500).json({ error: error });

    }
});

app.delete('/api/deleteTicket/:ticketId', async function (req, res) {
    try {
        const ccp = buildCCP();
        const wallet = await buildWallet(Wallets, walletPath);
        const identity = await wallet.get('user');
        if (!identity) {
            console.log('An identity for the user "user" does not exist in the wallet');
            console.log('Run the registerUser.js application before retrying');
            return;
        }
        const gateway = new Gateway();
        await gateway.connect(ccp, { wallet, identity: 'user', discovery: { enabled: true, asLocalhost: true } });
        const network = await gateway.getNetwork('ticketchannel');
        const contract = network.getContract('ticketing_chaincode');
        await contract.submitTransaction('DeleteTicket', req.params.ticketId);
        console.log("Ticket Deleted Sucessfully");
        res.send('Ticket Deleted Sucessfully');

    } catch (error) {
        console.error(`Failed to evaluate transaction: ${error}`);
        res.status(500).json({ error: error });

    }

});

app.put('/api/transferTicket/:ticketId', async function (req, res) {
    try {
        const ccp = buildCCP();
        const wallet = await buildWallet(Wallets, walletPath);
        const identity = await wallet.get('user');
        if (!identity) {
            console.log('An identity for the user "user" does not exist in the wallet');
            console.log('Run the registerUser.js application before retrying');
            return;
        }
        const gateway = new Gateway();
        await gateway.connect(ccp, { wallet, identity: 'user', discovery: { enabled: true, asLocalhost: true } });
        const network = await gateway.getNetwork('ticketchannel');
        const contract = network.getContract('ticketing_chaincode');
        await contract.submitTransaction('TransferTicket', req.params.ticketId, req.body.owner);
        console.log("Ticket Owner Changed");
        res.send('Ticket Owner Changed');

    } catch (error) {
        console.error(`Failed to evaluate transaction: ${error}`);
        res.status(500).json({ error: error });
    }
});

app.patch('/api/buyTicket', async function (req, res) {
    try {
        const ccp = buildCCP();
        const wallet = await buildWallet(Wallets, walletPath);
        const identity = await wallet.get('user');
        if (!identity) {
            console.log('An identity for the user "user" does not exist in the wallet');
            console.log('Run the registerUser.js application before retrying');
            return;
        }
        const gateway = new Gateway();
        await gateway.connect(ccp, { wallet, identity: 'user', discovery: { enabled: true, asLocalhost: true } });
        const network = await gateway.getNetwork('ticketchannel');
        const contract = network.getContract('ticketing_chaincode');
        console.log("Result :",req.body);
        const result=await contract.submitTransaction('BuyTicket', req.body.EventId, req.body.Owner,req.body.TicketType,req.body.TotalTickets);
        console.log(result);
        console.log("Ticket successfully buy.",req.body.ticketId);
        res.send('Ticket successfully buy.');
    } catch (error) {
        console.error(`Failed to evaluate transaction: ${error}`);
        res.status(500).json({ error: error });
    }
});

app.put('/api/resaleTicket/:ticketId', async function (req, res) {
    try {
        const ccp = buildCCP();
        const wallet = await buildWallet(Wallets, walletPath);
        const identity = await wallet.get('user');
        if (!identity) {
            console.log('An identity for the user "user" does not exist in the wallet');
            console.log('Run the registerUser.js application before retrying');
            return;
        }
        const gateway = new Gateway();
        await gateway.connect(ccp, { wallet, identity: 'user', discovery: { enabled: true, asLocalhost: true } });
        const network = await gateway.getNetwork('ticketchannel');
        const contract = network.getContract('ticketing_chaincode');
        await contract.submitTransaction('ResaleTicket', req.params.ticketId, req.body.Price);
        console.log("Ticket successfully resale.");
        res.send('Ticket successfully resale.');

    } catch (error) {
        console.error(`Failed to evaluate transaction: ${error}`);
        res.status(500).json({ error: error });
    }
});

app.put('/api/redeemTicket/:ticketId', async function (req, res) {
    try {
        const ccp = buildCCP();
        const wallet = await buildWallet(Wallets, walletPath);
        const identity = await wallet.get('user');
        if (!identity) {
            console.log('An identity for the user "user" does not exist in the wallet');
            console.log('Run the registerUser.js application before retrying');
            return;
        }
        const gateway = new Gateway();
        await gateway.connect(ccp, { wallet, identity: 'user', discovery: { enabled: true, asLocalhost: true } });
        const network = await gateway.getNetwork('ticketchannel');
        const contract = network.getContract('ticketing_chaincode');
        await contract.submitTransaction('RedeemTicket', req.params.ticketId);
        console.log("Ticket Redeemed");
        res.send('Ticket Redeemed');

    } catch (error) {
        console.error(`Failed to evaluate transaction: ${error}`);
        res.status(500).json({ error: error });
    }
});
let port = process.env.PORT || 8080;
app.listen(port, () => console.log(`server listening on port ${port}....`));
