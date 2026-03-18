terraform {
  required_providers {
    docker = {
      source  = "kreuzwerker/docker"
      version = "~> 3.0"
    }
  }
}

provider "docker" {}

# Pull NGINX image
resource "docker_image" "nginx" {
  name = "nginx:latest"
}

# Create container
resource "docker_container" "nginx_container" {
  name  = "nginx-server"
  image = docker_image.nginx.image_id

  ports {
    internal = 80
    external = var.external_port
  }
}