<walkthrough-metadata>
  <meta name="title" content="Serverless Journey>: Extend your code to call Cloud APIs" />
  <meta name="description" content="Learn how to call Cloud APIs directly from your code and use service accounts to secure your app." />
  <meta name="keywords" content="deploy, containers, console, run" />
</walkthrough-metadata>

# Serverless Journey: Extend your code to call Cloud APIs

![Tutorial header image](https://storage.googleapis.com/gweb-cloudblog-publish/images/databases_2022_HTRs5Tr.max-700x700.jpg)

In this tutorial we'll learn how to extend the existing code to call Cloud APIs directly. Currently, the deployed application uses a library which contains a static set of jokes. Whenever the library is used it randomly selects a joke and returns it. After a while we will surely start to see the same jokes again and the only way to see new jokes is when a human would actually implement them in the library.

Luckily, there is a thing called _generative AI_ now. Google Cloud Vertex AI contains a Google-built, pre-trained, PaLM model which is a general-purpose, generative LLM that can be query with free-texts prompts to generate all sorts of text-based outputs. In this tutorial we'll implement the `model:predict` endpoint of Vertex AI to execute this model in order to new dad jokes in a generative matter.

Additionally, we'll learn a little bit about custom service accounts, IAM permissions and how to use the principle of least privilege to secure our services on Cloud Run.

<walkthrough-tutorial-difficulty difficulty="3"></walkthrough-tutorial-difficulty>

Estimated time:
<walkthrough-tutorial-duration duration="45"></walkthrough-tutorial-duration>

To get started, click **Start**.

## Project setup

First, let's make sure we got the correct project selected. Go ahead and select the provided project ID.

<walkthrough-project-setup billing="true"></walkthrough-project-setup>

Run the following to make sure all required APIs are enabled. Note that `aiplatform.googleapis.com` is added.

<walkthrough-enable-apis apis="cloudbuild.googleapis.com,
run.googleapis.com,aiplatform.googleapis.com,
artifactregistry.googleapis.com">
</walkthrough-enable-apis>

## About the different ways to call Google APIs

In the cloud we often need to implement API calls to interact with other services. Platforms like Google Cloud are a large collection of APIs; if we learn how to use them we can use the entire wealth of Google Cloud's functionality in our applications and build almost anything!

There are typically three different ways to interact with Google APIs programmatically and we should choose them in the following order:

1. **Cloud Client Libraries**: These are the recommended option. Cloud Client Libraries are SDK that you can use in a language-native, idiomatic style. They give you a high-level interface to the most important operation and allow you to quickly and comfortably get the job done. An example in the Go eco-system would be the `cloud.google.com/go/storage` package, which implements the most commonly used operations of Google Cloud Storage. Have a look at the [documentation of the package](https://pkg.go.dev/cloud.google.com/go/storage) and see how it respects and implements native language concepts like the `io.Writer` interface of the Go programming language.

2. **Google API Client Libraries**: Should no Cloud Client Library be available for what you are trying to accomplish you can fall back to using a Google API Client Library. These libraries are auto-generated and should be available for almost all Google APIs.

3. **Direct Implementation**: You can always choose to implement exposed APIs directly with your own client code. While most APIs are available as traditional RESTful services communicating JSON-serialized data, some APIs also implement [gRPC](https://grpc.io).

Have a look at [this page of the Google Cloud documentation](https://cloud.google.com/apis/docs/client-libraries-explained) to learn more about the differences between the available libraries.

Let's configure `gcloud` for the default regions:

```bash
# Set default locations for Cloud Run and Artifact Registry to europe-north1, Finland.

gcloud config set run/region europe-north1 
gcloud config set artifacts/location europe-north1
```

Okay, all set!

## Extending the code

In order to be able to replace the statically created jokes with jokes generated by Vertex AI, the code needs to be extended in the following ways:

1. **Create a client to execute the remote model**: The first steps is to safely instantiate the correct client type, that we are going to use later on to interact with the API. This type holds all the configuration for endpoints and handles authentication, too. Have a look at the [documentation of the package](https://pkg.go.dev/cloud.google.com/go/aiplatform@v1.45.0/apiv1beta1#NewPredictionClient) for `aiplatform.PredictionClient`. The client will automatically look for credentials according to the rules of [Google's Application Default Credentials scheme](https://cloud.google.com/docs/authentication/application-default-credentials).

2. **Form a request object**: Executing a call against `model:predict` on Vertex AI requires the caller to send an object with some configuration needed by the API to run the prediction. Take a look at the [documentation of the type](https://pkg.go.dev/cloud.google.com/go/aiplatform@v1.45.0/apiv1beta1/aiplatformpb#PredictRequest) `aiplatformpb.PredictRequest`. One field of the type is asking to supply and `Endpoint string`. In order to use the Google-built _PaLM 2 for Text_ model, we'll set this property like this:

```golang
endpoint := fmt.Sprintf(
	"projects/%s/locations/us-central1/publishers/google/models/text-bison",
	project,
)
```

Have a look at the [Vertex AI Model Garden](https://console.cloud.google.com/vertex-ai/publishers/google/model-garden/text-bison) to learn more about _PaLM 2 for Text_ (`text-bison`) and other available models.

The call also requires other inputs, like the actual text prompt and parameters to tune the output of the model.

4. **Execute the call against the API to run the model prediction**: Finally, we'll use the previously instantiated client to actually send the `aiplatformpb.PredictRequest` and execute the API call. Take a look at the [documentation of the method](https://pkg.go.dev/cloud.google.com/go/aiplatform@v1.45.0/apiv1beta1#PredictionClient.Predict) `Predict` on the client type.

Now, it's time to make some changes to the code.

<walkthrough-info-message>You may now attempt to **implement the above code changes yourself** for which you should have a good understanding of the Go programming language. If you choose to do so, you should stop reading now and give it your best shot.</walkthrough-info-message>

Alternatively, you can **use a prebuilt library** to accomplish the same. If that's more your cup of tea, go hit 'Next' and proceed to the next section.jjourney.svg

## Using a library to call the PaLM 2 for Text model

To get started, first have a look at the Go module `github.com/helloworlddan/tortuneai`. The module provides the package `github.com/helloworlddan/tortuneai/tortuneai` which implements the Cloud Client Library for Vertex AI to call the _PaLM2 for Text_ model. [Read through the code](https://github.com/helloworlddan/tortuneai/blob/main/tortuneai/tortuneai.go) of `tortuneai.HitMe` and see how it implements the aforementioned steps to interact with the API.
 
In order to use the package we need to first get the module like this:

```bash
go get github.com/helloworlddan/tortuneai 
```

Once that is completed, we can update
<walkthrough-editor-open-file filePath="cloudshell_open/serverless/main.go">
the main application source file main.go
</walkthrough-editor-open-file>
and change all references to the previously used package `tortune` with `tortuneai`.

Notice that the signature for `tortuneai.HitMe()` is different from the previous `tortune.HitMe()`. While the original function did not require any parameters, you are required to pass two `string` values into the new one: One with an actual text prompt for the PaLM 2 model and one with your Google Cloud project ID. Additionally, the function now returns multiple return values: a `string` containing the response from the API and an `error`. If everything goes well, the error will be `nil`, if not it will contain information about what went wrong.

Here is a possible implementation:

```golang 
joke, err := tortuneai.HitMe("tell me something about the year of the rabbit", "<walkthrough-project-id/>")
if err != nil {
    fmt.Fprintf(w, "error: %v\n", err)
    return
}
fmt.Fprint(w, joke)
```

Let's check if the modified code compiles by running it:

```bash
go run main.go
```

This recompiles and starts the web server. Let's check the application with the Web Preview <walkthrough-web-preview-icon></walkthrough-web-preview-icon> at the top right in Cloud Shell and see if we can successfully interact with the PaLM 2 model.

If you are satisfied you can focus the terminal again and terminate the web server with `Ctrl-C`.

It's good practice to clean up old dependencies from the `go.mod` file. You can do this automatically my running:

```bash
go mod tidy
```

If you like you can stay at this point for a moment, change the prompt (the first argument to `tortuneai.HitMe()`), re-run with `go run main.go` and use the Web Preview <walkthrough-web-preview-icon></walkthrough-web-preview-icon> to have a look at how the change in prompt affected the model's output.

## Creating a custom service account for the Cloud Run service

Good! The code changes we made seem to work, now it's time to deploy the changes to the cloud. 

When running the the code from Cloud Shell, the underlying implementation used Google ADC to find credentials. In this case it was using the credentials of the Cloud Shell user identity (yours).

Cloud Run can be configured to use a service account, which exposes credentials to the code running in your container. Your application can then make authenticated requests against Google APIs.

Per default, Cloud Run uses the [Compute Engine default service account](https://cloud.google.com/compute/docs/access/service-accounts#default_service_account). This service account has wide permissions and should generally be replaced by a service account with the least amount of permissions required to do whatever your service needs to do. The Compute Engine default service has a lot of permissions, but it does not have the required permissions to execute [the API call](
https://cloud.google.com/vertex-ai/docs/reference/rest/v1beta1/projects.locations.publishers.models/predict).

When it comes to identifying the correct IAM roles to attach to identities [this reference page on the IAM documentation](https://cloud.google.com/iam/docs/understanding-roles) is an extremely useful resources. On the page you can check the section for [Vertex AI roles](https://cloud.google.com/iam/docs/understanding-roles#vertex-ai-roles) and learn that [_Vertex AI User_](https://cloud.google.com/iam/docs/understanding-roles#aiplatform.user) (`roles/aiplatform.user`) is a suitable role for our Cloud Run service, because this role contains the permission `aiplatform.endpoints.predict`. If you are unsure which permission you require, you can always check the [API reference for the required operation](https://cloud.google.com/vertex-ai/docs/reference/rest/v1beta1/projects.locations.publishers.models/predict).

<walkthrough-info-message>**Note**: You could argue that the role _Vertex AI User_ has too many permissions for the Cloud Run service and a security-conscious person would probably agree with you. If you really wanted to make sure that the Cloud Run service only had least amount of privilege to execute the absolutely required permissions, you would have to create a [IAM custom role](https://cloud.google.com/iam/docs/creating-custom-roles) to achieve this.</walkthrough-info-message>

For now, we'll stick with _Vertex AI User_.

Next, let's create a brand new customer IAM service account like this:

```bash
gcloud iam service-accounts create tortune
```

We can bind the identified role to the various resource levels. For the sake of simplicity, let's attach it at the project level by executing:

```bash
gcloud projects add-iam-policy-binding $(gcloud config get-value project) \
    --member serviceAccount:tortune@$(gcloud config get-value project).iam.gserviceaccount.com  \
    --role "roles/aiplatform.user"
```

The service account will now be able to use all the permissions in _Vertex AI User_ on all resources in our current project. Finally, we need to deploy a need Cloud Run revision by updating the service configuration so that our Cloud Run service will use the newly-created service account:

```bash
gcloud run services update jokes \
    --service-account tortune@$(gcloud config get-value project).iam.gserviceaccount.com 
```

Now, that all IAM resources and configurations are in place, we can commit changes and push the to the remote repository to trigger a new Cloud Build execution to deploy changes in a CI/CD fashion, like this:

```bash
git add .
git commit -m 'upgrade to PaLM 2 for generative jokes'
git push origin main
```

Navigate to [Cloud Build's dashboard](https://console.cloud.google.com/cloud-build/dashboard) and click into the active build to monitor it's progress.

Once completed, you should be able to get fresh generated content by cURLing the endpoint of the Cloud Run Service:

```bash
curl $(gcloud run services describe jokes --format 'status.url')
```

Amazing!

## Summary

You now know how to call Google APIs directly from your code and understand how to secure your services with least-privilege service accounts.

<walkthrough-conclusion-trophy></walkthrough-conclusion-trophy>

```bash
cloudshell launch-tutorial .journey/03-operate.neos.md
```

