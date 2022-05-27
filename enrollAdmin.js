'use strict'

const FabricCAServices = require('fabric-ca-client');
const { wallets } = require('fabric-network');
const fs = require('fs');
const path = require('path');

async function enrollAdmin(){
    try{
        const ccpPath = path.resolve(__dirname,'..','..','test-network','organizations','peerOrganizations','org1.example.com','connection-org1.json');
        const ccp = JSON.parse(fs.readFileSync(ccpPath,'utf8'));

        const caInfo = ccp.certificateAuthorities['ca.org1.example.com'];
        const caTLSCACerts = caInfo.tlsCACerts.pem;
        const ca = new FabricCAServices(caInfo.url,{trustedRoots:caTLSCACerts,verify:false},caInfo.caName);

        const walletPath = path.join(process.cwd(),'wallet');
        const wallet = await wallets.newFileSystemWallet(walletPath);
        

        const identity = await wallet.get('admin');
        if(identity)
        {
            console.log('An identity for the admin already exists in the wallet');
            return;
        }

        const enrollment = await ca.enroll({enrollmentId: 'admin',enrollmentSecret:'adminpw'});
        const x509Identity = {
            credentials : 
                {
                    certificate:enrollment.certificate,
                    privateKey: enrollment.key.toBytes(),
                },
                mspId:'Org1MSP',
                type:'x.509',
        };
        await wallet.put('admin',x509Identity);

    }catch(error){
        console.error('Failed to enroll admin : ${error}');
        process.exit(1);
    }
}

enrollAdmin();