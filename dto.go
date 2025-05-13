package main

type TransfersListResponse struct {
	Code        int    `json:"code"`
	Data        Data   `json:"data"`
	GeneratedAt int64  `json:"generated_at"`
	Message     string `json:"message"`
}

type TransfersListRequest struct {
	Address        *string   `json:"address"`
	AfterID        []*int    `json:"after_id"`
	AssetSymbol    *string   `json:"asset_symbol"`
	AssetUniqueID  *string   `json:"asset_unique_id"`
	BlockRange     *string   `json:"block_range"`
	Currency       *string   `json:"currency"`
	Direction      *string   `json:"direction"`
	ExtrinsicIndex *string   `json:"extrinsic_index"`
	FilterNFT      *bool     `json:"filter_nft"`
	IncludeTotal   *bool     `json:"include_total"`
	ItemID         *int      `json:"item_id"`
	MaxAmount      *string   `json:"max_amount"`
	MinAmount      *string   `json:"min_amount"`
	Order          *string   `json:"order"`
	Page           *int      `json:"page"`
	Row            *int      `json:"row"`
	Success        *bool     `json:"success"`
	Timeout        *int      `json:"timeout"`
	TokenCategory  []*string `json:"token_category"`
}

type SymbolTokenListResponse struct {
	Code        *int           `json:"code"`
	Data        *TokenListData `json:"data"`
	GeneratedAt *int64         `json:"generated_at"`
	Message     *string        `json:"message"`
}

type AccountResponse struct {
	Code        *int                 `json:"code,omitempty"`
	Message     *string              `json:"message,omitempty"`
	GeneratedAt *int64               `json:"generated_at,omitempty"`
	Data        *AccountResponseData `json:"data,omitempty"`
}

type AccountResponseData struct {
	Account *Account `json:"account,omitempty"`
}

type Account struct {
	Address            *string                 `json:"address,omitempty"`
	Balance            *string                 `json:"balance,omitempty"`
	Lock               *string                 `json:"lock,omitempty"`
	BalanceLock        *string                 `json:"balance_lock,omitempty"`
	IsEvmContract      *bool                   `json:"is_evm_contract,omitempty"`
	AccountDisplay     *AccountDisplay1        `json:"account_display,omitempty"`
	EvmAccount         *string                 `json:"evm_account,omitempty"`
	CountExtrinsic     *int                    `json:"count_extrinsic,omitempty"`
	NftAmount          *string                 `json:"nft_amount,omitempty"`
	Display            *string                 `json:"display,omitempty"`
	Web                *string                 `json:"web,omitempty"`
	Riot               *string                 `json:"riot,omitempty"`
	Email              *string                 `json:"email,omitempty"`
	Legal              *string                 `json:"legal,omitempty"`
	Twitter            *string                 `json:"twitter,omitempty"`
	Github             *string                 `json:"github,omitempty"`
	Matrix             *string                 `json:"matrix,omitempty"`
	Discord            *string                 `json:"discord,omitempty"`
	Reserved           *string                 `json:"reserved,omitempty"`
	Bonded             *string                 `json:"bonded,omitempty"`
	Unbonding          *string                 `json:"unbonding,omitempty"`
	DemocracyLock      *string                 `json:"democracy_lock,omitempty"`
	ConvictionLock     *string                 `json:"conviction_lock,omitempty"`
	ElectionLock       *string                 `json:"election_lock,omitempty"`
	Nonce              *int                    `json:"nonce,omitempty"`
	Role               *string                 `json:"role,omitempty"`
	Stash              *string                 `json:"stash,omitempty"`
	IsCouncilMember    *bool                   `json:"is_council_member,omitempty"`
	IsTechcommMember   *bool                   `json:"is_techcomm_member,omitempty"`
	IsRegistrar        *bool                   `json:"is_registrar,omitempty"`
	IsFellowshipMember *bool                   `json:"is_fellowship_member,omitempty"`
	IsModuleAccount    *bool                   `json:"is_module_account,omitempty"`
	IsERC20            *bool                   `json:"is_erc20,omitempty"`
	IsERC721           *bool                   `json:"is_erc721,omitempty"`
	Proxy              *map[string]interface{} `json:"proxy,omitempty"`
	Multisig           *map[string]interface{} `json:"multisig,omitempty"`
}

type AccountDisplay1 struct {
	Address *string                 `json:"address,omitempty"`
	People  *map[string]interface{} `json:"people,omitempty"`
}

type TokenListData struct {
	Detail map[string]*TokenDetail `json:"detail"`
	Token  []*string               `json:"token"`
}

type TokenDetail struct {
	AssetType              *string       `json:"asset_type"`
	AvailableBalance       *string       `json:"available_balance"`
	BondedLockedBalance    *string       `json:"bonded_locked_balance"`
	ConvictionLockBalance  *string       `json:"conviction_lock_balance"`
	DemocracyLockedBalance *string       `json:"democracy_locked_balance"`
	DisplayName            *string       `json:"display_name"`
	ElectionLockedBalance  *string       `json:"election_locked_balance"`
	ExternalData           *ExternalData `json:"external_data"`
	FreeBalance            *string       `json:"free_balance"`
	Inflation              *string       `json:"inflation"`
	LockedBalance          *string       `json:"locked_balance"`
	NominatorBonded        *string       `json:"nominator_bonded"`
	Price                  *string       `json:"price"`
	PriceChange            *string       `json:"price_change"`
	ReservedBalance        *string       `json:"reserved_balance"`
	Symbol                 *string       `json:"symbol"`
	TokenDecimals          *int          `json:"token_decimals"`
	TotalIssuance          *string       `json:"total_issuance"`
	TreasuryBalance        *string       `json:"treasury_balance"`
	UnbondedLockedBalance  *string       `json:"unbonded_locked_balance"`
	UniqueID               *string       `json:"unique_id"`
	ValidatorBonded        *string       `json:"validator_bonded"`
	VestingBalance         *string       `json:"vesting_balance"`
}

type ExternalData struct {
	AuthorizationSource *string `json:"authorization_source"`
	CirculatingSupply   *string `json:"circulating_supply"`
	Source              *string `json:"source"`
}

type Data struct {
	Count     int                 `json:"count"`
	Total     map[string]Property `json:"total"`
	Transfers []Transfer          `json:"transfers"`
}

type Property struct {
	Received string `json:"received"`
	Sent     string `json:"sent"`
	Total    string `json:"total"`
}

type Transfer struct {
	Amount                *string         `json:"amount"`
	AmountV2              *string         `json:"amount_v2"`
	AssetSymbol           *string         `json:"asset_symbol"`
	AssetType             *string         `json:"asset_type"`
	AssetUniqueID         *string         `json:"asset_unique_id"`
	BlockNum              *int            `json:"block_num"`
	BlockTimestamp        *int64          `json:"block_timestamp"`
	CurrencyAmount        *string         `json:"currency_amount"`
	CurrentCurrencyAmount *string         `json:"current_currency_amount"`
	EventIdx              *int            `json:"event_idx"`
	ExtrinsicIndex        *string         `json:"extrinsic_index"`
	Fee                   *string         `json:"fee"`
	From                  *string         `json:"from"`
	FromAccountDisplay    *AccountDisplay `json:"from_account_display"`
	Hash                  *string         `json:"hash"`
	IsLock                *bool           `json:"is_lock"`
	ItemDetail            *ItemDetail     `json:"item_detail"`
	ItemID                *string         `json:"item_id"`
	Module                *string         `json:"module"`
	Nonce                 *int            `json:"nonce"`
	Success               *bool           `json:"success"`
	To                    *string         `json:"to"`
	ToAccountDisplay      *AccountDisplay `json:"to_account_display"`
	TransferID            *int            `json:"transfer_id"`
}

type AccountDisplay struct {
	AccountIndex *string      `json:"account_index"`
	Address      *string      `json:"address"`
	Display      *string      `json:"display"`
	EvmAddress   *string      `json:"evm_address"`
	EvmContract  *EvmContract `json:"evm_contract"`
	Identity     *bool        `json:"identity"`
	Judgements   []*Judgement `json:"judgements"`
	Merkle       *Merkle      `json:"merkle"`
	Parent       *Parent      `json:"parent"`
	People       *People      `json:"people"`
}

type EvmContract struct {
	ContractName *string `json:"contract_name"`
	VerifySource *string `json:"verify_source"`
}

type Judgement struct {
	Index     *int    `json:"index"`
	Judgement *string `json:"judgement"`
}

type Merkle struct {
	AddressType *string `json:"address_type"`
	TagName     *string `json:"tag_name"`
	TagSubtype  *string `json:"tag_subtype"`
	TagType     *string `json:"tag_type"`
}

type Parent struct {
	Address   *string `json:"address"`
	Display   *string `json:"display"`
	Identity  *bool   `json:"identity"`
	SubSymbol *string `json:"sub_symbol"`
}

type People struct {
	Display    *string      `json:"display"`
	Identity   *bool        `json:"identity"`
	Judgements []*Judgement `json:"judgements"`
	Parent     *Parent      `json:"parent"`
}

type ItemDetail struct {
	CollectionSymbol *string  `json:"collection_symbol"`
	FallbackImage    *string  `json:"fallback_image"`
	Image            *string  `json:"image"`
	LocalImage       *string  `json:"local_image"`
	Media            []*Media `json:"media"`
	Symbol           *string  `json:"symbol"`
	Thumbnail        *string  `json:"thumbnail"`
}

type Media struct {
	Types *string `json:"types"`
	URL   *string `json:"url"`
}
