# Fiona
<img align="right" src="https://vignette.wikia.nocookie.net/muppet/images/2/22/Fiona.jpg/revision/latest/scale-to-width-down/200?cb=20081201050027">

## What is it?

Fiona is a http based service to conveniently set up users and a standard set of policies for a minio based S3 bucket. 
The main purpose is to provide a setup of separately available "folders" within a bucket for specific 
users to make simple object storage available for clients. 

The most important endpoint is the `/createuser` endpoint, described in [the API documentation](./API.md). 

The component is named after the Fiona Fraggle (https://muppet.fandom.com/wiki/Fiona_Fraggle). 

## Building Fiona

Fiona is a go application, using Go modules. Fiona has been developed with go version 1.13.6. 

```
git clone https://github.com/Skatteetaten/fiona.git
cd fiona
make
```

## S3 server - a prerequisite for testing and deployment

Fiona has been developed with a basic minio server for S3 server. Since the purpose of Fiona is to set up users on 
such a server, a running S3 server (minio) is needed to use Fiona.

## Deployment

There are some configuration needed for deploying and running Fiona

### Configuration settings

Fiona need to be configured to connect to the S3 server. This is done by environment variables. All variables have 
defined defaults, so for very basic testing, fiona can start without them, but only as long as the S3 server conforms 
to the defaults.

Here is a summary of the environment variables used by Fiona:

| Environment variable | Default | Description |
| ---| ---| ---|
| FIONA_S3_HOST | localhost | The host name of the S3 server |
| FIONA_S3_PORT | 9000 | The port of the S3 server |
| FIONA_S3_USESSL | false | Set to true if the S3 server uses SSL |
| FIONA_S3_REGION | us-east-1 | The region of the S3 server, also used for the bucket |
| FIONA_RANDOMPASS | false | Set to true if each user should get a separate password (recommended)|
| FIONA_DEFAULTPASS | S3userpass | The returned userpass if FIONA_RANDOMPASS is false |
| FIONA_ACCESS_KEY | aurora | Access key for the S3 server admin (recommended to override) |
| FIONA_SECRET_KEY | fragleberget | Access secret for the S3 server admin (recommended to override) |
| FIONA_DEBUG | false | Set to true to enable debug logging |
| FIONA_AURORATOKENLOCATION | ./aurora-token | The location of a file for authentication token see [the API](./API.md) for information |

### Aurora token

To authenticate endpoint requests, a token is used.  This token is stored in a file as indicated by the 
FIONA_AURORATOKENLOCATION configuration, and is mandatory for Fiona to work (an error will occur on startup if missing).

## Using Fiona - API

Fiona provides an http based API as a service.  [The API is described here](./API.md)

## Versioning

We use [Semantic versioning](http://semver.org/) for our release versions. 

## Authors

* **Rune Offerdal** - *Initial work*

## License

This project is licensed under the Apache 2.0 License - see the [LICENSE](./LICENSE) file for details
