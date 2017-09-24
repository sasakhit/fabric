package main

import (
  "encoding/json"
	"errors"
	"fmt"
  "bytes"

	"github.com/hyperledger/fabric/core/chaincode/shim"
  pb "github.com/hyperledger/fabric/protos/peer"
)

var logger = shim.NewLogger("chaincode_sl")

type SlChaincode struct {
}

type Transaction struct {
  BRInd string `json:"brInd"`
  Borrower string `json:"borrower"`
  Lender string `json:"lender"`
  TradeDate string `json:"tradeDate"`
  SettleDate string `json:"settleDate"`
  SecCode string `json:"secCode"`
  Qty float64 `json:"qty"`
  Ccy string `json:"ccy"`
  Amt float64 `json:"amt"`
}

type Outstanding struct {
  Borrower string `json:"borrower"`
  Lender string `json:"lender"`
  SecCode string `json:"secCode"`
  Qty float64 `json:"qty"`
  Price float64 `json:"price"`
  Mtm float64 `json:"mtm"`
}

func (cc *SlChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Info("########### chaincode_sl Init ###########")

	return shim.Success(nil)
}

func (cc *SlChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Info("########### chaincode_sl Invoke ###########")

	function, args := stub.GetFunctionAndParameters()

	if function == "getOutstandings" {
		return cc.getOutstandings(stub, args)
	} else if function == "getTransactions" {
    return cc.getTransactions(stub, args)
  } else if function == "tradeSl" {
    return cc.tradeSl(stub, args)
  } else if function == "calcMarginCall" {
    return cc.calcMarginCall(stub, args)
  } else if function == "offsetOutstandings" {
    return cc.offsetOutstandings(stub, args)
  } else if function == "revaluateMtm" {
    return cc.revaluateMtm(stub, args)
  }

	logger.Errorf("Unknown action, check the first argument: %v", args[0])
	return shim.Error(fmt.Sprintf("Unknown action, check the first argument: %v", args[0]))
}

func (cc *SlChaincode) getOutstanding(stub shim.ChaincodeStubInterface, key string) (Outstanding, bool, error) {
  var outstanding Outstanding
  found := false

  outstandingBytes, err := stub.GetState(key)
  if err != nil {
    logger.Error("Error retrieving outstanding " + key)
    return outstanding, found, errors.New("Error retrieving outstanding " + key)
  }

	if outstandingBytes != nil {
		logger.Info("Outstanding found " + key)
		found = true

    err = json.Unmarshal(outstandingBytes, &outstanding)
    if err != nil {
      logger.Error("Error unmarshalling outstanding " + key)
      return outstanding, found, errors.New("Error unmarshalling outstanding " + key)
    }
    logger.Info("Outstanding found: ", outstanding)
	}

  return outstanding, found, nil
}

func (cc *SlChaincode) addKeys(stub shim.ChaincodeStubInterface, newKeys []string) (error) {
  keysBytes, err := stub.GetState("OutstandingKeys")
  if err != nil {
    logger.Error("Error retrieving Outstanding keys")
    return errors.New("Error retrieving Outstanding keys")
  }

  var keys []string
  if keysBytes != nil {
    err = json.Unmarshal(keysBytes, &keys)
    if err != nil {
      logger.Error("Error unmarshalling Outstanding keys")
      return errors.New("Error unmarshalling Outstanding keys")
    }
  }

  keys = append(keys, newKeys...)

  keysBytesToWrite, err := json.Marshal(&keys)
  if err != nil {
    logger.Error("Error marshalling keys")
    return errors.New("Error marshalling the keys")
  }

  logger.Info("Put state on OutstandingKeys")
  err = stub.PutState("OutstandingKeys", keysBytesToWrite)
  if err != nil {
    logger.Error("Error writting keys")
    return errors.New("Error writing the keys")
  }

  return nil
}

func (cc *SlChaincode) deleteKey(stub shim.ChaincodeStubInterface, oldKey string) (error) {
  keysBytes, err := stub.GetState("OutstandingKeys")
  if err != nil {
    logger.Error("Error retrieving Outstanding keys")
    return errors.New("Error retrieving Outstanding keys")
  }

  var keys []string
  if keysBytes != nil {
    err = json.Unmarshal(keysBytes, &keys)
    if err != nil {
      logger.Error("Error unmarshalling Outstanding keys")
      return errors.New("Error unmarshalling Outstanding keys")
    }
  }

  var newKeys []string
  for _, key := range keys {
    if key != oldKey {
      newKeys = append(newKeys, key)
    }
  }

  keysBytesToWrite, err := json.Marshal(&newKeys)
  if err != nil {
    logger.Error("Error marshalling keys")
    return errors.New("Error marshalling the keys")
  }

  logger.Info("Put state on OutstandingKeys")
  err = stub.PutState("OutstandingKeys", keysBytesToWrite)
  if err != nil {
    logger.Error("Error writting keys")
    return errors.New("Error writing the keys")
  }

  return nil
}

func (cc *SlChaincode) deleteOutstanding(stub shim.ChaincodeStubInterface, key string) (error) {
  logger.Info("Delete state")

  err := stub.DelState(key)
  if err != nil {
    logger.Error("Error deleting outstanding " + key)
    return errors.New("Error deleting " + key)
  }

  err = cc.deleteKey(stub, key)
  if err != nil {
    logger.Error("Error deleting outstanding " + key)
    return errors.New("Error deleting " + key)
  }

  return nil
}

func (cc *SlChaincode) writeOutstanding(stub shim.ChaincodeStubInterface, key string, outstanding Outstanding) (error) {
  if &outstanding == nil || outstanding == (Outstanding{}) {
    return nil
  }

  if outstanding.Qty == 0 {
    err := cc.deleteOutstanding(stub, key)
    if err != nil {
      logger.Error("Error deleting " + key)
      return errors.New("Error deleting " + key)
    }
    return nil
  }

  outstandingBytesToWrite, err := json.Marshal(&outstanding)
  if err != nil {
    logger.Error("Error marshalling outstanding " + key)
    return errors.New("Error marshalling outstanding " + key)
  }

  logger.Info("Put state on outstanding: ", outstanding)
  err = stub.PutState(key, outstandingBytesToWrite)
  if err != nil {
    logger.Error("Error writing outstanding " + key)
    return errors.New("Error writing outstanding " + key)
  }

  return nil
}

func (cc *SlChaincode) writeTransaction(stub shim.ChaincodeStubInterface, transaction Transaction) (error) {
  var allTransactions []Transaction

  if &transaction == nil || transaction == (Transaction{}) {
    return nil
  }

  allTransactionsBytes, err := stub.GetState("AllTransactions")
  if err != nil {
    logger.Error("Error retrieving AllTransactions")
    return errors.New("Error retrieving AllTransactions")
  }

  // In case of no transactions
  if allTransactionsBytes != nil {
    err = json.Unmarshal(allTransactionsBytes, &allTransactions)
    if err != nil {
      logger.Error("Error unmarshalling AllTransactions")
      return errors.New("Error unmarshalling AllTransactions")
    }
  }

  allTransactions = append(allTransactions, transaction)

  allTransactionBytesToWrite, err := json.Marshal(&allTransactions)
  if err != nil {
    logger.Error("Error marshalling allTransaction")
    return errors.New("Error marshalling allTransaction")
  }

  logger.Info("Put state on transaction")
  err = stub.PutState("AllTransactions", allTransactionBytesToWrite)
  if err != nil {
    logger.Error("Error writing allTransaction")
    return errors.New("Error writing allTransaction")
  }

  return nil
}

func (cc *SlChaincode) tradeSl(stub shim.ChaincodeStubInterface, args []string) pb.Response {
  logger.Info("*** Trading Stockloan ***")

  err := cc.bookTrade(stub, args, Transaction{})
  if err != nil {
    logger.Error("Error writing Stockloan trades")
    return shim.Error("Error writing Stockloan trades")
  }

  return shim.Success(nil)
}

func (cc *SlChaincode) bookTrade(stub shim.ChaincodeStubInterface, args []string, transaction Transaction) (error) {
  logger.Info("*** Booking Trade (bookTrade function) ***")

  var tr Transaction
  outstandingSl := Outstanding{}
  outstandingColl := Outstanding{}
  var err error
  var newKeys []string

  if transaction != (Transaction{}) {
    tr = transaction
  } else {
    if len(args) != 1 {
      logger.Error("Incorrect number of arguments. Expecting Stockloan transaction record")
      return errors.New("Incorrect number of arguments. Expecting Stockloan transaction record")
    }

    logger.Info("Unmarshalling Transaction " + args[0])
	  err = json.Unmarshal([]byte(args[0]), &tr)
	  if err != nil {
		  logger.Error("Error Unmarshalling Transaction")
		  return errors.New("Invalid Stockloan transaction")
    }
  }

  // For Return, change sign for quantity and amount
  if tr.BRInd == "R" {
    tr.Qty = -1 * tr.Qty
    tr.Amt = -1 * tr.Amt
  }

  slKey := tr.Borrower + "-" + tr.Lender + "-" + tr.SecCode

  logger.Info("Getting State on Borrower-Lender-SecCode " + slKey)
  outstandingSl, outstandingSlFound, err := cc.getOutstanding(stub, slKey)
  if err != nil {
    logger.Error("Error getting outstanding " + slKey)
    return errors.New("Error getting outstanding " + slKey)
  }
  logger.Info("outstandingSl: ", outstandingSl)
  logger.Info("outstandingSlFound: ", outstandingSlFound)

  // Get outstanding collateral
  collKey := tr.Lender + "-" + tr.Borrower + "-" + tr.Ccy

  logger.Info("Getting State on Lender-Borrower-Ccy " + collKey)
  outstandingColl, outstandingCollFound, err := cc.getOutstanding(stub, collKey)
  if err != nil {
    logger.Error("Error getting outstanding " + collKey)
    return errors.New("Error getting outstanding " + collKey)
  }
  logger.Info("outstandingColl: ", outstandingColl)
  logger.Info("outstandingCollFound: ", outstandingCollFound)

  if tr.BRInd == "R" && ( outstandingSl.Qty + tr.Qty < 0 || outstandingColl.Qty + tr.Amt < 0) {
    return errors.New("Not enough outstandings to return")
  }

  logger.Info("Transfering outstanding")
  if tr.Qty != 0 {
    if outstandingSlFound == true {
      outstandingSl.Qty += tr.Qty
    } else {
      outstandingSl = Outstanding{Borrower: tr.Borrower, Lender: tr.Lender, SecCode: tr.SecCode, Qty: tr.Qty, Price: 0, Mtm: 0}
      newKeys = append(newKeys, slKey)
    }
  }

  if tr.Amt != 0 {
    if outstandingCollFound == true {
      outstandingColl.Qty += tr.Amt
    } else {
      outstandingColl = Outstanding{Borrower: tr.Lender, Lender: tr.Borrower, SecCode: tr.Ccy, Qty: tr.Amt, Price: 0, Mtm: 0}
      newKeys = append(newKeys, collKey)
    }
  }

  // Write to World State
  // Transaction
  err = cc.writeTransaction(stub, tr)
  if err != nil {
    logger.Error("Error writing transaction")
    return errors.New("Error writing transaction")
  }

  // Add new OutstandingKeys to World State
  if len(newKeys) > 0 {
    err = cc.addKeys(stub, newKeys)
    if err != nil {
  		logger.Error("Error adding new OutstandingKeys")
  		return errors.New("Error adding new OutstandingKeys")
  	}
  }

  // Write to World State
  // Stockloan
  logger.Info("New outstandingSl: ", outstandingSl)
  err = cc.writeOutstanding(stub, slKey, outstandingSl)
  if err != nil {
	  logger.Error("Error writing outstanding " + slKey)
	  return errors.New("Error writing outstanding " + slKey)
  }

  // Collateral
  logger.Info("New outstandingColl: ", outstandingColl)
  err = cc.writeOutstanding(stub, collKey, outstandingColl)
  if err != nil {
		logger.Error("Error writing outstanding " + collKey)
		return errors.New("Error writing outstanding " + collKey)
	}

  logger.Info("Successfully booked trade")
  return nil
}

func (cc *SlChaincode) getPrice(stub shim.ChaincodeStubInterface, secCode string) (float64, error) {
  var price float64
  if secCode == "JPY" {
    price = 1
  } else {
    price = 1000
  }

  return price, nil
}

func (cc *SlChaincode) offsetOutstandings(stub shim.ChaincodeStubInterface, args []string) pb.Response {
  logger.Info("*** Offsetting outstandings (offsetOutstandings function) ***")

  // Get list of all the keys
  keys, err := cc.getAllKeys(stub, "OutstandingKeys")
  if err != nil {
    logger.Error("Error getting all keys")
    return shim.Error("Error writing outstanding")
  }

  // In case of no outstandings
  if keys == nil {
    return shim.Success(nil)
  }

  offsetKey := [3]string{}
  offsetTotals := map[[3]string]float64{}
  var qty float64

  for _, key := range keys {
    outstanding, outstandingFound, err := cc.getOutstanding(stub, key)
    if err != nil || outstandingFound == false {
      logger.Error("Error getting outstanding " + key)
      return shim.Error("Error getting outstanding " + key)
    }

    if outstanding.Borrower < outstanding.Lender {
      offsetKey = [3]string{outstanding.Borrower, outstanding.Lender, outstanding.SecCode}
      qty = outstanding.Qty
    } else {
      offsetKey = [3]string{outstanding.Lender, outstanding.Borrower, outstanding.SecCode}
      qty = -1 * outstanding.Qty
    }

    value, ok := offsetTotals[offsetKey]
    if ok == true {
      key1 := offsetKey[0] + "-" + offsetKey[1] + "-" + offsetKey[2]
      key2 := offsetKey[1] + "-" + offsetKey[0] + "-" + offsetKey[2]
      if qty + value == 0 {
        err = cc.deleteOutstanding(stub, key1)
        if err != nil {
      		logger.Error("Error deleting " + key1)
      		return shim.Error("Error deleting " + key1)
      	}
        err = cc.deleteOutstanding(stub, key2)
        if err != nil {
      		logger.Error("Error deleting " + key2)
      		return shim.Error("Error deleting " + key2)
      	}
      } else if qty + value > 0 {
        outstanding.Qty = qty + value
        err = cc.writeOutstanding(stub, key1, outstanding)
        if err != nil {
      		logger.Error("Error writing outstanding " + key1)
      		return shim.Error("Error writing outstanding " + key1)
      	}
        err = cc.deleteOutstanding(stub, key2)
        if err != nil {
      		logger.Error("Error deleting " + key2)
      		return shim.Error("Error deleting " + key2)
      	}
      } else { // qty + value < 0
        err = cc.deleteOutstanding(stub, key1)
        if err != nil {
      		logger.Error("Error deleting " + key1)
      		return shim.Error("Error deleting " + key1)
      	}
        outstanding.Borrower = offsetKey[1]
        outstanding.Lender = offsetKey[0]
        outstanding.Qty = -1 * (value + qty)
        err = cc.writeOutstanding(stub, key2, outstanding)
        if err != nil {
      		logger.Error("Error writing outstanding " + key2)
      		return shim.Error("Error writing outstanding " + key2)
      	}
      }
    } else {
      offsetTotals[offsetKey] = qty
    }
	}

  return shim.Success(nil)
}

func (cc *SlChaincode) revaluateMtm(stub shim.ChaincodeStubInterface, args []string) pb.Response {
  logger.Info("*** Revaluating MTM (revaluateMtm function) ***")

  // Get list of all the keys
  keys, err := cc.getAllKeys(stub, "OutstandingKeys")
  if err != nil {
    logger.Error("Error getting all keys")
    return shim.Error("Error writing outstanding")
  }

  // In case of no outstandings
  if keys == nil {
    return shim.Success(nil)
  }

  for _, key := range keys {
    outstanding, outstandingFound, err := cc.getOutstanding(stub, key)
    if err != nil || outstandingFound == false {
      logger.Error("Error getting outstanding " + key)
      return shim.Error("Error getting outstanding " + key)
    }

    price, err := cc.getPrice(stub, outstanding.SecCode)
    if err != nil {
			logger.Error("Error getting price " + outstanding.SecCode)
			return shim.Error("Error getting price " + outstanding.SecCode)
		}
    outstanding.Price = price
    outstanding.Mtm = price * outstanding.Qty

    err = cc.writeOutstanding(stub, key, outstanding)
    if err != nil {
  		logger.Error("Error writing outstanding " + key)
  		return shim.Error("Error writing outstanding " + key)
  	}
  }

  return shim.Success(nil)
}

func (cc *SlChaincode) calcMarginCall(stub shim.ChaincodeStubInterface, args []string) pb.Response {
  logger.Info("Calculating margin call")

  if len(args) != 1 {
    return shim.Error("Incorrect number of arguments. Expecting revaluation date")
  }
  revalDate := args[0]
  logger.Info("revalDate: ", revalDate)

  // Get list of all the keys
  keys, err := cc.getAllKeys(stub, "OutstandingKeys")
  if err != nil {
    logger.Error("Error getting all keys")
    return shim.Error("Error writing outstanding")
  }

  // In case of no outstandings
  if keys == nil {
    return shim.Success(nil)
  }

  // Margin calculation
  outTotalKey := [2]string{}
  outTotals := map[[2]string]float64{}
  var mtm float64
  marginCall := Transaction{}

	for _, key := range keys {
    outstanding, outstandingFound, err := cc.getOutstanding(stub, key)
    if err != nil || outstandingFound == false {
      logger.Error("Error getting outstanding " + key)
      return shim.Error("Error getting outstanding " + key)
    }

    if outstanding.Borrower < outstanding.Lender {
      outTotalKey = [2]string{outstanding.Borrower, outstanding.Lender}
      mtm = outstanding.Mtm
    } else {
      outTotalKey = [2]string{outstanding.Lender, outstanding.Borrower}
      mtm = -1 * outstanding.Mtm
    }

    value, ok := outTotals[outTotalKey]
    if ok == true {
      delete(outTotals, outTotalKey)
      outTotals[outTotalKey] = value + mtm
    } else {
      outTotals[outTotalKey] = mtm
    }
	}

  threshold := 100.0
  for key, value := range outTotals {
    logger.Info("value: ", value)
    if value > threshold {
      marginCall = Transaction{BRInd: "M", Borrower: key[0], Lender: key[1], TradeDate: revalDate, SettleDate: revalDate, SecCode: "", Qty: 0, Ccy: "JPY", Amt: value}
    } else if value < -1 * threshold {
      marginCall = Transaction{BRInd: "M", Borrower: key[1], Lender: key[0], TradeDate: revalDate, SettleDate: revalDate, SecCode: "", Qty: 0, Ccy: "JPY", Amt: -1 * value}
    }

    if marginCall != (Transaction{}) {
      err := cc.bookTrade(stub, nil, marginCall)
      if err != nil {
  		  logger.Error("Error writing marginCall")
  		  return shim.Error("Error writing marginCall")
      }
    }

  }

  logger.Info("Successfully completed Invoke")
  return shim.Success(nil)
}

func (cc *SlChaincode) getAllKeys(stub shim.ChaincodeStubInterface, keyName string) ([]string, error) {
  // Get list of all the keys
  keysBytes, err := stub.GetState(keyName)
  if err != nil {
    logger.Error("Error retrieving Outstanding keys " + keyName)
    return nil, errors.New("Error retrieving Outstanding keys " + keyName)
  }

  // In case of no outstandings
  if keysBytes == nil {
    return nil, nil
  }

  var keys []string
  err = json.Unmarshal(keysBytes, &keys)
  if err != nil {
    logger.Error("Error unmarshalling Outstanding keys")
    return nil, errors.New("Error unmarshalling Outstanding keys")
  }

  return keys, nil
}

func (cc *SlChaincode) getOutstandings(stub shim.ChaincodeStubInterface, args []string) pb.Response {
  var allOutstandings []Outstanding

  // Get list of all the keys
  keys, err := cc.getAllKeys(stub, "OutstandingKeys")
  if err != nil {
    logger.Error("Error getting all keys")
    return shim.Error("Error getting all keys")
  }

  // Sample code for rich query
  queryString := "{\"selector\":{\"secCode\":\"8756\"}}"
  logger.Info("queryString: ", queryString)
	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
  logger.Info("queryResults: ", queryResults)

  // In case of no outstandings
  if keys == nil {
    noOutstandings := []Outstanding{Outstanding{Borrower: "No outstandings", Lender: "", SecCode: "", Qty: 0, Price: 0, Mtm: 0}}

    noOutstandingsBytes, err := json.Marshal(&noOutstandings)
    if err != nil {
      logger.Error("Error marshalling noOutstandings")
      return shim.Error("Error marshalling noOutstandings")
    }

    return shim.Success(noOutstandingsBytes)
  }

	// Get all the outstandings
	for _, key := range keys {
		outstandingBytes, err := stub.GetState(key)

		var outstanding Outstanding
		err = json.Unmarshal(outstandingBytes, &outstanding)
		if err != nil {
			logger.Error("Error retrieving outstanding " + key)
			return shim.Error("Error retrieving outstanding " + key)
		}

		logger.Info("Appending outstanding " + key)
		allOutstandings = append(allOutstandings, outstanding)
	}

  allOutstandingsBytes, err := json.Marshal(&allOutstandings)
	if err != nil {
		logger.Error("Error marshalling allOutstandings")
		return shim.Error("Error marshalling allOutstandings")
	}

  logger.Info("All success, returning allOutstandings")
	return shim.Success(allOutstandingsBytes)
}

func (cc *SlChaincode) getTransactions(stub shim.ChaincodeStubInterface, args []string) pb.Response {
  allTransactionsBytes, err := stub.GetState("AllTransactions")
  if err != nil {
    logger.Error("Error retrieving AllTransactions")
    return shim.Error("Error retrieving AllTransactions")
  }

  if allTransactionsBytes == nil {
    noTransactions := []Transaction{Transaction{BRInd: "No transactions", Borrower: "", Lender: "", TradeDate: "", SettleDate: "", SecCode: "", Qty: 0, Ccy: "", Amt: 0}}

    noTransactionsBytes, err := json.Marshal(&noTransactions)
    if err != nil {
      logger.Error("Error marshalling noTransactions")
      return shim.Error("Error marshalling noTransactions")
    }

    return shim.Success(noTransactionsBytes)
  }

  return shim.Success(allTransactionsBytes)
}

// =========================================================================================
// getQueryResultForQueryString executes the passed in query string.
// Result set is built and returned as a byte array containing the JSON results.
// =========================================================================================
func getQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	logger.Info("- getQueryResultForQueryString queryString: ", queryString)

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	logger.Info("- getQueryResultForQueryString queryResult: ", buffer.String())

	return buffer.Bytes(), nil
}

func main() {
	err := shim.Start(new(SlChaincode))
	if err != nil {
		logger.Errorf("Error starting Stockloan chaincode: %s", err)
	}
}
