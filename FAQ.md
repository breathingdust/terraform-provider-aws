# Frequently Asked Questions

## Who are the maintainers?

The HashiCorp Terraform AWS provider team is :

* Mary Curtrali - Product Manager [github](https://github.com/maryelizbeth) [twitter](https://twitter.com/marycutrali)
* Brian Flad - Engineering Lead [github](https://github.com/bflad) 
* Graham Davison - Engineer [github](https://github.com/gdavison)
* Angie Pinilla - Engineer [github](https://github.com/angie44)
* Simon Davis - Engineering Manager [github](https://github.com/maryelizbeth)
* Kerim Satirli - Developer Advocate [github](https://github.com/ksatirli)

## Why isn’t my PR merged yet?

Unfortunately, due to the volume of issues and new pull requests we receive, we are unable to give each one the full attention that we would like. We always focus on the contributions that provide the most value to the most community members.

## How do you decide what gets merged for each release?

The number one factor we look at when deciding what issues to look at are your reactions, comments, and upvotes on the issues or PR’s. The items with the most support are always on our radar, and we commit to keep the community updated on their status and potential timelines.

We also are investing time to improve the contributing experience by improving documentation, adding more linter coverage to ensure that incoming PR's can be in as good shape as possible. This will allow us to get through them quicker.

## How often do you release?

We release weekly on Thursday. We release often to ensure we can bring value to the community at a frequent cadence and to ensure we are in a good place to react to AWS region launches and service announcements.

## Backward Compatibility Promise

Our policy is described on the Terraform website [here](https://www.terraform.io/docs/extend/best-practices/versioning.html). 

## AWS just announced a new region, when will I see it in the provider.

Normally pretty quickly. We usually see the region appear within the `aws-go-sdk` within a couple days of the announcement. Depending on when it lands, we can often get it out within the current or following weekly release. Comparatively, adding support for a new  region in the S3 backend can take a little longer, as it is shipped as part of Terraform Core and not via the AWS Provider. 

Please note that this new region requires a manual process to enable in your account. Once enabled in the console, it takes a few minutes for everything to work properly.

If the region is not enabled properly, or the enablement process is still in progress, you may receive errors like these:

```
$ terraform apply

Error: error validating provider credentials: error calling sts:GetCallerIdentity: InvalidClientTokenId: The security token included in the request is invalid.
    status code: 403, request id: 142f947b-b2c3-11e9-9959-c11ab17bcc63

  on main.tf line 1, in provider "aws":
   1: provider "aws" {

To use this new region before support has been added to the Terraform AWS Provider, you can disable the provider's automatic region validation via:
provider "aws" {
  # ... potentially other configuration ...

  region                 = "af-south-1"
  skip_region_validation = true
}

```

## How can I help?

Great question, if you have contributed before check out issues with the `help-wanted` label. These are normally issues with support, but we are currently unable to field resources to work on. If you are just getting started, take a look at issues with the `good-first-issue` label. Items with these labels will always be given priority for response.

## How can I become a maintainer?

This is an area under active research. Stay tuned!