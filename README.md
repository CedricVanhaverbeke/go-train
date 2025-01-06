# Go-train

<img src="./gobike.png" alt="logo" width="200"/>

this project aims to control my indoor bike trainer with a training and display the training as it would on Zwift or other platforms. Since I don't care about those virtual environments, but I'd like to watch a movie, the training should be shown as an overlay.

## Goal

The ultimate goal of this app is to start a predefined training, finish it on my bike trainer at home and upload it to Strava.

## Running the project

This is a `go 1.23.0` project

```bash
# After cloning
1. go mod tidy

# Looks for a suitable bluetooth trainer, starts a training and controls the trainer.
2a. go run main.go 

# Starts a training on a mock bluetooth trainer. It mocks incoming data from the trainer
2b. go run main.go -m true
```

## Done

- [x] Create a random training with durations and power
- [x] Show an overlay with ebitenengine
- [x] Create basic sprites
- [x] Make the sprites update every second with data from the random trainer
- [x] Make the training have a sense of time, it should progress every second
- [x] Connect to a bluetooth trainer
- [x] Control the bluetooth trainer (EGC power mode)
- [x] Get the power from the trainer
- [x] Have a mock setting
- [x] Create a mock route


## TODO

- [] Add elevation metric to the route
- [] Determine distance / determine the progression of the route
- [] Get the speed from the trainer
- [] Write data to a gpx file
