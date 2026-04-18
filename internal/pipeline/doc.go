// Package pipeline provides an end-to-end processing stage that connects the
// port scanner to the notification system.
//
// Each call to Run performs one scan cycle:
//  1. Scan current port bindings via portscanner.
//  2. Diff against the previous snapshot via alerting.Alerter.
//  3. Apply filter rules to drop uninteresting entries.
//  4. Check suppression rules to silence known-noisy alerts.
//  5. Deduplicate repeated events within a TTL window.
//  6. Dispatch remaining alerts through the notify.Notifier.
//
// All stages are optional except Scanner and Notifier; omitting a stage
// simply skips that processing step.
package pipeline
