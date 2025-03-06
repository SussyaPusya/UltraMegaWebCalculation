package jsonutils

import (
	"encoding/json"

	"github.com/SussyaPusya/UltraMegaWebCalculation/internal/domain"
)

func ParseJson(bytew []byte) (domain.JsonReq, error) {

	var jsonStrr domain.JsonReq

	err := json.Unmarshal(bytew, &jsonStrr)

	if err != nil {
		return domain.JsonReq{}, err
	}

	return jsonStrr, nil

}
