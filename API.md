# API for Fiona

Fiona is a http based admin service for setting up multi tenant minio users on a limited number of S3 buckets.  

## Endpoints

The endpoints use JSON both for POST payloads (body) and for returned information. 

Errors may return content as plain, non-JSON strings.  

### Create User with Policy for a Path

  Creates a user with a policy on a specific path for a bucket and returns access information.
  
  Precondition: The named bucket must exist

* **URL**

  /buckets/{bucketname}/paths/{path}/userpolicies/

* **Method:**
  
  `POST`
  
*  **URL Params**
    
   None

* **Data Params**

  Input is provided as JSON
  
  **Required**
  
  `"username": <username>`
  
  `"access": <list of access specifiers>`
  
  **Optional**
  
  None
  
  **Example**
  
  `{"username":"testuser", "access":["READ", "WRITE", "DELETE"]}`
  
* **Authorization**

  Yes, see [Access control](#access-control)

* **Success Response:**
  
  The user is created with access policy to the specified basepath to create, read and delete objects. 
  A JSON structure is returned with information necessary to use the S3 bucket.

  * **Code:** 201 CREATED <br />
    **Content:** `{"accessKey":"aUserName","secretKey":"someSecretKey","host":"https://localhost:9000"}`
 
* **Error Response:**

  * **Code:** 401 UNAUTHORIZED <br />
    **Content:** `Unauthorized`

  OR

  * **Code:** 422 UNPROCESSABLE ENTITY <br />
    **Content:** `Could not unmarshal body` or `Error creating user: Bucket does not exist`

  OR

  * **Code:** 401 FORBIDDEN <br />
    **Content:** `Missing required input`
  
  OR

  * **Code:** 400 BAD REQUEST <br />
    **Content:** `Could not read request body`

  
* **Sample Call:**

```
  curl -d '{"username":"testuser", "access":["READ", "WRITE", "DELETE"]}' -H 'Content-Type: application/json' -H 'Authorization: aurora-token token' http://localhost:8080/buckets/abucketname/paths/apath/userpolicies/
```
  
### List users

  Lists users policy name and status.

* **URL**

  /listusers

* **Method:**
  
  `GET`
  
*  **URL Params**
    
  None

* **Data Params**
  
  None
    
* **Authorization**

  Yes, see [Access control](#access-control)

* **Success Response:**
  
  A full list of users are returned with policy names and activity status. 

  * **Code:** 200 OK <br />
    **Content:** 
    `{"testuser":{"policyName":"RWDutvtestpath_137","status":"enabled"}}`
 
* **Error Response:**

  * **Code:** 401 UNAUTHORIZED <br />
    **Content:** `Unauthorized`
  
* **Sample Call:**

```
  curl -H 'Content-Type: application/json' -H 'Authorization: aurora-token token' http://localhost:9000/listusers
```

### Server info

  This is a passthrough call to the connected S3 (minio) server that provides status information.

* **URL**

  /serverinfo

* **Method:**
  
  `GET`
  
*  **URL Params**
    
  None

* **Data Params**
  
  None
    
* **Authorization**

  Yes, see [Access control](#access-control)

* **Success Response:**
  
  Status information for the connected S3 server is returned in JSON format.

  * **Code:** 200 OK <br />
    **Content:** 
    `{"mode":"online","deploymentID":"a274ba71-d4a5-448f-ad96-4132aadc3461","buckets":{"count":1},"objects":{},"usage":{},"services":{"vault":{"status":"disabled"},"ldap":{}},"backend":{"backendType":"FS"},"servers":[{"state":"ok","endpoint":"minio-aurora-dev.utv.paas.skead.no:80","uptime":174805,"version":"2020-02-20T22:51:23Z","commitID":"d4dcf1d7225a38ecf94abe7cbe7c69a93dc7c0b0","network":{"minio-aurora-dev.utv.paas.skead.no:80":"online"},"disks":[{"path":"/data","state":"ok","totalspace":246950133760,"usedspace":50852524032}]}]}`
 
* **Error Response:**

  * **Code:** 401 UNAUTHORIZED <br />
    **Content:** `Unauthorized`
  
* **Sample Call:**

```
  curl -H 'Content-Type: application/json' -H 'Authorization: aurora-token token' http://localhost:9000/serverinfo
```

## Access control
  
Fiona uses an HTTP Authorization request header for access control. 

### Syntax

`Authorization: <type> <credentials>`
 
* \<type\>

  `aurora-token`
  
* \<credentials\>

  a secret string stored with the application

## Management interface

Fiona provides a management-interface for health check and environment variables. The endpoints are made available on a separate port, 
usually 8081.

### Discovery endpoint

  The management interface has a discovery endpoint that provides the endpoints for health and env. Use this to discover 
  the other endpoints.

* **URL**

  /management

* **Method:**

  `GET`

*  **URL Params**

  None

* **Data Params**
  
  None

* **Authorization**

  None

* **Success Response:**

  Returns a JSON structure of the configured management interface endpoints.

  * **Code:** 200 OK <br />
    **Content:** `{"_links":{"env":{"href":"http://localhost:8081/env"},"health":{"href":"http://localhost:8081/health"}}}`

* **Sample Call:**

```
  curl -H 'Content-Type: application/json' http://localhost:8081/management
```
