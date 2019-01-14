version: '2'
services:
  default-website:
    image: flaccid/rancher-services-gateway:latest
    stdin_open: true
    tty: true
    command:
    - --ui
    - -d
    labels:
      io.rancher.scheduler.affinity:host_label: tier=private
      io.rancher.container.agent.role: environment
      label: private
      io.rancher.container.create_agent: 'true'
      io.rancher.container.pull_image: always
  service-discovery:
    image: flaccid/rancher-services-gateway:latest
    stdin_open: true
    tty: true
    command:
    - -t
    - '15'
    labels:
      io.rancher.container.agent.role: environment
      io.rancher.container.create_agent: 'true'
      io.rancher.container.pull_image: always
      io.rancher.scheduler.affinity:host_label_soft: tier=private
  router-tcp-80:
    image: jamessharp/docker-nginx-https-redirect
    stdin_open: true
    tty: true
    ports:
    - 80:80/tcp
    labels:
      io.rancher.container.pull_image: always
      io.rancher.scheduler.global: 'true'
      io.rancher.scheduler.affinity:host_label_soft: tier=private
  router-tcp-443:
    image: rancher/lb-service-haproxy:v0.7.15
    ports:
    - 443:443/tcp
    labels:
      io.rancher.container.agent.role: environmentAdmin,agent
      io.rancher.container.agent_service.drain_provider: 'true'
      io.rancher.scheduler.affinity:host_label: tier=private
      io.rancher.container.create_agent: 'true'
      services_gateway: 'true'
      io.rancher.scheduler.global: 'true'
  dns-updater:
    image: flaccid/ranch53:latest
    stdin_open: true
    tty: true
    command:
    - --sync-host-pools
    - --sync-lb-services
    environment:
      AWS_ACCESS_KEY_ID: ${AWS_ACCESS_KEY_ID}
      AWS_SECRET_ACCESS_KEY: ${AWS_SECRET_ACCESS_KEY}
      POLL_INTERVAL: '30'
{{- if (.Values.RANCHER_URL)}}
      RANCHER_URL: ${RANCHER_URL}
{{- end}}
{{- if (.Values.RANCHER_ACCESS_KEY)}}
      RANCHER_ACCESS_KEY: ${RANCHER_ACCESS_KEY}
{{- end}}
{{- if (.Values.RANCHER_SECRET_KEY)}}
      RANCHER_SECRET_KEY: ${RANCHER_SECRET_KEY}
{{- end}}
{{- if (.Values.FORWARD_PROXY)}}
      http_proxy: ${FORWARD_PROXY}
{{- end}}
    labels:
      io.rancher.container.agent.role: environment
      io.rancher.container.create_agent: 'true'
      io.rancher.container.pull_image: always
      io.rancher.scheduler.affinity:host_label_soft: tier=private
{{- if eq .Values.USE_LETSENCRYPT "Yes"}}
  letsencrypt:
    image: vxcontrol/rancher-letsencrypt:v1.0.0
    environment:
      API_VERSION: Production
      AWS_ACCESS_KEY: '${AWS_ACCESS_KEY_LETSENCRYPT}'
      AWS_SECRET_KEY: '${AWS_SECRET_KEY_LETSENCRYPT}'
      CERT_NAME: '${CERT_NAME}'
      DOMAINS: '${DOMAINS}'
      EMAIL: '${EMAIL}'
      EULA: 'Yes'
      PROVIDER: Route53
      PUBLIC_KEY_TYPE: RSA-2048
      RENEWAL_PERIOD_DAYS: '20'
      RENEWAL_TIME: '12'
      RUN_ONCE: 'false'
    volumes:
    - /var/lib/rancher:/var/lib/rancher
    labels:
      io.rancher.scheduler.affinity:host_label: tier=public
      io.rancher.container.agent.role: environment
      io.rancher.container.create_agent: 'true'
      tier: public
{{- end}}
