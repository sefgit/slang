package daemon

import (
	"net/http"
	"encoding/json"
	"io"
	"github.com/Bitspark/slang/pkg/api"
	"github.com/Bitspark/slang/pkg/core"
	"strings"
	"github.com/Bitspark/slang/pkg/builtin"
	"path/filepath"
	"log"
)

type DaemonService struct {
	Routes map[string]*DaemonEndpoint
}

type DaemonEndpoint struct {
	Handle func(w http.ResponseWriter, r *http.Request)
}

func readJSON(r io.Reader, dat interface{}) error {
	dec := json.NewDecoder(r)
	return dec.Decode(&dat)
}

func writeJSON(w io.Writer, dat interface{}) error {
	return json.NewEncoder(w).Encode(dat)
}

type outJSON struct {
	Objects []OperatorDefJSON `json:"objects"`
	Status  string            `json:"status"`
	Error   *Error            `json:"error,omitempty"`
}

type OperatorDefJSON struct {
	Name string           `json:"name"`
	Def  core.OperatorDef `json:"def"`
	Type string           `json:"type"`
}

type Error struct {
	Msg  string `json:"msg"`
	Code string `json:"code"`
}

var OperatorDefService = &DaemonService{map[string]*DaemonEndpoint{
	"/": {func(w http.ResponseWriter, r *http.Request) {
		var dataOut outJSON
		var opDefList []OperatorDefJSON
		var err error
		cwd := r.FormValue("cwd")

		e := api.NewEnviron(cwd)

		opNames, err := e.ListOperatorNames()

		if err == nil {
			builtinOpNames := builtin.GetBuiltinNames()

			// Gather builtin/elementary opDefs
			for _, opFQName := range builtinOpNames {
				opDef, err := builtin.GetOperatorDef(opFQName)

				if err != nil {
					break
				}

				opDefList = append(opDefList, OperatorDefJSON{
					Name: opFQName,
					Type: "elementary",
					Def:  opDef,
				})
			}

			if err == nil {
				// Gather opDefs from local & lib
				for _, opFQName := range opNames {
					opDefFilePath, err := e.GetOperatorDefFilePath(strings.Replace(opFQName, ".", string(filepath.Separator), -1), "")
					if err != nil {
						continue
					}

					opDef, err := e.ReadOperatorDef(opDefFilePath, nil)
					if err != nil {
						continue
					}

					opType := "lib"
					if e.IsLocalOperator(opFQName) {
						opType = "local"
					}

					opDefList = append(opDefList, OperatorDefJSON{
						Name: opFQName,
						Type: opType,
						Def:  opDef,
					})
				}
			}
		}

		if err == nil {
			dataOut = outJSON{Status: "success", Objects: opDefList}
		} else {
			dataOut = outJSON{Status: "error", Error: &Error{err.Error(), "E0001"}}
		}

		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.WriteHeader(200)
		err = writeJSON(w, dataOut)
		if err != nil {
			log.Print(err)
		}
	}},
}}
