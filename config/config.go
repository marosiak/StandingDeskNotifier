package config

import (
	"encoding/json"
	"os"
	"time"
)

// Custom type for duration
type HumanFriendlyDuration time.Duration

func (h HumanFriendlyDuration) Duration() time.Duration {
	return time.Duration(h)
}

// UnmarshalJSON for HumanFriendlyDuration to parse strings like "1h5m15s"
func (d *HumanFriendlyDuration) UnmarshalJSON(b []byte) error {
	// Remove the quotes around the value
	strValue := string(b)
	if len(strValue) >= 2 && strValue[0] == '"' && strValue[len(strValue)-1] == '"' {
		strValue = strValue[1 : len(strValue)-1]
	}

	duration, err := time.ParseDuration(strValue)
	if err != nil {
		return err
	}
	*d = HumanFriendlyDuration(duration)
	return nil
}

// MarshalJSON for HumanFriendlyDuration to convert it back to a string format
func (d HumanFriendlyDuration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

// String method to represent HumanFriendlyDuration as a string
func (d HumanFriendlyDuration) String() string {
	return time.Duration(d).String()
}

type Config struct {
	RangeSensorTriggerPin int8                  `json:"range_sensor_trigger_pin"`
	RangeSensorEchoPin    int8                  `json:"range_sensor_echo_pin"`
	BuzzerPin             uint8                 `json:"buzzer_pin"`
	DeskBottomPosition    float32               `json:"desk_bottom_position"`
	DeskTopPosition       float32               `json:"desk_top_position"`
	DurationToStand       HumanFriendlyDuration `json:"duration_to_stand"`
	DurationToSit         HumanFriendlyDuration `json:"duration_to_sit"`
	NotifyToSit           bool                  `json:"notify_to_sit"` // make "beep" after passing minimum time to stand
}

var config *Config
var lastUpdatedAt time.Time

func getDefault() Config {
	return Config{
		RangeSensorTriggerPin: 17,
		RangeSensorEchoPin:    27,
		BuzzerPin:             23,
		DeskBottomPosition:    95,
		DeskTopPosition:       120,
		DurationToStand:       HumanFriendlyDuration(time.Minute * 60),
		DurationToSit:         HumanFriendlyDuration(time.Minute * 10),
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
