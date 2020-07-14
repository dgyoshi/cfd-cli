## command

### getpubkeyfromprivkey
```
go run ./ getpubkeyfromprivkey -privkey <privatekey> -comp
go run ./ getpubkeyfromprivkey -wif <private wif>
```

### genprivkeyfromstrings
```
go run ./ genprivkeyfromstrings -text <message>
go run ./ genprivkeyfromstrings -text "<message1|message2|message3>"
```

### getextkeypairfromseed
```
go run ./ getextkeypairfromseed -seed <seed> -network <network>
go run ./ getextkeypairfromseed -seed <seed> -network <network> -path <bip32path>
```

### getextkeypairfrommnemonic
```
go run ./ getextkeypairfrommnemonic -mnemonic <mnemonic> -network <network>
go run ./ getextkeypairfrommnemonic -mnemonic <mnemonic> -network <network> -path <bip32paths>
```

### decoderawtransaction
```
go run ./ decoderawtransaction -tx <tx> -network <network>
go run ./ decoderawtransaction -tx <tx> -network <network> -elements
go run ./ decoderawtransaction -file <filename> -network <network> -elements
```

### encodedersignature
```
go run ./ encodedersignature -signature <signature> -sighashtype <sighashtype>
go run ./ encodedersignature -signature <signature> -sighashtype <sighashtype> -anyonecanpay
```

### verifysigntransaction
```
go run ./ verifysigntransaction -tx <tx> -txid <txid> -vout <vout> -address <address> -addresstype <addressType> -amount <amount>
go run ./ verifysigntransaction -tx <tx> -elements -txid <txid> -vout <vout> -address <address> -addresstype <addressType> -amount <amount>
go run ./ verifysigntransaction -file <filename> -elements -txid <txid> -vout <vout> -address <address> -addresstype <addressType> -commitment <amountCommitment>
go run ./ verifysigntransaction -file <filename> -elements -txid <txid> -vout <vout> -descriptor <descriptor> -commitment <amountCommitment>
```

### verifysignature
```
go run ./ verifysignature -tx <tx> -txid <txid> -vout <vout> -signature <signature> -pubkey <pubkey> -addresstype <addressType> -sighashtype <sighashtype> -amount <amount>
go run ./ verifysignature -tx <tx> -elements -txid <txid> -vout <vout> -signature <signature> -script <redeemScript> -addresstype <addressType> -sighashtype <sighashtype> -anyonecanpay -amount <amount>
go run ./ verifysignature -file <filename> -elements -txid <txid> -vout <vout> -signature <signature> -descriptor <descriptor> -sighashtype <sighashtype> -commitment <amountCommitment>
```

### initializetransaction
```
go run ./ initializetransaction -tx <tx> -version <version> -locktime <locktime>
go run ./ initializetransaction -file <filename> -elements -version <version> -locktime <locktime>
```

### appendtxin
```
go run ./ appendtxin -tx <tx> -txid <txid> -vout <vout> -sequence <sequence>
go run ./ appendtxin -file <filename> -elements -txid <txid> -vout <vout> -sequence <sequence>
(save utxo data)
go run ./ appendtxin -file <filename> -elements -txid <txid> -vout <vout> -sequence <sequence> -amount <amount> -asset <asset> -assetblinder <assetblinder> -assetcommitment <assetcommitment> -blinder <blinder> -amountcommitment <amountcommitment> -descriptor <descriptor>
```

### appendtxout
```
go run ./ appendtxout -tx <tx> -amount <amount> -address <address>
go run ./ appendtxout -file <filename> -elements -amount <amount> -address <address> -asset <asset>
go run ./ appendtxout -file <filename> -elements -amount <amount> -asset <asset> -fee
go run ./ appendtxout -file <filename> -elements -amount <amount> -asset <asset> -destroy
go run ./ appendtxout -file <filename> -elements -amount <amount> -lockingscript <lockingScript> -asset <asset>
```

### setrawreissueasset
(utxo setting is call appendtxin)
```
go run ./ setrawreissueasset -tx <tx> -txid <txid> -vout <vout> -amount <amount> -entropy <entropy> -address <address>
go run ./ setrawreissueasset -file <filename> -txid <txid> -vout <vout> -amount <amount> -entropy <entropy> -assetblinder <assetBlinder> -lockingscript <lockingscript>
```

### estimatefee
```
go run ./ estimatefee -tx <tx> -feerate <feerate>
go run ./ estimatefee -file <filename> -elements -feerate <feerate> -asset <asset>
```

### blindrawtransaction
(utxo setting is call appendtxin)
```
go run ./ blindrawtransaction -tx <tx>
go run ./ blindrawtransaction -file <filename> -blindingkeys <issuanceKey1|issuanceKey2|...>
```

### createsignaturehash
```
go run ./ createsignaturehash -tx <tx> -elements -txid <txid> -vout <vout> -sighashtype <sighashtype> -anyonecanpay -amount <amount>
go run ./ createsignaturehash -file <filename> -elements -txid <txid> -vout <vout> -sighashtype <sighashtype> -amountcommitment <amountcommitment>
go run ./ createsignaturehash -file <filename> -elements -txid <txid> -vout <vout> -pubkey <pubkey> -sighashtype <sighashtype> -amountcommitment <amountcommitment>
go run ./ createsignaturehash -file <filename> -elements -txid <txid> -vout <vout> -script <redeemScript> -sighashtype <sighashtype> -amountcommitment <amountcommitment>
```

### getsignature
```
go run ./ getsignature -sighash <sighash> -privkey <privkey> -grindr
go run ./ getsignature -sighash <sighash> -extpriv <extpriv> -bip32path <bip32path>
```

### addsigntransaction
```
go run ./ addsigntransaction -tx <tx> -elements -txid <txid> -vout <vout> -signature <signature> -pubkey <pubkey> -addresstype <addresstype> -sighashtype <sighashtype> -anyonecanpay
go run ./ addsigntransaction -file <filename> -elements -txid <txid> -vout <vout> -signature <signature1|signature2|...> -pubkey <pubkey1|pubkey2|...> -script <redeemScript> -addresstype <addresstype> -sighashtype <sighashtype>
go run ./ addsigntransaction -file <filename> -elements -txid <txid> -vout <vout> -signature <signature1|signature2|...> -pubkey <pubkey1|pubkey2|...> -sighashtype <sighashtype>
```

### signwithprivkey
```
go run ./ signwithprivkey -tx <tx> -elements -txid <txid> -vout <vout> -privkey <privkey> -grindr -addresstype <addresstype> -sighashtype <sighashtype> -anyonecanpay
go run ./ signwithprivkey -file <filename> -elements -txid <txid> -vout <vout> -extpriv <extpriv> -bip32path <bip32path> -sighashtype <sighashtype> -anyonecanpay
```

### getcommitment
```
go run ./ getcommitment -asset <asset> -amount <amount> -assetblinder <assetBlinder> -blinder <blinder>
```

### parsedescriptor
```
go run ./ parsedescriptor -network <network> -childnum <childnumber> -descriptor <outputDescriptor>
```
