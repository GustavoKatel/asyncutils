package supervisor

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	mocks "github.com/GustavoKatel/asyncutils/supervisor/mocks"
)

func TestLifeCycle1(t *testing.T) {
	assert := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mocks.NewMockService(ctrl)

	ctx := context.Background()

	spv, err := New(ctx)
	assert.Nil(err)

	token, err := spv.AddService(service)
	assert.Nil(err)

	service.EXPECT().Init(token.Context()).Return(nil)
	service.EXPECT().Run(token.Context()).Return(nil)
	service.EXPECT().Clean(token.Context()).DoAndReturn(func(ctx context.Context) error {
		token.Stop()
		return nil
	})

	assert.Nil(spv.Start())

	<-token.Context().Done()
}
