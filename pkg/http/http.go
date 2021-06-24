package http

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/consensus/clique"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/signer/core"
	"github.com/pefish/ether-clef/pkg/global"
	"github.com/pefish/ether-clef/pkg/internal/ethapi"
	"github.com/pefish/ether-clef/version"
	go_crypto "github.com/pefish/go-crypto"
	go_logger "github.com/pefish/go-logger"
	go_mysql "github.com/pefish/go-mysql"
	"github.com/pkg/errors"
	"math/big"
	"mime"
	"strings"
)

// SignerAPI defines the actual implementation of ExternalAPI
type SignerAPI struct {
	chainID *big.Int
}

func NewSignerAPI(chainID int64) *SignerAPI {
	return &SignerAPI{big.NewInt(chainID)}
}

func (api *SignerAPI) List(ctx context.Context) ([]common.Address, error) {
	addresses := make([]common.Address, 0)
	results := make([]struct {
		Address string `json:"address"`
	}, 0)
	err := go_mysql.MysqlInstance.Select(&results, "address", "address", map[string]interface{}{
		"is_ban": 0,
	})
	if err != nil {
		return nil, err
	}
	for _, result := range results {
		addresses = append(addresses, common.HexToAddress(result.Address))
	}
	go_logger.Logger.Debug("List")
	return addresses, nil
}

// New creates a new password protected Account. The private key is protected with
// the given password. Users are responsible to backup the private key that is stored
// in the keystore location thas was specified when this API was created.
func (api *SignerAPI) New(ctx context.Context) (common.Address, error) {
	go_logger.Logger.Debug("New")
	return common.Address{}, errors.New("not implement")
}

// SignTransaction signs the given Transaction and returns it both as json and rlp-encoded form
func (api *SignerAPI) SignTransaction(ctx context.Context, args core.SendTxArgs) (*ethapi.SignTransactionResult, error) {
	go_logger.Logger.DebugF("SignTransaction. args: %#v", args)

	// 验证 method
	var payload *hexutil.Bytes
	if args.Data != nil {
		payload = args.Data
	}
	if args.Input != nil {
		payload = args.Input
	}
	if payload == nil {
		return nil, errors.New("method not be allowed")
	}
	inputStr := payload.String()
	if strings.HasPrefix(inputStr, "0x") {
		inputStr = inputStr[2:]
	}
	if len(inputStr) < 8 {
		return nil, errors.New("method not be allowed")
	}
	methodId := inputStr[:8]
	if _, ok := global.AllowedMethod.Load(methodId); !ok {
		queryResults := make([]struct {
			MethodId string `json:"method_id"`
		}, 0)
		err := go_mysql.MysqlInstance.Select(&queryResults, "method", "method_id", map[string]interface{}{
			"is_ban": 0,
		})
		if err != nil {
			return nil, err
		}
		for _, queryResult := range queryResults {
			global.AllowedMethod.Store(queryResult.MethodId, true)
		}
		if _, ok := global.AllowedMethod.Load(methodId); !ok {
			return nil, errors.New("method not be allowed")
		}
	}
	// 验证 chain id
	if args.ChainID != nil {
		requestedChainId := (*big.Int)(args.ChainID)
		if api.chainID.Cmp(requestedChainId) != 0 {
			return nil, fmt.Errorf("requested chainid %d does not match the configuration of the signer", requestedChainId)
		}
	}

	// Convert fields into a real transaction
	txArgs := ethapi.TransactionArgs{
		Gas:                  &args.Gas,
		GasPrice:             args.GasPrice,
		MaxFeePerGas:         args.MaxFeePerGas,
		MaxPriorityFeePerGas: args.MaxPriorityFeePerGas,
		Value:                &args.Value,
		Nonce:                &args.Nonce,
		Data:                 args.Data,
		Input:                args.Input,
		AccessList:           args.AccessList,
		ChainID:              args.ChainID,
	}
	// Add the To-field, if specified
	if args.To != nil {
		to := args.To.Address()
		txArgs.To = &to
	}
	var unsignedTx = txArgs.ToTransaction()

	privateKeyECDSA, err := api.privateKey(strings.ToLower(args.From.Original()))
	if err != nil {
		return nil, err
	}
	signedTx, err := types.SignTx(unsignedTx, types.NewEIP155Signer(api.chainID), privateKeyECDSA)
	if err != nil {
		return nil, err
	}

	data, err := signedTx.MarshalBinary()
	if err != nil {
		return nil, err
	}
	response := ethapi.SignTransactionResult{Raw: data, Tx: signedTx}

	return &response, nil

}

func (api *SignerAPI) privateKey(from string) (*ecdsa.PrivateKey, error) {
	privateKeyInter, ok := global.Addresses.Load(from)
	privateKey := ""
	if !ok {
		var queryResult struct{
			Priv string `json:"priv"`
		}
		notFound, err := go_mysql.MysqlInstance.SelectFirst(&queryResult, "address", "priv", map[string]interface{}{
			"is_ban": 0,
			"address": from,
		})
		if err != nil {
			return nil, err
		}
		if notFound {
			return nil, errors.New("Not be allowed to sign")
		}
		privateKey = go_crypto.Crypto.MustAesCbcDecrypt(global.Password, queryResult.Priv)
		global.Addresses.Store(from, privateKey)
	} else {
		privateKey = privateKeyInter.(string)
	}
	privateKeyBuf, err := hex.DecodeString(privateKey)
	if err != nil {
		return nil, err
	}
	return crypto.ToECDSA(privateKeyBuf)
}

func (api *SignerAPI) determineSignatureFormat(ctx context.Context, contentType string, addr common.MixedcaseAddress, data interface{}) (*core.SignDataRequest, bool, error) {
	var (
		req          *core.SignDataRequest
		useEthereumV = true
	)
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return nil, useEthereumV, err
	}

	switch mediaType {
	case core.IntendedValidator.Mime:
		// Data with an intended validator
		validatorData, err := core.UnmarshalValidatorData(data)
		if err != nil {
			return nil, useEthereumV, err
		}
		sighash, msg := core.SignTextValidator(validatorData)
		messages := []*core.NameValueType{
			{
				Name:  "This is a request to sign data intended for a particular validator (see EIP 191 version 0)",
				Typ:   "description",
				Value: "",
			},
			{
				Name:  "Intended validator address",
				Typ:   "address",
				Value: validatorData.Address.String(),
			},
			{
				Name:  "Application-specific data",
				Typ:   "hexdata",
				Value: validatorData.Message,
			},
			{
				Name:  "Full message for signing",
				Typ:   "hexdata",
				Value: fmt.Sprintf("0x%x", msg),
			},
		}
		req = &core.SignDataRequest{ContentType: mediaType, Rawdata: []byte(msg), Messages: messages, Hash: sighash}
	case core.ApplicationClique.Mime:
		// Clique is the Ethereum PoA standard
		stringData, ok := data.(string)
		if !ok {
			return nil, useEthereumV, fmt.Errorf("input for %v must be an hex-encoded string", core.ApplicationClique.Mime)
		}
		cliqueData, err := hexutil.Decode(stringData)
		if err != nil {
			return nil, useEthereumV, err
		}
		header := &types.Header{}
		if err := rlp.DecodeBytes(cliqueData, header); err != nil {
			return nil, useEthereumV, err
		}
		// The incoming clique header is already truncated, sent to us with a extradata already shortened
		if len(header.Extra) < 65 {
			// Need to add it back, to get a suitable length for hashing
			newExtra := make([]byte, len(header.Extra)+65)
			copy(newExtra, header.Extra)
			header.Extra = newExtra
		}
		// Get back the rlp data, encoded by us
		if len(header.Extra) < 65 {
			return nil, useEthereumV, fmt.Errorf("clique header extradata too short, %d < 65", len(header.Extra))
		}
		cliqueRlp := clique.CliqueRLP(header)
		sighash := clique.SealHash(header).Bytes()

		messages := []*core.NameValueType{
			{
				Name:  "Clique header",
				Typ:   "clique",
				Value: fmt.Sprintf("clique header %d [0x%x]", header.Number, header.Hash()),
			},
		}
		// Clique uses V on the form 0 or 1
		useEthereumV = false
		req = &core.SignDataRequest{ContentType: mediaType, Rawdata: cliqueRlp, Messages: messages, Hash: sighash}
	default: // also case TextPlain.Mime:
		// Calculates an Ethereum ECDSA signature for:
		// hash = keccak256("\x19${byteVersion}Ethereum Signed Message:\n${message length}${message}")
		// We expect it to be a string
		if stringData, ok := data.(string); !ok {
			return nil, useEthereumV, fmt.Errorf("input for text/plain must be an hex-encoded string")
		} else {
			if textData, err := hexutil.Decode(stringData); err != nil {
				return nil, useEthereumV, err
			} else {
				sighash, msg := accounts.TextAndHash(textData)
				messages := []*core.NameValueType{
					{
						Name:  "message",
						Typ:   accounts.MimetypeTextPlain,
						Value: msg,
					},
				}
				req = &core.SignDataRequest{ContentType: mediaType, Rawdata: []byte(msg), Messages: messages, Hash: sighash}
			}
		}
	}
	req.Address = addr
	req.Meta = core.MetadataFromContext(ctx)
	return req, useEthereumV, nil
}

func (api *SignerAPI) SignData(ctx context.Context, contentType string, addr common.MixedcaseAddress, data interface{}) (hexutil.Bytes, error) {
	var req, transformV, err = api.determineSignatureFormat(ctx, contentType, addr, data)
	if err != nil {
		return nil, err
	}
	privateKeyECDSA, err := api.privateKey(strings.ToLower(addr.Original()))
	if err != nil {
		return nil, err
	}
	hashBuf := crypto.Keccak256(req.Rawdata)
	signature, err := crypto.Sign(hashBuf, privateKeyECDSA)
	if err != nil {
		return nil, err
	}
	if transformV {
		signature[64] += 27
	}
	go_logger.Logger.Debug("SignData")
	return signature, nil
}

func (api *SignerAPI) SignTypedData(ctx context.Context, addr common.MixedcaseAddress, typedData core.TypedData) (hexutil.Bytes, error) {
	domainSeparator, err := typedData.HashStruct("EIP712Domain", typedData.Domain.Map())
	if err != nil {
		return nil, err
	}
	typedDataHash, err := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	if err != nil {
		return nil, err
	}
	rawData := []byte(fmt.Sprintf("\x19\x01%s%s", string(domainSeparator), string(typedDataHash)))
	privateKeyECDSA, err := api.privateKey(strings.ToLower(addr.Original()))
	if err != nil {
		return nil, err
	}
	signature, err := crypto.Sign(crypto.Keccak256(rawData), privateKeyECDSA)
	if err != nil {
		return nil, err
	}
	signature[64] += 27
	go_logger.Logger.Debug("SignTypedData")
	return signature, nil
}

func (api *SignerAPI) SignGnosisSafeTx(ctx context.Context, signerAddress common.MixedcaseAddress, gnosisTx core.GnosisSafeTx, methodSelector *string) (*core.GnosisSafeTx, error) {
	go_logger.Logger.Debug("SignGnosisSafeTx")
	return nil, errors.New("not implement")
}

func (api *SignerAPI) EcRecover(ctx context.Context, data hexutil.Bytes, sig hexutil.Bytes) (common.Address, error) {
	if len(sig) != 65 {
		return common.Address{}, fmt.Errorf("signature must be 65 bytes long")
	}
	if sig[64] != 27 && sig[64] != 28 {
		return common.Address{}, fmt.Errorf("invalid Ethereum signature (V is not 27 or 28)")
	}
	sig[64] -= 27 // Transform yellow paper V from 27/28 to 0/1
	hash := accounts.TextHash(data)
	rpk, err := crypto.SigToPub(hash, sig)
	if err != nil {
		return common.Address{}, err
	}
	go_logger.Logger.Debug("EcRecover")
	return crypto.PubkeyToAddress(*rpk), nil
}

// Returns the external api version. This method does not require user acceptance. Available methods are
// available via enumeration anyway, and this info does not contain user-specific data
func (api *SignerAPI) Version(ctx context.Context) (string, error) {
	go_logger.Logger.Debug("Version")
	return version.Version, nil
}
