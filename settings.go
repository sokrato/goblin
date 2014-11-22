package goblin


type Settings map[string]interface{}


var defaultSettings = Settings{
    "debug": false,
    "env": "dev",
}
