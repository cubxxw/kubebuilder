/*
Copyright 2023 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cronjob

const cronjobSpecExplaination = `

// +kubebuilder:docs-gen:collapse=Imports

/*
 First, let's take a look at our spec.  As we discussed before, spec holds
 *desired state*, so any "inputs" to our controller go here.

 Fundamentally a CronJob needs the following pieces:

 - A schedule (the *cron* in CronJob)
 - A template for the Job to run (the
 *job* in CronJob)

 We'll also want a few extras, which will make our users' lives easier:

 - A deadline for starting jobs (if we miss this deadline, we'll just wait till
   the next scheduled time)
 - What to do if multiple jobs would run at once (do we wait? stop the old one? run both?)
 - A way to pause the running of a CronJob, in case something's wrong with it
 - Limits on old job history

 Remember, since we never read our own status, we need to have some other way to
 keep track of whether a job has run.  We can use at least one old job to do
 this.

 We'll use several markers (` + "`" + `// +comment` + "`" + `) to specify additional metadata.  These
 will be used by [controller-tools](https://github.com/kubernetes-sigs/controller-tools) when generating our CRD manifest.
 As we'll see in a bit, controller-tools will also use GoDoc to form descriptions for
 the fields.
*/
`

const cronjobSpecStruct = `
	// schedule in Cron format, see https://en.wikipedia.org/wiki/Cron.
	// +kubebuilder:validation:MinLength=0
	// +required
	Schedule string` + " `" + `json:"schedule"` + "`" + `

	// startingDeadlineSeconds defines in seconds for starting the job if it misses scheduled
	// time for any reason.  Missed jobs executions will be counted as failed ones.
	// +optional
	// +kubebuilder:validation:Minimum=0
	StartingDeadlineSeconds *int64` + " `" + `json:"startingDeadlineSeconds,omitempty"` + "`" + `

	// concurrencyPolicy specifies how to treat concurrent executions of a Job.
	// Valid values are:
	// - "Allow" (default): allows CronJobs to run concurrently;
	// - "Forbid": forbids concurrent runs, skipping next run if previous run hasn't finished yet;
	// - "Replace": cancels currently running job and replaces it with a new one
	// +optional
	// +kubebuilder:default:=Allow
	ConcurrencyPolicy ConcurrencyPolicy` + " `" + `json:"concurrencyPolicy,omitempty"` + "`" + `

	// suspend tells the controller to suspend subsequent executions, it does
	// not apply to already started executions.  Defaults to false.
	// +optional
	Suspend *bool` + " `" + `json:"suspend,omitempty"` + "`" + `

	// jobTemplate defines the job that will be created when executing a CronJob.
	// +required
	JobTemplate batchv1.JobTemplateSpec` + " `" + `json:"jobTemplate"` + "`" + `

	// successfulJobsHistoryLimit defines the number of successful finished jobs to retain.
	// This is a pointer to distinguish between explicit zero and not specified.
	// +optional
	// +kubebuilder:validation:Minimum=0
	SuccessfulJobsHistoryLimit *int32` + " `" + `json:"successfulJobsHistoryLimit,omitempty"` + "`" + `

	// failedJobsHistoryLimit defines the number of failed finished jobs to retain.
	// This is a pointer to distinguish between explicit zero and not specified.
	// +optional
	// +kubebuilder:validation:Minimum=0
	FailedJobsHistoryLimit *int32` + " `" + `json:"failedJobsHistoryLimit,omitempty"` + "`" + `
}

/*
 We define a custom type to hold our concurrency policy.  It's actually
 just a string under the hood, but the type gives extra documentation,
 and allows us to attach validation on the type instead of the field,
 making the validation more easily reusable.
*/

// ConcurrencyPolicy describes how the job will be handled.
// Only one of the following concurrent policies may be specified.
// If none of the following policies is specified, the default one
// is AllowConcurrent.
// +kubebuilder:validation:Enum=Allow;Forbid;Replace
type ConcurrencyPolicy string

const (
	// AllowConcurrent allows CronJobs to run concurrently.
	AllowConcurrent ConcurrencyPolicy = "Allow"

	// ForbidConcurrent forbids concurrent runs, skipping next run if previous
	// hasn't finished yet.
	ForbidConcurrent ConcurrencyPolicy = "Forbid"

	// ReplaceConcurrent cancels currently running job and replaces it with a new one.
	ReplaceConcurrent ConcurrencyPolicy = "Replace"
)

/*
 Next, let's design our status, which holds observed state.  It contains any information
 we want users or other controllers to be able to easily obtain.

 We'll keep a list of actively running jobs, as well as the last time that we successfully
 ran our job.  Notice that we use` + " `" + `metav1.Time` + "`" + ` instead of` + " `" + `time.Time` + "`" + ` to get the stable
 serialization, as mentioned above.
*/`

const cronjobList = `

	// active defines a list of pointers to currently running jobs.
	// +optional
	// +listType=atomic
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=10
	Active []corev1.ObjectReference` + " `" + `json:"active,omitempty"` + "`" + `

	// lastScheduleTime defines when was the last time the job was successfully scheduled.
	// +optional
	LastScheduleTime *metav1.Time` + " `" + `json:"lastScheduleTime,omitempty"` + "`" + `
`

const docCommentStatusSub = `
/*
 Finally, we have the rest of the boilerplate that we've already discussed.
 As previously noted, we don't need to change this, except to mark that
 we want a status subresource, so that we behave like built-in kubernetes types.
*/

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
`
