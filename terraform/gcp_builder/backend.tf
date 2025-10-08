terraform {
  backend "gcs" {
    bucket  = "tfstate-grpc-ecom"
    prefix  = "jenkins/builder"
  }
}