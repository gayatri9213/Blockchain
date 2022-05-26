/*
 SPDX-License-Identifier: Apache-2.0
*/
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

const index = "Event~Id"

var ticketCount int
var contract SimpleChaincode

type SimpleChaincode struct {
	contractapi.Contract
}

type Ticket struct {
	TicketId   string  `json:"ticketId"`
	TicketType string  `json:"ticketType"`
	Price      float32 `json:"price"`
	Seats      int     `json:"seats"`
	Owner      string  `json:"owner"`
	EventId    string  `json:"eventId"`
	Status     string  `json:"status"`
	Minter     string  `json:"minter"`
}

type HistoryQueryResult struct {
	Record    *Ticket   `json:"record"`
	TxId      string    `json:"txId"`
	Timestamp time.Time `json:"timestamp"`
	IsDelete  bool      `json:"isDelete"`
}

type PaginatedQueryResult struct {
	Records             []*Ticket `json:"records"`
	FetchedRecordsCount int32     `json:"fetchedRecordsCount"`
	Bookmark            string    `json:"bookmark"`
}

func (t *SimpleChaincode) IDGenerator(doctype string, count int) (string, error) {

	docSubstring := doctype[0:6]
	s := []string{docSubstring, strconv.Itoa(count)}
	return strings.Join(s, ""), nil
}

func (t *SimpleChaincode) CreateTicket(ctx contractapi.TransactionContextInterface, ticketType string, price float32, seats int, owner string, eventId string) error {
	id, _ := contract.IDGenerator("ticket", ticketCount)
	ticketBytes, err := ctx.GetStub().GetState(id)
	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}
	//check if ID already exists (return the state of the ID by checking the world state)
	if ticketBytes != nil {
		return fmt.Errorf("the ticket already exists for user")
	}

	ticket := &Ticket{
		TicketId:   id,
		TicketType: ticketType,
		Price:      price,
		Seats:      seats,
		Owner:      owner,
		EventId:    eventId,
		Status:     "NEW",
		Minter:     owner,
	}

	ticketBytes, err = json.Marshal(ticket)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(id, ticketBytes)
	if err != nil {
		return err
	}
	eventIdIndexKey, err := ctx.GetStub().CreateCompositeKey(index, []string{ticket.EventId, ticket.TicketId})
	if err != nil {
		return err
	}
	ticketCount += 1

	value := []byte{0x00}
	return ctx.GetStub().PutState(eventIdIndexKey, value)
}

func (t *SimpleChaincode) ReadTicket(ctx contractapi.TransactionContextInterface, ticketId string) (*Ticket, error) {
	ticketBytes, err := ctx.GetStub().GetState(ticketId)
	if err != nil {
		return nil, fmt.Errorf("failed to get asset %s: %v", ticketId, err)
	}
	if ticketBytes == nil {
		return nil, fmt.Errorf("asset %s does not exist", ticketId)
	}

	var ticket Ticket
	err = json.Unmarshal(ticketBytes, &ticket)
	if err != nil {
		return nil, err
	}

	return &ticket, nil
}

func (t *SimpleChaincode) DeleteTicket(ctx contractapi.TransactionContextInterface, ticketId string) error {
	ticket, err := t.ReadTicket(ctx, ticketId)
	if err != nil {
		return err
	}

	err = ctx.GetStub().DelState(ticketId)
	if err != nil {
		return fmt.Errorf("failed to delete ticket %s: %v", ticketId, err)
	}

	eventIdIndexKey, err := ctx.GetStub().CreateCompositeKey(index, []string{ticket.EventId, ticket.TicketId})
	if err != nil {
		return err
	}

	return ctx.GetStub().DelState(eventIdIndexKey)
}

func (t *SimpleChaincode) TransferTicket(ctx contractapi.TransactionContextInterface, ticketId, newOwner string) error {
	ticket, err := t.ReadTicket(ctx, ticketId)
	if err != nil {
		return err
	}

	ticket.Owner = newOwner

	ticketAsBytes, _ := json.Marshal(ticket)

	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(ticketId, ticketAsBytes)
}

func constructQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) ([]*Ticket, error) {
	var tickets []*Ticket
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
		tickets = append(tickets, &ticket)
	}

	return tickets, nil
}

func (t *SimpleChaincode) GetTicketsByRange(ctx contractapi.TransactionContextInterface, startKey, endKey string) ([]*Ticket, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	return constructQueryResponseFromIterator(resultsIterator)
}

func (t *SimpleChaincode) QueryAllTickets(ctx contractapi.TransactionContextInterface) ([]*Ticket, error) {
	startKey := ""
	endKey := ""

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()
	return constructQueryResponseFromIterator(resultsIterator)
}

func (t *SimpleChaincode) BuyTicket(ctx contractapi.TransactionContextInterface, eventId string, newOwner string, ticketType string, totalTickets int) error {

	ticketResultsIterator, err := ctx.GetStub().GetStateByPartialCompositeKey(index, []string{eventId})
	if err != nil {
		return err
	}
	defer ticketResultsIterator.Close()

	for ticketResultsIterator.HasNext() {
		responseRange, err := ticketResultsIterator.Next()
		if err != nil {
			return err
		}
		_, compositeKeyParts, err := ctx.GetStub().SplitCompositeKey(responseRange.Key)
		if err != nil {
			return err
		}

		if len(compositeKeyParts) > 1 {

			returnedTicketID := compositeKeyParts[1]

			ticket, err := t.ReadTicket(ctx, returnedTicketID)
			if err != nil {
				return fmt.Errorf("error in reading ticket %v", err)
			}

			if (ticket.Status == "NEW" || ticket.Status == "RESALE") && ticket.TicketType == ticketType {
				ticket.Owner = newOwner
				ticket.Status = "SOLD"
				fmt.Println("Ticket successfully sold.")
			} else {
				if err != nil {
					return fmt.Errorf("you cannot buy this ticket. %s: %v", returnedTicketID, err)
				}
			}

			ticketBytes, err := json.Marshal(ticket)
			if err != nil {
				return err
			}
			err = ctx.GetStub().PutState(returnedTicketID, ticketBytes)
			fmt.Println("returnedTicketID :", returnedTicketID)
			if err != nil {
				return fmt.Errorf("transfer failed for ticket %s: %v", returnedTicketID, err)
			}

		}

	}

	return nil
}

func (t *SimpleChaincode) ResaleTicket(ctx contractapi.TransactionContextInterface, ticketId string, newPrice float32) error {
	ticket, err := t.ReadTicket(ctx, ticketId)
	if err != nil {
		return err
	}
	if ticket.Status == "SOLD" {
		ticket.Price = newPrice
		ticket.Status = "RESALE"
		fmt.Println("Ticket successfully resale.")
	}
	if ticket.Status != "SOLD" {
		fmt.Println("You cannot resale this ticket.")
	}
	ticketAsBytes, _ := json.Marshal(ticket)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(ticketId, ticketAsBytes)
}

func (t *SimpleChaincode) RedeemTicket(ctx contractapi.TransactionContextInterface, ticketId string) error {
	ticket, err := t.ReadTicket(ctx, ticketId)
	if err != nil {
		return err
	}
	if ticket.Status == "SOLD" {

		ticket.Status = "REDEEM"
		fmt.Println("Ticket successfully redeem.")
	}
	if ticket.Status != "SOLD" {
		fmt.Println("You cannot redeem this ticket.")
	}
	ticketAsBytes, _ := json.Marshal(ticket)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(ticketId, ticketAsBytes)
}

func (t *SimpleChaincode) QueryTicketsByOwner(ctx contractapi.TransactionContextInterface, owner string) ([]*Ticket, error) {

	queryString := fmt.Sprintf(`{"selector":{"status":"SOLD","owner":"%s"}}`, owner)
	return getQueryResultForQueryString(ctx, queryString)
}

func (t *SimpleChaincode) QueryTickets(ctx contractapi.TransactionContextInterface, queryString string) ([]*Ticket, error) {
	return getQueryResultForQueryString(ctx, queryString)
}

func getQueryResultForQueryString(ctx contractapi.TransactionContextInterface, queryString string) ([]*Ticket, error) {
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	return constructQueryResponseFromIterator(resultsIterator)
}

func (t *SimpleChaincode) GetTicketsByRangeWithPagination(ctx contractapi.TransactionContextInterface, startKey string, endKey string, pageSize int, bookmark string) ([]*Ticket, error) {

	resultsIterator, _, err := ctx.GetStub().GetStateByRangeWithPagination(startKey, endKey, int32(pageSize), bookmark)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	return constructQueryResponseFromIterator(resultsIterator)
}

func (t *SimpleChaincode) QueryTicketsWithPagination(ctx contractapi.TransactionContextInterface, queryString string, pageSize int, bookmark string) (*PaginatedQueryResult, error) {

	return getQueryResultForQueryStringWithPagination(ctx, queryString, int32(pageSize), bookmark)
}

func getQueryResultForQueryStringWithPagination(ctx contractapi.TransactionContextInterface, queryString string, pageSize int32, bookmark string) (*PaginatedQueryResult, error) {

	resultsIterator, responseMetadata, err := ctx.GetStub().GetQueryResultWithPagination(queryString, pageSize, bookmark)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	tickets, err := constructQueryResponseFromIterator(resultsIterator)
	if err != nil {
		return nil, err
	}

	return &PaginatedQueryResult{
		Records:             tickets,
		FetchedRecordsCount: responseMetadata.FetchedRecordsCount,
		Bookmark:            responseMetadata.Bookmark,
	}, nil
}

func (t *SimpleChaincode) GetTicketHistory(ctx contractapi.TransactionContextInterface, ticketId string) ([]HistoryQueryResult, error) {
	log.Printf("GetTicketHistory: ID %v", ticketId)

	resultsIterator, err := ctx.GetStub().GetHistoryForKey(ticketId)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()
	var i int
	var records []HistoryQueryResult
	for i = 0; resultsIterator.HasNext(); i++ {
		response, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var ticket Ticket
		if len(response.Value) > 0 {
			err = json.Unmarshal(response.Value, &ticket)
			if err != nil {
				return nil, err
			}
		} else {
			ticket = Ticket{
				TicketId: ticketId,
			}
		}

		timestamp, err := ptypes.Timestamp(response.Timestamp)
		if err != nil {
			return nil, err
		}

		record := HistoryQueryResult{
			TxId:      response.TxId,
			Timestamp: timestamp,
			Record:    &ticket,
			IsDelete:  response.IsDelete,
		}
		records = append(records, record)
	}

	return records, nil
}

func (t *SimpleChaincode) Init(ctx contractapi.TransactionContextInterface) error {
	ticketCount = 1

	return nil
}

func (t *SimpleChaincode) ShowTicketByEventId(ctx contractapi.TransactionContextInterface, eventId string) ([]*Ticket, error) {

	queryString := fmt.Sprintf(`{"selector":{"eventId":"%s"}}`, eventId)
	return getQueryResultForQueryString(ctx, queryString)
}

func main() {
	// chaincode, err := contractapi.NewChaincode(&SimpleChaincode{})
	// if err != nil {
	// 	log.Panicf("Error creating asset chaincode: %v", err)
	// }

	// if err := chaincode.Start(); err != nil {
	// 	log.Panicf("Error starting asset chaincode: %v", err)
	// }

}

/*
====CHAINCODE EXECUTION SAMPLES (CLI) ==================


peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C ticketchannel -n ticketing_chaincode $PEER_CONN_PARAMS -c '{"function":"CreateTicket","Args":["Gold","1500","1","Gayatri","3"]}'

peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile $ORDERER_CA -C ticketchannel -n ticketing_chaincode $PEER_CONN_PARAMS -c '{"function":"ReadTicket","Args":["ticket2"]}'

peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile $ORDERER_CA -C ticketchannel -n ticketing_chaincode $PEER_CONN_PARAMS -c '{"function":"DeleteTicket","Args":["ticket2"]}'

peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile $ORDERER_CA -C ticketchannel -n ticketing_chaincode $PEER_CONN_PARAMS -c '{"function":"GetTicketsByRange","Args":["ticket4","ticket7"]}'

peer chaincode query -C ticketchannel -n ticketing_chaincode -c '{"Args":["QueryAllTickets"]}'

peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C ticketchannel -n ticketing_chaincode $PEER_CONN_PARAMS -c '{"function":"BuyTicket","Args":["3","Max","Gold","1"]}'

peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile $ORDERER_CA -C ticketchannel -n ticketing_chaincode $PEER_CONN_PARAMS -c '{"function":"ResaleTicket","Args":["ticket4",0.25]}'

peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile $ORDERER_CA -C ticketchannel -n ticketing_chaincode $PEER_CONN_PARAMS -c '{"function":"RedeemTicket","Args":["ticket2"]}'

peer chaincode query -C ticketchannel -n ticketing_chaincode -c '{"Args":["QueryTicketsByOwner","Aniruddha"]}'

peer chaincode query -C ticketchannel -n ticketing_chaincode -c '{"Args":["QueryTickets","{\"selector\":{\"owner\":\"Max\"}}"]}'

peer chaincode query -C ticketchannel -n ticketing_chaincode -c '{"Args":["GetTicketHistory","ticket1"]}'

peer chaincode query -C ticketchannel -n ticketing_chaincode -c '{"Args":["GetTicketsByRangeWithPagination","ticket6","ticket7","2","g1AAAAA-eJzLYWBgYMpgSmHgKy5JLCrJTq2MT8lPzkzJBYqzl2QmZ6eWmIOkOWDSyBJZAB9lEjk"]}'

peer chaincode query -C ticketchannel -n ticketing_chaincode -c '{"Args":["QueryTicketsWithPagination","{\"selector\":{\"owner\":\"Max\"}}","3",""]}'

peer chaincode query -C ticketchannel -n ticketing_chaincode -c '{"Args":["ShowTicketByEventId","31"]}'

*/
