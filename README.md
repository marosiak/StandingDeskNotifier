# 🏋️‍♂️ Standing Desk Notifier
Bought a standing desk just to sit in front of it like it’s a regular desk? Guilty as charged? 😅 No worries! I’ve got your back (literally). Let’s get those glutes off the chair with this nifty little notifier that *politely* insists you stand up. 

## How It Works: A User’s Tale 🎢

1. **You're sitting**... comfortably.
2. **You're still sitting**... and getting *way* too cozy.
3. **You hear a soft "beep beep"**... just a gentle nudge.
4. **You ignore it**... it’s just a beep, right?
5. **40 seconds later**... BEEP BEEP BEEP! Every 3 seconds! 🔊 
6. **You surrender**... reluctantly raise your desk, stand up, and bask in the sweet sound of silence.
7. **1 minute later**... you’re tired already? Seriously?
8. **You sit back down**... but wait...
9. **BEEP BEEP BEEP**... Time to stand up again, buddy! 14 more minutes, you can do it! 💪

## 🛠️ Hardware
Keeping it simple and straightforward, here’s what you’ll need:
1. **Raspberry Pi Zero** (with Raspbian and charging cable, of course).
2. **HC-SR04** - The trusty ultrasonic sensor measuring the distance from desk to floor.
3. **Active Buzzer** - The source of your new productivity (and maybe a bit of frustration).
4. **1x resistor 2.7k or 2.2k**
5. **1x resistor 4.7k**

## 🛠️ Hardware installation
// TODO

You will need to modify HC-SR04 in order to work with 3.3v

**Tutorial**: https://www.instructables.com/Modify-Ultrasonic-Sensors-for-3-Volts-Logic-prepar/

## 💰 Hardware Prices
1. **Amazon.com**: ~ $27.50
2. **Botland.pl**: ~ $19.68 (~ 74 PLN)

## 🚀 Installation
It’s as easy as running `chmod +x install.sh && ./install.sh`. The script does all the heavy lifting. Once it’s done, it’ll drop a few logs telling you where it stashed the goodies—like the binary and your shiny new `config.json`.

## 🔧 Configuration
Your `config.json` is where the magic happens. You can edit it on the fly while the app is running, but remember, changes will take about 10 minutes to kick in. If patience isn’t your virtue, just restart the app after tweaking the config.
```json
{
  "range_sensor_trigger_pin": 5,
  "range_sensor_echo_pin": 6,
  "buzzer_pin": 7,

  "desk_bottom_position": 50.5,
  "desk_top_position": 120.0,

  "duration_to_stand": "15m",
  "duration_to_sit": "1h",
  "notify_to_sit": true
}
