// Package all imports all hardware drivers.
package all

import (
	_ "github.com/tehmaze/benjamin/deck/infinitton" // Infinitton iDisplay support
	_ "github.com/tehmaze/benjamin/deck/streamdeck" // Corsair (Elgato) Stream Deck support
)
