version: '2'
catalog:
  name: "Services Gateway"
  version: "0.0.2"
  description: "Public exposure in a non-candid way."
  questions:
  - variable: "AWS_ACCESS_KEY_ID"
    label: "AWS Access Key ID"
    description: "AWS Access Key ID to update Route 53."
    type: "string"
    required: true
  - variable: "AWS_SECRET_ACCESS_KEY"
    label: "AWS Secret Access Key"
    description: "AWS Secret Access Key to update Route 53."
    type: "string"
    required: true
  - variable: "FORWARD_PROXY"
    label: "Forward Proxy"
    description: "(optional) - a forward proxy to use (in format, user:password@hostname:port)."
    type: "string"
    required: false
  - variable: "DEFAULT_TLS_CERTIFICATE"
    label: "Default TLS Certificate"
    description: "The X509 certificate to use for TLS connections."
    type: certificate
    required: true
services:
  default-website:
    scale: 2
    start_on_create: true
  dns-updater:
    scale: 1
    start_on_create: true
  router-tcp-80:
    start_on_create: true
  router-tcp-443:
    start_on_create: true
    lb_config:
      default_cert: '${DEFAULT_TLS_CERTIFICATE}'
      port_rules:
      - hostname: ''
        path: ''
        priority: 1
        protocol: https
        service: default-website
        source_port: 443
        target_port: 8080
    health_check:
      healthy_threshold: 2
      response_timeout: 2000
      port: 42
      unhealthy_threshold: 3
      initializing_timeout: 60000
      interval: 2000
      strategy: recreate
      reinitializing_timeout: 60000
  service-discovery:
    scale: 1
    start_on_create: true
