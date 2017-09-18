package main

import (
  "encoding/json"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type SlChaincode struct {
}

type Transaction struct {
  BRInd string `json:"brInd"`
  Borrower string `json:"borrower"`
  Lender string `json:"lender"`
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

func (cc *SlChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
  fmt.Println("Initializing")

  return nil, nil
}

func (cc *SlChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
  fmt.Println("Invoking")
  // Handle by function name
  if function == "tradeSl" {
    return nil, cc.tradeSl(stub, args, Transaction{})
  } else if function == "calcMarginCall" {
    return nil, cc.calcMarginCall(stub, args)
  } else if function == "offsetOutstandings" {
    return nil, cc.offsetOutstandings(stub, args)
  }

  return nil, errors.New("Received unknown function")
}

func (cc *SlChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
  fmt.Println("Querying")
  // Handle by function name
  if function == "getOutstandings" {
    return cc.getOutstandings(stub, args)
  } else if function == "getTransactions" {
    return cc.getTransactions(stub, args)
  }

  return nil, errors.New("Received unknown function")
}

func (cc *SlChaincode) getOutstanding(stub shim.ChaincodeStubInterface, key string) (Outstanding, bool, error) {
  var outstanding Outstanding
  found := false

  outstandingBytes, err := stub.GetState(key)
  if err != nil {
    fmt.Println("Error retrieving outstanding " + key)
    return outstanding, found, errors.New("Error retrieving outstanding " + key)
  }

	if outstandingBytes != nil {
		fmt.Println("Outstanding found " + key)
		found = true

    fmt.Println("Unmarshalling Outstanding")
    err = json.Unmarshal(outstandingBytes, &outstanding)
    if err != nil {
      fmt.Println("Error unmarshalling outstanding " + key)
      return outstanding, found, errors.New("Error unmarshalling outstanding " + key)
    }
	}

  return outstanding, found, nil
}

func (cc *SlChaincode) addKeys(stub shim.ChaincodeStubInterface, newKeys []string) (error) {
  keysBytes, err := stub.GetState("OutstandingKeys")
  if err != nil {
    fmt.Println("Error retrieving Outstanding keys")
    return errors.New("Error retrieving Outstanding keys")
  }

  var keys []string
  if keysBytes != nil {
    err = json.Unmarshal(keysBytes, &keys)
    if err != nil {
      fmt.Println("Error unmarshalling Outstanding keys")
      return errors.New("Error unmarshalling Outstanding keys")
    }
  }

  keys = append(keys, newKeys...)

  keysBytesToWrite, err := json.Marshal(&keys)
  if err != nil {
    fmt.Println("Error marshalling keys")
    return errors.New("Error marshalling the keys")
  }

  fmt.Println("Put state on OutstandingKeys")
  err = stub.PutState("OutstandingKeys", keysBytesToWrite)
  if err != nil {
    fmt.Println("Error writting keys")
    return errors.New("Error writing the keys")
  }

  return nil
}

func (cc *SlChaincode) deleteKey(stub shim.ChaincodeStubInterface, oldKey string) (error) {
  keysBytes, err := stub.GetState("OutstandingKeys")
  if err != nil {
    fmt.Println("Error retrieving Outstanding keys")
    return errors.New("Error retrieving Outstanding keys")
  }

  var keys []string
  if keysBytes != nil {
    err = json.Unmarshal(keysBytes, &keys)
    if err != nil {
      fmt.Println("Error unmarshalling Outstanding keys")
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
    fmt.Println("Error marshalling keys")
    return errors.New("Error marshalling the keys")
  }

  fmt.Println("Put state on OutstandingKeys")
  err = stub.PutState("OutstandingKeys", keysBytesToWrite)
  if err != nil {
    fmt.Println("Error writting keys")
    return errors.New("Error writing the keys")
  }

  return nil
}

func (cc *SlChaincode) deleteOutstanding(stub shim.ChaincodeStubInterface, key string) (error) {
  fmt.Println("Delete state")

  err := stub.DelState(key)
  if err != nil {
    fmt.Println("Error deleting outstanding " + key)
    return errors.New("Error deleting " + key)
  }

  err = cc.deleteKey(stub, key)
  if err != nil {
    fmt.Println("Error deleting outstanding " + key)
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
      fmt.Println("Error deleting " + key)
      return errors.New("Error deleting " + key)
    }
    return nil
  }

  outstandingBytesToWrite, err := json.Marshal(&outstanding)
  if err != nil {
    fmt.Println("Error marshalling outstanding " + key)
    return errors.New("Error marshalling outstanding " + key)
  }

  fmt.Println("Put state on outstanding")
  err = stub.PutState(key, outstandingBytesToWrite)
  if err != nil {
    fmt.Println("Error writing outstanding " + key)
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
    fmt.Println("Error retrieving AllTransactions")
    return errors.New("Error retrieving AllTransactions")
  }

  // In case of no transactions
  if allTransactionsBytes != nil {
    err = json.Unmarshal(allTransactionsBytes, &allTransactions)
    if err != nil {
      fmt.Println("Error unmarshalling AllTransactions")
      return errors.New("Error unmarshalling AllTransactions")
    }
  }

  allTransactions = append(allTransactions, transaction)

  allTransactionBytesToWrite, err := json.Marshal(&allTransactions)
  if err != nil {
    fmt.Println("Error marshalling allTransaction")
    return errors.New("Error marshalling allTransaction")
  }

  fmt.Println("Put state on transaction")
  err = stub.PutState("AllTransactions", allTransactionBytesToWrite)
  if err != nil {
    fmt.Println("Error writing allTransaction")
    return errors.New("Error writing allTransaction")
  }

  return nil
}

func (cc *SlChaincode) tradeSl(stub shim.ChaincodeStubInterface, args []string, transaction Transaction) (error) {
  fmt.Println("Trading Stockloan")

  var tr Transaction
  outstandingSl := Outstanding{}
  outstandingColl := Outstanding{}
  var err error
  var newKeys []string

  if transaction != (Transaction{}) {
    tr = transaction
  } else {
    if len(args) != 1 {
      return errors.New("Incorrect number of arguments. Expecting Stockloan transaction record")
    }

    fmt.Println("Unmarshalling Transaction")
	  err = json.Unmarshal([]byte(args[0]), &tr)
	  if err != nil {
		  fmt.Println("Error Unmarshalling Transaction")
		  return errors.New("Invalid Stockloan transaction")
    }
  }

  // For Return, change sign for quantity and amount
  if tr.BRInd == "R" {
    tr.Qty = -1 * tr.Qty
    tr.Amt = -1 * tr.Amt
  }

  slKey := tr.Borrower + "-" + tr.Lender + "-" + tr.SecCode

  fmt.Println("Getting State on Borrower-Lender-SecCode " + slKey)
  outstandingSl, outstandingSlFound, err := cc.getOutstanding(stub, slKey)
  if err != nil {
    fmt.Println("Error getting outstanding " + slKey)
    return errors.New("Error getting outstanding " + slKey)
  }

  // Get outstanding collateral
  collKey := tr.Lender + "-" + tr.Borrower + "-" + tr.Ccy

  fmt.Println("Getting State on Lender-Borrower-Ccy " + collKey)
  outstandingColl, outstandingCollFound, err := cc.getOutstanding(stub, collKey)
  if err != nil {
    fmt.Println("Error getting outstanding " + collKey)
    return errors.New("Error getting outstanding " + collKey)
  }

  if tr.BRInd == "R" && ( outstandingSl.Qty + tr.Qty < 0 || outstandingColl.Qty + tr.Amt < 0) {
    return errors.New("Not enough outstandings to return")
  }

  fmt.Println("Transfering outstanding")
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
    fmt.Println("Error writing transaction")
    return errors.New("Error writing transaction")
  }

  // Add new OutstandingKeys to World State
  if len(newKeys) > 0 {
    err = cc.addKeys(stub, newKeys)
    if err != nil {
  		fmt.Println("Error adding new OutstandingKeys")
  		return errors.New("Error adding new OutstandingKeys")
  	}
  }

  // Write to World State
  // Stockloan
  err = cc.writeOutstanding(stub, slKey, outstandingSl)
  if err != nil {
	  fmt.Println("Error writing outstanding " + slKey)
	  return errors.New("Error writing outstanding " + slKey)
  }

  // Collateral
  err = cc.writeOutstanding(stub, collKey, outstandingColl)
  if err != nil {
		fmt.Println("Error writing outstanding " + collKey)
		return errors.New("Error writing outstanding " + collKey)
	}

  // Offset Outstandings
  err = cc.offsetOutstandings(stub, nil)
  if err != nil {
		fmt.Println("Error offseting outstandings")
		return errors.New("Error offseting outstandings")
	}

  // Revaluating MTM
  err = cc.revaluateMtm(stub, nil)
  if err != nil {
    fmt.Println("Error revaluation MTM")
    return errors.New("Error revaluation MTM")
  }

  fmt.Println("Successfully completed Invoke")
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

func (cc *SlChaincode) offsetOutstandings(stub shim.ChaincodeStubInterface, args []string) (error) {
  fmt.Println("Offsetting outstandings")

  // Get list of all the keys
  keys, err := cc.getAllKeys(stub, "OutstandingKeys")
  if err != nil {
    fmt.Println("Error getting all keys")
    return errors.New("Error writing outstanding")
  }

  // In case of no outstandings
  if keys == nil {
    return nil
  }

  offsetKey := [3]string{}
  offsetTotals := map[[3]string]float64{}
  var qty float64

  for _, key := range keys {
    outstanding, outstandingFound, err := cc.getOutstanding(stub, key)
    if err != nil || outstandingFound == false {
      fmt.Println("Error getting outstanding " + key)
      return errors.New("Error getting outstanding " + key)
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
      		fmt.Println("Error deleting " + key1)
      		return errors.New("Error deleting " + key1)
      	}
        err = cc.deleteOutstanding(stub, key2)
        if err != nil {
      		fmt.Println("Error deleting " + key2)
      		return errors.New("Error deleting " + key2)
      	}
      } else if qty + value > 0 {
        outstanding.Qty = qty + value
        err = cc.writeOutstanding(stub, key1, outstanding)
        if err != nil {
      		fmt.Println("Error writing outstanding " + key1)
      		return errors.New("Error writing outstanding " + key1)
      	}
        err = cc.deleteOutstanding(stub, key2)
        if err != nil {
      		fmt.Println("Error deleting " + key2)
      		return errors.New("Error deleting " + key2)
      	}
      } else { // qty + value < 0
        err = cc.deleteOutstanding(stub, key1)
        if err != nil {
      		fmt.Println("Error deleting " + key1)
      		return errors.New("Error deleting " + key1)
      	}
        outstanding.Borrower = offsetKey[1]
        outstanding.Lender = offsetKey[0]
        outstanding.Qty = -1 * (value + qty)
        err = cc.writeOutstanding(stub, key2, outstanding)
        if err != nil {
      		fmt.Println("Error writing outstanding " + key2)
      		return errors.New("Error writing outstanding " + key2)
      	}
      }
    } else {
      offsetTotals[offsetKey] = qty
    }
	}

  return nil
}

func (cc *SlChaincode) revaluateMtm(stub shim.ChaincodeStubInterface, args []string) (error) {
  fmt.Println("Revaluating MTM")

  // Get list of all the keys
  keys, err := cc.getAllKeys(stub, "OutstandingKeys")
  if err != nil {
    fmt.Println("Error getting all keys")
    return errors.New("Error writing outstanding")
  }

  // In case of no outstandings
  if keys == nil {
    return nil
  }

  for _, key := range keys {
    outstanding, outstandingFound, err := cc.getOutstanding(stub, key)
    if err != nil || outstandingFound == false {
      fmt.Println("Error getting outstanding " + key)
      return errors.New("Error getting outstanding " + key)
    }

    price, err := cc.getPrice(stub, outstanding.SecCode)
    if err != nil {
			fmt.Println("Error getting price " + outstanding.SecCode)
			return errors.New("Error getting price " + outstanding.SecCode)
		}
    outstanding.Price = price
    outstanding.Mtm = price * outstanding.Qty

    err = cc.writeOutstanding(stub, key, outstanding)
    if err != nil {
  		fmt.Println("Error writing outstanding " + key)
  		return errors.New("Error writing outstanding " + key)
  	}
  }

  return nil
}

func (cc *SlChaincode) calcMarginCall(stub shim.ChaincodeStubInterface, args []string) (error) {
  fmt.Println("Calculating margin call")

  err := cc.revaluateMtm(stub, nil)
  if err != nil {
    fmt.Println("Error revaluation MTM")
    return errors.New("Error revaluation MTM")
  }

  // Get list of all the keys
  keys, err := cc.getAllKeys(stub, "OutstandingKeys")
  if err != nil {
    fmt.Println("Error getting all keys")
    return errors.New("Error writing outstanding")
  }

  // In case of no outstandings
  if keys == nil {
    return nil
  }

  // Margin calculation
  outTotalKey := [2]string{}
  outTotals := map[[2]string]float64{}
  var mtm float64
  marginCall := Transaction{}

	for _, key := range keys {
    outstanding, outstandingFound, err := cc.getOutstanding(stub, key)
    if err != nil || outstandingFound == false {
      fmt.Println("Error getting outstanding " + key)
      return errors.New("Error getting outstanding " + key)
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
    if value > threshold {
      marginCall = Transaction{BRInd: "M", Borrower: key[0], Lender: key[1], SecCode: "", Qty: 0, Ccy: "JPY", Amt: value}
    } else if value < -1 * threshold {
      marginCall = Transaction{BRInd: "M", Borrower: key[1], Lender: key[0], SecCode: "", Qty: 0, Ccy: "JPY", Amt: -1 * value}
    }

    if marginCall != (Transaction{}) {
      err := cc.tradeSl(stub, nil, marginCall)
      if err != nil {
  		  fmt.Println("Error writing marginCall")
  		  return errors.New("Error writing marginCall")
      }
    }

  }

  fmt.Println("Successfully completed Invoke")
  return nil
}

func (cc *SlChaincode) getAllKeys(stub shim.ChaincodeStubInterface, keyName string) ([]string, error) {
  // Get list of all the keys
  keysBytes, err := stub.GetState(keyName)
  if err != nil {
    fmt.Println("Error retrieving Outstanding keys " + keyName)
    return nil, errors.New("Error retrieving Outstanding keys " + keyName)
  }

  // In case of no outstandings
  if keysBytes == nil {
    return nil, nil
  }

  var keys []string
  err = json.Unmarshal(keysBytes, &keys)
  if err != nil {
    fmt.Println("Error unmarshalling Outstanding keys")
    return nil, errors.New("Error unmarshalling Outstanding keys")
  }

  return keys, nil
}

func (cc *SlChaincode) getOutstandings(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
  var allOutstandings []Outstanding

  // Get list of all the keys
  keys, err := cc.getAllKeys(stub, "OutstandingKeys")
  if err != nil {
    fmt.Println("Error getting all keys")
    return nil, errors.New("Error getting all keys")
  }

  // In case of no outstandings
  if keys == nil {
    noOutstandings := []Outstanding{Outstanding{Borrower: "No outstandings", Lender: "", SecCode: "", Qty: 0, Price: 0, Mtm: 0}}

    noOutstandingsBytes, err := json.Marshal(&noOutstandings)
    if err != nil {
      fmt.Println("Error marshalling noOutstandings")
      return nil, err
    }

    return noOutstandingsBytes, nil
  }

	// Get all the outstandings
	for _, key := range keys {
		outstandingBytes, err := stub.GetState(key)

		var outstanding Outstanding
		err = json.Unmarshal(outstandingBytes, &outstanding)
		if err != nil {
			fmt.Println("Error retrieving outstanding " + key)
			return nil, errors.New("Error retrieving outstanding " + key)
		}

		fmt.Println("Appending outstanding " + key)
		allOutstandings = append(allOutstandings, outstanding)
	}

  allOutstandingsBytes, err := json.Marshal(&allOutstandings)
	if err != nil {
		fmt.Println("Error marshalling allOutstandings")
		return nil, err
	}

  fmt.Println("All success, returning allOutstandings")
	return allOutstandingsBytes, nil
}

func (cc *SlChaincode) getTransactions(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
  allTransactionsBytes, err := stub.GetState("AllTransactions")
  if err != nil {
    fmt.Println("Error retrieving AllTransactions")
    return nil, errors.New("Error retrieving AllTransactions")
  }

  if allTransactionsBytes == nil {
    noTransactions := []Transaction{Transaction{BRInd: "No transactions", Borrower: "", Lender: "", SecCode: "", Qty: 0, Ccy: "", Amt: 0}}

    noTransactionsBytes, err := json.Marshal(&noTransactions)
    if err != nil {
      fmt.Println("Error marshalling noTransactions")
      return nil, err
    }

    return noTransactionsBytes, nil
  }

  return allTransactionsBytes, nil
}

func main() {
  err := shim.Start(new(SlChaincode))
  if err != nil {
    fmt.Printf("Error starting chaincode: %s", err)
  }
}

//stub *shim.ChaincodeStub
//stub shim.ChaincodeStubInterface
