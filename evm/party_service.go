package evm

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
)

func GetUSDTPricePerGasUnit() (*big.Int, error) {
	// Get the price of USDT through an API
	resp, err := http.Get("https://api.coingecko.com/api/v3/simple/price?ids=tether&vs_currencies=eth")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]map[string]float64
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	priceInETH, ok := result["tether"]["eth"]
	if !ok {
		return nil, errors.New("could not retrieve price from API response")
	}

	priceBigInt := new(big.Int)
	priceBigInt.SetString(fmt.Sprintf("%.0f", priceInETH*1e18), 10) // 假设价格单位为 wei

	return priceBigInt, nil
}
