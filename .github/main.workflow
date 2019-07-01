workflow "Main" {
  on = "push"
  resolves = ["GitHub Action for Google Cloud-1"]
}

action "cedrickring/golang-action@1.3.0" {
  uses = "cedrickring/golang-action@1.3.0"
  args = "go get -v && go test ./..."
  env = {
    GO111MODULE = "on"
  }
}

action "GitHub Action for Google Cloud" {
  uses = "actions/gcloud/auth@ba93088eb19c4a04638102a838312bb32de0b052"
  needs = ["cedrickring/golang-action@1.3.0"]
  secrets = ["GCLOUD_AUTH"]
}

action "cloud build" {
  uses = "actions/gcloud/cli@ba93088eb19c4a04638102a838312bb32de0b052"
  needs = ["GitHub Action for Google Cloud"]
  args = "builds submit --tag gcr.io/carbonoffsets/backend --project carbonoffsets"
}

action "GitHub Action for Google Cloud-1" {
  uses = "actions/gcloud/cli@ba93088eb19c4a04638102a838312bb32de0b052"
  needs = ["cloud build"]
  args = "beta run deploy --image gcr.io/carbonoffsets/backend --project carbonoffsets --platform managed"
}
