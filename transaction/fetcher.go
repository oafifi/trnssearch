package transaction

import (
  "os"
  "log"
  "errors"
)

func fetchProviderTransactions(provider string) ([]Transaction,error) {

  var t []Transaction

  path, exist := providerPathList[provider]
  if exist == false {
    log.Printf("Error: %v provider data file path is missing from providerPathList",provider)
    return t,errors.New("Provider data file path is missing")
  }

  providerDecoder, err := CreateProviderDecoder(provider)
  if err != nil {
    log.Printf("Error creating %v provider decoder. | Error: %v",provider, err)
    return t,err
  }

  // Open our jsonFile
  jsonFile, err := os.Open(path)
  if err != nil {
    log.Printf("Error opening %v provider data file. Path: %v | Error: %v",provider, path, err)
    return t,err
  }

  // defer the closing of our jsonFile so that we can parse it later on
  defer jsonFile.Close()

  t,err = providerDecoder.Decode(jsonFile)
  if err != nil {
    log.Printf("Error decoding %v provider data file. Path: %v | Error: %v",provider, path, err)
    return t,err
  }

  return t,nil
}

func fetchAllProvidersTransactions() ([]Transaction,error) {
  var transactions []Transaction

  for provider := range providerPathList {
    t,err := fetchProviderTransactions(provider)
    if err != nil {
      return t,err
    }
    transactions = append(transactions, t...)
  }

  return transactions,nil
}

func fetchProviderTransactionsConcurrent(provider string, c chan Result) {

  t,err := fetchProviderTransactions(provider)

  c <- Result{t,err}

  close(c)
}

func fetchAllProvidersTransactionsConcurrent() ([]Transaction,error) {
  var transactions []Transaction

  //This should be implemented in a dynamic way using providers list and single channel

  c1 := make(chan Result)
  c2 := make(chan Result)

  go fetchProviderTransactionsConcurrent(ProviderFlypayA, c1)
  go fetchProviderTransactionsConcurrent(ProviderFlypayB, c2)

  tA := <- c1
  tB := <- c2
  if tA.Error != nil || tB.Error != nil {
    return transactions, errors.New("Error retreiving data")
  }
  transactions = append(transactions, tA.Transactions...)
  transactions = append(transactions, tB.Transactions...)

  return transactions,nil
}
