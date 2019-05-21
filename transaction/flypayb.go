package transaction

import (
  "io"
  //"log"
  "encoding/json"
  "errors"
  "fmt"
  "strings"
)

//FlypayB is the format of FlypayB transaction
type FlypayB struct {
  Value float64
  TransactionCurrency string
  StatusCode int
  OrderInfo string
  PaymentID string
}

func (flypayB FlypayB) status() (string,error) {
  switch flypayB.StatusCode {
  case 100:
    return StatusAuthorised,nil
  case 200:
    return StatusDecline,nil
  case 300:
    return StatusRefunded,nil
  }

  return "", errors.New("Invalid status code")
}

//FlypayBDecoder is a struct that implements Decoder interface
type FlypayBDecoder struct{}

//Decodes the next json value in the stream to Transaction
func (flypayBDecoder FlypayBDecoder) decodeTransaction(d *json.Decoder) (Transaction,error) {

  var t FlypayB
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
  stdT = Transaction{ProviderFlypayB,tStatus,
    t.Value,t.TransactionCurrency,t.PaymentID,
    t.OrderInfo}

  return stdT,nil
}

func (flypayBDecoder FlypayBDecoder) Decode(reader io.Reader) ([]Transaction,error) {

  var transactionsList []Transaction

  decoder := json.NewDecoder(reader)

  // read open curly {
	t, err := decoder.Token()
	if err != nil {
		return transactionsList, err
	}
  d, validD := t.(json.Delim)
  if validD != true || d.String() != "{" {
    return transactionsList, fmt.Errorf("Invalid FlypayB json format, expected { got %v",t)
  }

  //read transactions key
  t, err = decoder.Token()
  if err != nil {
		return transactionsList, err
	}
  s, validS := t.(string)
  if validS != true || strings.ToLower(s) != "transactions" {
    return transactionsList, fmt.Errorf("Invalid FlypayB json format, expected transactions key got %v",t)
  }

  //read transactions open bracket [
	t, err = decoder.Token()
  if err != nil {
		return transactionsList, err
	}
  d, validD = t.(json.Delim)
  if validD != true || d.String() != "[" {
    return transactionsList, fmt.Errorf("Invalid FlypayB json format, expected [ got %v",t)
  }

	// while the transactions array contains more transaction items
	for decoder.More() {
		stdTransaction, err := flypayBDecoder.decodeTransaction(decoder)
    if err == nil {
      transactionsList = append(transactionsList, stdTransaction)
    }
	}

  return transactionsList,nil
}
