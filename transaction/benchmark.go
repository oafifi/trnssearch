package transaction

import(
  "os"
  "encoding/json"
  //"strings"
  "io"
  "io/ioutil"
)

func fetchProviderTransactions_Stream(path string) []Transaction {

  //var t []Transaction
  //path, _ := providerPathList[provider]


  providerDecoder := FlypayADecoder{}

  // Open our jsonFile
  jsonFile, _ := os.Open(path)

  // defer the closing of our jsonFile so that we can parse it later on
  defer jsonFile.Close()

  t,_ := providerDecoder.Decode(jsonFile)
  return t
}

func fetchProviderTransactions_Unmarshal(path string) []Transaction {

  //var t []Transaction
  //path, _ := providerPathList[provider]


  providerDecoder := FlypayADecoder_Unmarshal{}

  // Open our jsonFile
  jsonFile, _ := os.Open(path)

  // defer the closing of our jsonFile so that we can parse it later on
  defer jsonFile.Close()

  t,_ := providerDecoder.Decode(jsonFile)
  return t
}


////////////////////////////////////////////////////////////////////////////////////////////////
//Define decoder that uses json.Unmarshal to read whole file

//FlypayADecoder1 with decode method unmarshal implementation reading all file at once
type FlypayADecoder_Unmarshal struct{}

type FlypayATransactions struct{
  Transactions []FlypayA
}

func (flypayADecoder FlypayADecoder_Unmarshal) Decode(reader io.Reader) ([]Transaction,error) {

var transactionsList []Transaction

byteValue, _ := ioutil.ReadAll(reader)

var tr FlypayATransactions

json.Unmarshal(byteValue, &tr)
//fmt.Println(tr)
for _,t := range tr.Transactions {
  tStatus,_ := t.status()

  transactionsList = append(transactionsList, Transaction{ProviderFlypayA,tStatus,
    t.Amount,t.Currency,t.TransactionID,
    t.OrderReference})

}

  return transactionsList,nil
}
