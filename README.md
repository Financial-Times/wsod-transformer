# Wall St On Demand - Concept Transformer

[![CircleCI](https://circleci.com/gh/Financial-Times/wsod-transformer.svg?style=svg)](https://circleci.com/gh/Financial-Times/wsod-transformer)

Loads TME data for Wall St. On Demand financial instrument codes, and transforms the series to the internal UP json model.
The service exposes endpoints for getting all the series and for getting series by uuid.

# Usage

To get source code:
```
go get -u github.com/Financial-Times/wsod-transformer
```

To run:
```
$GOPATH/bin/wsod-transformer --port=8080 --base-url="http://localhost:8080/transformers/wsod/" --tme-base-url="https://tme.ft.com" --tme-username="user" --tme-password="pass" --token="token" --tme-taxonomy-name="WSODClassification"`

export|set PORT=8080  
export|set BASE_URL="http://localhost:8080/transformers/wsod/"  
export|set TME_BASE_URL="https://tme.ft.com"  
export|set TME_USERNAME="user"  
export|set TME_PASSWORD="pass"  
export|set TOKEN="token"  
export|set CACHE_FILE_NAME="cache.db"  
$GOPATH/bin/wsod-transformer  
```

### Docker

Docker build:
```
docker build -t coco/wsod-transformer .
```

To run:

```
docker run -ti --env BASE_URL=<base url> --env TME_BASE_URL=<structure service url> --env TME_USERNAME=<user> --env TME_PASSWORD=<pass> --env TOKEN=<token> --env CACHE_FILE_NAME=<file> --env "TME_TAXONOMY_NAME=WSODClassification" coco/wsod-transformer
```

# Endpoints

* `/transformers/wsod` - Get all WSOD as APIURLs
* `/transformers/wsod/{uuid}` - Get WSOD data of the given uuid
* `/transformers/wsod/__ids` - Get a stream of WSOD IDs in this format {id : uuid}
* `/transformers/wsod/__count` - Get count of WSOD
