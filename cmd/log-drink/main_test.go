package main

import (
	"testing"
	"time"
	"uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_getResponseSK(t *testing.T) {
	userID := uuid.NewV4().String()
	tm := new(time.Now())

	sk, err := getCompositeID(userID, tm)
	require.NoError(t, err)
	assert.Len(t, sk, 32)
}
