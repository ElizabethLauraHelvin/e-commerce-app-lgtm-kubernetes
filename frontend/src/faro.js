import {
  initializeFaro,
  getWebInstrumentations,
} from '@grafana/faro-web-sdk'

import { TracingInstrumentation } from '@grafana/faro-web-tracing'

let faroInstance = null

export function initFaro() {
  if (faroInstance) return faroInstance

  // Ambil runtime config jika tersedia
  const config = window.__CONFIG__ || {}

  // Gunakan FARO_COLLECTOR_URL dari runtime config.
  // Jika tidak tersedia, gunakan endpoint Faro Alloy sebagai fallback.
  const collectorUrl =
    config.FARO_COLLECTOR_URL ||
    'http://20.239.125.172/faro/collect'

  faroInstance = initializeFaro({
    url: collectorUrl,

    app: {
      name: config.APP_NAME || 'ecommerce-frontend',
      version: config.APP_VERSION || '1.0.0',
      environment: config.APP_ENV || 'production',
    },

    sessionTracking: {
      enabled: true,
      persistent: true,
    },

    instrumentations: [
      ...getWebInstrumentations({
        captureConsole: true,
        captureConsoleDisabledLevels: [],
        enablePerformanceInstrumentation: true,
      }),

      new TracingInstrumentation({
        instrumentationOptions: {
          propagateTraceHeaderCorsUrls: [/.*/],
        },
      }),
    ],
  })

  console.log('Faro initialized, collector:', collectorUrl)

  return faroInstance
}

export function getFaro() {
  return faroInstance
}