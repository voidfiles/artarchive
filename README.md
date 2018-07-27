# artfinder

Will find an index art entries

## Plan

- Download feeds output Feed Items
  - Feed Items should resolve the URLs
  - Find best content etc
- Convert Feed Items into Slides
- Download current Slide info for all Slides
- If slide has image, but doesn't have google photo url
  - Upload image to photo album
- Push slides back to s3


## What is needed

- Should download and index entries on to s3
  - For each entry
    - Find all the images, download, and push to s3
    - Rewrite image links in HTML
  - Save entries normalized to s3
- For each entry create an HTML page
  - /art/{year}/{month}/{day}/{entry-id}-{title}/index.html
- Re-write daily index for each date affected
- Re-write /index.html to reflect all dates known


## Contributing

1. Fork
2. Commit
3. Make sure tests pass
4. Create PR
5. Get code review

## Getting Started

This will install dependencies. You'll be able to run tests and
build the codebase afterwards.

```bash
make init
```

Run the tests

```bash
make test
```

To build

```bash
make build
```
