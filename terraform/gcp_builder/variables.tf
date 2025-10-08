variable "gcp_project_id" {
  description = "ID Proyek Google Cloud Anda."
  type        = string
}

variable "gcp_region" {
  description = "Region GCP untuk membuat VM."
  type        = string
  default     = "asia-southeast1"
}

variable "gcp_zone" {
  description = "Zone GCP untuk membuat VM."
  type        = string
  default     = "asia-southeast1-b"
}

variable "ssh_user" {
  description = "Username untuk SSH ke VM."
  type        = string
  default     = "jenkins"
}

variable "ssh_public_key_path" {
  description = "Path ke file kunci SSH publik."
  type        = string
  default     = "~/.ssh/id_rsa.pub"
}