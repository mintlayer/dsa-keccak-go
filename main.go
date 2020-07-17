// This package is just a test demonstration of DSA 

package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/sha3"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

type Stake struct {
	Id     int
	Weight int
	Age    int
	PubKey string
}

type Slot struct {
	Id     int
	PubKey string
}

const MIN_STAKE = 440000
const MAX_VAL = 4294967295

func main() {
	// Example of BTC hash 
	btc_hash, _ := hex.DecodeString("0000000000000000000d9ed0f796aeee51b200c7293a6e31c101a0e4159bf310")

	// Sample struct of stakes
	stakes := []Stake{
		{0, MIN_STAKE, 0, "PK0"},
		{1, MIN_STAKE, 0, "PK1"},
		{2, MIN_STAKE + 1, 0, "PK2"},
		{3, MIN_STAKE * 2, 0, "PK3"},
		{4, MIN_STAKE * 2, 0, "PK4"},
		{5, MIN_STAKE * 4, 0, "PK5"},
		{6, 10000000, 0, "PK6"},
		{7, 10000001, 0, "PK7"},
		{8, 20000000, 0, "PK8"},
	}

	// Sample list of slots
	slots := []Slot{
		{0, ""},
		{1, ""},
		{2, ""},
		{3, ""},
		{4, ""},
		{5, ""},
		{6, ""},
		{7, ""},
		{8, ""},
		{9, ""},
	}

	// Calc sum of weights
	var stake_total uint32
	for _, s := range stakes {
		stake_total += uint32(s.Weight)
	}

	// This is actual DSA
	//  assumption: size of stakes/ min_stake > size of slots
	var stakes_copy = make([]Stake, len(stakes))
	copy(stakes_copy, stakes)
	for i, _ := range slots {
		rnd := rand(i, btc_hash, stake_total)
		for el_id, s := range stakes_copy {
			if uint32(s.Weight) >= rnd {
				if s.Weight >= MIN_STAKE {
					s.Weight -= MIN_STAKE
				} else {
					stakes_copy = append(stakes_copy[:el_id], stakes_copy[el_id+1:]...)
				}
				slots[i].PubKey = s.PubKey
				break
			} else {
				rnd -= uint32(s.Weight)
			}
		}
	}

	// Print the result
	fmt.Println("Initial Stake:")
	for _, s := range stakes {
		fmt.Println(s.PubKey, " - ", s.Weight)
	}
	fmt.Println("Final Result:")
	for i, s := range slots {
		fmt.Println("SLOT ", i, " -> ", s.PubKey)
	}
}

// calculate u as keccak("k",i) << hash
func rand(i int, hash []byte, max uint32) uint32 {
	var u uint32
	out := make([]byte, 4)
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, uint32(i))
	c1 := sha3.NewCShake256([]byte("Mintlayer"), buf.Bytes())
	c1.Write(hash)
	c1.Read(out)
	binary.Read(bytes.NewReader(out), binary.LittleEndian, &u)
	// Its ok for the poc to use float 
	normalized := uint32((float64(max) / float64(MAX_VAL)) * float64(u))
	return normalized
}
