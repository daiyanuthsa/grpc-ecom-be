output "instance_ip" {
  description = "Alamat IP publik dari VM build agent."
  value       = google_compute_instance.build_agent.network_interface[0].access_config[0].nat_ip
}