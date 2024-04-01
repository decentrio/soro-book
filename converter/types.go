package converter

import (
	"time"
)

type Hash []byte
type PoolId Hash
type ScSymbol string

type TransactionEnvelope struct {
	V0      *TransactionV0Envelope      `json:"v0,omitempty"`
	V1      *TransactionV1Envelope      `json:"v1,omitempty"`
	FeeBump *FeeBumpTransactionEnvelope `json:"feebump,omitempty"`
}

type TransactionV0Envelope struct {
	Tx         TransactionV0        `json:"tx,omitempty"`
	Signatures []DecoratedSignature `json:"signatures,omitempty"`
}

type TransactionV0 struct {
	SourceAccountEd25519 string           `json:"source_account_ed25519,omitempty"`
	Fee                  uint32           `json:"fee,omitempty"`
	SeqNum               int64            `json:"seq_num,omitempty"`
	TimeBounds           *TimeBounds      `json:"time_bounds,omitempty"`
	Memo                 Memo             `json:"memo,omitempty"`
	Operations           []Operation      `json:"operations,omitempty"`
	Ext                  TransactionV0Ext `json:"ext,omitempty"`
}

type DecoratedSignature struct {
	Hint      []byte `json:"hint,omitempty"`
	Signature []byte `json:"signature,omitempty"`
}

type TimeBounds struct {
	MinTime uint64 `json:"min_time,omitempty"`
	MaxTime uint64 `json:"max_time,omitempty"`
}

type Memo struct {
	Text    *string `json:"text,omitempty"`
	Id      *uint64 `json:"id,omitempty"`
	Hash    *string `json:"hash,omitempty"`
	RetHash *string `json:"rethash,omitempty"`
}

type Operation struct {
	SourceAccount *MuxedAccount `json:"source_account,omitempty"`
	Body          OperationBody `json:"body,omitempty"`
}

type MuxedAccount struct {
	Ed25519  *string               `json:"ed25519,omitempty"`
	Med25519 *MuxedAccountMed25519 `json:"med25519,omitempty"`
}

type MuxedAccountMed25519 struct {
	Id      uint64  `json:"id,omitempty"`
	Ed25519 *string `json:"ed25519,omitempty"`
}

type OperationBody struct {
	CreateAccountOp                 *CreateAccountOp                 `json:"create_account_op,omitempty"`
	PaymentOp                       *PaymentOp                       `json:"payment_op,omitempty"`
	PathPaymentStrictReceiveOp      *PathPaymentStrictReceiveOp      `json:"path_payment_strict_receive_op,omitempty"`
	ManageSellOfferOp               *ManageSellOfferOp               `json:"manage_sell_offer_op,omitempty"`
	CreatePassiveSellOfferOp        *CreatePassiveSellOfferOp        `json:"create_passive_sell_offer_op,omitempty"`
	SetOptionsOp                    *SetOptionsOp                    `json:"set_options_op,omitempty"`
	ChangeTrustOp                   *ChangeTrustOp                   `json:"change_trust_op,omitempty"`
	AllowTrustOp                    *AllowTrustOp                    `json:"allow_trust_op,omitempty"`
	Destination                     *MuxedAccount                    `json:"muxed_account,omitempty"`
	ManageDataOp                    *ManageDataOp                    `json:"manage_data_op,omitempty"`
	BumpSequenceOp                  *BumpSequenceOp                  `json:"bump_sequence_op,omitempty"`
	ManageBuyOfferOp                *ManageBuyOfferOp                `json:"manage_buy_offer_op,omitempty"`
	PathPaymentStrictSendOp         *PathPaymentStrictSendOp         `json:"path_payment_strict_send_op,omitempty"`
	CreateClaimableBalanceOp        *CreateClaimableBalanceOp        `json:"create_claimable_balance_op,omitempty"`
	ClaimClaimableBalanceOp         *ClaimClaimableBalanceOp         `json:"claim_claimable_balance_op,omitempty"`
	BeginSponsoringFutureReservesOp *BeginSponsoringFutureReservesOp `json:"begin_sponsoring_future_reserves_op,omitempty"`
	RevokeSponsorshipOp             *RevokeSponsorshipOp             `json:"revoke_sponsorship_op,omitempty"`
	ClawbackOp                      *ClawbackOp                      `json:"clawback_op,omitempty"`
	ClawbackClaimableBalanceOp      *ClawbackClaimableBalanceOp      `json:"clawback_claimable_balance_op,omitempty"`
	SetTrustLineFlagsOp             *SetTrustLineFlagsOp             `json:"set_trust_line_flags_op,omitempty"`
	LiquidityPoolDepositOp          *LiquidityPoolDepositOp          `json:"liquidity_pool_deposit_op,omitempty"`
	LiquidityPoolWithdrawOp         *LiquidityPoolWithdrawOp         `json:"liquidity_pool_withdraw_op,omitempty"`
	InvokeHostFunctionOp            *InvokeHostFunctionOp            `json:"invoke_host_function_op,omitempty"`
	ExtendFootprintTtlOp            *ExtendFootprintTtlOp            `json:"extend_footprint_ttl_op,omitempty"`
	RestoreFootprintOp              *RestoreFootprintOp              `json:"restore_footprint_op,omitempty"`
}

type PublicKey struct {
	Ed25519 string `json:"ed25519,omitempty"`
}

type Asset struct {
	AssetType string    `json:"asset_type,omitempty"`
	AssetCode []byte    `json:"asset_code,omitempty"`
	Issuer    PublicKey `json:"issuer,omitempty"`
}

type CreateAccountOp struct {
	Destination     PublicKey `json:"public_key,omitempty"`
	StartingBalance int64     `json:"starting_balance,omitempty"`
}

type PaymentOp struct {
	Destination MuxedAccount `json:"muxed_account,omitempty"`
	Asset       Asset        `json:"asset,omitempty"`
	Amount      int64        `json:"amount,omitempty"`
}

type PathPaymentStrictReceiveOp struct {
	SendAsset   Asset        `json:"send_asset,omitempty"`
	SendMax     int64        `json:"send_max,omitempty"`
	Destination MuxedAccount `json:"destination,omitempty"`
	DestAsset   Asset        `json:"dest_asset,omitempty"`
	DestAmount  int64        `json:"dest_amount,omitempty"`
	Path        []Asset      `json:"path,omitempty"`
}

type Price struct {
	N int32 `json:"n,omitempty"`
	D int32 `json:"d,omitempty"`
}

type ManageSellOfferOp struct {
	Selling   Asset `json:"selling,omitempty"`
	Buying    Asset `json:"buying,omitempty"`
	BuyAmount int64 `json:"buy_amount,omitempty"`
	Price     Price `json:"price,omitempty"`
	OfferId   int64 `json:"offer_id,omitempty"`
}

type CreatePassiveSellOfferOp struct {
	Selling Asset `json:"selling,omitempty"`
	Buying  Asset `json:"buying,omitempty"`
	Amount  int64 `json:"amount,omitempty"`
	Price   Price `json:"price,omitempty"`
}

type SetOptionsOp struct {
	InflationDest *PublicKey `json:"inflation_dest,omitempty"`
	ClearFlags    *uint32    `json:"clear_flags,omitempty"`
	SetFlags      *uint32    `json:"set_flags,omitempty"`
	MasterWeight  *uint32    `json:"master_weight,omitempty"`
	LowThreshold  *uint32    `json:"low_threshold,omitempty"`
	MedThreshold  *uint32    `json:"med_threshold,omitempty"`
	HighThreshold *uint32    `json:"high_threshold,omitempty"`
	HomeDomain    *string    `json:"home_domain,omitempty"`
	Signer        *Signer    `json:"signer,omitempty"`
}

type Signer struct {
	Key    SignerKey `json:"signer_key,omitempty"`
	Weight uint32    `json:"weight,omitempty"`
}

type SignerKey struct {
	Ed25519              *string                        `json:"ed25519,omitempty"`
	PreAuthTx            *string                        `json:"pre_auth_tx,omitempty"`
	HashX                *string                        `json:"hash_x,omitempty"`
	Ed25519SignedPayload *SignerKeyEd25519SignedPayload `json:"ed25519_signed_payload,omitempty"`
}

type SignerKeyEd25519SignedPayload struct {
	Ed25519 string `json:"ed25519,omitempty"`
	Payload []byte `json:"payload,omitempty"`
}

type ChangeTrustOp struct {
	Line  ChangeTrustAsset `json:"change_trust_asset,omitempty"`
	Limit int64            `json:"int64,omitempty"`
}

type ChangeTrustAsset struct {
	Asset         *Asset                   `json:"asset,omitempty"`
	LiquidityPool *LiquidityPoolParameters `json:"liquidity_pool,omitempty"`
}

type LiquidityPoolParameters struct {
	ConstantProduct *LiquidityPoolConstantProductParameters `json:"constant_product,omitempty"`
}

type LiquidityPoolConstantProductParameters struct {
	AssetA Asset `json:"asset_a,omitempty"`
	AssetB Asset `json:"asset_b,omitempty"`
	Fee    int32 `json:"fee,omitempty"`
}

type AllowTrustOp struct {
	Trustor   PublicKey `json:"trustor,omitempty"`
	AssetCode []byte    `json:"asset_code,omitempty"`
	Authorize uint32    `json:"authorize,omitempty"`
}

type ManageDataOp struct {
	DataName  string `json:"data_name,omitempty"`
	DataValue []byte `json:"data_value,omitempty"`
}

type BumpSequenceOp struct {
	BumpTo int64 `json:"bump_to,omitempty"`
}

type ManageBuyOfferOp struct {
	Selling   Asset `json:"selling,omitempty"`
	Buying    Asset `json:"buying,omitempty"`
	BuyAmount int64 `json:"buy_amount,omitempty"`
	Price     Price `json:"price,omitempty"`
	OfferId   int64 `json:"offer_id,omitempty"`
}

type PathPaymentStrictSendOp struct {
	SendAsset   Asset        `json:"send_asset,omitempty"`
	SendAmount  int64        `json:"send_amount,omitempty"`
	Destination MuxedAccount `json:"destination,omitempty"`
	DestAsset   Asset        `json:"dest_asset,omitempty"`
	DestMin     int64        `json:"dest_min,omitempty"`
	Path        []Asset      `json:"path,omitempty"`
}

type CreateClaimableBalanceOp struct {
	Asset     Asset      `json:"asset,omitempty"`
	Amount    int64      `json:"amount,omitempty"`
	Claimants []Claimant `json:"claimants,omitempty"`
}

type Claimant struct {
	V0 *ClaimantV0 `json:"v0,omitempty"`
}

type ClaimantV0 struct {
	Destination PublicKey      `json:"destination,omitempty"`
	Predicate   ClaimPredicate `json:"predicate,omitempty"`
}

type ClaimPredicate struct {
	AndPredicates  *[]ClaimPredicate `json:"and_predicates,omitempty"`
	OrPredicates   *[]ClaimPredicate `json:"or_predicates,omitempty"`
	NotPredicate   *ClaimPredicate   `json:"not_predicates,omitempty"`
	AbsBefore      *time.Time        `json:"abs_before,omitempty"`
	AbsBeforeEpoch *int64            `json:"abs_before_epoch,omitempty"`
	RelBefore      *int64            `json:"rel_before,omitempty"`
}

type ClaimClaimableBalanceOp struct {
	BalanceId ClaimableBalanceId `json:"balance_id,omitempty"`
}

type ClaimableBalanceId struct {
	V0 *string `json:"v0,omitempty"`
}

type BeginSponsoringFutureReservesOp struct {
	SponsoredId PublicKey `json:"sponsored_id,omitempty"`
}

type RevokeSponsorshipOp struct {
	LedgerKey *LedgerKey                 `json:"ledger_key,omitempty"`
	Signer    *RevokeSponsorshipOpSigner `json:"signer,omitempty"`
}

type RevokeSponsorshipOpSigner struct {
	AccountId PublicKey `json:"account_id,omitempty"`
	SignerKey SignerKey `json:"signer_key,omitempty"`
}

type LedgerKey struct {
	Account          *LedgerKeyAccount          `json:"account,omitempty"`
	TrustLine        *LedgerKeyTrustLine        `json:"trust_line,omitempty"`
	Offer            *LedgerKeyOffer            `json:"offer,omitempty"`
	Data             *LedgerKeyData             `json:"data,omitempty"`
	ClaimableBalance *LedgerKeyClaimableBalance `json:"claimable_balance,omitempty"`
	LiquidityPool    *LedgerKeyLiquidityPool    `json:"liquidity_pool,omitempty"`
	ContractData     *LedgerKeyContractData     `json:"contract_data,omitempty"`
	ContractCode     *LedgerKeyContractCode     `json:"contract_code,omitempty"`
	ConfigSetting    *LedgerKeyConfigSetting    `json:"config_setting,omitempty"`
	Ttl              *LedgerKeyTtl              `json:"ttl,omitempty"`
}

type LedgerKeyAccount struct {
	AccountId PublicKey `json:"account_id,omitempty"`
}

type LedgerKeyTrustLine struct {
	AccountId PublicKey      `json:"account_id,omitempty"`
	Asset     TrustLineAsset `json:"asset,omitempty"`
}

type TrustLineAsset struct {
	Asset           *Asset  `json:"asset,omitempty"`
	LiquidityPoolId *PoolId `json:"liquidity_pool_id,omitempty"`
}

type LedgerKeyOffer struct {
	SellerId PublicKey `json:"seller_id,omitempty"`
	OfferId  int64     `json:"offer_id,omitempty"`
}

type LedgerKeyData struct {
	AccountId PublicKey `json:"account_id,omitempty"`
	DataName  string    `json:"data_name,omitempty"`
}

type LedgerKeyClaimableBalance struct {
	BalanceId ClaimableBalanceId `json:"balance_id,omitempty"`
}

type LedgerKeyLiquidityPool struct {
	LiquidityPoolId PoolId `json:"liquidity_pool_id,omitempty"`
}

type LedgerKeyContractData struct {
	Contract   ScAddress `json:"contract,omitempty"`
	Key        ScVal     `json:"key,omitempty"`
	Durability int32     `json:"durability,omitempty"`
}

type ScAddress struct {
	AccountId  *PublicKey `json:"account_id,omitempty"`
	ContractId *string    `json:"contract_id,omitempty"`
}

type ScVal struct {
	B         *bool               `json:"b,omitempty"`
	Error     *ScError            `json:"error,omitempty"`
	U32       *uint32             `json:"u32,omitempty"`
	I32       *int32              `json:"i32,omitempty"`
	U64       *uint64             `json:"u64,omitempty"`
	I64       *int64              `json:"i64,omitempty"`
	Timepoint *uint64             `json:"timepoint,omitempty"`
	Duration  *uint64             `json:"duration,omitempty"`
	U128      *UInt128Parts       `json:"u128,omitempty"`
	I128      *Int128Parts        `json:"i128,omitempty"`
	U256      *UInt256Parts       `json:"u256,omitempty"`
	I256      *Int256Parts        `json:"i256,omitempty"`
	Bytes     *ScBytes            `json:"bytes,omitempty"`
	Str       *string             `json:"str,omitempty"`
	Sym       *ScSymbol           `json:"sym,omitempty"`
	Vec       *[]ScVal            `json:"vec,omitempty"`
	Map       *ScMap              `json:"map,omitempty"`
	Address   *ScAddress          `json:"address,omitempty"`
	NonceKey  *ScNonceKey         `json:"nonce_key,omitempty"`
	Instance  *ScContractInstance `json:"instance,omitempty"`
}

type ContractExecutable struct {
	WasmHash *string `json:"wasm_hash,omitempty"`
}

type ScContractInstance struct {
	Executable ContractExecutable `json:"executable,omitempty"`
	Storage    *ScMap             `json:"storage,omitempty"`
}

type ScNonceKey struct {
	Nonce int64 `json:"nonce,omitempty"`
}

type ScMap []ScMapEntry

type ScMapEntry struct {
	Key ScVal `json:"key,omitempty"`
	Val ScVal `json:"val,omitempty"`
}

type UInt128Parts struct {
	Hi uint64 `json:"hi,omitempty"`
	Lo uint64 `json:"lo,omitempty"`
}

type Int128Parts struct {
	Hi int64  `json:"hi,omitempty"`
	Lo uint64 `json:"lo,omitempty"`
}

type UInt256Parts struct {
	HiHi uint64 `json:"hihi,omitempty"`
	HiLo uint64 `json:"hilo,omitempty"`
	LoHi uint64 `json:"lohi,omitempty"`
	LoLo uint64 `json:"lolo,omitempty"`
}

type Int256Parts struct {
	HiHi int64  `json:"hihi,omitempty"`
	HiLo uint64 `json:"hilo,omitempty"`
	LoHi uint64 `json:"lohi,omitempty"`
	LoLo uint64 `json:"lolo,omitempty"`
}

type ScBytes []byte

type ScError struct {
	ContractCode *uint32 `json:"contract_code,omitempty"`
	Code         *int32  `json:"code,omitempty"`
}

type LedgerKeyContractCode struct {
	Hash string `json:"hash,omitempty"`
}

type LedgerKeyConfigSetting struct {
	ConfigSettingId int32 `json:"config_setting_id,omitempty"`
}

type LedgerKeyTtl struct {
	KeyHash string `json:"key_hash,omitempty"`
}

type ClawbackOp struct {
	Asset  Asset        `json:"asset,omitempty"`
	From   MuxedAccount `json:"from,omitempty"`
	Amount int64        `json:"amount,omitempty"`
}

type ClawbackClaimableBalanceOp struct {
	BalanceId ClaimableBalanceId `json:"balance_id,omitempty"`
}

type SetTrustLineFlagsOp struct {
	Trustor    PublicKey `json:"trustor,omitempty"`
	Asset      Asset     `json:"asset,omitempty"`
	ClearFlags uint32    `json:"clear_flags,omitempty"`
	SetFlags   uint32    `json:"set_flags,omitempty"`
}

type LiquidityPoolDepositOp struct {
	LiquidityPoolId PoolId `json:"liquidity_pool_id,omitempty"`
	MaxAmountA      int64  `json:"max_amount_a,omitempty"`
	MaxAmountB      int64  `json:"max_amount_b,omitempty"`
	MinPrice        Price  `json:"min_price,omitempty"`
	MaxPrice        Price  `json:"max_price,omitempty"`
}

type LiquidityPoolWithdrawOp struct {
	LiquidityPoolId PoolId `json:"liquidity_pool_id,omitempty"`
	Amount          int64  `json:"amount,omitempty"`
	MinAmountA      int64  `json:"min_amount_a,omitempty"`
	MinAmountB      int64  `json:"min_amount_b,omitempty"`
}

type InvokeHostFunctionOp struct {
	HostFunction HostFunction                `json:"host_function,omitempty"`
	Auth         []SorobanAuthorizationEntry `json:"auth,omitempty"`
}

type SorobanAuthorizationEntry struct {
	Credentials    SorobanCredentials          `json:"credentials,omitempty"`
	RootInvocation SorobanAuthorizedInvocation `json:"root_invocation,omitempty"`
}

type SorobanCredentials struct {
	Address *SorobanAddressCredentials `json:"address,omitempty"`
}

type SorobanAddressCredentials struct {
	Address                   ScAddress `json:"address,omitempty"`
	Nonce                     int64     `json:"nonce,omitempty"`
	SignatureExpirationLedger uint32    `json:"signature_expiration_ledger,omitempty"`
	Signature                 ScVal     `json:"signature,omitempty"`
}

type SorobanAuthorizedInvocation struct {
	Function       SorobanAuthorizedFunction     `json:"function,omitempty"`
	SubInvocations []SorobanAuthorizedInvocation `json:"sub_invocations,omitempty"`
}

type SorobanAuthorizedFunction struct {
	ContractFn           *InvokeContractArgs `json:"contract_fn,omitempty"`
	CreateContractHostFn *CreateContractArgs `json:"create_contract_host_fn,omitempty"`
}

type HostFunction struct {
	InvokeContract *InvokeContractArgs `json:"invoke_contract,omitempty"`
	CreateContract *CreateContractArgs `json:"create_contract,omitempty"`
	Wasm           *[]byte             `json:"wasm,omitempty"`
}

type InvokeContractArgs struct {
	ContractAddress ScAddress `json:"contract_address,omitempty"`
	FunctionName    ScSymbol  `json:"function_name,omitempty"`
	Args            []ScVal   `json:"args,omitempty"`
}

type CreateContractArgs struct {
	ContractIdPreimage ContractIdPreimage `json:"contract_id_preimage,omitempty"`
	Executable         ContractExecutable `json:"executable,omitempty"`
}

type ContractIdPreimage struct {
	FromAddress *ContractIdPreimageFromAddress `json:"from_address,omitempty"`
	FromAsset   *Asset                         `json:"from_asset,omitempty"`
}

type ContractIdPreimageFromAddress struct {
	Address ScAddress `json:"address,omitempty"`
	Salt    string    `json:"salt,omitempty"`
}

type ExtendFootprintTtlOp struct {
	Ext      ExtensionPoint `json:"ext,omitempty"`
	ExtendTo uint32         `json:"extend_to,omitempty"`
}

type RestoreFootprintOp struct {
	Ext ExtensionPoint `json:"ext,omitempty"`
}

type ExtensionPoint struct {
	V int32 `json:"v,omitempty"`
}

type TransactionV0Ext struct {
	V int32 `json:"v,omitempty"`
}

type TransactionV1Envelope struct {
	Tx         Transaction          `json:"tx,omitempty"`
	Signatures []DecoratedSignature `json:"signatures,omitempty"`
}

type Transaction struct {
	SourceAccount MuxedAccount   `json:"source_account,omitempty"`
	Fee           uint32         `json:"fee,omitempty"`
	SeqNum        int64          `json:"seq_num,omitempty"`
	Cond          Preconditions  `json:"cond,omitempty"`
	Memo          Memo           `json:"memo,omitempty"`
	Operations    []Operation    `json:"operations,omitempty"`
	Ext           TransactionExt `json:"ext,omitempty"`
}

type TransactionExt struct {
	V           int32                   `json:"V,omitempty"`
	SorobanData *SorobanTransactionData `json:"soroban_data,omitempty"`
}

type SorobanTransactionData struct {
	Ext         ExtensionPoint   `json:"ext,omitempty"`
	Resources   SorobanResources `json:"resources,omitempty"`
	ResourceFee int64            `json:"resource_fee,omitempty"`
}

type SorobanResources struct {
	Footprint    LedgerFootprint `json:"footprint,omitempty"`
	Instructions uint32          `json:"instructions,omitempty"`
	ReadBytes    uint32          `json:"read_bytes,omitempty"`
	WriteBytes   uint32          `json:"write_bytes,omitempty"`
}

type LedgerFootprint struct {
	ReadOnly  []LedgerKey `json:"read_only,omitempty"`
	ReadWrite []LedgerKey `json:"read_write,omitempty"`
}

type Preconditions struct {
	TimeBounds *TimeBounds      `json:"time_bounds,omitempty"`
	V2         *PreconditionsV2 `json:"v2,omitempty"`
}

type PreconditionsV2 struct {
	TimeBounds      *TimeBounds   `json:"time_bounds,omitempty"`
	LedgerBounds    *LedgerBounds `json:"ledger_bounds,omitempty"`
	MinSeqNum       *int64        `json:"min_seq_num,omitempty"`
	MinSeqAge       uint64        `json:"min_seq_age,omitempty"`
	MinSeqLedgerGap uint32        `json:"min_seq_ledger_gap,omitempty"`
	ExtraSigners    []SignerKey   `json:"extra_signers,omitempty"`
}

type LedgerBounds struct {
	MinLedger uint32 `json:"min_ledger,omitempty"`
	MaxLedger uint32 `json:"max_ledger,omitempty"`
}

type FeeBumpTransactionEnvelope struct {
	Tx         FeeBumpTransaction   `json:"tx,omitempty"`
	Signatures []DecoratedSignature `json:"signatures,omitempty"`
}

type FeeBumpTransaction struct {
	FeeSource MuxedAccount              `json:"fee_source,omitempty"`
	Fee       int64                     `json:"fee,omitempty"`
	InnerTx   FeeBumpTransactionInnerTx `json:"inner_tx,omitempty"`
	Ext       FeeBumpTransactionExt     `json:"ext,omitempty"`
}

type FeeBumpTransactionInnerTx struct {
	V1 *TransactionV1Envelope `json:"v1,omitempty"`
}

type FeeBumpTransactionExt struct {
	V int32 `json:"v,omitempty"`
}

type TransactionResultPair struct {
	TransactionHash string            `json:"transaction_hash,omitempty"`
	Result          TransactionResult `json:"result,omitempty"`
}

type TransactionResult struct {
	FeeCharged int64                   `json:"fee_charged,omitempty"`
	Result     TransactionResultResult `json:"result,omitempty"`
	Ext        TransactionResultExt    `json:"ext,omitempty"`
}

type TransactionResultResult struct {
	Code            int32                       `json:"code,omitempty"`
	InnerResultPair *InnerTransactionResultPair `json:"inner_result_pair,omitempty"`
	Results         *[]OperationResult          `json:"results,omitempty"`
}

type InnerTransactionResultPair struct {
	TransactionHash string                 `json:"transaction_hash,omitempty"`
	Result          InnerTransactionResult `json:"result,omitempty"`
}

type InnerTransactionResult struct {
	FeeCharged int64                        `json:"fee_charged,omitempty"`
	Result     InnerTransactionResultResult `json:"result,omitempty"`
	Ext        InnerTransactionResultExt    `json:"ext,omitempty"`
}

type InnerTransactionResultResult struct {
	Code    int32              `json:"code,omitempty"`
	Results *[]OperationResult `json:"results,omitempty"`
}

type InnerTransactionResultExt struct {
	V int32 `json:"v,omitempty"`
}

type OperationResult struct {
	Code int32              `json:"code,omitempty"`
	Tr   *OperationResultTr `json:"tr,omitempty"`
}

type OperationResultTr struct {
	CreateAccountResult                 *CreateAccountResult                 `json:"create_account_result,omitempty"`
	PaymentResult                       *PaymentResult                       `json:"payment_result,omitempty"`
	PathPaymentStrictReceiveResult      *PathPaymentStrictReceiveResult      `json:"path_payment_strict_receive_result,omitempty"`
	ManageSellOfferResult               *ManageSellOfferResult               `json:"manage_sell_offer_result,omitempty"`
	CreatePassiveSellOfferResult        *ManageSellOfferResult               `json:"create_passive_sell_offer_result,omitempty"`
	SetOptionsResult                    *SetOptionsResult                    `json:"set_options_result,omitempty"`
	ChangeTrustResult                   *ChangeTrustResult                   `json:"change_trust_result,omitempty"`
	AllowTrustResult                    *AllowTrustResult                    `json:"allow_trust_result,omitempty"`
	AccountMergeResult                  *AccountMergeResult                  `json:"account_merge_result,omitempty"`
	InflationResult                     *InflationResult                     `json:"inflation_result,omitempty"`
	ManageDataResult                    *ManageDataResult                    `json:"manage_data_result,omitempty"`
	BumpSeqResult                       *BumpSequenceResult                  `json:"bump_seq_result,omitempty"`
	ManageBuyOfferResult                *ManageBuyOfferResult                `json:"manage_buy_offer_result,omitempty"`
	PathPaymentStrictSendResult         *PathPaymentStrictSendResult         `json:"path_payment_strict_send_result,omitempty"`
	CreateClaimableBalanceResult        *CreateClaimableBalanceResult        `json:"create_claimable_balance_result,omitempty"`
	ClaimClaimableBalanceResult         *ClaimClaimableBalanceResult         `json:"claim_claimable_balance_result,omitempty"`
	BeginSponsoringFutureReservesResult *BeginSponsoringFutureReservesResult `json:"begin_sponsoring_future_reserves_result,omitempty"`
	EndSponsoringFutureReservesResult   *EndSponsoringFutureReservesResult   `json:"end_sponsoring_future_reserves_result,omitempty"`
	RevokeSponsorshipResult             *RevokeSponsorshipResult             `json:"revoke_sponsorship_result,omitempty"`
	ClawbackResult                      *ClawbackResult                      `json:"clawback_result,omitempty"`
	ClawbackClaimableBalanceResult      *ClawbackClaimableBalanceResult      `json:"clawback_claimable_balance_result,omitempty"`
	SetTrustLineFlagsResult             *SetTrustLineFlagsResult             `json:"set_trust_line_flags_result,omitempty"`
	LiquidityPoolDepositResult          *LiquidityPoolDepositResult          `json:"liquidity_pool_deposit_result,omitempty"`
	LiquidityPoolWithdrawResult         *LiquidityPoolWithdrawResult         `json:"liquidity_pool_withdraw_result,omitempty"`
	InvokeHostFunctionResult            *InvokeHostFunctionResult            `json:"invoke_host_function_result,omitempty"`
	ExtendFootprintTtlResult            *ExtendFootprintTtlResult            `json:"extend_footprint_ttl_result,omitempty"`
	RestoreFootprintResult              *RestoreFootprintResult              `json:"restore_footprint_result,omitempty"`
}

type RestoreFootprintResult struct {
	Code int32 `json:"code,omitempty"`
}

type ExtendFootprintTtlResult struct {
	Code int32 `json:"code,omitempty"`
}

type InvokeHostFunctionResult struct {
	Code    int32   `json:"code,omitempty"`
	Success *string `json:"success,omitempty"`
}

type LiquidityPoolWithdrawResult struct {
	Code int32 `json:"code,omitempty"`
}

type LiquidityPoolDepositResult struct {
	Code int32 `json:"code,omitempty"`
}

type SetTrustLineFlagsResult struct {
	Code int32 `json:"code,omitempty"`
}

type ClawbackClaimableBalanceResult struct {
	Code int32 `json:"code,omitempty"`
}

type ClawbackResult struct {
	Code int32 `json:"code,omitempty"`
}

type RevokeSponsorshipResult struct {
	Code int32 `json:"code,omitempty"`
}

type EndSponsoringFutureReservesResult struct {
	Code int32 `json:"code,omitempty"`
}

type BeginSponsoringFutureReservesResult struct {
	Code int32 `json:"code,omitempty"`
}

type ClaimClaimableBalanceResult struct {
	Code int32 `json:"code,omitempty"`
}

type CreateClaimableBalanceResult struct {
	Code      int32               `json:"code,omitempty"`
	BalanceId *ClaimableBalanceId `json:"balance_id,omitempty"`
}

type PathPaymentStrictSendResult struct {
	Code     int32                               `json:"code,omitempty"`
	Success  *PathPaymentStrictSendResultSuccess `json:"success,omitempty"`
	NoIssuer *Asset                              `json:"no_issuer,omitempty"`
}

type PathPaymentStrictSendResultSuccess struct {
	Offers []ClaimAtom         `json:"offers,omitempty"`
	Last   SimplePaymentResult `json:"last,omitempty"`
}

type ManageBuyOfferResult struct {
	Code    int32                     `json:"code,omitempty"`
	Success *ManageOfferSuccessResult `json:"success,omitempty"`
}

type BumpSequenceResult struct {
	Code int32 `json:"code,omitempty"`
}

type ManageDataResult struct {
	Code int32 `json:"code,omitempty"`
}

type InflationPayout struct {
	Destination PublicKey `json:"destination,omitempty"`
	Amount      int64     `json:"amount,omitempty"`
}

type InflationResult struct {
	Code    int32              `json:"code,omitempty"`
	Payouts *[]InflationPayout `json:"payouts,omitempty"`
}

type AccountMergeResult struct {
	Code                 int32  `json:"code,omitempty"`
	SourceAccountBalance *int64 `json:"source_account_balance,omitempty"`
}

type AllowTrustResult struct {
	Code int32 `json:"code,omitempty"`
}

type ChangeTrustResult struct {
	Code int32 `json:"code,omitempty"`
}

type SetOptionsResult struct {
	Code int32 `json:"code,omitempty"`
}

type ManageSellOfferResult struct {
	Code    int32                     `json:"code,omitempty"`
	Success *ManageOfferSuccessResult `json:"success,omitempty"`
}

type ManageOfferSuccessResult struct {
	OffersClaimed []ClaimAtom                   `json:"offers_claimed,omitempty"`
	Offer         ManageOfferSuccessResultOffer `json:"offer,omitempty"`
}

type OfferEntry struct {
	SellerId PublicKey     `json:"seller_id,omitempty"`
	OfferId  int64         `json:"offer_id,omitempty"`
	Selling  Asset         `json:"selling,omitempty"`
	Buying   Asset         `json:"buying,omitempty"`
	Amount   int64         `json:"amount,omitempty"`
	Price    Price         `json:"price,omitempty"`
	Flags    uint32        `json:"flags,omitempty"`
	Ext      OfferEntryExt `json:"ext,omitempty"`
}

type OfferEntryExt struct {
	V int32 `json:"v,omitempty"`
}

type ManageOfferSuccessResultOffer struct {
	Effect int32       `json:"effect,omitempty"`
	Offer  *OfferEntry `json:"offer,omitempty"`
}

type CreateAccountResult struct {
	Code int32 `json:"code,omitempty"`
}

type PaymentResult struct {
	Code int32 `json:"code,omitempty"`
}

type PathPaymentStrictReceiveResult struct {
	Code     int32                                  `json:"code,omitempty"`
	Success  *PathPaymentStrictReceiveResultSuccess `json:"success,omitempty"`
	NoIssuer *Asset                                 `json:"no_issuer,omitempty"`
}

type PathPaymentStrictReceiveResultSuccess struct {
	Offers []ClaimAtom         `json:"offers,omitempty"`
	Last   SimplePaymentResult `json:"last,omitempty"`
}

type ClaimAtom struct {
	V0            *ClaimOfferAtomV0   `json:"v0,omitempty"`
	OrderBook     *ClaimOfferAtom     `json:"order_book,omitempty"`
	LiquidityPool *ClaimLiquidityAtom `json:"liquidity_pool,omitempty"`
}

type ClaimOfferAtomV0 struct {
	SellerEd25519 string `json:"seller_ed25519,omitempty"`
	OfferId       int64  `json:"offer_id,omitempty"`
	AssetSold     Asset  `json:"asset_sold,omitempty"`
	AmountSold    int64  `json:"amount_sold,omitempty"`
	AssetBought   Asset  `json:"asset_bought,omitempty"`
	AmountBought  int64  `json:"amount_bought,omitempty"`
}

type ClaimOfferAtom struct {
	SellerId     PublicKey `json:"seller_id,omitempty"`
	OfferId      int64     `json:"offer_id,omitempty"`
	AssetSold    Asset     `json:"asset_sold,omitempty"`
	AmountSold   int64     `json:"amount_sold,omitempty"`
	AssetBought  Asset     `json:"asset_bought,omitempty"`
	AmountBought int64     `json:"amount_bought,omitempty"`
}

type ClaimLiquidityAtom struct {
	LiquidityPoolId PoolId `json:"liquidity_pool_id,omitempty"`
	AssetSold       Asset  `json:"asset_sold,omitempty"`
	AmountSold      int64  `json:"amount_sold,omitempty"`
	AssetBought     Asset  `json:"asset_bought,omitempty"`
	AmountBought    int64  `json:"amount_bought,omitempty"`
}

type SimplePaymentResult struct {
	Destination PublicKey `json:"destination,omitempty"`
	Asset       Asset     `json:"asset,omitempty"`
	Amount      int64     `json:"amount,omitempty"`
}

type TransactionResultExt struct {
	V int32 `json:"v,omitempty"`
}

type LedgerEntryChanges []LedgerEntryChange

type LedgerEntryChange struct {
	Created *LedgerEntry `json:"created,omitempty"`
	Updated *LedgerEntry `json:"updated,omitempty"`
	Removed *LedgerKey   `json:"removed,omitempty"`
	State   *LedgerEntry `json:"state,omitempty"`
}

type LedgerEntry struct {
	LastModifiedLedgerSeq uint32          `json:"last_modified_ledger_seq,omitempty"`
	Data                  LedgerEntryData `json:"data,omitempty"`
	Ext                   LedgerEntryExt  `json:"ext,omitempty"`
}

type LedgerEntryData struct {
	Account          *AccountEntry          `json:"account,omitempty"`
	TrustLine        *TrustLineEntry        `json:"trust_line,omitempty"`
	Offer            *OfferEntry            `json:"offer,omitempty"`
	Data             *DataEntry             `json:"data,omitempty"`
	ClaimableBalance *ClaimableBalanceEntry `json:"claimable_balance,omitempty"`
	LiquidityPool    *LiquidityPoolEntry    `json:"liquidity_pool,omitempty"`
	ContractData     *ContractDataEntry     `json:"contract_data,omitempty"`
	ContractCode     *ContractCodeEntry     `json:"contract_code,omitempty"`
	ConfigSetting    *ConfigSettingEntry    `json:"config_setting,omitempty"`
	Ttl              *TtlEntry              `json:"ttl,omitempty"`
}

type TtlEntry struct {
	KeyHash            string `json:"key_hash,omitempty"`
	LiveUntilLedgerSeq uint32 `json:"live_until_ledger_seq,omitempty"`
}

type ConfigSettingEntry struct {
	ConfigSettingId            int32                                  `json:"config_setting_id,omitempty"`
	ContractMaxSizeBytes       *uint32                                `json:"contract_max_size_bytes,omitempty"`
	ContractCompute            *ConfigSettingContractComputeV0        `json:"contract_compute,omitempty"`
	ContractLedgerCost         *ConfigSettingContractLedgerCostV0     `json:"contract_ledger_cost,omitempty"`
	ContractHistoricalData     *ConfigSettingContractHistoricalDataV0 `json:"contract_historical_data,omitempty"`
	ContractEvents             *ConfigSettingContractEventsV0         `json:"contract_events,omitempty"`
	ContractBandwidth          *ConfigSettingContractBandwidthV0      `json:"contract_bandwith,omitempty"`
	ContractCostParamsCpuInsns *ContractCostParams                    `json:"contract_cost_params_cpu_insns,omitempty"`
	ContractCostParamsMemBytes *ContractCostParams                    `json:"contract_cost_params_mem_bytes,omitempty"`
	ContractDataKeySizeBytes   *uint32                                `json:"contract_data_key_size_bytes,omitempty"`
	ContractDataEntrySizeBytes *uint32                                `json:"contract_data_entry_size_bytes,omitempty"`
	StateArchivalSettings      *StateArchivalSettings                 `json:"state_archival_settings,omitempty"`
	ContractExecutionLanes     *ConfigSettingContractExecutionLanesV0 `json:"contract_execute_lanes,omitempty"`
	BucketListSizeWindow       *[]uint64                              `json:"bucket_list_size_window,omitempty"`
	EvictionIterator           *EvictionIterator                      `json:"eviction_iterator,omitempty"`
}

type EvictionIterator struct {
	BucketListLevel  uint32 `json:"bucket_list_level,omitempty"`
	IsCurrBucket     bool   `json:"is_curr_bucket,omitempty"`
	BucketFileOffset uint64 `json:"bucket_file_offset,omitempty"`
}

type ConfigSettingContractExecutionLanesV0 struct {
	LedgerMaxTxCount uint32 `json:"ttl,omitempty"`
}

type StateArchivalSettings struct {
	MaxEntryTtl                    uint32 `json:"max_entry_ttl,omitempty"`
	MinTemporaryTtl                uint32 `json:"min_temporary_ttl,omitempty"`
	MinPersistentTtl               uint32 `json:"min_persistent_ttl,omitempty"`
	PersistentRentRateDenominator  int64  `json:"persistent_rent_rate_denominator,omitempty"`
	TempRentRateDenominator        int64  `json:"temp_rent_rate_denominator,omitempty"`
	MaxEntriesToArchive            uint32 `json:"max_entries_to_archive,omitempty"`
	BucketListSizeWindowSampleSize uint32 `json:"bucket_list_size_window_sample_size,omitempty"`
	BucketListWindowSamplePeriod   uint32 `json:"bucket_list_window_sample_period,omitempty"`
	EvictionScanSize               uint32 `json:"eviction_scan_size,omitempty"`
	StartingEvictionScanLevel      uint32 `json:"starting_eviction_scan_level,omitempty"`
}

type ContractCostParams []ContractCostParamEntry

type ContractCostParamEntry struct {
	Ext        ExtensionPoint `json:"ext,omitempty"`
	ConstTerm  int64          `json:"const_term,omitempty"`
	LinearTerm int64          `json:"linear_term,omitempty"`
}

type ConfigSettingContractBandwidthV0 struct {
	LedgerMaxTxsSizeBytes uint32 `json:"ledger_max_txs_size_bytes,omitempty"`
	TxMaxSizeBytes        uint32 `json:"tx_max_size_bytes,omitempty"`
	FeeTxSize1Kb          int64  `json:"fee_tx_size_1kb,omitempty"`
}

type ConfigSettingContractEventsV0 struct {
	TxMaxContractEventsSizeBytes uint32 `json:"tx_max_contract_events_size_bytes,omitempty"`
	FeeContractEvents1Kb         int64  `json:"fee_contract_events_1kb,omitempty"`
}

type ConfigSettingContractHistoricalDataV0 struct {
	FeeHistorical1Kb int64 `json:"fee_historical_1kb,omitempty"`
}

type ConfigSettingContractLedgerCostV0 struct {
	LedgerMaxReadLedgerEntries     uint32 `json:"ledger_max_read_ledger_entries,omitempty"`
	LedgerMaxReadBytes             uint32 `json:"ledger_max_read_bytes,omitempty"`
	LedgerMaxWriteLedgerEntries    uint32 `json:"ledger_max_write_ledger_entries,omitempty"`
	LedgerMaxWriteBytes            uint32 `json:"ledger_max_write_bytes,omitempty"`
	TxMaxReadLedgerEntries         uint32 `json:"tx_max_read_ledger_entries,omitempty"`
	TxMaxReadBytes                 uint32 `json:"tx_max_read_bytes,omitempty"`
	TxMaxWriteLedgerEntries        uint32 `json:"tx_max_write_ledger_entries,omitempty"`
	TxMaxWriteBytes                uint32 `json:"tx_max_write_bytes,omitempty"`
	FeeReadLedgerEntry             int64  `json:"fee_read_ledger_entry,omitempty"`
	FeeWriteLedgerEntry            int64  `json:"fee_write_ledger_entry,omitempty"`
	FeeRead1Kb                     int64  `json:"fee_read_1kb,omitempty"`
	BucketListTargetSizeBytes      int64  `json:"bucket_list_target_size_bytes,omitempty"`
	WriteFee1KbBucketListLow       int64  `json:"write_fee_1kb_bucket_list_low,omitempty"`
	WriteFee1KbBucketListHigh      int64  `json:"write_fee_1kb_bucket_list_high,omitempty"`
	BucketListWriteFeeGrowthFactor uint32 `json:"bucket_list_write_fee_growth_factor,omitempty"`
}

type ConfigSettingContractComputeV0 struct {
	LedgerMaxInstructions           int64  `json:"ledger_max_instructions,omitempty"`
	TxMaxInstructions               int64  `json:"tx_max_instructions,omitempty"`
	FeeRatePerInstructionsIncrement int64  `json:"fee_rate_per_instructions_increment,omitempty"`
	TxMemoryLimit                   uint32 `json:"tx_memory_limit,omitempty"`
}

type ContractCodeEntry struct {
	Ext  ExtensionPoint `json:"ext,omitempty"`
	Hash string         `json:"hash,omitempty"`
	Code []byte         `json:"code,omitempty"`
}

type ContractDataEntry struct {
	Ext        ExtensionPoint `json:"ext,omitempty"`
	Contract   ScAddress      `json:"contract,omitempty"`
	Key        ScVal          `json:"key,omitempty"`
	Durability int32          `json:"durability,omitempty"` //ContractDataDurability
	Val        ScVal          `json:"val,omitempty"`
}

type LiquidityPoolEntry struct {
	LiquidityPoolId PoolId                 `json:"liquidity_pool_id,omitempty"`
	Body            LiquidityPoolEntryBody `json:"body,omitempty"`
}

type LiquidityPoolEntryBody struct {
	ConstantProduct *LiquidityPoolEntryConstantProduct `json:"constant_product,omitempty"`
}

type LiquidityPoolEntryConstantProduct struct {
	Params                   LiquidityPoolConstantProductParameters `json:"params,omitempty"`
	ReserveA                 int64                                  `json:"reserve_a,omitempty"`
	ReserveB                 int64                                  `json:"reserve_b,omitempty"`
	TotalPoolShares          int64                                  `json:"total_pool_shares,omitempty"`
	PoolSharesTrustLineCount int64                                  `json:"pool_shares_trust_line_count,omitempty"`
}

type ClaimableBalanceEntry struct {
	BalanceId ClaimableBalanceId       `json:"balance_id,omitempty"`
	Claimants []Claimant               `json:"claimants,omitempty"`
	Asset     Asset                    `json:"asset,omitempty"`
	Amount    int64                    `json:"amount,omitempty"`
	Ext       ClaimableBalanceEntryExt `json:"ext,omitempty"`
}

type ClaimableBalanceEntryExt struct {
	V  int32                             `json:"v,omitempty"`
	V1 *ClaimableBalanceEntryExtensionV1 `json:"v1,omitempty"`
}

type ClaimableBalanceEntryExtensionV1 struct {
	Flags uint32                              `json:"flags,omitempty"`
	Ext   ClaimableBalanceEntryExtensionV1Ext `json:"ext,omitempty"`
}

type ClaimableBalanceEntryExtensionV1Ext struct {
	V int32 `json:"v,omitempty"`
}

type DataEntry struct {
	AccountId PublicKey    `json:"account_id,omitempty"`
	DataName  string       `json:"data_name,omitempty"`
	DataValue []byte       `json:"data_value,omitempty"`
	Ext       DataEntryExt `json:"ext,omitempty"`
}

type DataEntryExt struct {
	V int32 `json:"v,omitempty"`
}

type TrustLineEntry struct {
	AccountId PublicKey         `json:"account_id,omitempty"`
	Asset     TrustLineAsset    `json:"asset,omitempty"`
	Balance   int64             `json:"balance,omitempty"`
	Limit     int64             `json:"limit,omitempty"`
	Flags     uint32            `json:"flags,omitempty"`
	Ext       TrustLineEntryExt `json:"ext,omitempty"`
}

type TrustLineEntryExt struct {
	V  int32             `json:"v,omitempty"`
	V1 *TrustLineEntryV1 `json:"v1,omitempty"`
}

type TrustLineEntryV1 struct {
	Liabilities Liabilities         `json:"liabilities,omitempty"`
	Ext         TrustLineEntryV1Ext `json:"ext,omitempty"`
}

type TrustLineEntryV1Ext struct {
	V  int32                      `json:"v,omitempty"`
	V2 *TrustLineEntryExtensionV2 `json:"v2,omitempty"`
}

type TrustLineEntryExtensionV2 struct {
	LiquidityPoolUseCount int32                        `json:"liquidity_pool_use_count,omitempty"`
	Ext                   TrustLineEntryExtensionV2Ext `json:"ext,omitempty"`
}

type TrustLineEntryExtensionV2Ext struct {
	V int32 `json:"v,omitempty"`
}

type AccountEntry struct {
	AccountId     PublicKey       `json:"account_id,omitempty"`
	Balance       int64           `json:"balance,omitempty"`
	SeqNum        int64           `json:"seq_num,omitempty"`
	NumSubEntries uint32          `json:"num_sub_entries,omitempty"`
	InflationDest *PublicKey      `json:"inflation_dest,omitempty"`
	Flags         uint32          `json:"flags,omitempty"`
	HomeDomain    string          `json:"home_domain,omitempty"`
	Thresholds    []byte          `json:"thresholds,omitempty"`
	Signers       []Signer        `json:"signers,omitempty"`
	Ext           AccountEntryExt `json:"ext,omitempty"`
}

type AccountEntryExt struct {
	V  int32                    `json:"v,omitempty"`
	V1 *AccountEntryExtensionV1 `json:"v1,omitempty"`
}

type AccountEntryExtensionV1 struct {
	Liabilities Liabilities                `json:"liabilities,omitempty"`
	Ext         AccountEntryExtensionV1Ext `json:"ext,omitempty"`
}

type AccountEntryExtensionV1Ext struct {
	V  int32                    `json:"v,omitempty"`
	V2 *AccountEntryExtensionV2 `json:"v2,omitempty"`
}

type AccountEntryExtensionV2 struct {
	NumSponsored        uint32                     `json:"num_sponsored,omitempty"`
	NumSponsoring       uint32                     `json:"num_sponsoring,omitempty"`
	SignerSponsoringIDs []PublicKey                `json:"signer_sponsoring_ids,omitempty"`
	Ext                 AccountEntryExtensionV2Ext `json:"ext,omitempty"`
}

type AccountEntryExtensionV2Ext struct {
	V  int32                    `json:"v,omitempty"`
	V3 *AccountEntryExtensionV3 `json:"v3,omitempty"`
}

type AccountEntryExtensionV3 struct {
	Ext       ExtensionPoint `json:"ext,omitempty"`
	SeqLedger uint32         `json:"seq_ledger,omitempty"`
	SeqTime   uint64         `json:"seq_time,omitempty"`
}

type Liabilities struct {
	Buying  int64 `json:"buying,omitempty"`
	Selling int64 `json:"selling,omitempty"`
}

type LedgerEntryExt struct {
	V  int32                   `json:"v,omitempty"`
	V1 *LedgerEntryExtensionV1 `json:"v1,omitempty"`
}

type LedgerEntryExtensionV1 struct {
	SponsoringId PublicKey                 `json:"sponsoring_id,omitempty"`
	Ext          LedgerEntryExtensionV1Ext `json:"ext,omitempty"`
}

type LedgerEntryExtensionV1Ext struct {
	V int32 `json:"v,omitempty"`
}

type TransactionMeta struct {
	V          int32              `json:"v,omitempty"`
	Operations *[]OperationMeta   `json:"operations,omitempty"`
	V1         *TransactionMetaV1 `json:"v1,omitempty"`
	V2         *TransactionMetaV2 `json:"v2,omitempty"`
	V3         *TransactionMetaV3 `json:"v3,omitempty"`
}

type OperationMeta struct {
	Changes LedgerEntryChanges `json:"changes,omitempty"`
}

type TransactionMetaV1 struct {
	TxChanges  LedgerEntryChanges `json:"tx_changes,omitempty"`
	Operations []OperationMeta    `json:"operations,omitempty"`
}

type TransactionMetaV2 struct {
	TxChangesBefore LedgerEntryChanges `json:"tx_changes_before,omitempty"`
	Operations      []OperationMeta    `json:"operations,omitempty"`
	TxChangesAfter  LedgerEntryChanges `json:"tx_changes_after,omitempty"`
}

type TransactionMetaV3 struct {
	Ext             ExtensionPoint          `json:"ext,omitempty"`
	TxChangesBefore LedgerEntryChanges      `json:"tx_changes_before,omitempty"`
	Operations      []OperationMeta         `json:"operations,omitempty"`
	TxChangesAfter  LedgerEntryChanges      `json:"tx_changes_after,omitempty"`
	SorobanMeta     *SorobanTransactionMeta `json:"soroban_meta,omitempty"`
}

type SorobanTransactionMeta struct {
	Ext              ExtensionPoint    `json:"ext,omitempty"`
	Events           []ContractEvent   `json:"events,omitempty"`
	ReturnValue      ScVal             `json:"return_value,omitempty"`
	DiagnosticEvents []DiagnosticEvent `json:"diagnostic_events,omitempty"`
}

type DiagnosticEvent struct {
	InSuccessfulContractCall bool          `json:"in_successful_contract_call,omitempty"`
	Event                    ContractEvent `json:"event,omitempty"`
}

type ContractEvent struct {
	Ext               ExtensionPoint `json:"ext,omitempty"`
	ContractId        *string        `json:"contract_id,omitempty"`
	ContractEventType int32          `json:"contract_event_type,omitempty"`
	EventType         string         `json:"event_type,omitempty"`
	Transfer          *TransferEvent `json:"transfer,omitempty"`
	Mint              *MintEvent     `json:"mint,omitempty"`
	Clawback          *ClawbackEvent `json:"claw_back,omitempty"`
	Burn              *BurnEvent     `json:"burn,omitempty"`
}

type TransferEvent struct {
	From   string      `json:"from,omitempty"`
	To     string      `json:"to,omitempty"`
	Amount Int128Parts `json:"amount,omitempty"`
}

type MintEvent struct {
	Admin  string      `json:"admin,omitempty"`
	To     string      `json:"to,omitempty"`
	Amount Int128Parts `json:"amount,omitempty"`
}

type ClawbackEvent struct {
	Admin  string      `json:"admin,omitempty"`
	From   string      `json:"from,omitempty"`
	Amount Int128Parts `json:"amount,omitempty"`
}

type BurnEvent struct {
	From   string      `json:"from,omitempty"`
	Amount Int128Parts `json:"amount,omitempty"`
}

type TransactionResultMeta struct {
	Result            TransactionResultPair `json:"result,omitempty"`
	FeeProcessing     LedgerEntryChanges    `json:"fee_processing,omitempty"`
	TxApplyProcessing TransactionMeta       `json:"tx_apply_processing,omitempty"`
}
