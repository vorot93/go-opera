package ethapi

import (
	"context"

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
	profiles := es.ValidatorProfiles
	for _, vid := range es.Validators.IDs() {
		validatorProfiles[hexutil.Uint64(vid)] = map[string]interface{}{
			"weight": (*hexutil.Big)(profiles[vid].Weight),
			"pubkey": profiles[vid].PubKey.String(),
		}
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

	validators := make(map[hexutil.Uint64]hexutil.Uint64, len(es.Validators.Values))
	for k, v := range es.Validators.Values {
		validators[hexutil.Uint64(k)] = hexutil.Uint64(v)
	}

	r := es.Rules
	rules := map[string]interface{}{
		"rlp": map[string]interface{}{
			"name":      r.Name,
			"networkId": hexutil.Uint64(r.NetworkID),
			"dag": map[string]interface{}{
				"maxParents":     hexutil.Uint64(r.Dag.MaxParents),
				"maxFreeParents": hexutil.Uint64(r.Dag.MaxFreeParents),
				"maxExtraData":   hexutil.Uint64(r.Dag.MaxExtraData),
			},
			"epochs": map[string]interface{}{
				"maxEpochGas":      hexutil.Uint64(r.Epochs.MaxEpochGas),
				"maxEpochDuration": hexutil.Uint64(r.Epochs.MaxEpochDuration),
			},
			"blocks": map[string]interface{}{
				"maxBlockGas":              hexutil.Uint64(r.Blocks.MaxBlockGas),
				"maxEmptyBlockSkipPeriods": hexutil.Uint64(r.Blocks.MaxEmptyBlockSkipPeriod),
			},
		},
		"upgrades": r.Upgrades,
	}

	_ = rules

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
