package config

import (
	"encoding/json"
	"os"
	"time"
)

type Config struct {
	RangeSensorTriggerPin int8    `json:"range_sensor_trigger_pin"`
	RangeSensorEchoPin    int8    `json:"range_sensor_echo_pin"`
	BuzzerPin             uint8   `json:"buzzer_pin"`
	DeskBottomPosition    float32 `json:"desk_bottom_position"`
	DeskTopPosition       float32 `json:"desk_top_position"`

	DurationToStand time.Duration `json:"duration_to_stand"`
	DurationToSit   time.Duration `json:"duration_to_sit"`
	NotifyToSit     bool          `json:"notify_to_sit"` // make "beep" after passing minimum time to stand
}

var config *Config
var lastUpdatedAt time.Time

func getDefault() Config {
	return Config{
		RangeSensorTriggerPin: -99,
		RangeSensorEchoPin:    -99,
		BuzzerPin:             99,
		DeskBottomPosition:    50,
		DeskTopPosition:       120,
		DurationToStand:       time.Minute * 60,
		DurationToSit:         time.Minute * 10,
		NotifyToSit:           true,
	}

}

func NewConfig() *Config {
	if config == nil || time.Since(lastUpdatedAt) > time.Minute*10 {
		// open from config.json if exists, else create with default
		config = &Config{}
		lastUpdatedAt = time.Now()

		by, err := os.ReadFile("config.json")
		if err == nil {
			err = json.Unmarshal(by, config)
			if err == nil {
				return config
			}
		}

		defaultConfigBy, err := json.Marshal(getDefault())
		if err != nil {
			return nil
		}

		err = os.WriteFile("config.json", defaultConfigBy, 0644)
		if err != nil {
			panic(err)
		}
		return config
	}
	return config
}
