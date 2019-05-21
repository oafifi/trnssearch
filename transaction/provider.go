package transaction

import (
  "io"
  "errors"
)

//constants containing providers names
const (
  //ProviderFlypayA is const defining FlypayA transaction provider name
  ProviderFlypayA = "FlypayA"
  //ProviderFlypayB is const defining FlypayB transaction provider name
  ProviderFlypayB = "FlypayB"
)

//A simple holder for file paths, better to create config
var providerPathList = map[string]string {
  ProviderFlypayA: "transaction/testcases/flypayA.json",
  ProviderFlypayB: "transaction/testcases/flypayB.json",
}

//ProviderDecoder is an interface that has method to decode provider input stream
//into a list of standard transaction format
//Each provider must create his own ProviderDecoder by implementing this interface
//and add his provider to create provider
type ProviderDecoder interface {
  Decode(io.Reader) ([]Transaction,error)
}

//CreateProviderDecoder Each provider must add its ProviderDecoder to this simple factory function
//A better implementation is ProviderDecoder factory interface and it will ease testing fetching functions as well
func CreateProviderDecoder(provider string) (ProviderDecoder,error) {
  var decoder ProviderDecoder

  switch provider {
  case ProviderFlypayA:
    return FlypayADecoder{},nil
  case ProviderFlypayB:
    return FlypayBDecoder{},nil
  }

  return decoder,errors.New("Provider decoder not found")
}
