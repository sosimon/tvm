# Token Vending Machine

The vending machine generates temporary credentials to your AWS account - access and secret keys for APIs and a URL for the AWS console. The generated temporary credentials has access permissions that is the _intersection_ of

    1. The policy that is passed into the `GetFederationToken` call, and
    2. Policies that are attached to the IAM user whose credentials were used to all `GetFederationToken`.

See [docs](https://docs.aws.amazon.com/IAM/latest/UserGuide/id_credentials_temp_control-access_getfederationtoken.html) for more details. Currently, the policy passed into `GetFederationToken` is hardcoded to allow full access to EC2 and to manage SSH key pairs.

The temporary credentials are defaulted to expire after 15 minutes (900 seconds), which is the minimum session duration. A longer expiration can be specified, in seconds, using the `-x` flag, up to 36 hours (129600 seconds).

## Setup

Setup AWS credentials the same way you would for AWS CLI. See [SDK Configuration](http://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html#id2) and [AWS CLI Configuration](http://docs.aws.amazon.com/cli/latest/userguide/cli-chap-getting-started.html#config-settings-and-precedence) for more information.

## Usage

```
tvm [-u tempUsername] [-x sessionDuration]
```

`tempUsername` defaults to `temp-user` if none is specified.
`sessionDuration` is in seconds and defaults to 900 seconds if none are
specified.

`tvm` tries to read a policy file `policy.json` in the same directory. If the policy file is not present, it will default to allow all on all resources. The generated credentials will still end up having permissions that are the intersection of this and the permissions of the IAM user used to call this tool (you can't have more permissions than the IAM user used!) 

## Example

```
tvm -u foobar -x 3600
```

## References

[Creating a URL that Enables Federated Users to Access the AWS Management Console](http://docs.aws.amazon.com/IAM/latest/UserGuide/id_roles_providers_enable-console-custom-url.html)
