package xdr

type Hash []byte
type PoolId Hash
type ScSymbol string

type Envelope struct {
	V0 *TransactionV0Envelope `json:"v0,omitempty"`
	V1 *TransactionV1Envelope `json:"v1,omitempty"`
}

type TransactionV0Envelope struct {
	Tx         TransactionV0        `json:"tx,omitempty"`
	Signatures []DecoratedSignature `json:"signatures,omitempty"`
}

type TransactionV0 struct {
	SourceAccountEd25519 []byte           `json:"source_account_ed25519,omitempty"`
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
	Ed25519  *[]byte               `json:"ed25519,omitempty"`
	Med25519 *MuxedAccountMed25519 `json:"med25519,omitempty"`
}

type MuxedAccountMed25519 struct {
	Id      uint64  `json:"id,omitempty"`
	Ed25519 *[]byte `json:"ed25519,omitempty"`
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
	Ed25519 []byte `json:"ed25519,omitempty"`
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
	Selling Asset `json:"selling,omitempty"`
	Buying  Asset `json:"buying,omitempty"`
	Amount  int64 `json:"amount,omitempty"`
	Price   Price `json:"price,omitempty"`
	OfferId int64 `json:"offer_id,omitempty"`
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
	Ed25519              *[]byte                        `json:"ed25519,omitempty"`
	PreAuthTx            *[]byte                        `json:"pre_auth_tx,omitempty"`
	HashX                *[]byte                        `json:"hash_x,omitempty"`
	Ed25519SignedPayload *SignerKeyEd25519SignedPayload `json:"ed25519_signed_payload,omitempty"`
}

type SignerKeyEd25519SignedPayload struct {
	Ed25519 []byte `json:"ed25519,omitempty"`
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
	Asset     []byte    `json:"asset,omitempty"`
	Authorize uint32    `json:"authorize,omitempty"`
}

type ManageDataOp struct {
	DataName  string `json:"data_name,omitempty"`
	DataValue []byte `json:"data_value,omitempty"`
}

type BumpSequenceOp struct {
	BumpTo int64
}

type ManageBuyOfferOp struct {
	Selling   Asset
	Buying    Asset
	BuyAmount int64
	Price     Price
	OfferId   int64
}

type PathPaymentStrictSendOp struct {
	SendAsset   Asset
	SendAmount  int64
	Destination MuxedAccount
	DestAsset   Asset
	DestMin     int64
	Path        []Asset
}

type CreateClaimableBalanceOp struct {
	Asset     Asset
	Amount    int64
	Claimants []Claimant `xdrmaxsize:"10"`
}

type Claimant struct {
	V0 *ClaimantV0
}

type ClaimantV0 struct {
	Destination PublicKey
	Predicate   ClaimPredicate
}

type ClaimPredicate struct {
	AndPredicates *[]ClaimPredicate `xdrmaxsize:"2"`
	OrPredicates  *[]ClaimPredicate `xdrmaxsize:"2"`
	NotPredicate  **ClaimPredicate
	AbsBefore     *int64
	RelBefore     *int64
}

type ClaimClaimableBalanceOp struct {
	BalanceId ClaimableBalanceId
}

type ClaimableBalanceId struct {
	V0 *string
}

type BeginSponsoringFutureReservesOp struct {
	SponsoredId PublicKey
}

type RevokeSponsorshipOp struct {
	LedgerKey *LedgerKey
	Signer    *RevokeSponsorshipOpSigner
}

type RevokeSponsorshipOpSigner struct {
	AccountId PublicKey
	SignerKey SignerKey
}

type LedgerKey struct {
	Account          *LedgerKeyAccount
	TrustLine        *LedgerKeyTrustLine
	Offer            *LedgerKeyOffer
	Data             *LedgerKeyData
	ClaimableBalance *LedgerKeyClaimableBalance
	LiquidityPool    *LedgerKeyLiquidityPool
	ContractData     *LedgerKeyContractData
	ContractCode     *LedgerKeyContractCode
	ConfigSetting    *LedgerKeyConfigSetting
	Ttl              *LedgerKeyTtl
}

type LedgerKeyAccount struct {
	AccountId PublicKey
}

type LedgerKeyTrustLine struct {
	AccountId PublicKey
	Asset     TrustLineAsset
}

type TrustLineAsset struct {
	Asset           *Asset  `json:"asset,omitempty"`
	LiquidityPoolId *PoolId `json:"liquidity_pool_id,omitempty"`
}

type LedgerKeyOffer struct {
	SellerId PublicKey
	OfferId  int64
}

type LedgerKeyData struct {
	AccountId PublicKey
	DataName  string
}

type LedgerKeyClaimableBalance struct {
	BalanceId ClaimableBalanceId
}

type LedgerKeyLiquidityPool struct {
	LiquidityPoolId PoolId
}

type LedgerKeyContractData struct {
	Contract   ScAddress
	Key        ScVal
	Durability int32
}

type ScAddress struct {
	AccountId  *PublicKey
	ContractId *Hash
}

type ScVal struct {
	B         *bool
	Error     *ScError
	U32       *uint32
	I32       *int32
	U64       *uint64
	I64       *int64
	Timepoint *uint64
	Duration  *uint64
	U128      *UInt128Parts
	I128      *Int128Parts
	U256      *UInt256Parts
	I256      *Int256Parts
	Bytes     *ScBytes
	Str       *string
	Sym       *ScSymbol
	Vec       **ScVal
	Map       **ScMap
	Address   *ScAddress
	NonceKey  *ScNonceKey
	Instance  *ScContractInstance
}

type ContractExecutable struct {
	WasmHash *Hash
}

type ScContractInstance struct {
	Executable ContractExecutable
	Storage    *ScMap
}

type ScNonceKey struct {
	Nonce int64
}

type ScMap []ScMapEntry

type ScMapEntry struct {
	Key ScVal
	Val ScVal
}

type UInt128Parts struct {
	Hi uint64
	Lo uint64
}

type Int128Parts struct {
	Hi int64  `json:"hi,omitempty"`
	Lo uint64 `json:"lo,omitempty"`
}

type UInt256Parts struct {
	HiHi uint64
	HiLo uint64
	LoHi uint64
	LoLo uint64
}

type Int256Parts struct {
	HiHi int64
	HiLo uint64
	LoHi uint64
	LoLo uint64
}

type ScBytes []byte

type ScError struct {
	ContractCode *uint32
	Code         *uint32
}

type LedgerKeyContractCode struct {
	Hash Hash
}

type LedgerKeyConfigSetting struct {
	ConfigSettingId int32
}

type LedgerKeyTtl struct {
	KeyHash Hash
}

type ClawbackOp struct {
	Asset  Asset
	From   MuxedAccount
	Amount int64
}

type ClawbackClaimableBalanceOp struct {
	BalanceId ClaimableBalanceId
}

type SetTrustLineFlagsOp struct {
	Trustor    PublicKey
	Asset      Asset
	ClearFlags uint32
	SetFlags   uint32
}

type LiquidityPoolDepositOp struct {
	LiquidityPoolId PoolId
	MaxAmountA      int64
	MaxAmountB      int64
	MinPrice        Price
	MaxPrice        Price
}

type LiquidityPoolWithdrawOp struct {
	LiquidityPoolId PoolId
	Amount          int64
	MinAmountA      int64
	MinAmountB      int64
}

type InvokeHostFunctionOp struct {
	HostFunction HostFunction
	Auth         []SorobanAuthorizationEntry
}

type SorobanAuthorizationEntry struct {
	Credentials    SorobanCredentials
	RootInvocation SorobanAuthorizedInvocation
}

type SorobanCredentials struct {
	Address *int32
}

type SorobanAuthorizedInvocation struct {
	Function       SorobanAuthorizedFunction
	SubInvocations []SorobanAuthorizedInvocation
}

type SorobanAuthorizedFunction struct {
	ContractFn           *InvokeContractArgs
	CreateContractHostFn *CreateContractArgs
}

type HostFunction struct {
	InvokeContract *InvokeContractArgs
	CreateContract *CreateContractArgs
	Wasm           *[]byte
}

type InvokeContractArgs struct {
	ContractAddress ScAddress
	FunctionName    ScSymbol
	Args            []ScVal
}

type CreateContractArgs struct {
	ContractIdPreimage ContractIdPreimage
	Executable         ContractExecutable
}

type ContractIdPreimage struct {
	FromAddress *ContractIdPreimageFromAddress
	FromAsset   *Asset
}

type ContractIdPreimageFromAddress struct {
	Address ScAddress
	Salt    []byte
}

type ExtendFootprintTtlOp struct {
	Ext      ExtensionPoint
	ExtendTo uint32
}

type RestoreFootprintOp struct {
	Ext ExtensionPoint
}

type ExtensionPoint struct {
	V int32
}

type TransactionV0Ext struct {
	V int32
}

type TransactionV1Envelope struct {
	Tx         Transaction
	Signatures []DecoratedSignature `xdrmaxsize:"20"`
}

type Transaction struct {
	SourceAccount MuxedAccount
	Fee           uint32
	SeqNum        int64
	Cond          Preconditions
	Memo          Memo
	Operations    []Operation `xdrmaxsize:"100"`
	Ext           TransactionExt
}

type TransactionExt struct {
	SorobanData *SorobanTransactionData
}

type SorobanTransactionData struct {
	Ext         ExtensionPoint
	Resources   SorobanResources
	ResourceFee int64
}

type SorobanResources struct {
	Footprint    LedgerFootprint
	Instructions uint32
	ReadBytes    uint32
	WriteBytes   uint32
}

type LedgerFootprint struct {
	ReadOnly  []LedgerKey
	ReadWrite []LedgerKey
}

type Preconditions struct {
	TimeBounds *TimeBounds
	V2         *PreconditionsV2
}

type PreconditionsV2 struct {
	TimeBounds      *TimeBounds
	LedgerBounds    *LedgerBounds
	MinSeqNum       *int64
	MinSeqAge       uint64
	MinSeqLedgerGap uint32
	ExtraSigners    []SignerKey `xdrmaxsize:"2"`
}

type LedgerBounds struct {
	MinLedger uint32
	MaxLedger uint32
}
