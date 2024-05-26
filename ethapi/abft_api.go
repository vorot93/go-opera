package ethapi

import (
	"context"

	"github.com/Fantom-foundation/go-opera/opera"
	"github.com/Fantom-foundation/lachesis-base/inter/idx"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
)

// PublicAbftAPI provides an API to access consensus related information.
// It offers only methods that operate on public data that is freely available to anyone.
type PublicAbftAPI struct {
	b Backend
}

// NewPublicAbftAPI creates a new SFC protocol API.
func NewPublicAbftAPI(b Backend) *PublicAbftAPI {
	return &PublicAbftAPI{b}
}

type EventInfo struct {
	ID           common.Hash    `json:"id"`
	GasPowerLeft GasPowerLeft   `json:"gasPowerLeft"`
	Time         hexutil.Uint64 `json:"time"`
}

type ValidatorEpochState struct {
	GasRefund      hexutil.Uint64 `json:"gasRefund"`
	PrevEpochEvent EventInfo      `json:"prevEpochEvent"`
}

type DagRules struct {
	MaxParents     hexutil.Uint64 `json:"maxParents"`
	MaxFreeParents hexutil.Uint64 `json:"maxFreeParents"`
	MaxExtraData   hexutil.Uint64 `json:"maxExtraData"`
}

type EpochRules struct {
	MaxEpochGas      hexutil.Uint64 `json:"maxEpochGas"`
	MaxEpochDuration hexutil.Uint64 `json:"maxEpochDuration"`
}

type BlocksRules struct {
	MaxBlockGas             hexutil.Uint64 `json:"maxBlockGas"`
	MaxEmptyBlockSkipPeriod hexutil.Uint64 `json:"maxEmptyBlockSkipPeriod"`
}

type RulesRLP struct {
	Name      string         `json:"name"`
	NetworkID hexutil.Uint64 `json:"networkId"`
	Dag       DagRules       `json:"dag"`
	Epochs    EpochRules     `json:"epochs"`
	Blocks    BlocksRules    `json:"blocks"`
	Economy   EconomyRules   `json:"economy"`
}

type GasRules struct {
	MaxEventGas          hexutil.Uint64 `json:"maxEventGas"`
	EventGas             hexutil.Uint64 `json:"eventGas"`
	ParentGas            hexutil.Uint64 `json:"parentGas"`
	ExtraDataGas         hexutil.Uint64 `json:"extraDataGas"`
	BlockVotesBaseGas    hexutil.Uint64 `json:"blockVotesBaseGas"`
	BlockVoteGas         hexutil.Uint64 `json:"blockVoteGas"`
	EpochVoteGas         hexutil.Uint64 `json:"epochVoteGas"`
	MisbehaviourProofGas hexutil.Uint64 `json:"misbehaviourProofGas"`
}

type EconomyRules struct {
	BlockMissedSlack hexutil.Uint64 `json:"blockMissedSlack"`
	Gas              GasRules       `json:"gas"`
	MinGasPrice      *hexutil.Big   `json:"minGasPrice"`
	ShortGasPower    GasPowerRules  `json:"shortGasPower"`
	LongGasPower     GasPowerRules  `json:"longGasPower"`
}

type GasPowerRules struct {
	AllocPerSec        hexutil.Uint64 `json:"allocPerSec"`
	MaxAllocPeriod     hexutil.Uint64 `json:"maxAllocPeriod"`
	StartupAllocPeriod hexutil.Uint64 `json:"startupAllocPeriod"`
	MinStartupGas      hexutil.Uint64 `json:"minStartupGas"`
}

type Rules struct {
	Rlp      RulesRLP       `json:"rlp"`
	Upgrades opera.Upgrades `json:"upgrades"`
}

type GasPowerLeft [2]hexutil.Uint64

func (s *PublicAbftAPI) GetFullEpochState(ctx context.Context, epoch rpc.BlockNumber) (map[string]interface{}, error) {
	_, es, err := s.b.GetEpochBlockState(ctx, epoch)
	if err != nil {
		return nil, err
	}
	if es == nil {
		return nil, nil
	}

	validatorProfiles := map[hexutil.Uint64]interface{}{}
	validators := make(map[hexutil.Uint64]hexutil.Uint64)
	profiles := es.ValidatorProfiles
	for _, vid := range es.Validators.IDs() {
		validatorProfiles[hexutil.Uint64(vid)] = map[string]interface{}{
			"weight": (*hexutil.Big)(profiles[vid].Weight),
			"pubkey": profiles[vid].PubKey.String(),
		}
		validators[hexutil.Uint64(vid)] = hexutil.Uint64(es.Validators.Get(vid))
	}

	validatorState := make([]ValidatorEpochState, len(es.ValidatorStates))
	for i, vs := range es.ValidatorStates {
		gas_power_left := GasPowerLeft{}
		for j, gas := range vs.PrevEpochEvent.GasPowerLeft.Gas {
			gas_power_left[j] = hexutil.Uint64(gas)
		}
		validatorState[i] = ValidatorEpochState{
			GasRefund: hexutil.Uint64(vs.GasRefund),
			PrevEpochEvent: EventInfo{
				ID:           common.Hash(vs.PrevEpochEvent.ID),
				GasPowerLeft: gas_power_left,
				Time:         hexutil.Uint64(vs.PrevEpochEvent.Time),
			},
		}
	}

	r := es.Rules
	rules := Rules{
		Rlp: RulesRLP{
			Name:      r.Name,
			NetworkID: hexutil.Uint64(r.NetworkID),
			Dag: DagRules{
				MaxParents:     hexutil.Uint64(r.Dag.MaxParents),
				MaxFreeParents: hexutil.Uint64(r.Dag.MaxFreeParents),
				MaxExtraData:   hexutil.Uint64(r.Dag.MaxExtraData),
			},
			Epochs: EpochRules{
				MaxEpochGas:      hexutil.Uint64(r.Epochs.MaxEpochGas),
				MaxEpochDuration: hexutil.Uint64(r.Epochs.MaxEpochDuration),
			},
			Blocks: BlocksRules{
				MaxBlockGas:             hexutil.Uint64(r.Blocks.MaxBlockGas),
				MaxEmptyBlockSkipPeriod: hexutil.Uint64(r.Blocks.MaxEmptyBlockSkipPeriod),
			},
			Economy: EconomyRules{
				BlockMissedSlack: hexutil.Uint64(r.Economy.BlockMissedSlack),
				Gas: GasRules{
					MaxEventGas:          hexutil.Uint64(r.Economy.Gas.MaxEventGas),
					EventGas:             hexutil.Uint64(r.Economy.Gas.EventGas),
					ParentGas:            hexutil.Uint64(r.Economy.Gas.ParentGas),
					ExtraDataGas:         hexutil.Uint64(r.Economy.Gas.ExtraDataGas),
					BlockVotesBaseGas:    hexutil.Uint64(r.Economy.Gas.BlockVotesBaseGas),
					BlockVoteGas:         hexutil.Uint64(r.Economy.Gas.BlockVoteGas),
					EpochVoteGas:         hexutil.Uint64(r.Economy.Gas.EpochVoteGas),
					MisbehaviourProofGas: hexutil.Uint64(r.Economy.Gas.MisbehaviourProofGas),
				},
				MinGasPrice: (*hexutil.Big)(r.Economy.MinGasPrice),
				ShortGasPower: GasPowerRules{
					AllocPerSec:        hexutil.Uint64(r.Economy.ShortGasPower.AllocPerSec),
					MaxAllocPeriod:     hexutil.Uint64(r.Economy.ShortGasPower.MaxAllocPeriod),
					StartupAllocPeriod: hexutil.Uint64(r.Economy.ShortGasPower.StartupAllocPeriod),
					MinStartupGas:      hexutil.Uint64(r.Economy.ShortGasPower.MinStartupGas),
				},
				LongGasPower: GasPowerRules{
					AllocPerSec:        hexutil.Uint64(r.Economy.LongGasPower.AllocPerSec),
					MaxAllocPeriod:     hexutil.Uint64(r.Economy.LongGasPower.MaxAllocPeriod),
					StartupAllocPeriod: hexutil.Uint64(r.Economy.LongGasPower.StartupAllocPeriod),
					MinStartupGas:      hexutil.Uint64(r.Economy.LongGasPower.MinStartupGas),
				},
			},
		},
		Upgrades: es.Rules.Upgrades,
	}

	res := map[string]interface{}{
		"epoch":             hexutil.Uint64(epoch),
		"epochStart":        hexutil.Uint64(es.EpochStart),
		"prevEpochStart":    hexutil.Uint64(es.PrevEpochStart),
		"epochStateRoot":    es.EpochStateRoot.String(),
		"validators":        validators,
		"validatorStates":   validatorState,
		"validatorProfiles": validatorProfiles,
		"rules":             rules,
	}

	return res, nil

}

func (s *PublicAbftAPI) GetValidators(ctx context.Context, epoch rpc.BlockNumber) (map[hexutil.Uint64]interface{}, error) {
	bs, es, err := s.b.GetEpochBlockState(ctx, epoch)
	if err != nil {
		return nil, err
	}
	if es == nil {
		return nil, nil
	}
	res := map[hexutil.Uint64]interface{}{}
	for _, vid := range es.Validators.IDs() {
		profiles := es.ValidatorProfiles
		if epoch == rpc.PendingBlockNumber {
			profiles = bs.NextValidatorProfiles
		}
		res[hexutil.Uint64(vid)] = map[string]interface{}{
			"weight": (*hexutil.Big)(profiles[vid].Weight),
			"pubkey": profiles[vid].PubKey.String(),
		}
	}
	return res, nil
}

// GetDowntime returns validator's downtime.
func (s *PublicAbftAPI) GetDowntime(ctx context.Context, validatorID hexutil.Uint) (map[string]interface{}, error) {
	blocks, period, err := s.b.GetDowntime(ctx, idx.ValidatorID(validatorID))
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"offlineBlocks": hexutil.Uint64(blocks),
		"offlineTime":   hexutil.Uint64(period),
	}, nil
}

// GetEpochUptime returns validator's epoch uptime in nanoseconds.
func (s *PublicAbftAPI) GetEpochUptime(ctx context.Context, validatorID hexutil.Uint) (hexutil.Uint64, error) {
	v, err := s.b.GetUptime(ctx, idx.ValidatorID(validatorID))
	if err != nil {
		return 0, err
	}
	if v == nil {
		return 0, nil
	}
	return hexutil.Uint64(v.Uint64()), nil
}

// GetOriginatedEpochFee returns validator's originated epoch fee.
func (s *PublicAbftAPI) GetOriginatedEpochFee(ctx context.Context, validatorID hexutil.Uint) (*hexutil.Big, error) {
	v, err := s.b.GetOriginatedFee(ctx, idx.ValidatorID(validatorID))
	if err != nil {
		return nil, err
	}
	return (*hexutil.Big)(v), nil
}
