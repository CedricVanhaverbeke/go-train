package bluetooth

import (
	bluetooth_mocks "overlay/pkg/bluetooth/mocks"
	"testing"

	"go.uber.org/mock/gomock"
)

func TestVerifyDevice(t *testing.T) {
	ctrl := gomock.NewController()
	device := bluetooth_mocks.NewMockbluetootdevice()

	adapter.EXPECT().Connect(gomock.Any().String(), gomock.Any()).Return(device)
}
