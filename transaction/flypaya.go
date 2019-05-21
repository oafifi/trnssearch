package transaction

import (
  "io"
  //"log"
  "encoding/json"
  "errors"
  "fmt"
  "strings"
)

//FlypayA is the format of FlypayA transaction
type FlypayA struct {
  Amount float64
  Currency string
  StatusCode int
  OrderReference string
  TransactionID string
}

func (flypayA FlypayA) status() (string, error) {
  switch flypayA.StatusCode {
  case 1:
    return StatusAuthorised,nil
  case 2:
    return StatusDecline,nil
  case 3:
    return StatusRefunded,nil
  }

  return "", errors.New("Invalid status code")
}

//FlypayADecoder is a struct that implements Decoder interface
type FlypayADecoder struct{}

//Decodes the next json value in the stream to Transaction
func (flypayADecoder FlypayADecoder) decodeTransaction(d *json.Decoder) (Transaction,error) {

  var t FlypayA
  var stdT Transaction

  // decode an array value (Message)
  err := d.Decode(&t)
  if err != nil {
    return stdT,err
  }

  tStatus,err := t.status()
  if err != nil {
    return stdT,err
  }
  stdT = Transaction{ProviderFlypayA,tStatus,
    t.Amount,t.Currency,t.TransactionID,
    t.OrderReference}

  return stdT,nil
}

func (flypayADecoder FlypayADecoder) Decode(reader io.Reader) ([]Transaction,error) {

  var transactionsList []Transaction

  decoder := json.NewDecoder(reader)

	// read open curly {
	t, err := decoder.Token()
	if err != nil {
		return transactionsList, err
	}
  d, validD := t.(json.Delim)
  if validD != true || d.String() != "{" {
    return transactionsList, fmt.Errorf("Invalid FlypayA json format, expected { got %v",t)
  }

  //read transactions key
  t, err = decoder.Token()
  if err != nil {
		return transactionsList, err
	}
  s, validS := t.(string)
  if validS != true || strings.ToLower(s) != "transactions" {
    return transactionsList, fmt.Errorf("Invalid FlypayA json format, expected transactions key got %v",t)
  }

  //read transactions open bracket [
	t, err = decoder.Token()
  if err != nil {
		return transactionsList, err
	}
  d, validD = t.(json.Delim)
  if validD != true || d.String() != "[" {
    return transactionsList, fmt.Errorf("Invalid FlypayA json format, expected [ got %v",t)
  }

	// while the transactions array contains more transaction items
	for decoder.More() {
		stdTransaction, err := flypayADecoder.decodeTransaction(decoder)
    if err == nil {
      transactionsList = append(transactionsList, stdTransaction)
    }
	}

  return transactionsList,nil
}
