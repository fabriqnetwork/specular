// Copyright 2023, Specular contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Adapted from Optimism's `op-bindings`

package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/specularL2/specular/ops/bindings/ast"
	"github.com/specularL2/specular/ops/bindings/gen/hardhat"
)

type flags struct {
	Contracts   string
	SourceMaps  string
	OutDir      string
	Package     string
	ContractDir string
	AbigenBin   string
}

type data struct {
	Name          string
	StorageLayout string
	DeployedBin   string
	Package       string
}

func main() {
	var f flags
	flag.StringVar(&f.OutDir, "out", "", "Output directory to put code in")
	flag.StringVar(&f.Contracts, "contracts", "artifacts.json", "Path to file containing list of contracts to generate bindings for")
	flag.StringVar(&f.Package, "package", "artifacts", "Go package name")
	flag.StringVar(&f.ContractDir, "contract-dir", "", "Contract directory")
	flag.StringVar(&f.AbigenBin, "abigen-bin", "", "Abigen binary")
	flag.Parse()

	if f.ContractDir == "" {
		log.Fatal("must provide -contract-dir")
	}
	log.Printf("Using contract directory %s\n", f.ContractDir)

	if f.AbigenBin == "" {
		log.Fatal("must provide -abigen-bin")
	}
	log.Printf("Using abigen %s\n", f.AbigenBin)

	contractData, err := os.ReadFile(f.Contracts)
	if err != nil {
		log.Fatal("error reading contract list: %w\n", err)
	}
	contracts := []string{}
	if err := json.Unmarshal(contractData, &contracts); err != nil {
		log.Fatal("error parsing contract list: %w\n", err)
	}

	if len(contracts) == 0 {
		log.Fatalf("must define a list of contracts")
	}

	t := template.Must(template.New("artifact").Parse(tmpl))

	// Make a temp dir to hold all the inputs for abigen
	dir, err := os.MkdirTemp("", "specular-bindings")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Using package %s\n", f.Package)

	defer os.RemoveAll(dir)
	log.Printf("created temp dir %s\n", dir)

	buildInfos, err := hardhat.ReadBuildInfos(filepath.Join(f.ContractDir, "artifacts"))
	if err != nil {
		log.Fatalf("error reading build infos: %v\n", err)
	}

	for _, contract := range contracts {
		log.Printf("generating code for %s\n", contract)

		name := contract[strings.LastIndex(contract, "/")+1 : len(contract)-4] // remove .sol

		artifactPath := path.Join(f.ContractDir, "artifacts", contract, name+".json")
		artifact, err := hardhat.ReadArtifact(artifactPath)
		if err != nil {
			log.Fatalf("error reading artifact: %v\n", err)
		}

		abi := artifact.Abi
		rawAbi, err := json.Marshal(abi)
		if err != nil {
			log.Fatalf("error marshaling abi: %v\n", err)
		}
		abiFile := path.Join(dir, name+".abi")
		if err := os.WriteFile(abiFile, rawAbi, 0o600); err != nil {
			log.Fatalf("error writing file: %v\n", err)
		}
		rawBytecode := artifact.Bytecode.String()
		bytecodeFile := path.Join(dir, name+".bin")
		if err := os.WriteFile(bytecodeFile, []byte(rawBytecode), 0o600); err != nil {
			log.Fatalf("error writing file: %v\n", err)
		}

		cwd, err := os.Getwd()
		if err != nil {
			log.Fatalf("error getting cwd: %v\n", err)
		}

		outFile := path.Join(cwd, f.Package, name+".go")

		cmd := exec.Command(f.AbigenBin, "--abi", abiFile, "--bin", bytecodeFile, "--pkg", f.Package, "--type", name, "--out", outFile)
		cmd.Stdout = os.Stdout

		if err := cmd.Run(); err != nil {
			log.Fatalf("error running abigen: %v\n", err)
		}

		storageLayout, err := hardhat.GetStorageLayout(name, buildInfos[name])
		if err != nil {
			log.Fatalf("error getting storage layout: %v\n", err)
		}
		canonicalStorage := ast.CanonicalizeASTIDs(storageLayout, filepath.Join(f.ContractDir, ".."))
		ser, err := json.Marshal(canonicalStorage)
		if err != nil {
			log.Fatalf("error marshaling storage: %v\n", err)
		}
		serStr := strings.Replace(string(ser), "\"", "\\\"", -1)

		d := data{
			Name:          name,
			StorageLayout: serStr,
			DeployedBin:   artifact.DeployedBytecode.String(),
			Package:       f.Package,
		}

		fname := filepath.Join(f.OutDir, name+"_more.go")
		outfile, err := os.OpenFile(
			fname,
			os.O_RDWR|os.O_CREATE|os.O_TRUNC,
			0o600,
		)
		if err != nil {
			log.Fatalf("error opening %s: %v\n", fname, err)
		}

		if err := t.Execute(outfile, d); err != nil {
			log.Fatalf("error writing template %s: %v", outfile.Name(), err)
		}
		outfile.Close()
		log.Printf("wrote file %s\n", outfile.Name())
	}
}

var tmpl = `// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package {{.Package}}

import (
	"encoding/json"

	"github.com/specularL2/specular/ops/bindings/solc"
)

const {{.Name}}StorageLayoutJSON = "{{.StorageLayout}}"

var {{.Name}}StorageLayout = new(solc.StorageLayout)

var {{.Name}}DeployedBin = "{{.DeployedBin}}"
func init() {
	if err := json.Unmarshal([]byte({{.Name}}StorageLayoutJSON), {{.Name}}StorageLayout); err != nil {
		panic(err)
	}

	layouts["{{.Name}}"] = {{.Name}}StorageLayout
	deployedBytecodes["{{.Name}}"] = {{.Name}}DeployedBin
}
`
