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

variable "ssh_public_key_content" {
  description = "Konten dari kunci SSH publik."
  type        = string
  sensitive   = true // Tandai sebagai sensitif
}