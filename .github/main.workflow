workflow "Main" {
  on = "push"
  resolves = ["gcloud run deploy"]
}

action "go test" {
  uses = "cedrickring/golang-action@1.3.0"
  args = "go get -v && go test ./..."
  env = {
    GO111MODULE = "on"
  }
}

# Build and push docker image


action "docker build" {
  uses = "actions/docker/cli@master"
  args = ["build", "-t", "gcr.io/carbonoffsets/backend", "."]
  needs = ["go test"]
}

action "docker credential-helper" {
  needs = ["gcloud login"]
  uses = "actions/gcloud/cli@master"
  args = ["auth", "configure-docker", "--quiet"]
}

action "docker push" {
  needs = ["docker build", "docker credential-helper"]
  uses = "actions/gcloud/cli@master"
  runs = "sh -c"
  args = ["docker push gcr.io/carbonoffsets/backend"]
}


# Deploy to google cloud run

action "gcloud login" {
  uses = "actions/gcloud/auth@ba93088eb19c4a04638102a838312bb32de0b052"
  needs = ["go test"]
  secrets = ["GCLOUD_AUTH"]
}

action "gcloud run deploy" {
  uses = "actions/gcloud/cli@ba93088eb19c4a04638102a838312bb32de0b052"
  needs = ["docker push"]
  args = "beta run deploy --quiet --image gcr.io/carbonoffsets/backend --project carbonoffsets --platform managed"
}
