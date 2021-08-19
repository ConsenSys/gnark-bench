/*
Copyright © 2021 ConsenSys Software Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/consensys/gnark-bench/circuit"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/pkg/profile"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gnark-bench",
	Short: "runs benchmarks and profiles using gnark",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var (
	fCircuit     *string
	fCircuitSize *int
	fAlgo        *string
	fProfile     *string
	fCount       *int
	fCurve       *string
)

var (
	curveID ecc.ID
	p       func(p *profile.Profile)
	c       circuit.BenchCircuit
)

func init() {

	_curves := ecc.Implemented()
	curves := make([]string, len(_curves))
	for i := 0; i < len(_curves); i++ {
		curves[i] = strings.ToLower(_curves[i].String())
	}

	fCircuit = rootCmd.PersistentFlags().String("circuit", "expo", "name of the circuit to use")
	fCircuitSize = rootCmd.PersistentFlags().Int("size", 10000, "size of the circuit, parameter to circuit constructor")
	fCount = rootCmd.PersistentFlags().Int("count", 2, "bench count (time is averaged on number of executions)")
	fAlgo = rootCmd.PersistentFlags().String("algo", "prove", "name of the algorithm to benchmark. must be compile, setup, prove or verify")
	fProfile = rootCmd.PersistentFlags().String("profile", "none", "type of profile. must be none, trace, cpu or mem")
	fCurve = rootCmd.PersistentFlags().String("curve", "bn254", "curve name. must be "+fmt.Sprint(curves))
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// plonkCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// plonkCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func parseFlags() error {
	if *fCircuitSize <= 0 {
		return errors.New("circuit size must be >= 0")
	}
	if *fCount <= 0 {
		return errors.New("bench count must be >= 0")
	}

	switch *fAlgo {
	case "compile", "setup", "prove", "verify":
	default:
		return errors.New("invalid algo")
	}

	switch *fProfile {
	case "none":
	case "trace":
		p = profile.TraceProfile
	case "cpu":
		p = profile.CPUProfile
	case "mem":
		p = profile.MemProfile
	default:
		return errors.New("invalid profile")
	}

	curves := ecc.Implemented()
	for _, id := range curves {
		if *fCurve == strings.ToLower(id.String()) {
			curveID = id
		}
	}
	if curveID == ecc.UNKNOWN {
		return errors.New("invalid curve")
	}

	var ok bool
	c, ok = circuit.BenchCircuits[*fCircuit]
	if !ok {
		return errors.New("unknown circuit")
	}

	return nil
}
