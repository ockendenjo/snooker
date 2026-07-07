package main

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_getResponseSK(t *testing.T) {
	userID := uuid.NewString()
	tm := new(time.Now())

	sk, err := getCompositeID(userID, tm)
	require.NoError(t, err)
	assert.Len(t, sk, 32)
}
