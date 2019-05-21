package api

import (
  "github.com/oafifi/trnssearch/transaction"
  "net/http"
  "strconv"
  "encoding/json"
)


func searchTransactionHandler(w http.ResponseWriter, r *http.Request){
  query := r.URL.Query()
  params := transaction.Parameters{}

  param := query.Get("provider")
  params.Provider = param

  param = query.Get("statusCode")
  params.StatusCode = param

  param = query.Get("currency")
  params.Currency = param

  amountMin := query.Get("amountMin")
  if amount, err := strconv.ParseFloat(amountMin, 64); err == nil {
    params.AmountMin = amount
  }

  amountMax := query.Get("amountMax")
  if amount, err := strconv.ParseFloat(amountMax, 64); err == nil {
    params.AmountMax = amount
  }

  transactions,err := transaction.FindBy(params)

  w.Header().Set("Content-Type", "application/json")

  if err != nil {
    json.NewEncoder(w).Encode(struct{Status bool}{Status: false})
  } else {
    json.NewEncoder(w).Encode(transactions)
  }


}
