package transaction

import (
  "testing"
  "strings"
  "encoding/json"
  //"fmt"
)

func TestFlypayA_Status(t *testing.T) {

  // set non-existant StatusCode
  transaction := FlypayA{StatusCode: 9999}
  status,err := transaction.status()
  if err == nil {
    t.Errorf("StatusCode %v is a non-existent status code and must return error", transaction.StatusCode)
  }

  transaction.StatusCode = 1
  status,err = transaction.status()
  if err != nil {
    t.Errorf("StatusCode %v has no standard status mapping", transaction.StatusCode)
  } else if status != StatusAuthorised {
    t.Errorf("StatusCode %v has wrong standard status mapping %v",
      transaction.StatusCode, status)
  }

  transaction.StatusCode = 2
  status,err = transaction.status()
  if err != nil {
    t.Errorf("StatusCode %v has no standard status mapping", transaction.StatusCode)
  } else if status != StatusDecline {
    t.Errorf("StatusCode %v has wrong standard status mapping %v",
      transaction.StatusCode, status)
  }

  transaction.StatusCode = 3
  status,err = transaction.status()
  if err != nil {
    t.Errorf("StatusCode %v has no standard status mapping", transaction.StatusCode)
  } else if status != StatusRefunded {
    t.Errorf("StatusCode %v has wrong standard status mapping %v",
      transaction.StatusCode, status)
  }
}

func TestFlypayA_DecodeTransaction(t *testing.T) {
  fDecoder := FlypayADecoder{}

  s := `
    {
      "amount":1000,
      "currency":"AUD",
      "statusCode":1,
      "orderReference":"2e58bd43-0001",
      "transactionId": "flypay-a-0001"
    }
    {
      "amount":150,
      "currency":"AUD",
      "statusCode":999,
      "orderReference":"2e58bd43-0002",
      "transactionId": "flypay-a-0002"
    }
    {
      "val500,
      "curr":"AUD",
      "code":3,
      "order":"2e58bd43-0003",
      "Id": "flypay-a-0003"
    }
  `
  r := strings.NewReader(s)
  d := json.NewDecoder(r)

  transaction,err := fDecoder.decodeTransaction(d)
  validT := Transaction{ProviderFlypayA,StatusAuthorised,1000,"AUD","flypay-a-0001","2e58bd43-0001"}
  if err != nil {
    t.Errorf("err should be nil for successful decoding, error: %v",err)
  } else if transaction != validT {
    t.Error("Wrong transaction decoding")
  }

  transaction,err = fDecoder.decodeTransaction(d)
  if err == nil {
    t.Error("err should not be nil for decoding transaction with wrong status")
  }

  transaction,err = fDecoder.decodeTransaction(d)
  if err == nil {
    t.Error("err should not be nil for decoding transaction with different structure")
  }
}



func TestFlypayADecoder_Decode(t *testing.T) {
  fDecoder := FlypayADecoder{}

  //Test first delim cases

  s := `
  {
  "transactions":[
    {
      "amount":1000,
      "currency":"AUD",
      "statusCode":1,
      "orderReference":"2e58bd43-0001",
      "transactionId": "flypay-a-0001"
    }
    ]
  }
  `

  r := strings.NewReader(s)

  transactions,err := fDecoder.Decode(r)
  validT := Transaction{ProviderFlypayA,StatusAuthorised,1000,"AUD","flypay-a-0001","2e58bd43-0001"}
  if err != nil {
    t.Errorf("err should be nil for successful decoding, error: %v",err)
  } else if len(transactions) != 1 || transactions[0] != validT {
    t.Error("Wrong transactions list building and decoding")
  }

  s = `
  "transactions":[
    {
      "amount":1000,
      "currency":"AUD",
      "statusCode":1,
      "orderReference":"2e58bd43-0001",
      "transactionId": "flypay-a-0001"
    }
    ]
  }
    `

  r.Reset(s)
  transactions,err = fDecoder.Decode(r)
  validT = Transaction{ProviderFlypayA,StatusAuthorised,1000,"AUD","flypay-a-0001","2e58bd43-0001"}
  if err == nil {
    t.Errorf("Invalid json structure without opening {")
  }

  s = `
  `

  r.Reset(s)
  transactions,err = fDecoder.Decode(r)
  validT = Transaction{ProviderFlypayA,StatusAuthorised,1000,"AUD","flypay-a-0001","2e58bd43-0001"}
  if err == nil {
    t.Errorf("Invalid json structure without opening {")
  }

  //Test second delimeter cases

  s = `
  {
  "transact":[
    {
      "amount":1000,
      "currency":"AUD",
      "statusCode":1,
      "orderReference":"2e58bd43-0001",
      "transactionId": "flypay-a-0001"
    }
    ]
  }
  `

  r.Reset(s)
  transactions,err = fDecoder.Decode(r)
  validT = Transaction{ProviderFlypayA,StatusAuthorised,1000,"AUD","flypay-a-0001","2e58bd43-0001"}
  if err == nil {
    t.Errorf("Invalid json structure without transactions key")
  }

  s = `
  {
  {"transaction":[
    {
      "amount":1000,
      "currency":"AUD",
      "statusCode":1,
      "orderReference":"2e58bd43-0001",
      "transactionId": "flypay-a-0001"
    }
    ]
  }
  }
  `

  r.Reset(s)
  transactions,err = fDecoder.Decode(r)
  validT = Transaction{ProviderFlypayA,StatusAuthorised,1000,"AUD","flypay-a-0001","2e58bd43-0001"}
  if err == nil {
    t.Errorf("Invalid json structure without excess delim rather than transactions key")
  }

  s = `
  {
  `

  r.Reset(s)
  transactions,err = fDecoder.Decode(r)
  validT = Transaction{ProviderFlypayA,StatusAuthorised,1000,"AUD","flypay-a-0001","2e58bd43-0001"}
  if err == nil {
    t.Errorf("Invalid json structure with only {")
  }

  //Test third delim cases

  s = `
  {
  "transactions":{
    {
      "amount":1000,
      "currency":"AUD",
      "statusCode":1,
      "orderReference":"2e58bd43-0001",
      "transactionId": "flypay-a-0001"
    }
    ]
  }
    `

  r.Reset(s)
  transactions,err = fDecoder.Decode(r)
  validT = Transaction{ProviderFlypayA,StatusAuthorised,1000,"AUD","flypay-a-0001","2e58bd43-0001"}
  if err == nil {
    t.Errorf("Invalid json structure without opening [")
  }

  s = `
  {
  "transactions":
    {
      "amount":1000,
      "currency":"AUD",
      "statusCode":1,
      "orderReference":"2e58bd43-0001",
      "transactionId": "flypay-a-0001"
    }
    ]
  }
    `

  r.Reset(s)
  transactions,err = fDecoder.Decode(r)
  validT = Transaction{ProviderFlypayA,StatusAuthorised,1000,"AUD","flypay-a-0001","2e58bd43-0001"}
  if err == nil {
    t.Errorf("Invalid json structure without opening [")
  }

  s = `
  {
  "transactions":
  `

  r.Reset(s)
  transactions,err = fDecoder.Decode(r)
  validT = Transaction{ProviderFlypayA,StatusAuthorised,1000,"AUD","flypay-a-0001","2e58bd43-0001"}
  if err == nil {
    t.Errorf("Invalid json structure without opening opening [")
  }
}
