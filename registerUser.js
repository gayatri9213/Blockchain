'use strict'

const FabricCAServices = require('fabric-ca-client');
const { wallets } = require('fabric-network');
const fs = require('fs');
const path = require('path');

async function registerUser(){
    try{
        const ccpPath = path.resolve(__dirname,'..','..','test-network','organizations','peerOrganizations','org1.example.com','connection-org1.json');
        const ccp = JSON.parse(fs.readFileSync(ccpPath,'utf8'));

        const caURL = ccp.certificateAuthorities['ca.org1.example.com'].url;
        const ca = new FabricCAServices(caURL);

        const walletPath = path.join(process.cwd(),'wallet');
        const wallet = await wallets.newFileSystemWallet(walletPath);
        
        const identity = await wallet.get('appUser');
        if(identity){
            console.log('An identity for the user "appUser" already exists in the wallet');
            return;
        }

        const adminIdentity = await wallet.get('admin');
        if(!adminIdentity){
            console.log('An identity for the admin does not exists in the wallet');
            console.log('Run the enrollAdmin.js before retrying');
            return;
        }

        const provider = wallet.getProviderRegistry().getProvider(adminIdentity.type);
        const adminUser = await provider.getUserContest(adminIdentity,'admin');

        const secret = await ca.register({
            affiliation : 'org1.department1',
            enrollmentId : 'appUser',
            role: 'client',
        },adminUser);

        const enrollment = await ca.enroll({
            enrollmentId : 'appUser',
            enrollmentSecret : secret
        })

        const x509Identity = {
            credentials :{
                certificate : enrollment.certificate,
                privateKey : enrollment.key.toBytes(),
            },
            mspId : 'Org1MSP',
            type : 'x.509'
        };

        await wallet.put('appUser',x509Identity);
        console.log('Successfully registered user "appUser" and imported it into the wallet');
    }catch(error){
        console.error('Failed to register user "appUser":${error}');
        process.exit(1);
    }
}

registerUser();