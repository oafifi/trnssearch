package api

import (
  "net/http"
  "log"
)

func HandleRequests() {
    http.HandleFunc("/api/payment/transaction", searchTransactionHandler)
    log.Fatal(http.ListenAndServe(":8081", nil))
}
