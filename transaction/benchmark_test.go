package transaction

import(
  "testing"
)

var testFilePath = "testcases/flypayA.json"

//
func BenchmarkFlypayADecoder_FetchProviderTransactions_Stream(b *testing.B) {

  for i := 0; i < b.N; i++ {
    fetchProviderTransactions_Stream(testFilePath)
  }

}

func BenchmarkFlypayADecoder_FetchProviderTransactions_Unmarshal(b *testing.B) {

  for i := 0; i < b.N; i++ {
    fetchProviderTransactions_Unmarshal(testFilePath)
  }

}

func BenchmarkFetchAllProvidersTransactions(b *testing.B) {
  //override global variable
  providerPathList = map[string]string {
    ProviderFlypayA: "testcases/flypayA.json",
    ProviderFlypayB: "testcases/flypayB.json",
  }

  for i := 0; i < b.N; i++ {
    fetchAllProvidersTransactions()
  }

}

func BenchmarkFetchAllProvidersTransactionsConcurrent(b *testing.B) {
  //override global variable
  providerPathList = map[string]string {
    ProviderFlypayA: "testcases/flypayA.json",
    ProviderFlypayB: "testcases/flypayB.json",
  }

  for i := 0; i < b.N; i++ {
    fetchAllProvidersTransactionsConcurrent()
  }

}
