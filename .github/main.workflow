workflow "Main" {
  on = "push"
  resolves = ["cedrickring/golang-action@1.3.0"]
}

action "cedrickring/golang-action@1.3.0" {
  uses = "cedrickring/golang-action@1.3.0"
  args = "go test ./..."
}
