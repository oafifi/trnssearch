package transaction

const (
  //StatusAuthorised is a constant representing transaction status "authorised"
  StatusAuthorised = "authorised"
  //StatusDecline is a constant representing transaction status "decline"
  StatusDecline = "decline"
  //StatusRefunded is a constant representing transaction status "refunded"
  StatusRefunded = "refunded"
)

//Transaction is the standard transaction structure
type Transaction struct {
  Provider string
  StatusCode string
  Amount float64
  Currency string
  TransactionID string
  OrderReference string
}

type Parameters struct {
  Provider string
  StatusCode string
  AmountMin float64
  AmountMax float64
  Currency string
}

type Result struct {
  Transactions []Transaction
  Error error
}

func FindBy(parameters Parameters) ([]Transaction,error) {

  var transactions []Transaction
  var err error

  if parameters.Provider == "" {
    transactions,err = fetchAllProvidersTransactionsConcurrent()
  } else {
    transactions,err = fetchProviderTransactions(parameters.Provider)
  }

  if err != nil {
    return transactions,err
  }

  return filter(transactions, parameters),nil
}

func isTransactionCriteriaSatisfied(t Transaction, p Parameters) bool {
  if p.StatusCode != "" && p.StatusCode != t.StatusCode {
    return false
  }

  if p.AmountMin != 0 && p.AmountMin > t.Amount {
    return false
  }
  if p.AmountMax != 0 && p.AmountMax < t.Amount {
    return false
  }

  if p.Currency != "" && p.Currency != t.Currency {
    return false
  }

  return true
}

func filter(transactionList []Transaction, p Parameters) []Transaction {

  var filteredList []Transaction

  for _, transaction := range transactionList {

    if isTransactionCriteriaSatisfied(transaction, p) {
      filteredList = append(filteredList, transaction)
    }

  }

  return filteredList
}
