package consumer

import "os"

// SubWalletService Represent the sub wallet service contract
type SubWalletService interface {
	Process(topics []string, signals chan os.Signal)
}
