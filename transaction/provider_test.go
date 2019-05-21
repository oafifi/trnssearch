package transaction

import (
  "testing"
)

func TestCreateProviderDecoder(t *testing.T) {
  decoder,err := CreateProviderDecoder(ProviderFlypayA)
  if err != nil {
    t.Error("FlypayA decoder is missing")
  }
  _,valid := decoder.(FlypayADecoder)
  if valid == false {
    t.Errorf("Decoder expected was FlypayADEcoder, returned %T",decoder)
  }

  decoder,err = CreateProviderDecoder(ProviderFlypayB)
  if err != nil {
    t.Error("FlypayB decoder is missing")
  }
  _,valid = decoder.(FlypayBDecoder)
  if valid == false {
    t.Errorf("Decoder expected was FlypayBDecoder, returned %T",decoder)
  }

  _,err = CreateProviderDecoder("Non-existant provider")
  if err == nil {
    t.Error("Should return error on trying to create non-existant decoder")
  }
}
