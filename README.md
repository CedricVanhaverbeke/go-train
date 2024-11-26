# Go-train

this project aims to control my indoor bike trainer with a training and display the training as it would on Zwift or other platforms. Since I don't care about those virtual environments, but I'd like to watch a movie, the training should be shown as an overlay.

## Running the project

This is a `go 1.23.0` project

```bash
# After cloning
1. go mod tidy
2. go run main.go # A fake training will start here
```

## Done

- [x] Create a random training with durations and power
- [x] Show an overlay with ebitenengine
- [x] Create basic sprites
- [x] Make the sprites update every second with data from the random trainer
- [x] Make the training have a sense of time, it should progress every second


## TODO

- [ ] Comminucate with my bike trainer which includes consuming data and feeding the trainer with information.
- [ ] Display a training with increasing or descreasing power
