package loginhistory

import (
	"context"
	"fmt"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	loginhistorycrud "github.com/NpoolPlatform/login-gateway/pkg/crud/loginhistory"
	npool "github.com/NpoolPlatform/message/npool/logingateway"

	"github.com/go-resty/resty/v2"
)

func GetByAppUser(ctx context.Context, in *npool.GetLoginHistoriesRequest) (*npool.GetLoginHistoriesResponse, error) {
	resp, err := loginhistorycrud.GetByAppUser(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail get login history: %v", err)
	}

	infos := []*npool.LoginHistory{}
	ipLocations := map[string]string{}

	for _, info := range resp.Infos {
		if info.Location != "" {
			infos = append(infos, info)
			ipLocations[info.ClientIP] = info.Location
			continue
		}

		if location, ok := ipLocations[info.ClientIP]; ok {
			info.Location = location
			infos = append(infos, info)

			err = loginhistorycrud.Update(ctx, info)
			if err != nil {
				logger.Sugar().Warnf("fail update login info: %v", err)
				continue
			}

			continue
		}

		type ipResp struct {
			Error   bool   `json:"error"`
			City    string `json:"city"`
			Country string `json:"country_name"`
			IP      string `json:"ip"`
			Reason  string `json:"reason"`
		}

		resp, err := resty.New().R().SetResult(&ipResp{}).Get(fmt.Sprintf("https://ipapi.co/%v/json", info.ClientIP))
		if err != nil {
			logger.Sugar().Warnf("fail get login location: %v", err)
			infos = append(infos, info)
			continue
		}

		rc, ok := resp.Result().(*ipResp)
		if !ok {
			logger.Sugar().Warnf("fail get login location: %v", rc.Reason)
			infos = append(infos, info)
			continue
		}

		if rc.Error {
			logger.Sugar().Warnf("fail get login location: %v", rc.Reason)
			info.Location = rc.Reason
			infos = append(infos, info)
			continue
		}

		ipLocations[info.ClientIP] = fmt.Sprintf("%v, %v", rc.City, rc.Country)
		info.Location = ipLocations[info.ClientIP]
		infos = append(infos, info)
	}

	return &npool.GetLoginHistoriesResponse{
		Infos: infos,
	}, nil
}
