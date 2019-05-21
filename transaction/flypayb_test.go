package transaction

import (
  "testing"
  "strings"
  "encoding/json"
)

func TestFlypayBStatus(t *testing.T) {
  // set non-existant StatusCode
  transaction := FlypayB{StatusCode: 9999}
  status,err := transaction.status()
  if err == nil {
    t.Errorf("StatusCode %v is a non-existent status code and must return error", transaction.StatusCode)
  }

  transaction.StatusCode = 100
  status,err = transaction.status()
  if err != nil {
    t.Errorf("StatusCode %v has no standard status mapping", transaction.StatusCode)
  } else if status != StatusAuthorised {
    t.Errorf("StatusCode %v has wrong standard status mapping %v",
      transaction.StatusCode, status)
  }

  transaction.StatusCode = 200
  status,err = transaction.status()
  if err != nil {
    t.Errorf("StatusCode %v has no standard status mapping", transaction.StatusCode)
  } else if status != StatusDecline {
    t.Errorf("StatusCode %v has wrong standard status mapping %v",
      transaction.StatusCode, status)
  }

  transaction.StatusCode = 300
  status,err = transaction.status()
  if err != nil {
    t.Errorf("StatusCode %v has no standard status mapping", transaction.StatusCode)
  } else if status != StatusRefunded {
    t.Errorf("StatusCode %v has wrong standard status mapping %v",
      transaction.StatusCode, status)
  }
}

func TestFlypayBDecodeTransaction(t *testing.T) {
  fDecoder := FlypayBDecoder{}

  s := `
  {
    "value":200,
    "transactionCurrency":"AUD",
    "statusCode":100,
    "orderInfo":"2e58bd43-0051",
    "paymentId": "flypay-b-0001"
  }
  {
    "value":671,
    "transactionCurrency":"AUD",
    "statusCode":9999,
    "orderInfo":"2e58bd43-0052",
    "paymentId": "flypay-b-0002"
  }
  {
    "valu2000,
    "transactionCurrency":"AUD",
    "statusCode":300,
    "orderInfo":"2e58bd43-0053",
    "paymentId": "flypay-b-0003"
  }
  `
  r := strings.NewReader(s)
  d := json.NewDecoder(r)

  transaction,err := fDecoder.decodeTransaction(d)
  validT := Transaction{ProviderFlypayB,StatusAuthorised,200,"AUD","flypay-b-0001","2e58bd43-0051"}
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

func TestFlypayBDecoder_Decode(t *testing.T) {
  fDecoder := FlypayBDecoder{}

  s := `
  {
  "transactions":[
    {
      "value":200,
      "transactionCurrency":"AUD",
      "statusCode":100,
      "orderInfo":"2e58bd43-0051",
      "paymentId": "flypay-b-0001"
    }
    ]
  }
  `

    r := strings.NewReader(s)

    transactions,err := fDecoder.Decode(r)
    validT := Transaction{ProviderFlypayB,StatusAuthorised,200,"AUD","flypay-b-0001","2e58bd43-0051"}
    if err != nil {
      t.Errorf("err should be nil for successful decoding, error: %v",err)
    } else if len(transactions) != 1 || transactions[0] != validT {
      t.Error("Wrong transactions list building and decoding")
    }

    s = `
    "transactions":[
    {
      "value":200,
      "transactionCurrency":"AUD",
      "statusCode":100,
      "orderInfo":"2e58bd43-0051",
      "paymentId": "flypay-b-0001"
    }
      ]
    }
      `

    r.Reset(s)
    transactions,err = fDecoder.Decode(r)
    if err == nil {
      t.Errorf("Invalid json structure without opening {")
    }

    s = `
    `

    r.Reset(s)
    transactions,err = fDecoder.Decode(r)
    if err == nil {
      t.Errorf("Invalid json structure without opening {")
    }

    //Test second delimeter cases

    s = `
    {
    "transact":[
    {
      "value":200,
      "transactionCurrency":"AUD",
      "statusCode":100,
      "orderInfo":"2e58bd43-0051",
      "paymentId": "flypay-b-0001"
    }
      ]
    }
    `

    r.Reset(s)
    transactions,err = fDecoder.Decode(r)
    if err == nil {
      t.Errorf("Invalid json structure without transactions key")
    }

    s = `
    {
    {"transaction":[
      {
        "value":200,
        "transactionCurrency":"AUD",
        "statusCode":100,
        "orderInfo":"2e58bd43-0051",
        "paymentId": "flypay-b-0001"
      }
      ]
    }
    }
    `

    r.Reset(s)
    transactions,err = fDecoder.Decode(r)
    if err == nil {
      t.Errorf("Invalid json structure without excess delim rather than transactions key")
    }

    s = `
    {
    `

    r.Reset(s)
    transactions,err = fDecoder.Decode(r)
    if err == nil {
      t.Errorf("Invalid json structure with only {")
    }

    //Test third delim cases

    s = `
    {
    "transactions":{
      {
        "value":200,
        "transactionCurrency":"AUD",
        "statusCode":100,
        "orderInfo":"2e58bd43-0051",
        "paymentId": "flypay-b-0001"
      }
      ]
    }
      `

    r.Reset(s)
    transactions,err = fDecoder.Decode(r)
    if err == nil {
      t.Errorf("Invalid json structure without opening [")
    }

    s = `
    {
    "transactions":
    {
      "value":200,
      "transactionCurrency":"AUD",
      "statusCode":100,
      "orderInfo":"2e58bd43-0051",
      "paymentId": "flypay-b-0001"
    }
      ]
    }
      `

    r.Reset(s)
    transactions,err = fDecoder.Decode(r)
    if err == nil {
      t.Errorf("Invalid json structure without opening [")
    }

    s = `
    {
    "transactions":
    `

    r.Reset(s)
    transactions,err = fDecoder.Decode(r)
    if err == nil {
      t.Errorf("Invalid json structure without opening opening [")
    }
}
