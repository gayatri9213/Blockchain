cd basic-ticket-network/
chmod -R 0755 networkScript.sh
./networkScript.sh

====Export variables for CLI queries ==================

export PATH=${PWD}/../bin:$PATH
export FABRIC_CFG_PATH=$PWD/../config/
source ./scripts/setPeerConnectionParam.sh 1 
source ./scripts/setOrgPeerContext.sh 1


====CHAINCODE EXECUTION SAMPLES (CLI) ==================

peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C ticketchannel -n ticketing_chaincode $PEER_CONN_PARAMS -c '{"function":"CreateTicket","Args":["Gold","1500","1","Gayatri","3"]}'

peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile $ORDERER_CA -C ticketchannel -n ticketing_chaincode $PEER_CONN_PARAMS -c '{"function":"ReadTicket","Args":["ticket2"]}'

peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile $ORDERER_CA -C ticketchannel -n ticketing_chaincode $PEER_CONN_PARAMS -c '{"function":"DeleteTicket","Args":["ticket2"]}'

peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile $ORDERER_CA -C ticketchannel -n ticketing_chaincode $PEER_CONN_PARAMS -c '{"function":"GetTicketsByRange","Args":["ticket4","ticket7"]}'

peer chaincode query -C ticketchannel -n ticketing_chaincode -c '{"Args":["QueryAllTickets"]}'

peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C ticketchannel -n ticketing_chaincode $PEER_CONN_PARAMS -c '{"function":"BuyTicket","Args":["3","Max","Gold","1"]}'

peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile $ORDERER_CA -C ticketchannel -n ticketing_chaincode $PEER_CONN_PARAMS -c '{"function":"ResaleTicket","Args":["ticket4",0.25]}'

peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile $ORDERER_CA -C ticketchannel -n ticketing_chaincode $PEER_CONN_PARAMS -c '{"function":"RedeemTicket","Args":["ticket2"]}'

peer chaincode query -C ticketchannel -n ticketing_chaincode -c '{"Args":["QueryTicketsByOwner","30"]}'

peer chaincode query -C ticketchannel -n ticketing_chaincode -c '{"Args":["QueryTickets","{\"selector\":{\"owner\":\"Max\"}}"]}'

peer chaincode query -C ticketchannel -n ticketing_chaincode -c '{"Args":["GetTicketHistory","ticket1"]}'

peer chaincode query -C ticketchannel -n ticketing_chaincode -c '{"Args":["GetTicketsByRangeWithPagination","ticket6","ticket7","2","g1AAAAA-eJzLYWBgYMpgSmHgKy5JLCrJTq2MT8lPzkzJBYqzl2QmZ6eWmIOkOWDSyBJZAB9lEjk"]}'

peer chaincode query -C ticketchannel -n ticketing_chaincode -c '{"Args":["QueryTicketsWithPagination","{\"selector\":{\"owner\":\"Max\"}}","3",""]}'



queryString := fmt.Sprintf("{\"selector\":{\"Status\":\"NEW\",\"eventId\":\"%s\"}}", eventId)
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}

	defer resultsIterator.Close()

	var ticketArray []*Ticket
	
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var ticket Ticket
		err = json.Unmarshal(queryResult.Value, &ticket)
		if err != nil {
			return nil, err
		}
		ticketArray = append(ticketArray, &ticket)
		

	}

	return ticketArray, nil

      var count int
  m := map[string][]int{}
  for _, t := range m {
    count += len(t)
  }