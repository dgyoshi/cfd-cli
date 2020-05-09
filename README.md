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
```
