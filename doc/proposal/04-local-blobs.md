# Local Blobs

Usually sources and resources are stored in some external storage, e.g. a docker image is stored in some OCI registry
and the *Component Descriptor* contains only the information how to access it. As an alternative *Component Repositories*
MAY provide the possibility to store technical artifacts together with the *Component Descriptors* in the
*Component Repository* itself as so-called *local blobs*. This allows to pack all component versions with their
technical artifacts in a *Component Repository* as a completely self-contained package. This is a typical requirement
if you need to deliver your product into a fenced landscape. This also allows storing e.g. configuration data together
with your *Component Descriptor*.

## Functions for Local Blobs

### UploadLocalBlob

**Description**: Allows uploading binary data. The binary data belong to a particular *Component Descriptor*
and can be referenced by the component descriptor in its *resources* or *sources* section.
*Component Descriptors* are not allowed to reference local blobs of other *Component Descriptors* in their resources.

When uploading a local blob, it is not REQUIRED that the corresponding *Component Descriptor* already exists.
Usually local blobs are uploaded first because it is not allowed to upload a *Component Descriptor* if its local
blobs not already exist.

The optional parameter *mediaType* provides information about the internal structure of the provided blob.

With the optional parameter *annotations* you could provide additional information about the blob. This information
could be used by the *Component Repository* itself or later if the local blob is stored again in some external
location, e.g. an OCI registry.

*LocalAccessInfo* provides the information how to access the blob data with the method *GetLocalBlob* (see below).

With the return value *globalAccessInfo*, the *Component Repository* could optionally provide an external reference to
the resource, e.g. if the blob contains the data of an OCI image it could provide an external OCI image reference.

**Inputs**:

- String name: Name of the *Component Descriptor*
- String version: Version of the *Component Descriptor*
- BinaryStream data: Binary stream containing the local blob data.
- String mediaType: media-type of the uploaded data (optional)
- map(string,string) annotations: Additional information about the uploaded artifact (optional)

**Outputs**:

- String localAccessInfo: The information how to access the source or resource as a *local blob*.
- String globalAccessInfo (optional): The information how to access the source or resource via a global reference.

**Errors**:

- invalidArgument: If one of the input parameters is empty or not valid
- repositoryError: If some error occurred in the *Component Repository*

**Example**:
Assume you want to upload an OCI image to your *Component Repository* with the *UploadLocalBlob* function with media type
*application/vnd.oci.image.manifest.v1+json* and the *annotations* "name: test/monitoring", and get the *localAccessInfo*:

```
"digest: sha256:b5733194756a0a4a99a4b71c4328f1ccf01f866b5c3efcb4a025f02201ccf623"
```

Then the entry in the *Component Descriptor* might look as follows. It is up to you, if you add the annotations
provided to the upload function and depends on the use case.

```
...
  resources:
  - name: example-image
    type: oci-image
    access:
      type: localOciBlob
      mediaType: application/vnd.oci.image.manifest.v1+json
      annotations:
          name: test/monitoring
      localAccess: "digest: sha256:b5733194756a0a4a99a4b71c4328f1ccf01f866b5c3efcb4a025f02201ccf623"
... 
```

The *Component Repository* could also provide some *globalAccessInfo* containing the location in an OCI registry:

```
imageReference: somePrefix/test/monitoring@sha:...
type: ociRegistry
```

An entry to this resource with this information in the *Component Descriptor* looks as the following:

```
...
  resources:
  - name: example-image
    type: oci-image
    access:
      type: localOciBlob
      mediaType: application/vnd.oci.image.manifest.v1+json
      annotations:
        name: test/monitoring
      localAccess: "digest: sha256:b5733194756a0a4a99a4b71c4328f1ccf01f866b5c3efcb4a025f02201ccf623"
      globalAccess: 
        imageReference: somePrefix/test/monitoring@sha:...
        type: ociRegistry
... 
```

### GetLocalBlob

**Description**: Fetches the binary data of a local blob. *blobIdentifier* is the *Component Repository* specific
access information you got when you uploaded the local blob.

**Inputs**:

- String name: Name of the *Component Descriptor*
- String version: Version of the *Component Descriptor*
- String blobIdentifier: Identifier of the local blob

**Outputs**:

- BinaryStream data: Binary stream containing the local blob data.

**Errors**:

- doesNotExist: If the local blob does not exist
- invalidArgument: If one of the input parameters is empty or invalid
- repositoryError: If some error occurred in the *Component Repository*

### DeleteLocalBlob

**Description**: Deletes a local blob. *blobIdentifier* is the *Component Repository* specific
information you got when you uploaded the local blob.

An error occurs if there is still an existing reference to the local blob.


**Inputs**:

- String name: Name of the *Component Descriptor*
- String version: Version of the *Component Descriptor*
- String blobIdentifier: Identifier of the local blob

**Outputs**:

**Errors**:

- doesNotExist: If the local blob does not exist
- existingReference: If the local blob is still referenced
- invalidArgument: If one of the input parameters is empty
- repositoryError: If some error occurred in the *Component Repository*