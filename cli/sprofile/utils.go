package sprofile

import (
	"encoding/json"
	"net/http"
	"runtime"

	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-client-electron/cli/errortypes"
	"github.com/pritunl/pritunl-client-electron/cli/profile"
	"github.com/pritunl/pritunl-client-electron/cli/service"
)

func GetAll() (sprfls []*Sprofile, err error) {
	reqUrl := service.GetAddress() + "/sprofile"

	authKey, err := service.GetAuthKey()
	if err != nil {
		return
	}

	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		err = errortypes.RequestError{
			errors.Wrap(err, "sprofile: Get request failed"),
		}
		return
	}

	if runtime.GOOS == "linux" {
		req.Host = "unix"
	}
	req.Header.Set("Auth-Key", authKey)
	req.Header.Set("User-Agent", "pritunl")
	req.Header.Set("Content-Type", "application/json")

	resp, err := service.GetClient().Do(req)
	if err != nil {
		err = errortypes.RequestError{
			errors.Wrap(err, "sprofile: Request failed"),
		}
		return
	}
	defer resp.Body.Close()

	sprfls = []*Sprofile{}
	err = json.NewDecoder(resp.Body).Decode(&sprfls)
	if err != nil {
		err = errortypes.ParseError{
			errors.Wrap(err, "sprofile: Failed to parse response"),
		}
		return
	}

	sprflsMap := map[string]*Sprofile{}
	for _, sprfl := range sprfls {
		sprflsMap[sprfl.Id] = sprfl
	}

	prfls, err := profile.GetAll()
	if err != nil {
		return
	}

	for _, prfl := range prfls {
		sprfl := sprflsMap[prfl.Id]
		if sprfl != nil {
			sprfl.Profile = prfl
		}
	}

	return
}
