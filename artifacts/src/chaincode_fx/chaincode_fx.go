package main

import (
  "encoding/json"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type FxChaincode struct {
}

type Transaction struct {
  FromId string `json:"fromId"`
  FromCcy string `json:"fromCcy"`
  FromAmt float64 `json:"fromAmt"`
  ToId string `json:"toId"`
  ToCcy string `json:"toCcy"`
  ToAmt float64 `json:"toAmt"`
}

type Account struct {
  Id string `json:"id"`
  Balances []Balance `json:"balances"`
}

type Balance struct {
  Ccy string `json:"ccy"`
  Amt float64 `json:"amt"`
}

func (cc *FxChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
  fmt.Println("Initializing")

  var accounts = []Account{
    Account{Id: "0895123456", Balances: []Balance{Balance{Ccy: "JPY", Amt: 1000000}}},
    Account{Id: "0895999999", Balances: []Balance{Balance{Ccy: "USD", Amt: 10000}}}
  }
  var keys []string

  for _, account := range accounts {
    // Convert to JSON
    accountBytes, err := json.Marshal(account)
    if err != nil {
      fmt.Println("error creating account " + account.Id)
      return nil, errors.New("Error creating account " + account.Id)
    }
    // Add to world state
    err = stub.PutState(account.Id, accountBytes)
    fmt.Println("created account" + account.Id)

    keys = append(keys, account.Id)
  }

	keysBytes, err := json.Marshal(&keys)
	if err != nil {
		fmt.Println("Error marshalling keys")
		return nil, errors.New("Error marshalling the keys")
	}
	fmt.Println("Put state on AccountKeys")
	err = stub.PutState("AccountKeys", keysBytes)
	if err != nil {
		fmt.Println("Error writting keys")
		return nil, errors.New("Error writing the keys")
	}

  return nil, nil
}


func (cc *FxChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
  fmt.Println("Invoking")
  // Handle by function name
  if function == "tradeFx" {
    return cc.tradeFx(stub, args)
  }

  return nil, errors.New("Received unknown function")
}


func (cc *FxChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
  fmt.Println("Querying")
  // Handle by function name
  if function == "getBalances" {
    return cc.getBalances(stub, args)
  }

  return nil, errors.New("Received unknown function")
}


func (cc *FxChaincode) tradeFx(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
  fmt.Println("Trading FX")

  if len(args) != 1 {
    return nil, errors.New("Incorrect number of arguments. Expecting FX transaction record")
  }

  var tr Transaction

  fmt.Println("Unmarshalling Transaction")
	err := json.Unmarshal([]byte(args[0]), &tr)
	if err != nil {
		fmt.Println("Error Unmarshalling Transaction")
		return nil, errors.New("Invalid FX transaction")
	}

  var fromId Account
	fmt.Println("Getting State on fromId " + tr.FromId)
	fromIdBytes, err := stub.GetState(tr.FromId)
	if err != nil {
		fmt.Println("Account not found " + tr.FromId)
		return nil, errors.New("Account not found " + tr.FromId)
	}

	fmt.Println("Unmarshalling FromId ")
	err = json.Unmarshal(fromIdBytes, &fromId)
	if err != nil {
		fmt.Println("Error unmarshalling account " + tr.FromId)
		return nil, errors.New("Error unmarshalling account " + tr.FromId)
	}

  var toId Account
	fmt.Println("Getting State on toId " + tr.ToId)
	toIdBytes, err := stub.GetState(tr.ToId)
	if err != nil {
		fmt.Println("Account not found " + tr.ToId)
		return nil, errors.New("Account not found " + tr.ToId)
	}

	fmt.Println("Unmarshalling ToId ")
	err = json.Unmarshal(toIdBytes, &toId)
	if err != nil {
		fmt.Println("Error unmarshalling account " + tr.ToId)
		return nil, errors.New("Error unmarshalling account " + tr.ToId)
	}

  fmt.Println("Checking fromId has enough amount")
  fromCcyFound := false
  fromAmt := 0.00
  for _, balance := range fromId.Balances {
    if balance.Ccy == tr.FromCcy {
      fromCcyFound = true
      fromAmt = balance.Amt
    }
  }

  // If fromId doesn't own this currency
  if fromCcyFound == false {
    fmt.Println("The account " + tr.FromId + " doesn't own " + tr.FromCcy)
    return nil, errors.New("The account " + tr.FromId + " doesn't own " + tr.FromCcy)
  } else {
    fmt.Println("The FromId does own this currency")
  }

  // If fromId doesn't own enough amount of this currency
  if fromAmt < tr.FromAmt {
    fmt.Println("The account " + tr.FromId + " doesn't own enough amount of " + tr.FromCcy)
    return nil, errors.New("The account " + tr.FromId + " doesn't own enough amount of " + tr.FromCcy)
  } else {
    fmt.Println("The FromId owns enough amount of this currency")
  }

  fmt.Println("Checking toId has enough amount")
	toCcyFound := false
	toAmt := 0.00
	for _, balance := range toId.Balances {
		if balance.Ccy == tr.ToCcy {
			toCcyFound = true
			toAmt = balance.Amt
		}
	}

  // If toId doesn't own this currency
	if toCcyFound == false {
		fmt.Println("The account " + tr.ToId + "doesn't own " + tr.ToCcy)
		return nil, errors.New("The account " + tr.ToId + "doesn't own " + tr.ToCcy)
	} else {
		fmt.Println("The ToId does own this currency")
	}

	// If toId doesn't own enough amount of this currency
	if toAmt < tr.ToAmt {
		fmt.Println("The account " + tr.ToId + " doesn't own enough amount of " + tr.ToCcy)
		return nil, errors.New("The account " + tr.ToId + " doesn't own enough amount of " + tr.ToCcy)
	} else {
		fmt.Println("The ToId owns enough amount of this currency")
	}

  fmt.Println("Transfering currencies for fromId")
  toCcyFound = false
  for key, balance := range fromId.Balances {
    if balance.Ccy == tr.FromCcy {
      fmt.Println("Reducing Amount from FromId")
			fromId.Balances[key].Amt -= tr.FromAmt
    }
    if balance.Ccy == tr.ToCcy {
      toCcyFound = true
      fmt.Println("Increasing Amount to FromId")
      fromId.Balances[key].Amt += tr.ToAmt
    }
  }

  if toCcyFound == false {
		var newBalanceFromId Balance
		fmt.Println("As ToCcy was not found, appending the currency to FromId")
		newBalanceFromId.Ccy = tr.ToCcy
		newBalanceFromId.Amt = tr.ToAmt
		fromId.Balances = append(fromId.Balances, newBalanceFromId)
	}

  fmt.Println("Transfering currencies for toId")
  fromCcyFound = false
  for key, balance := range toId.Balances {
    if balance.Ccy == tr.ToCcy {
      fmt.Println("Reducing Amount from ToId")
			toId.Balances[key].Amt -= tr.ToAmt
    }
    if balance.Ccy == tr.FromCcy {
      fromCcyFound = true
      fmt.Println("Increasing Amount to ToId")
      toId.Balances[key].Amt += tr.FromAmt
    }
  }

  if fromCcyFound == false {
		var newBalanceToId Balance
		fmt.Println("As FromCcy was not found, appending the currency to ToId")
		newBalanceToId.Ccy = tr.FromCcy
		newBalanceToId.Amt = tr.FromAmt
		toId.Balances = append(toId.Balances, newBalanceToId)
	}

  // Write everything back
  // FromId
  fromIdBytesToWrite, err := json.Marshal(&fromId)
	if err != nil {
		fmt.Println("Error marshalling the fromId")
		return nil, errors.New("Error marshalling the fromId")
	}
	fmt.Println("Put state on fromId")
	err = stub.PutState(tr.FromId, fromIdBytesToWrite)
	if err != nil {
		fmt.Println("Error writing the fromId back")
		return nil, errors.New("Error writing the fromId back")
	}

  // ToId
  toIdBytesToWrite, err := json.Marshal(&toId)
	if err != nil {
		fmt.Println("Error marshalling the toId")
		return nil, errors.New("Error marshalling the toId")
	}
	fmt.Println("Put state on toId")
	err = stub.PutState(tr.ToId, toIdBytesToWrite)
	if err != nil {
		fmt.Println("Error writing the toId back")
		return nil, errors.New("Error writing the toId back")
	}

  fmt.Println("Successfully completed Invoke")
  	return nil, nil
}

func (cc *FxChaincode) getBalances(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
  var allAccounts []Account

  // Get list of all the keys
	keysBytes, err := stub.GetState("AccountKeys")
	if err != nil {
		fmt.Println("Error retrieving account keys")
		return nil, errors.New("Error retrieving account keys")
	}
	var keys []string
	err = json.Unmarshal(keysBytes, &keys)
	if err != nil {
		fmt.Println("Error unmarshalling account keys")
		return nil, errors.New("Error unmarshalling account keys")
	}

	// Get all the accounts
	for _, value := range keys {
		accountBytes, err := stub.GetState(value)

		var account Account
		err = json.Unmarshal(accountBytes, &account)
		if err != nil {
			fmt.Println("Error retrieving account " + value)
			return nil, errors.New("Error retrieving account " + value)
		}

		fmt.Println("Appending Account" + value)
		allAccounts = append(allAccounts, account)
	}

  allAccountsBytes, err := json.Marshal(&allAccounts)
	if err != nil {
		fmt.Println("Error marshalling allAccounts")
		return nil, err
	}
  fmt.Println("All success, returning allAccounts")
	return allAccountsBytes, nil
}

func main() {
  err := shim.Start(new(FxChaincode))
  if err != nil {
    fmt.Printf("Error starting chaincode: %s", err)
  }
}
