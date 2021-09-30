# generate-metadata action

This action creates JSON file containing metadata information

## Usage

See [action.yaml](https://github.com/hashicorp/actions-generate-metadata/blob/main/action.yml)

### Basic usage to create metadata.json file

```yaml
- name: Generate metadata file
    uses: hashicorp/actions-generate-metadata@main
    id: execute
    with:
    repository: "consul-terraform-sync"
    version: 1.2.3
```

### Create metadata.json file using command in the version input
```yaml
- name: Generate metadata file
    uses: hashicorp/actions-generate-metadata@main
    id: execute
    with:
    repository: "consul-terraform-sync"
    version: "make version"
```

## Inputs
* **`repository`** - (required). The repository name for collecting the metadata
* **`version`** - (required). Indicates the version to be set in the metadata file. Can also accept the command which will set version(e.g "make version")
* **`filePath`** - (optional). Existing path that denotes the location of the metadata file to be created. The action will not create specified directory if it not exist. Default is set to Github action root path.
* **`repositoryOwner`** - (optional). The repository owner (organization or user). Default is set to "hashicorp" organization
* **`metadataFileName`** - (optional). The name of the file produced by the action. The generated file will have a JSON format. Default is set to "metadata.json"

Example command for the `version` input are provided [here](https://github.com/99/generate-metadata-action#create-metadatajson-file-using-command-in-the-version-input)

## Outputs
* `filepath` - The path where `metadataFileName` was created after action finished executing.

## Example of the generated metadata file
```json
{
"repository": "consul-terraform-sync",
"repositoryOwner": "hashicorp",
"sha": "4671a6594f2a2650f066489a4fbfe35c3b1e3d35",
"version": "1.0.3",
"buildWorkflowId": "1284662138"
}
```
