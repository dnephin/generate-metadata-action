# generate-metadata-action

This action creates metadata.json file containing release related information


## Usage

See [action.yaml](action.yaml)

### Usage

```yaml
- name: Generate metadata file
  uses: HashiCorp-RelEng-Dev/generate-metadata-action@main
  id: execute
  with:
    repository: ${{ github.event.repository.name }}
    eventPath: "."
  
```

## Inputs
* `repository` - The remote repository you want the sha from
* `organization` - (optional) the organization. Default "hashicorp"

## Outputs
* `generated-filepath` - The metadata.json path