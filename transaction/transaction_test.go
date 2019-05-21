package transaction

import (
  "testing"
)

func TestIsTransactionCriteriaSatisfied(t *testing.T) {
  tr := Transaction{ProviderFlypayA,StatusAuthorised,1000,"AUD","flypay-a-0001","2e58bd43-0001"}
  p := Parameters{StatusCode: "does-not-exist"}

  valid := isTransactionCriteriaSatisfied(tr, p)
  if valid == true {
    t.Error("Error validating status code")
  }

  p.StatusCode = StatusAuthorised
  valid = isTransactionCriteriaSatisfied(tr, p)
  if valid == false {
    t.Error("Error validating same status code")
  }

  p = Parameters{Currency: "does-not-exist"}
  valid = isTransactionCriteriaSatisfied(tr, p)
  if valid == true {
    t.Error("Error validating currency")
  }

  p = Parameters{Currency: "AUD"}
  valid = isTransactionCriteriaSatisfied(tr, p)
  if valid == false {
    t.Error("Error validating same currency")
  }

  p = Parameters{AmountMin: 900}
  valid = isTransactionCriteriaSatisfied(tr, p)
  if valid == false {
    t.Error("Error validating amount min less than case")
  }

  p = Parameters{AmountMin: 1000}
  valid = isTransactionCriteriaSatisfied(tr, p)
  if valid == false {
    t.Error("Error validating amount min equal case")
  }

  p = Parameters{AmountMin: 1100}
  valid = isTransactionCriteriaSatisfied(tr, p)
  if valid == true {
    t.Error("Error validating amount min more than case")
  }

  p = Parameters{AmountMax: 900}
  valid = isTransactionCriteriaSatisfied(tr, p)
  if valid == true {
    t.Error("Error validating amount max less than case")
  }

  p = Parameters{AmountMax: 1000}
  valid = isTransactionCriteriaSatisfied(tr, p)
  if valid == false {
    t.Error("Error validating amount max equal case")
  }

  p = Parameters{AmountMax: 1100}
  valid = isTransactionCriteriaSatisfied(tr, p)
  if valid == false {
    t.Error("Error validating amount max more than case")
  }
}

func TestFilter(t *testing.T) {

  t1 := Transaction{ProviderFlypayA,StatusAuthorised,1000,"AUD","flypay-a-0001","2e58bd43-0001"}
  t2 := Transaction{ProviderFlypayA,StatusAuthorised,500,"AUD","flypay-a-0001","2e58bd43-0001"}
  t3 := Transaction{ProviderFlypayA,StatusAuthorised,1000,"EGP","flypay-a-0001","2e58bd43-0001"}
  tList := []Transaction{t1,t2,t3}

  p := Parameters{AmountMin: 900, AmountMax: 1100, Currency: "AUD"}

  filteredList := filter(tList, p)

  l := len(filteredList)

  if l != 1 {
    t.Errorf("Expected a list with 1 transaction, got different value %v",l)
  } else if l == 1 && filteredList[0] != t1 {
    t.Error("Expected a list with 1 transaction t1, got different value")
  }
}
