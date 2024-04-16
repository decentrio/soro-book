# Sorobook
## Aggregation

The `aggregation` package takes responsibility for collecting all the `stellar` on-chain data and indexing corresponding to storage.

The package itself include 5 part:
- Aggregation Process
- Ledger Process
- Transaction Process
- Contract Data Process
- Contract Event Process

### Aggregation Process
The `Aggregation Proccess` use [go/ingest](https://github.com/stellar/go/tree/master/ingest) package for collecting the on-chain`ledger` information. The aggregation is defined in [getNewLedger()](https://github.com/decentrio/soro-book/blob/fad4719f4a7fd0cc8b0ce342b5faac9f6d2ad7ad/aggregation/ledger.go#L14) and push to `ledgerQueue` which would be used by `Ledger Process`

### Ledger Proccess
```go=
type LedgerWrapper struct {
    ledger models.Ledger
    txs    []TransactionWrapper
}
```
The `Ledger Process` is defined in [ledgerProcessing()](https://github.com/decentrio/soro-book/blob/fad4719f4a7fd0cc8b0ce342b5faac9f6d2ad7ad/aggregation/ledger.go#L79). In this process,`ledger`  will be stored into database and push all transaction information into `txQueue` which would be used by `Transaction Process`

### Transaction Process
```go=
type TransactionWrapper struct {
	LedgerSequence uint32
	Tx             ingest.LedgerTransaction
	Ops            []transactionOperationWrapper
}
```
The `Transaction Process` is defined in [transactionProcessing()](https://github.com/decentrio/soro-book/blob/fad4719f4a7fd0cc8b0ce342b5faac9f6d2ad7ad/aggregation/transactions.go#L17C24-L17C45). Whenever `txQueue` receive a new element:
- First, it will [store new transaction into database](https://github.com/decentrio/soro-book/blob/fad4719f4a7fd0cc8b0ce342b5faac9f6d2ad7ad/aggregation/transactions.go#L37-L41). 
- Then, the process will extract transaction information to get [Contract Data](https://github.com/decentrio/soro-book/blob/fad4719f4a7fd0cc8b0ce342b5faac9f6d2ad7ad/aggregation/contract_data.go#L39) and [Contract Events](https://github.com/decentrio/soro-book/blob/fad4719f4a7fd0cc8b0ce342b5faac9f6d2ad7ad/aggregation/contract_events.go#L91). The `Contract Data` and `Contract Events` that extracted from `transaction` will be pushed to corresponding channel `contractDataEntrysQueue` or `assetContractEventsQueue`, `wasmContractEventsQueue`

### Contract Data Process
The `Contract Data Process` is defined in [contractDataEntryProcessing()](https://github.com/decentrio/soro-book/blob/fad4719f4a7fd0cc8b0ce342b5faac9f6d2ad7ad/aggregation/contract_data.go#L13). In this process, `contract data` will be stored into database

### Contract Events Process
The `Contract Events Proccess` is defined in [contractEventsProcessing()](https://github.com/decentrio/soro-book/blob/fad4719f4a7fd0cc8b0ce342b5faac9f6d2ad7ad/aggregation/contract_events.go#L35). There are 2 type of contract events that will corresponding with 2 channel `assetContractEventsQueue` and `wasmContractEventsQueue`. In this process, we check the specific type of event and store it into database
