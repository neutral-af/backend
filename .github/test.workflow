workflow "Go test" {
  on = "push"
  resolves = ["GitHub Action for Docker"]
}

action "GitHub Action for Docker" {
  uses = "cedrickring/golang-action@1.3.0"
  args = "go test ./..."
}
