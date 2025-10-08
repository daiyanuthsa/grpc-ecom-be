provider "google" {
  project = var.gcp_project_id
  region  = var.gcp_region
  zone    = var.gcp_zone
}

resource "google_compute_firewall" "allow_ssh" {
  name    = "allow-ssh"
  network = "default"

  allow {
    protocol = "tcp"
    ports    = ["22"]
  }

  source_ranges = ["0.0.0.0/0"]
  target_tags   = ["allow-ssh"]
}

resource "google_compute_instance" "build_agent" {
  name         = "jenkins-build-agent-arm64-${random_id.id.hex}"
  # T2A adalah seri VM ARM64 yang paling terjangkau
  machine_type = "t2a-standard-2" 
  zone         = var.gcp_zone
  tags = ["allow-ssh"]

  boot_disk {
    initialize_params {
      # Menggunakan image Ubuntu yang dioptimalkan untuk GCP dan mendukung ARM64
      image = "debian-cloud/debian-12-arm64"
    }
  }

  network_interface {
    network = "default"
    access_config {
      // Diperlukan untuk memberikan IP publik agar Jenkins bisa terhubung
    }
  }

  // Startup script untuk menginstal Docker & Buildx secara otomatis saat VM pertama kali boot
  metadata_startup_script = <<-EOF
    #!/bin/bash
    sudo apt-get update
    sudo apt-get install -y ca-certificates curl gnupg
    sudo install -m 0755 -d /etc/apt/keyrings
    curl -fsSL https://download.docker.com/linux/debian/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
    sudo chmod a+r /etc/apt/keyrings/docker.gpg
    echo \
      "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/debian \
      $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
      sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
    sudo apt-get update
    sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin
    sudo usermod -aG docker ${var.ssh_user}
  EOF

  // Menambahkan kunci SSH publik agar Jenkins bisa login
  metadata = {
    ssh-keys = "${var.ssh_user}:${var.ssh_public_key_content}"
  }

  // Service account untuk memberikan izin (opsional, tapi direkomendasikan)
  service_account {
    scopes = ["cloud-platform"]
  }
}

// ID acak untuk memastikan nama VM unik setiap kali dibuat
resource "random_id" "id" {
  byte_length = 4
}