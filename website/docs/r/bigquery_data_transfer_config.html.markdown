---
# ----------------------------------------------------------------------------
#
#     ***     AUTO GENERATED CODE    ***    Type: MMv1     ***
#
# ----------------------------------------------------------------------------
#
#     This file is automatically generated by Magic Modules and manual
#     changes will be clobbered when the file is regenerated.
#
#     Please read more about how to change this file in
#     .github/CONTRIBUTING.md.
#
# ----------------------------------------------------------------------------
subcategory: "BigQuery Data Transfer"
layout: "google"
page_title: "Google: google_bigquery_data_transfer_config"
sidebar_current: "docs-google-bigquery-data-transfer-config"
description: |-
  Represents a data transfer configuration.
---

# google\_bigquery\_data\_transfer\_config

Represents a data transfer configuration. A transfer configuration
contains all metadata needed to perform a data transfer.


To get more information about Config, see:

* [API documentation](https://cloud.google.com/bigquery/docs/reference/datatransfer/rest/v1/projects.locations.transferConfigs/create)
* How-to Guides
    * [Official Documentation](https://cloud.google.com/bigquery/docs/reference/datatransfer/rest/)

~> **Warning:** All arguments including `sensitive_params.secret_access_key` will be stored in the raw
state as plain-text. [Read more about sensitive data in state](/docs/state/sensitive-data.html).

## Example Usage - Bigquerydatatransfer Config Scheduled Query


```hcl
data "google_project" "project" {
}

resource "google_project_iam_member" "permissions" {
  project = data.google_project.project.project_id
  role   = "roles/iam.serviceAccountShortTermTokenMinter"
  member = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-bigquerydatatransfer.iam.gserviceaccount.com"
}

resource "google_bigquery_data_transfer_config" "query_config" {
  depends_on = [google_project_iam_member.permissions]

  display_name           = "my-query"
  location               = "asia-northeast1"
  data_source_id         = "scheduled_query"
  schedule               = "first sunday of quarter 00:00"
  destination_dataset_id = google_bigquery_dataset.my_dataset.dataset_id
  params = {
    destination_table_name_template = "my_table"
    write_disposition               = "WRITE_APPEND"
    query                           = "SELECT name FROM tabl WHERE x = 'y'"
  }
}

resource "google_bigquery_dataset" "my_dataset" {
  depends_on = [google_project_iam_member.permissions]

  dataset_id    = "my_dataset"
  friendly_name = "foo"
  description   = "bar"
  location      = "asia-northeast1"
}
```

## Argument Reference

The following arguments are supported:


* `display_name` -
  (Required)
  The user specified display name for the transfer config.

* `data_source_id` -
  (Required)
  The data source id. Cannot be changed once the transfer config is created.

* `params` -
  (Required)
  These parameters are specific to each data source.


- - -


* `destination_dataset_id` -
  (Optional)
  The BigQuery target dataset id.

* `schedule` -
  (Optional)
  Data transfer schedule. If the data source does not support a custom
  schedule, this should be empty. If it is empty, the default value for
  the data source will be used. The specified times are in UTC. Examples
  of valid format: 1st,3rd monday of month 15:30, every wed,fri of jan,
  jun 13:15, and first sunday of quarter 00:00. See more explanation
  about the format here:
  https://cloud.google.com/appengine/docs/flexible/python/scheduling-jobs-with-cron-yaml#the_schedule_format
  NOTE: the granularity should be at least 8 hours, or less frequent.

* `schedule_options` -
  (Optional)
  Options customizing the data transfer schedule.
  Structure is [documented below](#nested_schedule_options).

* `email_preferences` -
  (Optional)
  Email notifications will be sent according to these preferences to the
  email address of the user who owns this transfer config.
  Structure is [documented below](#nested_email_preferences).

* `notification_pubsub_topic` -
  (Optional)
  Pub/Sub topic where notifications will be sent after transfer runs
  associated with this transfer config finish.

* `data_refresh_window_days` -
  (Optional)
  The number of days to look back to automatically refresh the data.
  For example, if dataRefreshWindowDays = 10, then every day BigQuery
  reingests data for [today-10, today-1], rather than ingesting data for
  just [today-1]. Only valid if the data source supports the feature.
  Set the value to 0 to use the default value.

* `disabled` -
  (Optional)
  When set to true, no runs are scheduled for a given transfer.

* `sensitive_params` -
  (Optional)
  Different parameters are configured primarily using the the `params` field on this
  resource. This block contains the parameters which contain secrets or passwords so that they can be marked
  sensitive and hidden from plan output. The name of the field, eg: secret_access_key, will be the key
  in the `params` map in the api request.
  Credentials may not be specified in both locations and will cause an error. Changing from one location
  to a different credential configuration in the config will require an apply to update state.
  Structure is [documented below](#nested_sensitive_params).

* `location` -
  (Optional)
  The geographic location where the transfer config should reside.
  Examples: US, EU, asia-northeast1. The default value is US.

* `service_account_name` -
  (Optional)
  Optional service account name. If this field is set, transfer config will
  be created with this service account credentials. It requires that
  requesting user calling this API has permissions to act as this service account.

* `project` - (Optional) The ID of the project in which the resource belongs.
    If it is not provided, the provider project is used.


<a name="nested_schedule_options"></a>The `schedule_options` block supports:

* `disable_auto_scheduling` -
  (Optional)
  If true, automatic scheduling of data transfer runs for this
  configuration will be disabled. The runs can be started on ad-hoc
  basis using transferConfigs.startManualRuns API. When automatic
  scheduling is disabled, the TransferConfig.schedule field will
  be ignored.

* `start_time` -
  (Optional)
  Specifies time to start scheduling transfer runs. The first run will be
  scheduled at or after the start time according to a recurrence pattern
  defined in the schedule string. The start time can be changed at any
  moment. The time when a data transfer can be triggered manually is not
  limited by this option.

* `end_time` -
  (Optional)
  Defines time to stop scheduling transfer runs. A transfer run cannot be
  scheduled at or after the end time. The end time can be changed at any
  moment. The time when a data transfer can be triggered manually is not
  limited by this option.

<a name="nested_email_preferences"></a>The `email_preferences` block supports:

* `enable_failure_email` -
  (Required)
  If true, email notifications will be sent on transfer run failures.

<a name="nested_sensitive_params"></a>The `sensitive_params` block supports:

* `secret_access_key` -
  (Required)
  The Secret Access Key of the AWS account transferring data from.
  **Note**: This property is sensitive and will not be displayed in the plan.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource with format `{{name}}`

* `name` -
  The resource name of the transfer config. Transfer config names have the
  form projects/{projectId}/locations/{location}/transferConfigs/{configId}.
  Where configId is usually a uuid, but this is not required.
  The name is ignored when creating a transfer config.


## Timeouts

This resource provides the following
[Timeouts](/docs/configuration/resources.html#timeouts) configuration options:

- `create` - Default is 4 minutes.
- `update` - Default is 4 minutes.
- `delete` - Default is 4 minutes.

## Import


Config can be imported using any of these accepted formats:

```
$ terraform import google_bigquery_data_transfer_config.default {{name}}
```

## User Project Overrides

This resource supports [User Project Overrides](https://www.terraform.io/docs/providers/google/guides/provider_reference.html#user_project_override).
