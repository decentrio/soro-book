package xdr

import "github.com/stellar/go/xdr"

// func ConvertTransactionJSON(tx models.Transaction) models.TransactionJSON {
// 	return models.TransactionJSON{}
// }

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

func ConvertTimeBounds(tb *xdr.TimeBounds) (*TimeBounds, error) {
	return &TimeBounds{}, nil
}

func ConvertMemo(memo xdr.Memo) (Memo, error) {
	return Memo{}, nil
}

func ConvertOperation(op xdr.Operation) (Operation, error) {
	return Operation{}, nil
}

func ConvertTxV0Ext(ext xdr.TransactionV0Ext) (TransactionV0Ext, error) {
	return TransactionV0Ext{}, nil
}
