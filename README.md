# SendToKobo

With Mozilla deprecating Pocket, I needed a replacement to sync articles I want to read to my Kobo eReader.

**Send to Kobo** is a microservice architecture application that takes advantage of Kobo's Dropbox integration, allowing users to minimize articles (i.e. strip out ads) and upload them to their eReader.

## Tech Stack

**Send to Kobo** is hosted on GCP, in a way which maximizes utilization of their free-tier services.

* (TODO unimplemented) **Send to Kobo (Extension)** is a browser extension, allowing one-click web page uploading. Also available for Firefox.
* (TODO needs to have API endpoint updated but almost done) **Send to Kobo (Script)** is a user script intended for users of Safari, but usable with any script manager (e.g Userscripts on Safari, Tampermonkey on Chrome or Firefox, etc), allowing one-click upload functionality in the restricted Safari permissions ecosystem
* (TODO unimplemented) **Gateway** is the gateway API service (hosted with GCP API Gateway)
* (TODO unimplemented) **TokenManager** is the service responsible for managing user sessions, tokens, and login information (hosted on Compute Engine, with data stored in Firestore)
* **HtmlExtractor** uses readability.js to extract just the main content of a webpage (hosted on Cloud Run Functions)
* **HtmlToEpub** uses go-pandoc to convert the minimized webpage to the eReader native Epub format (hosted on Cloud Run)
* **FileUploader** uses Dropbox's HTTP API to upload files to the user's Dropbox account and sync them to the user's Kobo eReader

(TODO unimplemented) Not sure how to coordinate retries, associate particular calls with particular users, etc

All of this work is orchestrated using Cloud Pub/Sub and Cloud Workflows. (TODO maybe I should use GKE?)

(TODO unimplemented) CI is done via Github Actions, with source-code available under MIT license on Github [here](github.com/danielwolber-wood/SendToKobo). 

A monolithic, self-hostable version is also available as **Kobox** [here](github.com/danielwolber-wood/kobox), also under an MIT license.

(TODO unimplemented) Automatic docker upload

(TODO unimplemented) Automatic deployment of docker images on GCP

