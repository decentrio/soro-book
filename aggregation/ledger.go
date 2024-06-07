package aggregation

import (
	"fmt"
	"io"
	"time"

	"github.com/decentrio/soro-book/database/models"
	"github.com/stellar/go/ingest"
	backends "github.com/stellar/go/ingest/ledgerbackend"
	"github.com/stellar/go/xdr"
)

var (
	MultiHopContract    = "CCLZRD4E72T7JCZCN3P7KNPYNXFYKQCL64ECLX7WP5GNVYPYJGU2IO2G"
	PHOUSDCPoolContract = "CAZ6W4WHVGQBGURYTUOLCUOOHW6VQGAAPSPCD72VEDZMBBPY7H43AYEC"
	//PHO key 1
	PHOTokenKey      = uint32(1)
	PHOTokenContract = "CBZ7M5B3Y4WWBZ5XK5UZCAFOEZ23KSSZXYECYX3IXM6E2JOLQC52DK32"
	//USDC key 2
	USDCTokenKey      = uint32(2)
	USDCTokenContract = "CCW67TSZV3SSS2HXMBQ5JFGCKJNXKZM7UQUWUZPUTHXSTZLEO7SJMI75"
	// token power
	TokenPowerReduction = 10000000
)

type LedgerWrapper struct {
	ledger models.Ledger
	txs    []TransactionWrapper
}

func (as *Aggregation) getNewLedger() {
	// prepare range
	fmt.Println("prepare")
	from, to := as.prepare()
	fmt.Println("preprare done")
	// get ledger
	if !as.isSync {
		for seq := from; seq < to; seq++ {
			ledgerCloseMeta, err := as.backend.GetLedger(as.ctx, seq)
			if err != nil {
				as.Logger.Error(fmt.Sprintf("error get ledger %s", err.Error()))
				return
			}

			go func(l xdr.LedgerCloseMeta) {
				as.ledgerQueue <- l
			}(ledgerCloseMeta)
		}
	} else {
		seq := as.StartLedgerSeq
		ledgerCloseMeta, err := as.backend.GetLedger(as.ctx, seq)
		if err != nil {
			as.Logger.Error(fmt.Sprintf("error get ledger %s", err.Error()))
			return
		}

		go func(l xdr.LedgerCloseMeta) {
			as.ledgerQueue <- l
		}(ledgerCloseMeta)
		as.StartLedgerSeq++
	}
}

// aggregation process
func (as *Aggregation) ledgerProcessing() {
	for {
		select {
		// Receive a new tx
		case ledger := <-as.ledgerQueue:
			as.handleReceiveNewLedger(ledger)
		// Terminate process
		case <-as.BaseService.Terminate():
			return
		default:
		}
		time.Sleep(time.Millisecond)
	}
}

// handleReceiveTx
func (as *Aggregation) handleReceiveNewLedger(l xdr.LedgerCloseMeta) {
	ledger := getLedgerFromCloseMeta(l)
	// get tx
	if l.LedgerSequence() != 51994989 {
		return
	}
	txReader, err := ingest.NewLedgerTransactionReaderFromLedgerCloseMeta(as.Cfg.NetworkPassphrase, l)
	panicIf(err)
	defer txReader.Close()

	// Read each transaction within the ledger, extract its operations, and
	// accumulate the statistics we're interested in.
	var txs []TransactionWrapper
	for {
		tx, err := txReader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			as.Logger.Error(fmt.Sprintf("error txReader %s", err.Error()))
		}

		txWrapper := NewTransactionWrapper(tx, l.LedgerSequence(), ledger.LedgerTime)
		txs = append(txs, txWrapper)
	}

	// extract tx to get information
	for _, tx := range txs {
		invokes, _, err := isInvokeHostFunctionTx(tx.Tx, l.LedgerSequence())
		if err != nil {
			continue
		}

		if len(invokes) == 0 {
			continue
		}

		var historicalTrade HistoricalTrades
		if invokes[0].ContractId == MultiHopContract {
			var argsXdr xdr.InvokeContractArgs
			argsXdr.UnmarshalBinary(invokes[0].Args)

			// args layout
			// sender
			// swap(Vec(Map))
			// max_spread_bps(option)(i64)
			// amount(i128)
			// we only care about arg1 and arg3

			// extract arg1
			arg1 := argsXdr.Args[1].MustVec()
			var askAsset string
			var offerAsset string
			for _, av := range *arg1 {
				am := av.MustMap()
				for _, e := range *am {
					key := string(e.Key.MustSym())
					switch key {
					case "ask_asset":
						val := e.Val.MustAddress()
						askAsset, _ = val.String()

						fmt.Printf("ask_asset: %s\n", askAsset)
					case "ask_asset_min_amount":
						// TODO

					case "offer_asset":
						val := e.Val.MustAddress()
						offerAsset, _ = val.String()

						fmt.Printf("offer_asset: %s\n", offerAsset)
					default:
						continue
					}
				}
			}

			if askAsset == PHOTokenContract && offerAsset == USDCTokenContract {
				historicalTrade.TradeType = "buy"
			} else if askAsset == USDCTokenContract && offerAsset == PHOTokenContract {
				historicalTrade.TradeType = "sell"
			} else {
				as.Logger.Errorf("unknown trading pair %s - %s", askAsset, offerAsset)
				continue
			}

			//extract arg3
			arg3 := argsXdr.Args[3].MustI128()
			amount := uint64(arg3.Lo)
			fmt.Println("amount ", amount)

			// retrive contract pool liquidity data
			var usdcLiquidity uint64
			var phoLiquidity uint64
			contractDatas := tx.GetModelsContractDataEntry()
			for _, cd := range contractDatas {
				if cd.ContractId == PHOUSDCPoolContract {
					var keyXdr xdr.ScVal
					keyXdr.UnmarshalBinary(cd.KeyXdr)
					key := uint32(keyXdr.MustU32())

					var valXdr xdr.ScVal
					valXdr.UnmarshalBinary(cd.ValueXdr)
					val := valXdr.MustI128()

					switch key {
					case PHOTokenKey:
						phoLiquidity = uint64(val.Lo)
					case USDCTokenKey:
						usdcLiquidity = uint64(val.Lo)
					default:
					}

				}
			}
			fmt.Println("Pho: ", phoLiquidity)
			fmt.Println("usdc: ", usdcLiquidity)
		}

	}
	fmt.Println("===================================")
}

func (as *Aggregation) prepare() (uint32, uint32) {
	if !as.isSync {
		from := as.StartLedgerSeq
		to := from + DefaultPrepareStep

		var ledgerRange backends.Range
		if to > as.CurrLedgerSeq {
			ledgerRange = backends.UnboundedRange(from)
		} else {
			ledgerRange = backends.BoundedRange(from, to)
		}

		fmt.Println(ledgerRange)
		err := as.backend.PrepareRange(as.ctx, ledgerRange)
		if err != nil {
			as.Logger.Errorf("error prepare %s", err.Error())
			return 0, 0 // if prepare error, we should skip here
		} else {
			if to > as.CurrLedgerSeq {
				as.isSync = true
			}
		}
		as.StartLedgerSeq += DefaultPrepareStep
		return from, to
	}

	return 0, 0
}

func getLedgerFromCloseMeta(ledgerCloseMeta xdr.LedgerCloseMeta) models.Ledger {
	var ledgerHeader xdr.LedgerHeaderHistoryEntry
	switch ledgerCloseMeta.V {
	case 0:
		ledgerHeader = ledgerCloseMeta.MustV0().LedgerHeader
	case 1:
		ledgerHeader = ledgerCloseMeta.MustV1().LedgerHeader
	default:
		panic(fmt.Sprintf("Unsupported LedgerCloseMeta.V: %d", ledgerCloseMeta.V))
	}

	timeStamp := uint64(ledgerHeader.Header.ScpValue.CloseTime)

	return models.Ledger{
		Hash:       ledgerCloseMeta.LedgerHash().HexString(),
		PrevHash:   ledgerCloseMeta.PreviousLedgerHash().HexString(),
		Seq:        ledgerCloseMeta.LedgerSequence(),
		LedgerTime: timeStamp,
	}
}
