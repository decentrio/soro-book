package xdr

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
