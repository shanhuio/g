// Copyright (C) 2023  Shanhu Tech Inc.
//
// This program is free software: you can redistribute it and/or modify it
// under the terms of the GNU Affero General Public License as published by the
// Free Software Foundation, either version 3 of the License, or (at your
// option) any later version.
//
// This program is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
// FITNESS FOR A PARTICULAR PURPOSE.  See the GNU Affero General Public License
// for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package srpc

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"shanhu.io/pub/errcode"
)

type rpcError struct {
	StatusCode int
	Status     string
	Body       string
}

func (e *rpcError) Error() string {
	if e.Body != "" {
		return fmt.Sprintf("%s - %s", e.Status, e.Body)
	}
	return e.Status
}

func isSuccessStatus(statusCode int) bool { return statusCode/100 == 2 }

func respError(resp *http.Response) error {
	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return errcode.Annotate(err, "read response body")
	}

	rpcErr := &rpcError{
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
		Body:       strings.TrimSpace(string(bs)),
	}

	switch rpcErr.StatusCode {
	case http.StatusUnauthorized, http.StatusForbidden:
		return errcode.Add(errcode.Unauthorized, rpcErr)
	case http.StatusNotFound, http.StatusBadRequest:
		return errcode.Add(errcode.InvalidArg, rpcErr)
	}
	return rpcErr
}
