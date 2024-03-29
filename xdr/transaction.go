package xdr

import (
	"github.com/stellar/go/support/errors"
	"github.com/stellar/go/xdr"
)

// func ConvertTransactionJSON(tx models.Transaction) models.TransactionJSON {
// 	return models.TransactionJSON{}
// }

func ConvertTransactionResultPair(r xdr.TransactionResultPair) (TransactionResultPair, error) {

}

func ConvertTransactionResult(r xdr.TransactionResult) (TransactionResult, error) {

}

func ConvertTransactionResultResult(r xdr.TransactionResultResult) (TransactionResultResult, error) {

}

func ConvertInnerTransactionResultPair(r xdr.InnerTransactionResultPair) (InnerTransactionResultPair, error) {

}

func ConvertInnerTransactionResult(r xdr.InnerTransactionResult) (InnerTransactionResult, error) {

}

func ConvertInnerTransactionResultResult(r xdr.InnerTransactionResultResult) (InnerTransactionResultResult, error) {

}

func ConvertInnerTransactionResultExt(e xdr.InnerTransactionResultExt) InnerTransactionResultExt {
	return InnerTransactionResultExt{V: e.V}
}

func ConvertTransactionResultExt(e xdr.TransactionResultExt) TransactionResultExt {
	return TransactionResultExt{V: e.V}
}

func ConvertFeeBumpTransaction(tx xdr.FeeBumpTransaction) (FeeBumpTransaction, error) {
	var result FeeBumpTransaction

	feeSource, err := ConvertMuxedAccount(tx.FeeSource)
	if err != nil {
		return result, err
	}

	innerTx, err := ConvertFeeBumpTransactionInnerTx(tx.InnerTx)
	if err != nil {
		return result, err
	}

	ext := ConvertFeeBumpTransactionExt(tx.Ext)

	result.FeeSource = feeSource
	result.Fee = int64(tx.Fee)
	result.InnerTx = innerTx
	result.Ext = ext

	return result, nil
}

func ConvertTransaction(tx xdr.Transaction) (Transaction, error) {
	var result Transaction

	sourceAccount, err := ConvertMuxedAccount(tx.SourceAccount)
	if err != nil {
		return result, err
	}

	cond, err := ConvertPreconditions(tx.Cond)
	if err != nil {
		return result, err
	}

	memo, err := ConvertMemo(tx.Memo)
	if err != nil {
		return result, err
	}

	var ops []Operation
	for _, xdrOp := range tx.Operations {
		op, err := ConvertOperation(xdrOp)
		if err != nil {
			return result, err
		}
		ops = append(ops, op)
	}

	ext, err := ConvertTxExt(tx.Ext)
	if err != nil {
		return result, err
	}

	result.SourceAccount = sourceAccount
	result.Fee = uint32(tx.Fee)
	result.SeqNum = int64(tx.SeqNum)
	result.Cond = cond
	result.Memo = memo
	result.Operations = ops
	result.Ext = ext

	return result, nil

}

// TODO: testing
func ConvertTransactionV0(tx xdr.TransactionV0) (TransactionV0, error) {
	var txV0 TransactionV0

	txV0.SourceAccountEd25519 = tx.SourceAccountEd25519.String()
	txV0.Fee = uint32(tx.Fee)
	txV0.SeqNum = int64(tx.SeqNum)

	tb, err := ConvertTimeBounds(tx.TimeBounds)
	if err != nil {
		return txV0, err
	}
	txV0.TimeBounds = tb

	memo, err := ConvertMemo(tx.Memo)
	if err != nil {
		return txV0, err
	}
	txV0.Memo = memo

	var ops []Operation
	for _, opXdr := range tx.Operations {
		op, err := ConvertOperation(opXdr)
		if err != nil {
			return txV0, err
		}

		ops = append(ops, op)
	}
	txV0.Operations = ops

	ext, err := ConvertTxV0Ext(tx.Ext)
	if err != nil {
		return txV0, err
	}
	txV0.Ext = ext

	return txV0, nil
}

// TODO: testing
func ConvertTimeBounds(tb *xdr.TimeBounds) (*TimeBounds, error) {
	return &TimeBounds{
		MinTime: uint64(tb.MinTime),
		MaxTime: uint64(tb.MaxTime),
	}, nil
}

// TODO: testing
func ConvertMemo(memo xdr.Memo) (Memo, error) {
	var result Memo

	switch memo.Type {
	case xdr.MemoTypeMemoNone:
		return result, nil
	case xdr.MemoTypeMemoText:
		text, ok := memo.GetText()
		if !ok {
			return result, errors.Errorf("error invalid memo type text %v", memo)
		}
		result.Text = &text

		return result, nil
	case xdr.MemoTypeMemoId:
		xdrId, ok := memo.GetId()
		if !ok {
			return result, errors.Errorf("error invalid memo type id %v", memo)
		}
		id := uint64(xdrId)
		result.Id = &id

		return result, nil
	case xdr.MemoTypeMemoHash:
		xdrHash, ok := memo.GetHash()
		if !ok {
			return result, errors.Errorf("error invalid memo type hash %v", memo)
		}
		hash := xdrHash.HexString()
		result.Hash = &hash

		return result, nil
	case xdr.MemoTypeMemoReturn:
		xdrRetHash, ok := memo.GetRetHash()
		if !ok {
			return result, errors.Errorf("error invalid memo type rethash%v", memo)
		}
		retHash := xdrRetHash.HexString()
		result.RetHash = &retHash

		return result, nil
	default:
		return result, errors.Errorf("error invalid memo %v", memo)
	}
}

// TODO: testing
func ConvertTxV0Ext(e xdr.TransactionV0Ext) (TransactionV0Ext, error) {
	return TransactionV0Ext{
		V: e.V,
	}, nil
}

func ConvertTxExt(e xdr.TransactionExt) (TransactionExt, error) {
	var result TransactionExt
	data, err := ConvertSorobanTransactionData(*e.SorobanData)
	if err != nil {
		return result, err
	}

	result.V = e.V
	result.SorobanData = &data

	return result, nil
}

func ConvertFeeBumpTransactionInnerTx(f xdr.FeeBumpTransactionInnerTx) (FeeBumpTransactionInnerTx, error) {
	var result FeeBumpTransactionInnerTx
	switch f.Type {
	case xdr.EnvelopeTypeEnvelopeTypeTx:
		v1, err := ConvertTransactionV1Envelope(f.V1)
		if err != nil {
			return result, err
		}
		result.V1 = &v1

		return result, nil
	}
	return result, errors.Errorf("error invalid FeeBumpTransactionInnerTx %v", f.Type)
}

func ConvertFeeBumpTransactionExt(f xdr.FeeBumpTransactionExt) FeeBumpTransactionExt {
	return FeeBumpTransactionExt{V: f.V}
}
