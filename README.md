# API Monitoring Platform

This repository runs a small Go service that probes a set of HTTP endpoints and exports Prometheus metrics. It also includes Prometheus, Alertmanager, and Grafana configuration for visualization and alerting.

## What it does

- `GET /healthz` health endpoint
- `GET /metrics` Prometheus metrics
- Periodically probes configured targets (hardcoded in the app)
- Ships metrics to Prometheus, dashboards to Grafana, and alerts via Alertmanager

## Default endpoints

- API Monitoring service: `http://localhost:8080/healthz` and `http://localhost:8080/metrics`
- Prometheus (Docker Compose): `http://localhost:9091`
- Alertmanager (Docker Compose): `http://localhost:9093`
- Grafana (Docker Compose): `http://localhost:3000`

## Prometheus metrics

The app exports these metrics (with a `target` label):

- `probe_requests_total{target="..."}`: total number of probe requests
- `probe_failed{target="..."}`: total number of probe failures
- `probe_latency_seconds{target="..."}`: probe latency histogram (seconds)

Note: the Prometheus alert rules and Grafana dashboard expressions in this repo reference failure metrics by a slightly different name than what the app exports. If you see “no data” or alerts not firing, compare the metric names used in `observability/prometheus/alert.rules.yml`, `observability/grafana/dashboards/api-monitoring.json`, and the live output of `GET /metrics`.

## Local development (Docker Compose)

Docker Compose configuration is in `deployments/docker/docker-compose.yml`.

1. Start everything:
   ```bash
   docker compose -f deployments/docker/docker-compose.yml up --build
   ```
2. Open:
   - Grafana: `http://localhost:3000` (default login is typically `admin` / `admin`)
   - Prometheus: `http://localhost:9091`
3. Verify the app metrics are being scraped:
   - Go to Prometheus and run a query like `probe_requests_total`

### Alertmanager email configuration (important)

`observability/alertmanager/alertmanager.yml` uses environment-variable placeholders like:

- `ALERT_EMAIL_TO`
- `ALERT_EMAIL_FROM`
- `GMAIL_USERNAME`
- `GMAIL_APP_PASSWORD`

If these are not set, you may need to either provide them to the `alertmanager` container or update the Alertmanager config to use a different receiver.

## Kubernetes (manifests in this repo)

### API Monitoring app

- Deployment: `deployments/k8s/api-monitoring/deployment.yaml`
- Service: `deployments/k8s/api-monitoring/service.yaml`

### Prometheus

- Deployment/Service/ConfigMap: `observability/prometheus/k8s/`

### Grafana provisioning (note)

This repo includes Grafana provisioning files and dashboards under `observability/grafana/`, but it does not include a Grafana Deployment manifest.

To use Grafana on Kubernetes, ensure you install Grafana separately and mount/copy:

- `observability/grafana/provisioning/` (datasource + dashboard provider)
- `observability/grafana/dashboards/api-monitoring.json` (the dashboard)

## Customizing probe targets

Probe targets are hardcoded in `cmd/server/main.go` in the `targets := []config.Target{ ... }` slice.

To monitor different endpoints:

1. Edit `cmd/server/main.go`
2. Rebuild your image (Docker) or redeploy (Kubernetes)

## Notes / gotchas

- Metric names must match your Prometheus alert queries and Grafana dashboard expressions. If you change metric definitions in `internal/metrics/metrics.go`, update:
  - `observability/prometheus/alert.rules.yml`
  - `observability/grafana/dashboards/api-monitoring.json`

