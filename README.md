# ServisBOT report list interface for emails

Simple endpoint that lists and gets emails from a specific s3 bucket 
The BOT processes the emails (from forms and direct emails), it pushes the result to a report s3 bucket and fires an event to the service to store report meta data in couchbase

## Update for openshift pipelines
- removed all references to GOCD
- updated memory limits for osp tasks
