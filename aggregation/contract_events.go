package aggregation

import (
	"fmt"
	"time"

	"github.com/decentrio/soro-book/database/models"
	"github.com/stellar/go/strkey"
	"github.com/stellar/go/xdr"
)

const (
	// Implemented
	EventTypeTransfer = "transfer"
	EventTypeMint     = "mint"
	EventTypeClawback = "clawback"
	EventTypeBurn     = "burn"
	// TODO: Not implemented
	EventTypeIncrAllow
	EventTypeDecrAllow
	EventTypeSetAuthorized
	EventTypeSetAdmin
)

var (
	STELLAR_ASSET_CONTRACT_TOPICS = map[xdr.ScSymbol]string{
		xdr.ScSymbol("transfer"): EventTypeTransfer,
		xdr.ScSymbol("mint"):     EventTypeMint,
		xdr.ScSymbol("clawback"): EventTypeClawback,
		xdr.ScSymbol("burn"):     EventTypeBurn,
	}
)

// aggregation process
func (as *Aggregation) contractEventsProcessing() {
	for {
		// Block until state have sync successful
		if as.isReSync {
			continue
		}

		select {
		// Receive a new tx
		case event := <-as.assetContractEventsQueue:
			eventType := event.GetType()
			switch eventType {
			case EventTypeTransfer:
				// Create AssetContractTransferEvent
				transferEvent := event.(*models.AssetContractTransferEvent)
				_, err := as.db.CreateAssetContractTransferEvent(transferEvent)
				if err != nil {
					as.Logger.Error(fmt.Sprintf("Error create asset contract transfer event tx %s: %s", transferEvent.TxHash, err.Error()))
				}
			case EventTypeMint:
				// Create AssetContractTransferEvent
				mintEvent := event.(*models.AssetContractMintEvent)
				_, err := as.db.CreateAssetContractMintEvent(mintEvent)
				if err != nil {
					as.Logger.Error(fmt.Sprintf("Error create asset contract mint event tx %s: %s", mintEvent.TxHash, err.Error()))
				}
			case EventTypeClawback:
				// Create AssetContractTransferEvent
				cbEvent := event.(*models.AssetContractClawbackEvent)
				_, err := as.db.CreateAssetContractClawbackEvent(cbEvent)
				if err != nil {
					as.Logger.Error(fmt.Sprintf("Error create asset contract clawback event tx %s: %s", cbEvent.TxHash, err.Error()))
				}
			case EventTypeBurn:
				// Create AssetContractTransferEvent
				burnEvent := event.(*models.AssetContractBurnEvent)
				_, err := as.db.CreateAssetContractBurnEvent(burnEvent)
				if err != nil {
					as.Logger.Error(fmt.Sprintf("Error create asset contract burn event tx %s: %s", burnEvent.TxHash, err.Error()))
				}
			}
			// as.handleReceiveNewLedger(ledger)
		case event := <-as.wasmContractEventsQueue:
			// Create WasmContractEvents
			_, err := as.db.CreateWasmContractEvent(&event)
			if err != nil {
				as.Logger.Error(fmt.Sprintf("Error create wasm contract event tx %s: %s", event.TxHash, err.Error()))
			}
		// Terminate process
		case <-as.BaseService.Terminate():
			return
		}
		time.Sleep(time.Millisecond)
	}
}

func (tx TransactionWrapper) GetContractEvents() ([]models.WasmContractEvent, []models.StellarAssetContractEvent, error) {
	var wasmContractevents []models.WasmContractEvent
	var assetContractEvents []models.StellarAssetContractEvent
	for _, op := range tx.Ops {
		var order = uint32(1)
		if op.OperationType() == xdr.OperationTypeInvokeHostFunction {
			diagnosticEvents, innerErr := tx.Tx.GetDiagnosticEvents()
			if innerErr != nil {
				return nil, nil, innerErr
			}
			evts := filterEvents(diagnosticEvents)

			for _, evt := range evts {
				isAssetEvent := isStellarAssetContractEvent(evt)
				if !isAssetEvent {
					wasmEvent, err := tx.GetWasmContractEvents(evt, op.ID(), &order)
					if err != nil {
						continue
					}

					wasmContractevents = append(wasmContractevents, wasmEvent)
				} else {
					assetEvent, err := tx.GetStellarAssetContractEvents(evt, op.ID(), &order)
					if err != nil {
						continue
					}

					assetContractEvents = append(assetContractEvents, assetEvent)
				}

			}
		}
	}

	return wasmContractevents, assetContractEvents, nil
}

func (tx TransactionWrapper) GetWasmContractEvents(event xdr.ContractEvent, id int64, order *uint32) (models.WasmContractEvent, error) {
	eventBodyXdr, err := event.Body.MarshalBinary()
	if err != nil {
		return models.WasmContractEvent{}, err
	}

	contractId, err := strkey.Encode(strkey.VersionByteContract, event.ContractId[:])
	if err != nil {
		return models.WasmContractEvent{}, err
	}

	evt := models.WasmContractEvent{
		Id:           fmt.Sprintf("%019d-%010d", id, *order), // ID should be combine from operation ID and event index
		ContractId:   contractId,
		TxHash:       tx.Tx.Result.TransactionHash.HexString(),
		EventBodyXdr: eventBodyXdr,
	}
	*order++

	return evt, nil
}

func (tx TransactionWrapper) GetStellarAssetContractEvents(event xdr.ContractEvent, id int64, order *uint32) (models.StellarAssetContractEvent, error) {
	topics := event.Body.V0.Topics
	value := event.Body.V0.Data

	// Get event type
	fn, _ := topics[0].GetSym()
	eventType := STELLAR_ASSET_CONTRACT_TOPICS[fn]

	// Get event Id
	eventId := fmt.Sprintf("%019d-%010d", id, *order)

	// get contract Id
	contractId, err := strkey.Encode(strkey.VersionByteContract, event.ContractId[:])
	if err != nil {
		return nil, err
	}

	// Get Tx Hash
	txHash := tx.GetTransactionHash()

	// Get event data
	switch eventType {
	case EventTypeTransfer:
		transferEvent := models.AssetContractTransferEvent{
			Id:         eventId,
			ContractId: contractId,
			TxHash:     txHash,
		}
		err := transferEvent.Parse(topics, value)
		if err != nil {
			return nil, err
		}
		*order++

		return &transferEvent, nil
	case EventTypeMint:
		mintEvent := models.AssetContractMintEvent{
			Id:         eventId,
			ContractId: contractId,
			TxHash:     txHash,
		}
		err := mintEvent.Parse(topics, value)
		if err != nil {
			return nil, err
		}
		*order++

		return &mintEvent, nil
	case EventTypeClawback:
		cbEvent := models.AssetContractClawbackEvent{
			Id:         eventId,
			ContractId: contractId,
			TxHash:     txHash,
		}
		err := cbEvent.Parse(topics, value)
		if err != nil {
			return nil, err
		}
		*order++

		return &cbEvent, nil
	case EventTypeBurn:
		burnEvent := models.AssetContractBurnEvent{
			Id:         eventId,
			ContractId: contractId,
			TxHash:     txHash,
		}
		err := burnEvent.Parse(topics, value)
		if err != nil {
			return nil, err
		}
		*order++

		return &burnEvent, nil
	default:
		return nil, fmt.Errorf("event type ('%s') unsupported", eventType)
	}
}

func isStellarAssetContractEvent(event xdr.ContractEvent) bool {
	if event.Type != xdr.ContractEventTypeContract || event.ContractId == nil || event.Body.V != 0 {
		return false
	}

	topics := event.Body.V0.Topics

	// No relevant SAC events have <= 2 topics
	if len(topics) <= 2 {
		return false
	}

	fn, ok := topics[0].GetSym()
	if !ok {
		return false
	}

	if _, found := STELLAR_ASSET_CONTRACT_TOPICS[fn]; !found {
		return false
	}

	return true
}

func filterEvents(diagnosticEvents []xdr.DiagnosticEvent) []xdr.ContractEvent {
	var filtered []xdr.ContractEvent
	for _, diagnosticEvent := range diagnosticEvents {
		if !diagnosticEvent.InSuccessfulContractCall || diagnosticEvent.Event.Type != xdr.ContractEventTypeContract {
			continue
		}
		filtered = append(filtered, diagnosticEvent.Event)
	}
	return filtered
}
