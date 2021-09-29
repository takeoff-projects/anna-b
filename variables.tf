variable "project_id" {
  description = "The GCP Project ID."
  type        = string
}

variable "region" {
  type = string
  default = "europe-west3"
}

variable "key" {
  description = "Path to your service account key json file."
  type = string
}